module github.com/Qovery/do-k8s-replace-notready-nodes

go 1.16

replace github.com/Qovery/do-k8s-replace-notready-nodes/cmd => ./cmd

replace github.com/Qovery/do-k8s-replace-notready-nodes/utils => ./utils

require (
	github.com/Qovery/do-k8s-replace-notready-nodes/cmd v0.0.0-00010101000000-000000000000
	github.com/digitalocean/godo v1.59.0 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/spf13/cobra v1.1.3 // indirect
	github.com/spf13/viper v1.7.1 // indirect
	k8s.io/apimachinery v0.20.5 // indirect
)
