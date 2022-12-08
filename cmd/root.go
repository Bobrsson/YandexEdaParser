package cmd

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"YandexEdaParser/manager"
	"YandexEdaParser/structs"
)

var cfgFile string
var config *structs.Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:               "YandexEdaParser",
	Short:             "Приложение для парсинга Яндекс Еды",
	Long:              ``,
	PersistentPreRunE: initConfig,
	Run: func(cmd *cobra.Command, args []string) {
		manager.ServerRun(*config)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is cfg/config.yaml)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig(_ *cobra.Command, _ []string) (err error) {
	if config, err = manager.NewConfig(cfgFile); err != nil {
		return errors.Wrap(err, "Config loading error: %s")
		return err
	}

	return nil
}
