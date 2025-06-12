package output

import (
	"encoding/json"
	"fmt"
	"github.com/pterm/pterm"
	"sort"
	"strconv"
)

type NestedData struct {
	FieldName string
	Data      interface{}
	IsArray   bool
	Count     int
}

func printTable(body []byte) {
	err := displayJSONAsTables(body)
	if err != nil {
		pterm.Error.Printf("Error displaying JSON: %v\n", err)
		return
	}
}

func displayJSONAsTables(body []byte) error {
	var data interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	switch v := data.(type) {
	case []interface{}:
		// Array of objects
		return displayArrayOfObjects(v)
	case map[string]interface{}:
		// Single object
		return displaySingleObject(v, 1, 1)
	default:
		return fmt.Errorf("unsupported JSON structure: expected object or array of objects")
	}
}

func displayArrayOfObjects(objects []interface{}) error {
	// Check if any object has nested fields
	hasNestedFields := false
	for _, obj := range objects {
		objMap, ok := obj.(map[string]interface{})
		if !ok {
			return fmt.Errorf("array item is not an object")
		}

		if objectHasNestedFields(objMap) {
			hasNestedFields = true
			break
		}
	}

	if !hasNestedFields {
		// Display all records in a single table
		return displayAllRecordsInSingleTable(objects)
	}

	// Display each record separately (original behavior for nested objects)
	for i, obj := range objects {
		objMap, ok := obj.(map[string]interface{})
		if !ok {
			return fmt.Errorf("array item %d is not an object", i)
		}

		if err := displaySingleObject(objMap, i+1, len(objects)); err != nil {
			return err
		}

		// Add spacing between records (except for the last one)
		if i < len(objects)-1 {
			fmt.Println()
		}
	}
	return nil
}

func objectHasNestedFields(obj map[string]interface{}) bool {
	for _, value := range obj {
		if value == nil {
			continue // Skip null values
		}
		switch v := value.(type) {
		case map[string]interface{}:
			if len(v) > 0 { // Only count non-empty objects
				return true
			}
		case []interface{}:
			if len(v) > 0 { // Only count non-empty arrays
				return true
			}
		}
	}
	return false
}

func displayAllRecordsInSingleTable(objects []interface{}) error {
	if len(objects) == 0 {
		return nil
	}

	// Collect all unique keys from all objects
	allKeys := make(map[string]bool)
	for _, obj := range objects {
		objMap := obj.(map[string]interface{})
		for key := range objMap {
			allKeys[key] = true
		}
	}

	// Convert to sorted slice
	keys := make([]string, 0, len(allKeys))
	for key := range allKeys {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Create table data
	tableData := [][]string{keys} // Header row

	// Add each object as a row
	for _, obj := range objects {
		objMap := obj.(map[string]interface{})
		row := make([]string, len(keys))

		for i, key := range keys {
			if value, exists := objMap[key]; exists {
				row[i] = formatValue(value)
			} else {
				row[i] = ""
			}
		}
		tableData = append(tableData, row)
	}

	// Render table
	return pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
}

func displaySingleObject(obj map[string]interface{}, recordNum, totalRecords int) error {
	if totalRecords > 1 {
		headerText := fmt.Sprintf("Record %d of %d", recordNum, totalRecords)
		pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgGreen)).
			WithTextStyle(pterm.NewStyle(pterm.FgBlack)).
			Println(headerText)
	}

	// Separate main fields from nested data
	mainFields := make(map[string]interface{})
	nestedData := []NestedData{}

	for key, value := range obj {
		switch v := value.(type) {
		case map[string]interface{}:
			// Nested object
			nestedData = append(nestedData, NestedData{
				FieldName: key,
				Data:      v,
				IsArray:   false,
				Count:     0,
			})
		case []interface{}:
			// Array
			nestedData = append(nestedData, NestedData{
				FieldName: key,
				Data:      v,
				IsArray:   true,
				Count:     len(v),
			})
		default:
			// Simple field
			mainFields[key] = value
		}
	}

	// Display main fields table
	if len(mainFields) > 0 {
		if err := displayMainFieldsTable(mainFields); err != nil {
			return err
		}
	}

	// Display nested data tables
	for _, nested := range nestedData {
		//fmt.Println() // Add spacing before nested tables
		if err := displayNestedTable(nested); err != nil {
			return err
		}
	}

	return nil
}

func displayMainFieldsTable(fields map[string]interface{}) error {
	// Get sorted keys for consistent ordering
	keys := make([]string, 0, len(fields))
	for key := range fields {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Create table data
	tableData := [][]string{keys} // Header row
	values := make([]string, len(keys))

	for i, key := range keys {
		values[i] = formatValue(fields[key])
	}
	tableData = append(tableData, values)

	// Render table
	return pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
}

func displayNestedTable(nested NestedData) error {
	// Create smaller header
	var headerText string
	if nested.IsArray {
		headerText = fmt.Sprintf("%s (%d items)", nested.FieldName, nested.Count)
	} else {
		headerText = nested.FieldName
	}

	pterm.DefaultSection.Print(headerText)

	switch data := nested.Data.(type) {
	case map[string]interface{}:
		return displayNestedObject(data)
	case []interface{}:
		return displayNestedArray(data)
	default:
		return fmt.Errorf("unexpected nested data type for field %s", nested.FieldName)
	}
}

func displayNestedObject(obj map[string]interface{}) error {
	// Check if this object has nested structures
	mainFields := make(map[string]interface{})
	hasNested := false

	for key, value := range obj {
		switch value.(type) {
		case map[string]interface{}, []interface{}:
			hasNested = true
		default:
			mainFields[key] = value
		}
	}

	if hasNested {
		// Treat as a full object with potential nesting
		return displaySingleObjectContent(obj)
	} else {
		// Simple object, display as key-value table
		return displayMainFieldsTable(obj)
	}
}

func displayNestedArray(arr []interface{}) error {
	for i, item := range arr {
		if i > 0 {
			fmt.Println() // Add spacing between array items
		}

		switch v := item.(type) {
		case map[string]interface{}:
			if err := displaySingleObjectContent(v); err != nil {
				return err
			}
		default:
			// Simple array items
			return displaySimpleArray(arr)
		}
	}
	return nil
}

func displaySimpleArray(arr []interface{}) error {
	// For arrays of simple values, display as a single column table
	tableData := [][]string{{"Value"}} // Header

	for _, item := range arr {
		tableData = append(tableData, []string{formatValue(item)})
	}

	return pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
}

func displaySingleObjectContent(obj map[string]interface{}) error {
	// Similar to displaySingleObject but without the main record header
	mainFields := make(map[string]interface{})
	nestedData := []NestedData{}

	for key, value := range obj {
		switch v := value.(type) {
		case map[string]interface{}:
			nestedData = append(nestedData, NestedData{
				FieldName: key,
				Data:      v,
				IsArray:   false,
				Count:     0,
			})
		case []interface{}:
			nestedData = append(nestedData, NestedData{
				FieldName: key,
				Data:      v,
				IsArray:   true,
				Count:     len(v),
			})
		default:
			mainFields[key] = value
		}
	}

	// Display main fields
	if len(mainFields) > 0 {
		if err := displayMainFieldsTable(mainFields); err != nil {
			return err
		}
	}

	// Display nested data
	for _, nested := range nestedData {
		//fmt.Println()
		if err := displayNestedTable(nested); err != nil {
			return err
		}
	}

	return nil
}

func formatValue(value interface{}) string {
	if value == nil {
		return "null"
	}

	switch v := value.(type) {
	case string:
		return v
	case bool:
		return strconv.FormatBool(v)
	case float64:
		// Check if it's actually an integer
		if v == float64(int64(v)) {
			return strconv.FormatInt(int64(v), 10)
		}
		return strconv.FormatFloat(v, 'f', -1, 64)
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	default:
		// For complex types, use JSON representation
		if jsonBytes, err := json.Marshal(v); err == nil {
			return string(jsonBytes)
		}
		return fmt.Sprintf("%v", v)
	}
}
