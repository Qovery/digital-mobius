package utils

import (
	"context"
	"fmt"
	"github.com/digitalocean/godo"
	log "github.com/sirupsen/logrus"
	v12 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"strings"
	"time"
)

type NodeInfos struct {
	ClusterId string
	PoolId string
	NodeId string
}

func AuthenticateInCluster() (*kubernetes.Clientset, dynamic.Interface, error) {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get client config: %v", err)
	}


	// creates the clientset
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate client set: %v", err)
	}


	// creates dynClient for kube watching
	dynClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate dynamic client : %v", err)
	}

	return clientSet, dynClient, nil
}

func AuthenticateOutOfCluster() (*kubernetes.Clientset, dynamic.Interface, error) {
	kubeconfig := os.Getenv("KUBECONFIG")

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get client config: %v", err)
	}

	// create the clientset
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate client set: %v", err)
	}

	// creates dynClient for kube watching
	dynClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate dynamic client : %v", err)
	}

	return clientSet, dynClient, nil
}

func GetKubeStuckNodesInfos(clientSet *kubernetes.Clientset, minutesDelay time.Duration) []NodeInfos {
	nodes, err := clientSet.CoreV1().Nodes().List(context.TODO(), v1.ListOptions{Limit: int64(1000)})

	if err != nil {
		log.Errorf("Can't list kubernetes nodes: %s", err.Error())
	}

	stuckNodes := getStuckNodes(nodes, minutesDelay)

	if len(stuckNodes) == len(nodes.Items) {
		log.Debug("Cluster isn't ready yet.")
		return []NodeInfos{}
	}

	return stuckNodes
}

/**
* Using three conditions to check if node is stock : node Ready status is at "unknown", it's not Unschedulable and it was
* created more than the creation delay duration ago.
 */
func checkNode(node v12.Node, creationDelay time.Duration) NodeInfos {
	var fullNode NodeInfos

	for _, condition := range node.Status.Conditions {
		if strings.EqualFold(string(condition.Type), "Ready") &&
			strings.EqualFold(string(condition.Status), "Unknown") &&
			!node.Spec.Unschedulable  &&
			node.CreationTimestamp.Add(creationDelay).Before(time.Now()) {
			return getNodesInfos(node)
		}
	}

	return fullNode
}

func getStuckNodes(nodelist *v12.NodeList, creationDelay time.Duration) []NodeInfos {
	var stuckNodesId []NodeInfos
	var emptyNode NodeInfos

	for _, node := range nodelist.Items {
		checkedNode := checkNode(node, creationDelay)
		if emptyNode != checkedNode {
			stuckNodesId = append(stuckNodesId, checkedNode)
		}
	}

	return stuckNodesId
}

func WatchNodes(clientSet *kubernetes.Clientset, dynClient dynamic.Interface, DOclient *godo.Client, creationDelay time.Duration) {
	log.Debug("Starting kubernetes nodes watch.")

	result, err := clientSet.CoreV1().Nodes().Watch(context.TODO(), v1.ListOptions{Limit: int64(1000)})
	if err != nil {
		log.Printf("Can't watch nodes: %s", err)
		time.Sleep(30 * time.Second)
		WatchNodes(clientSet, dynClient, DOclient, creationDelay)
	}


	for event := range result.ResultChan() {
		if event.Type == watch.Added || event.Type == watch.Modified {
			node := getNodesInfosFromRuntineObject(dynClient, "default", event.Object)
			checkedNode := checkNode(node, creationDelay)
			emptyNode := NodeInfos{}
			if checkedNode != emptyNode {
				log.Infof("Recyling node %s", checkedNode)
				RecyleNode(DOclient, checkedNode)
			}
		}
	}

	log.Debug("Closing watch channel.")
}

func getNodesInfosFromRuntineObject(dynClient dynamic.Interface, nameSpace string, obj runtime.Object) v12.Node{
	unstructuredObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		log.Errorf("Can't unstructure runtime object: %s", err)
		return v12.Node{}
	}

	var node v12.Node
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(unstructuredObj, &node)
	if err != nil {
		log.Errorf("Can't convert to node: %s", err)
		return v12.Node{}
	}

	return node
}

func getNodesInfos(node v12.Node) NodeInfos {
	nodeInfos := NodeInfos{
		NodeId: "",
		PoolId: "",
		ClusterId: os.Getenv("DIGITAL_OCEAN_CLUSTER_ID"),
	}


	for key, value := range node.Labels {
		if key == "doks.digitalocean.com/node-pool-id" {
			nodeInfos.PoolId = value
		}

		if key == "doks.digitalocean.com/node-id" {
			nodeInfos.NodeId = value
		}
	}

	return nodeInfos
}

