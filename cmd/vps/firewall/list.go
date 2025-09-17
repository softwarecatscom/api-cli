package firewall

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	//"github.com/hostinger/api-cli/api"
	"github.com/hostinger/api-cli/client"
	"github.com/hostinger/api-cli/output"
	"github.com/hostinger/api-cli/utils"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/securityprovider"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "Get firewall list",
	Long:  `This endpoint retrieves a list of all firewalls available.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Create raw client to bypass the broken generated parsing
		rawClient, err := createRawClient()
		if err != nil {
			log.Fatal(err)
		}

		// Call the raw HTTP endpoint
		resp, err := rawClient.VPSGetFirewallListV1(context.TODO(), firewallListRequestParameters(cmd))
		if err != nil {
			log.Fatal(err)
		}
		defer func() { _ = resp.Body.Close() }()

		// Read the raw response body
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		// Handle error status codes
		if resp.StatusCode >= http.StatusBadRequest {
			output.Format(cmd, bodyBytes, resp.StatusCode)
			return
		}

		// Validate that we can parse the response as a direct array
		// This confirms the API returns the correct format
		var firewalls []client.VPSV1FirewallFirewallResource
		if err := json.Unmarshal(bodyBytes, &firewalls); err != nil {
			log.Fatalf("Failed to parse firewall list response: %v", err)
		}

		// Pass the raw JSON to output formatter
		output.Format(cmd, bodyBytes, resp.StatusCode)
	},
}

func init() {
	ListCmd.Flags().IntP("page", "", 1, "Page number")
}

func createRawClient() (client.ClientInterface, error) {
	// Replicate the logic from api.Request() but return raw client
	bearerToken, err := securityprovider.NewSecurityProviderBearerToken(viper.GetString("api_token"))
	if err != nil {
		return nil, err
	}

	var apiUrl = viper.GetString("api_url")
	if apiUrl == "" {
		apiUrl = "https://developers.hostinger.com"
	}

	return client.NewClient(
		apiUrl,
		client.WithRequestEditorFn(bearerToken.Intercept),
		client.WithRequestEditorFn(addUserAgent),
	)
}

func addUserAgent(_ context.Context, req *http.Request) error {
	req.Header.Set("User-Agent", "hostinger-cli/"+utils.CLIVersion.String(false))
	return nil
}

func firewallListRequestParameters(cmd *cobra.Command) *client.VPSGetFirewallListV1Params {
	pageId, _ := cmd.Flags().GetInt("page")

	return &client.VPSGetFirewallListV1Params{
		Page: utils.IntPtrOrNil(pageId),
	}
}
