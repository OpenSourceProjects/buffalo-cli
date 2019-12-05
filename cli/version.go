package cli

import (
	"context"
	"encoding/json"
	"fmt"

	bufcli "github.com/gobuffalo/buffalo-cli"
	"github.com/gobuffalo/buffalo-cli/internal/cmdx"
)

type versionCmd struct {
	*Buffalo
	help bool
	json bool
}

func (vc *versionCmd) Name() string {
	return "version"
}

func (vc *versionCmd) Main(ctx context.Context, args []string) error {
	flags := cmdx.NewFlagSet("buffalo info", cmdx.Stderr(ctx))
	flags.BoolVarP(&vc.help, "help", "h", false, "print this help")
	flags.BoolVarP(&vc.json, "json", "j", false, "Print information in json format")
	if err := flags.Parse(args); err != nil {
		return err
	}

	if vc.help {
		flags.Usage()
		return nil
	}

	if !vc.json {
		fmt.Fprintln(vc.Stdout, bufcli.Version)
		return nil
	}

	enc := json.NewEncoder(vc.Stdout)
	enc.SetIndent("", "    ")
	return enc.Encode(map[string]string{
		"version": bufcli.Version,
	})

}
