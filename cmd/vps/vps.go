package vps

import (
	"github.com/hostinger/api-cli/cmd/vps/actions"
	"github.com/hostinger/api-cli/cmd/vps/backups"
	"github.com/hostinger/api-cli/cmd/vps/data_centers"
	"github.com/hostinger/api-cli/cmd/vps/firewall"
	"github.com/hostinger/api-cli/cmd/vps/monarx"
	"github.com/hostinger/api-cli/cmd/vps/post_install_scripts"
	"github.com/hostinger/api-cli/cmd/vps/ptr"
	"github.com/hostinger/api-cli/cmd/vps/public_keys"
	"github.com/hostinger/api-cli/cmd/vps/templates"
	"github.com/hostinger/api-cli/cmd/vps/virtual_machines"
	"github.com/spf13/cobra"
)

var VpsGroupCmd = &cobra.Command{
	Use:   "vps",
	Short: "VPS management",
	Long:  ``,
}

func init() {
	VpsGroupCmd.AddCommand(actions.GroupCmd)
	VpsGroupCmd.AddCommand(backups.GroupCmd)
	VpsGroupCmd.AddCommand(data_centers.GroupCmd)
	VpsGroupCmd.AddCommand(firewall.GroupCmd)
	VpsGroupCmd.AddCommand(ptr.GroupCmd)
	VpsGroupCmd.AddCommand(monarx.GroupCmd)
	VpsGroupCmd.AddCommand(post_install_scripts.GroupCmd)
	VpsGroupCmd.AddCommand(public_keys.GroupCmd)

	VpsGroupCmd.AddCommand(virtual_machines.GroupCmd)
	VpsGroupCmd.AddCommand(templates.GroupCmd)
}
