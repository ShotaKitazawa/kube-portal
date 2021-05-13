package main

import (
	"github.com/ShotaKitazawa/kube-portal/pkg/cmd/kubeportal"
	"github.com/ShotaKitazawa/kube-portal/pkg/server"
)

func main() {
	opts, err := kubeportal.Parse()
	if err != nil {
		panic(err)
	}
	if err := server.Run(opts); err != nil {
		panic(err)
	}
}
