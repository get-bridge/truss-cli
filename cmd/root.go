package cmd

import (
	"os"

	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "truss-cli",
	Short: "A CLI for use with Bridge Truss",
	// Long: `TODO`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Errorln(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.truss.yaml)")

	rootCmd.PersistentFlags().StringP("env", "e", "", "The environment to target")
	rootCmd.PersistentFlags().StringP("region", "g", "", "The region to target")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			log.Errorln(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".truss")
	}

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Errorln("Error loading config: ", err)
		}
	}
}
