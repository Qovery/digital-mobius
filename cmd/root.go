package cmd

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "do-k8s-replace-notready-nodes",
	Short: "Automatically recycle Kubernetes NotReady nodes",
	Long: `Automatically recycle Kubernetes NotReady nodes on DigitalOcean

Digital Ocean supports the replacement of VM when something goes wrong. However (even on DOKS),
when a node is in a NotReady state for a long time, nothing happens as Kubelet is not responding
anymore and Kubernetes doesn't know what to do. This is why do-k8s-replace-notready-nodes is born.`,
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
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "level", "info", "set log level")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".do-k8s-replace-notready-nodes" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".do-k8s-replace-notready-nodes")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func setLogLevel() error {
	// set log level
	logLevel, _ := rootCmd.Flags().GetString("level")

	lvl, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return err
	}

	logrus.SetLevel(lvl)

	// use timestamp
	formatter := &logrus.TextFormatter{
		FullTimestamp: true,
	}
	logrus.SetFormatter(formatter)
	return nil
}