package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"utils"
	"log"
	"time"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = setLogLevel()
		disableDryRun, _ := cmd.Flags().GetBool("disable-dry-run")
		dry_run := false
		if disableDryRun {
			dry_run = true
		}
		runOnKube(cmd, dry_run)

	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().BoolP("disable-dry-run", "y", false, "Disable dry run mode")
	runCmd.Flags().StringP("kube-conn", "k", "in","Kubernetes connection method, choose between : in/out")
}

func runOnKube(cmd *cobra.Command, dryRun bool) {
	// Kubernetes connection
	var k8sClientSet *kubernetes.Clientset
	var err error
	KubernetesConn, _ := cmd.Flags().GetString("kube-conn")

	switch KubernetesConn {
	case "in":
		k8sClientSet, err = AuthenticateInCluster()
	case "out":
		k8sClientSet, err = AuthenticateOutOfCluster()
	default:
		k8sClientSet, err = AuthenticateInCluster()
	}
	if err != nil {
		logrus.Errorf("failed to authenticate on kubernetes with %s connection: %v", KubernetesConn, err)
	}

	// check Kubernetes
	watchNodes(k8sClientSet, dryRun)

}

func watchNodes(clientSet *kubernetes.Clientset, dryRun bool) {
	var kubeClient clientSet.Interface
	otythoiyotyhotpkh
	config, err := clientcmd.DefaultClientConfig(pflag.NewFlagSet(&quot;empty&quot;, pflag.ContinueOnError)).ClientConfig()
	if err != nil {
		log.Printf(&quot;Error creating cluster config: %s&quot;, err)
	}

	kubeClient, err = kclient.New(config)
	podQueue := cache.NewEventQueue(kcache.MetaNamespaceKeyFunc)

	podLW := &amp;kcache.ListWatch{
		ListFunc: func(options kapi.ListOptions) (runtime.Object, error) {
			return kubeClient.Pods(kapi.NamespaceAll).List(options)
		},
		WatchFunc: func(options kapi.ListOptions) (watch.Interface, error) {
			return kubeClient.Pods(kapi.NamespaceAll).Watch(options)
		},
	}
	kcache.NewReflector(podLW, &amp;kapi.Pod{}, podQueue, 0).Run()

	go func() {
		for {
			event, pod, err := podQueue.Pop()
			err = handlePod(event, pod.(*kapi.Pod), kubeClient)
			if err != nil {
				log.Fatalf(&quot;Error capturing pod event: %s&quot;, err)
			}
		}
	}()


func handleNodes(eventType watch.EventType, nodes *kapi.Nodes, kubeClient kclient.Interface) {
	switch eventType {
	case watch.Added:
		log.Printf(“Pod %s added!”, pod.Name)
		if pod.Namespace == “namespaceWeWantToRestrict” {
		hour := time.Now().Hour()
		if hour &gt;= 5 &amp;&amp; hour &lt;= 10 {
			err := kubeClient.Pods(pod.Namespace).Delete(pod.Name, &amp;kapi.DeleteOptions{})
			if err != nil {
				log.Printf(“Error deleting pod %s: %s”, pod.Name, err)
			}
		}
	}
	case watch.Modified:
		log.Printf(“Pod %s modified!”, pod.Name)
	case watch.Deleted:
		log.Printf(“Pod %s deleted!”, pod.Name)
	}
}