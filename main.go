package main

import (
	"log"

	"github.com/ShotaKitazawa/kube-portal/backend"
	"github.com/ShotaKitazawa/kube-portal/flag"
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
	if err := backend.Run(opts); err != nil {
		log.Fatal(err)
	}
}
