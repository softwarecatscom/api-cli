package virtual_machines

import (
	"context"
	"github.com/hostinger/api-cli/api"
	"github.com/hostinger/api-cli/output"
	"github.com/spf13/cobra"
	"log"
	"strconv"
)

var StopCmd = &cobra.Command{
	Use:   "stop <virtual machine ID>",
	Short: "Stop a virtual machine",
	Args:  cobra.MatchAll(cobra.ExactArgs(1)),
	Run: func(cmd *cobra.Command, args []string) {
		vmId, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatal(err)
		}

		r, err := api.Request().VPSStopVirtualMachineV1WithResponse(context.TODO(), vmId)
		if err != nil {
			log.Fatal(err)
		}

		output.Format(cmd, r.Body, r.StatusCode())
	},
}
