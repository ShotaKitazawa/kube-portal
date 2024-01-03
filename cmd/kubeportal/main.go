package main

import (
	"log"

	"github.com/ShotaKitazawa/kube-portal/cmd/kubeportal/flag"
	"github.com/ShotaKitazawa/kube-portal/server"
)

var (
	appVersion string
	appCommit  string
)

func main() {
	opts, err := flag.Parse(appVersion, appCommit)
	if err != nil {
		log.Fatal(err)
	}
	if err := server.Run(opts); err != nil {
		log.Fatal(err)
	}
}
