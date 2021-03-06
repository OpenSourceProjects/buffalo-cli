package plugins

import (
	"github.com/gobuffalo/buffalo-cli/internal/v1/plugins/plugcmds"
	"github.com/spf13/cobra"
)

// Available used to manage all of the available commands
// for the plugins
var Available = plugcmds.NewAvailable()

// PluginsCmd is the "root" command for the plugin features.
var PluginsCmd = &cobra.Command{
	Use:   "plugins",
	Short: "tools for working with buffalo plugins",
}

func init() {
	PluginsCmd.AddCommand(addCmd)
	PluginsCmd.AddCommand(listCmd)
	PluginsCmd.AddCommand(generateCmd)
	PluginsCmd.AddCommand(removeCmd)
	PluginsCmd.AddCommand(installCmd)
	PluginsCmd.AddCommand(cacheCmd)

	Available.Add("generate", generateCmd)
	Available.ListenFor("buffalo:setup:.+", Listen)
}
