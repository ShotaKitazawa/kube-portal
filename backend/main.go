package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v3"

	"github.com/ShotaKitazawa/kube-portal/backend"
	"github.com/ShotaKitazawa/kube-portal/flag"
)

var (
	appVersion string = "unknown"
	appCommit  string = "unknown"
)

func main() {
	cmd := &cli.Command{}
	cmd.Flags = flag.Flags(appVersion, appCommit)
	cmd.Action = func(ctx context.Context, cmd *cli.Command) error {
		return backend.Run(ctx, cmd)
	}
	cmd.Version = fmt.Sprintf("%s (Commit: %s)\n", appVersion, appCommit)

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
