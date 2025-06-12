package virtual_machines

import (
	"context"
	"github.com/hostinger/api-cli/api"
	"github.com/hostinger/api-cli/client"
	"github.com/hostinger/api-cli/output"
	"github.com/spf13/cobra"
	"log"
)

var SetHostnameCmd = &cobra.Command{
	Use:   "set-hostname <virtual machine ID> <hostname>",
	Short: "Set hostname",
	Args:  cobra.MatchAll(cobra.ExactArgs(2)),
	Run: func(cmd *cobra.Command, args []string) {
		r, err := api.Request().VPSSetHostnameV1WithResponse(
			context.TODO(),
			17992,
			client.VPSV1VirtualMachineHostnameUpdateRequest{
				Hostname: args[0],
			},
		)

		if err != nil {
			log.Fatal(err)
		}

		output.Format(cmd, r.Body, r.StatusCode())
	},
}
