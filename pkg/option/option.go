package option

import (
	"os"

	"github.com/CDN-Security/Echo/pkg/version"
	"github.com/jessevdk/go-flags"
)

type Option struct {
	ConfigFilePath string `short:"c" long:"config" description:"Path to the config file" required:"true" default:"config.yaml"`
	Version        func() `long:"version" description:"print version and exit" json:"-"`
}

var Opt Option

func init() {
	Opt.Version = version.PrintVersion
	if _, err := flags.Parse(&Opt); err != nil {
		os.Exit(1)
	}
}
