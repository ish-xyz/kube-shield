package cmd

import (
	"fmt"
	"os"

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
	Cmd.Flags().BoolP("debug", "d", false, "Path to the directory with the policies")

	// Required flags
	Cmd.MarkFlagRequired("policies")
	Cmd.MarkFlagRequired("config")

	viper.BindPFlag("policies", Cmd.Flags().Lookup("policies"))
	viper.BindPFlag("debug", Cmd.Flags().Lookup("debug"))
}

// Start the admission controller here
func start(cmd *cobra.Command, args []string) {
	if printVersion {
		fmt.Println(Version)
		os.Exit(0)
	}

}
