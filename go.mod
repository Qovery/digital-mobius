module github.com/Qovery/digital-mobius

go 1.16

replace github.com/Qovery/digital-mobius/cmd => ./cmd

replace github.com/Qovery/digital-mobius/utils => ./utils

require (
	github.com/Qovery/digital-mobius/cmd v0.0.0-20210414080004-add85fc52139 // indirect
	github.com/Qovery/digital-mobius/utils v0.0.0-20210414080004-add85fc52139 // indirect
	github.com/digitalocean/godo v1.60.0 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/spf13/cobra v1.1.3 // indirect
	github.com/spf13/viper v1.7.1 // indirect
	k8s.io/client-go v0.21.0 // indirect
)
