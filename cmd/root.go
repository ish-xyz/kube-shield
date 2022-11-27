package cmd

import (
	"fmt"
	"os"

	"github.com/RedLabsPlatform/kube-shield/pkg/config"
	"github.com/RedLabsPlatform/kube-shield/pkg/server"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	printVersion bool
	configFile   string
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
	Cmd.Flags().StringVarP(&configFile, "config", "c", "", "Kube-shield configuration file")

	// Bound Flags
	Cmd.Flags().StringP("policies", "p", "", "Path to the directory with the policies")
	Cmd.Flags().String("web-address", "0.0.0.0:8000", "Address where the webhook webserver is exposed")
	Cmd.Flags().String("web-path", "/webhook", "Path where the webhook webserver is reachable")
	Cmd.Flags().String("tls-key", "/etc/kube-shield/tls/key.pem", "Path to the tls private key")
	Cmd.Flags().String("tls-cert", "/etc/kube-shield/tls/cert.pem", "Path to the tls certificate")
	Cmd.Flags().String("metrics-address", "0.0.0.0:3000", "Address where the metrics are exposed")
	Cmd.Flags().String("metrics-path", "/metrics", "Path where the metrics are exposed")
	Cmd.Flags().BoolP("register-webhook", "r", true, "create ValidatingWebhookConfiguration resource in the current Kubernetes")
	Cmd.Flags().BoolP("debug", "d", false, "debug mode")

	// Required flags
	Cmd.MarkFlagRequired("policies")

	viper.BindPFlag("web.address", Cmd.Flags().Lookup("web-address"))
	viper.BindPFlag("web.path", Cmd.Flags().Lookup("web-path"))
	viper.BindPFlag("web.tls.key", Cmd.Flags().Lookup("tls-key"))
	viper.BindPFlag("web.tls.cert", Cmd.Flags().Lookup("tls-cert"))

	viper.BindPFlag("policies", Cmd.Flags().Lookup("policies"))
	viper.BindPFlag("register", Cmd.Flags().Lookup("register-webhook"))
	viper.BindPFlag("debug", Cmd.Flags().Lookup("debug"))
	viper.BindPFlag("metrics.address", Cmd.Flags().Lookup("metrics-address"))
	viper.BindPFlag("metrics.path", Cmd.Flags().Lookup("metrics-path"))

}

// Start the admission controller here
func start(cmd *cobra.Command, args []string) {

	if printVersion {
		fmt.Println(Version)
		os.Exit(0)
	}

	if configFile != "" {
		viper.SetConfigFile(configFile)
		if viper.ReadInConfig() != nil {
			logrus.Fatal("error loading configuration")
		}
	}

	cfg := config.NewConfig(
		viper.GetString("policies"),
		viper.GetString("web.address"),
		viper.GetString("web.path"),
		viper.GetString("web.tls.key"),
		viper.GetString("web.tls.cert"),
		viper.GetBool("register"),
		viper.GetBool("debug"),
		viper.GetString("metrics.address"),
		viper.GetString("metrics.path"),
	)

	err := cfg.Validate()
	if err != nil {
		logrus.Fatalf("config validation failed: %v", err)
	}

	srv, err := server.NewServer(cfg)
	if err != nil {
		logrus.Fatalf("cannot create new server instance: %v", err)
	}
	srv.Run()

}