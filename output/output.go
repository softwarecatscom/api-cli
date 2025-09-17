package output

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
)

func Format(cmd *cobra.Command, body []byte, statusCode int) {
	format, err := cmd.Flags().GetString("format")
	if err != nil {
		log.Fatal(err)
	}

	if statusCode >= http.StatusBadRequest {
		fmt.Println(string(body))
		os.Exit(1)
	}
	if format == "" {
		val := viper.GetString("format")
		if val != "" {
			format = val
		}
	}
	switch format {
	default:
		printTable(body)
	case "tree":
		printTreeEnhanced(body)
	case "json":
		fmt.Println(string(body))
	}
}
