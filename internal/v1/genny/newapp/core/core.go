package core

import (
	"github.com/gobuffalo/buffalo-cli/internal/v1/genny/ci"
	"github.com/gobuffalo/buffalo-cli/internal/v1/genny/docker"
	"github.com/gobuffalo/buffalo-cli/internal/v1/genny/plugins/install"
	"github.com/gobuffalo/buffalo-cli/internal/v1/genny/refresh"
	"github.com/gobuffalo/buffalo-cli/internal/v1/plugins/plugdeps"
	pop "github.com/gobuffalo/buffalo-pop/genny/newapp"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/gogen/gomods"
	"github.com/gobuffalo/meta"
	"github.com/markbates/errx"
)

// New generator for creating a Buffalo application
func New(opts *Options) (*genny.Group, error) {
	gg := &genny.Group{}

	// add the root generator
	g, err := rootGenerator(opts)
	if err != nil {
		return gg, err
	}
	gg.Add(g)

	app := opts.App

	g, err = gomods.Init(app.PackagePkg, app.Root)
	if err != nil {
		return gg, err
	}
	gg.Add(g)

	plugs, err := plugdeps.List(app)
	if err != nil && (errx.Unwrap(err) != plugdeps.ErrMissingConfig) {
		return nil, err
	}

	if opts.Docker != nil {
		// add the docker generator
		g, err = docker.New(opts.Docker)
		if err != nil {
			return gg, err
		}
		gg.Add(g)
	}

	if opts.Pop != nil {
		// add the pop generator
		gg2, err := pop.New(opts.Pop)
		if err != nil {
			return gg, err
		}
		gg.Merge(gg2)

		// add the plugin
		plugs.Add(plugdeps.Plugin{
			Binary: "buffalo-pop",
			GoGet:  "github.com/gobuffalo/buffalo-pop",
		})
	}

	if opts.CI != nil {
		// add the CI generator
		g, err = ci.New(opts.CI)
		if err != nil {
			return gg, err
		}
		gg.Add(g)
	}

	if opts.Refresh != nil {
		g, err = refresh.New(opts.Refresh)
		if err != nil {
			return gg, err
		}
		gg.Add(g)
	}

	// ---

	// install all of the plugins
	iopts := &install.Options{
		App:     app,
		Plugins: plugs.List(),
	}
	if app.WithSQLite {
		iopts.Tags = meta.BuildTags{"sqlite"}
	}

	ig, err := install.New(iopts)
	if err != nil {
		return gg, err
	}
	gg.Merge(ig)

	return gg, nil
}
