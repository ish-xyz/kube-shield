package cmd

import (
	"fmt"
	"os"

	"github.com/RedLabsPlatform/kube-shield/pkg/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	printVersion bool
	Version      = "unset" // set at build time
	Cmd          = &cobra.Command{
		Use:   "kube-shield",
		Short: "kube-shield is a Kubernetes Admission Controller",
		Run:   start,
	}
)

func Execute() error {
	return Cmd.Execute()
}

func init() {
	// Non Bound Flags
	Cmd.Flags().BoolVarP(&printVersion, "version", "v", false, "Print version of kube-shield")
	Cmd.Flags().StringP("config", "c", "", "Kube-shield configuration file")

	// Bound Flags
	Cmd.Flags().StringP("policies", "p", "", "Path to the directory with the policies")
	Cmd.Flags().StringP("metricsAddress", "m", "", "Address where the metrics are exposed")
	Cmd.Flags().BoolP("registerWebhook", "r", true, "create ValidatingWebhookConfiguration resource in the current Kubernetes")
	Cmd.Flags().BoolP("debug", "d", false, "Path to the directory with the policies")

	// Required flags
	Cmd.MarkFlagRequired("policies")

	viper.BindPFlag("policies", Cmd.Flags().Lookup("policies"))
	viper.BindPFlag("debug", Cmd.Flags().Lookup("debug"))
}

// Start the admission controller here
func start(cmd *cobra.Command, args []string) {

	if printVersion {
		fmt.Println(Version)
		os.Exit(0)
	}

	cfg := config.NewConfig(
		viper.GetString("policies"),
		viper.GetBool("registerWebhook"),
		viper.GetBool("debug"),
		viper.GetString("metricsAddress"),
	)

	err := cfg.Validate()
	if err != nil {
		logrus.Fatalf("config validation failed: %v", err)
	}

}
