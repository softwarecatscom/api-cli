package api

import (
	"github.com/hostinger/api-cli/client"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/securityprovider"
	"github.com/spf13/viper"
	"log"
)

func Request() client.ClientWithResponsesInterface {
	bearerToken, err := securityprovider.NewSecurityProviderBearerToken(viper.GetString("api_token"))
	if err != nil {
		panic(err)
	}

	var apiUrl = viper.GetString("api_url")
	if apiUrl == "" {
		apiUrl = "https://developers.hostinger.com"
	}

	c, err := client.NewClientWithResponses(apiUrl, client.WithRequestEditorFn(bearerToken.Intercept))
	if err != nil {
		log.Fatal(err)
	}

	return c
}
