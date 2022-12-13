package cmd

import (
	"fmt"
	"os"

	"github.com/RedLabsPlatform/kube-shield/pkg/webhook/cache"
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
	Cmd.Flags().StringP("kubeconfig", "k", "", "Path to the kubeconfig file to run outside of the cluster")
	Cmd.Flags().BoolVarP(&printVersion, "version", "v", false, "Print version of kube-shield")
	viper.BindPFlag("kubeconfig", Cmd.Flags().Lookup("kubeconfig"))
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
	// ng := &engine.Engine{
	// 	CacheController: cachectrl,
	// }
	// srv := server.Server{
	// 	Engine: ng,
	// }

	cachectrl.Run(make(chan struct{}), make(chan struct{}))
	//srv.Start()
}
