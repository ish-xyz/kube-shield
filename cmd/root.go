package cmd

import (
	"fmt"
	"os"

	"github.com/RedLabsPlatform/kube-shield/pkg/cache"
	"github.com/RedLabsPlatform/kube-shield/pkg/engine"
	"github.com/RedLabsPlatform/kube-shield/pkg/server"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
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
	Cmd.Flags().StringP("kubeconfig", "k", "", "path to the kubeconfig file to run outside of the cluster")
	Cmd.Flags().StringP("tls-cert", "p", "/var/ssl/server.crt", "path to the server TLS certificate")
	Cmd.Flags().StringP("tls-key", "i", "/var/ssl/server.key", "path to the server TLS key")
	Cmd.Flags().StringP("address", "a", ":8000", "address of the web server")
	Cmd.Flags().Bool("ipv4", false, "Run web server on ipv4")
	Cmd.Flags().BoolVarP(&printVersion, "version", "v", false, "Print version of kube-shield")

	// Flags binding
	viper.BindPFlag("kubeconfig", Cmd.Flags().Lookup("kubeconfig"))
	viper.BindPFlag("tls-cert", Cmd.Flags().Lookup("tls-cert"))
	viper.BindPFlag("tls-key", Cmd.Flags().Lookup("tls-key"))
	viper.BindPFlag("address", Cmd.Flags().Lookup("address"))
	viper.BindPFlag("ipv4", Cmd.Flags().Lookup("ipv4"))
}

// Start the admission controller here
func start(cmd *cobra.Command, args []string) {

	if printVersion {
		fmt.Println(Version)
		os.Exit(0)
	}

	// try in-cluster kubeconfig
	kubecfg, err := rest.InClusterConfig()
	if err != nil {
		// try flag passed kubeconfig
		kubecfg, err = clientcmd.BuildConfigFromFlags("", viper.GetString("kubeconfig"))
	}
	if err != nil {
		logrus.Fatal("failed to load kubeconfig")
	}

	// cert, err := tls.LoadX509KeyPair("/tmp/server.crt", "/tmp/server.key")
	// if err != nil {
	// 	logrus.Fatal("failed to load certificates")
	// }

	dc, err := dynamic.NewForConfig(kubecfg)
	if err != nil {
		logrus.Fatal("failed to load Kubernetes dynamic client")
	}
	index := cache.NewCacheIndex()
	cachectrl := cache.NewCacheController(dc, index)
	ngin := &engine.Engine{
		CacheController: cachectrl,
	}
	srv := server.NewServer(
		viper.GetString("address"),
		viper.GetString("tls-cert"),
		viper.GetString("tls-key"),
		viper.GetBool("ipv4"),
		ngin,
	)

	go cachectrl.Run(make(chan struct{}), make(chan struct{}))
	srv.Start()
}
