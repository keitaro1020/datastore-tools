package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Options struct {
	OptProject   string
	OptKind      string
	OptNamespace string
	OptKeyFile   string
	OptFilter    string
	OptCount     bool
	OptTable     bool
}

var (
	o = &Options{}
)

var cfgFile string

var RootCmd = &cobra.Command{
	Use:           "datastore-tools",
	Short:         "CLI for Google Cloud Datastore",
	SilenceErrors: true,
	SilenceUsage:  false,
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default $HOME/.datastore-tools.yml)")

	viper.BindPFlag("url", RootCmd.PersistentFlags().Lookup("url"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".datastore-tools")
	viper.AddConfigPath("$HOME")
	viper.AutomaticEnv()

	viper.ReadInConfig()
}
