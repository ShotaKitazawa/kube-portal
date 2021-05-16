package main

import (
	"github.com/ShotaKitazawa/kube-portal/cmd/kubeportal"
	"github.com/ShotaKitazawa/kube-portal/server"
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
