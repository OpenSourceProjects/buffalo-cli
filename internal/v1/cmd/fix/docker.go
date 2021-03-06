package fix

import (
	"context"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	bufcli "github.com/gobuffalo/buffalo-cli"
	"github.com/gobuffalo/genny"
)

func fixDocker(r *Runner) error {
	app := r.App
	if !app.WithDocker {
		return nil
	}
	fmt.Println("~~~ Upgrading Dockerfile ~~~")
	run := genny.WetRunner(context.Background())
	run.WithRun(func(r *genny.Runner) error {
		dk, err := r.FindFile(filepath.Join(app.Root, "Dockerfile"))
		if err != nil {
			return err
		}

		ex := regexp.MustCompile(`(v[0-9.][\S]+)`)
		lines := strings.Split(dk.String(), "\n")
		for i, l := range lines {
			if strings.HasPrefix(strings.ToLower(l), "from gobuffalo/buffalo") {
				l = ex.ReplaceAllString(l, bufcli.Version)
				lines[i] = l
			}
		}
		return r.File(genny.NewFileS(dk.Name(), strings.Join(lines, "\n")))
	})
	return run.Run()
}
