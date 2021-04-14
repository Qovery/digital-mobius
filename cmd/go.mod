module github.com/Qovery/do-k8s-replace-notready-nodes/cmd

go 1.16

replace github.com/Qovery/do-k8s-replace-notready-nodes/utils => ../utils

require (
	github.com/Qovery/do-k8s-replace-notready-nodes/utils v0.0.0-00010101000000-000000000000
	github.com/mitchellh/go-homedir v1.1.0
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.7.1
	k8s.io/api v0.20.5
	k8s.io/apimachinery v0.20.5
	k8s.io/client-go v0.20.5
)
