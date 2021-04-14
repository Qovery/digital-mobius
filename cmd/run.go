package cmd

import (
	"github.com/Qovery/do-k8s-replace-notready-nodes/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"os"
	"time"
)

var runCmd = &cobra.Command{
	Use:   "do-nodes",
	Short: "Recycle all not ready digital ocean clusters nodes",
	Run: func(cmd *cobra.Command, args []string) {
		_ = setLogLevel()
		dryRun, _ := cmd.Flags().GetBool("disable-dry-run")
		log.Info("Starting DO nodes recycler.")

		if !dryRun {
			log.Info("Running DO nodes recycler in dry mode.")
		}

		kubernetesConn, _ := cmd.Flags().GetString("kube-conn")
		if kubernetesConn != "in" && kubernetesConn != "out" {
			log.Error("Choose a Kubernetes connection method between 'in' or 'out'")
			return
		}

		token, isTokenPresent := os.LookupEnv("DIGITAL_OCEAN_TOKEN")
		if !isTokenPresent || token == "" {
			log.Error("You need to add your digital ocean token to env DIGITAL_OCEAN_TOKEN")
			return
		}

		clusterId, isIdPresent := os.LookupEnv("DIGITAL_OCEAN_CLUSTER_ID")
		if !isIdPresent || clusterId == "" {
			log.Error("You need to add the digital ocean cluster id to env DIGITAL_OCEAN_CLUSTER_ID")
			return
		}

		DelayEnv, isDelay := os.LookupEnv("DELAY_NODE_CREATION")
		if !isDelay || DelayEnv == "" {
			log.Error("You need to add the delay in minutes which consider a node is stuck to env DELAY_NODE_CREATION")
			return
		}

		creationDelay, err := time.ParseDuration(DelayEnv)
		if err != nil {
			log.Errorf("Can't parse MINUTES_DELAY_NODE_CREATION env: %s", err.Error())
			return
		}

		runKubeCmd(cmd, kubernetesConn, creationDelay, dryRun)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().BoolP("disable-dry-run", "y", false, "Disable dry run mode")
	runCmd.Flags().StringP("kube-conn", "k", "in","Kubernetes connection method, choose between : in/out")
}

func runKubeCmd(cmd *cobra.Command, kubernetesConn string, creationDelay time.Duration, dryRun bool) {
	k8sClientSet, dynamicClient := getKubeClient(cmd, kubernetesConn)
	DOclient := utils.GetDOClient()
	stuckNodes :=  getStuckNodes(k8sClientSet, creationDelay)
	if dryRun {
		utils.RecycleStuckNodes(DOclient, stuckNodes)
		log.Debug("Starting kubernetes nodes watch.")
		utils.WatchNodes(k8sClientSet, dynamicClient, DOclient, creationDelay)
	}
}

func getKubeClient(cmd *cobra.Command, kubernetesConn string) (*kubernetes.Clientset,dynamic.Interface) {
	var k8sClientSet *kubernetes.Clientset
	var dynamicClient dynamic.Interface
	var err error

	switch kubernetesConn {
		case "out":
			k8sClientSet, dynamicClient, err = utils.AuthenticateOutOfCluster()
		default:
			k8sClientSet, dynamicClient, err = utils.AuthenticateInCluster()
	}
	if err != nil {
		log.Errorf("Failed to authenticate on kubernetes with %s connection: %v", kubernetesConn, err)
	}

	return k8sClientSet, dynamicClient
}

func getStuckNodes(clientSet *kubernetes.Clientset, creationDelay time.Duration) []utils.NodeInfos {
	stuckNodes := utils.GetKubeStuckNodesInfos(clientSet, creationDelay)
	if len(stuckNodes) == 0 {
		log.Debug("There is no stuck nodes to recycle")
		return []utils.NodeInfos{}
	}

	log.Debugf("There is %d stuck node(s) to recycle", len(stuckNodes))

	return stuckNodes
}




