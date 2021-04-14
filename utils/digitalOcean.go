package utils

import (
	"context"
	"github.com/digitalocean/godo"
	log "github.com/sirupsen/logrus"
	"os"
)

func GetDOClient() *godo.Client {
	client := godo.NewFromToken(os.Getenv("DIGITAL_OCEAN_TOKEN"))

	return client
}

func RecyleNode(params ...interface{}) interface {} {
	client, clientOk := params[0].(*godo.Client)
	node, nodeOk := params[1].(NodeInfos)

	if clientOk && nodeOk {
		response, err := client.Kubernetes.DeleteNode(context.TODO(), node.ClusterId, node.PoolId, node.NodeId, &godo.KubernetesNodeDeleteRequest{Replace: true, SkipDrain: true})

		if err != nil {
			log.Errorf("Can't recycle node %s : %s", node.NodeId, err.Error())
		}

		if response.Status != "200" {
			log.Errorf("Request error: %s", response.String())
		}
	}

	return nil
}

func RecycleStuckNodes(client *godo.Client, stuckedNodes []NodeInfos) {
	RunWorkers(2,"Nodes recycling complete.", 1, RecyleNode, client, stuckedNodes)
}


