package config

import (
	"errors"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const EnvPrefix = "WATCH_"

func Load(cmd *cobra.Command, args []string) (*Config, error) {
	conf, ok := FromContext(cmd.Context())
	if !ok {
		panic("command missing config")
	}
	return conf, conf.Load(cmd, args)
}

func (c *Config) Load(cmd *cobra.Command, args []string) error {
	var errs []error
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if !f.Changed {
			if val, ok := os.LookupEnv(EnvName(f.Name)); ok {
				if err := f.Value.Set(val); err != nil {
					errs = append(errs, err)
				}
			}
		}
	})
	InitLog(c)

	c.Domains = append(c.Domains, args...)
	if len(c.Domains) == 0 && c.Completion == "" {
		return errors.New("missing domain")
	}

	return errors.Join(errs...)
}

func EnvName(name string) string {
	name = strings.ToUpper(name)
	name = strings.ReplaceAll(name, "-", "_")
	return EnvPrefix + name
}
