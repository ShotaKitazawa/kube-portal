package kubeportal

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/thanhpk/randstr"
)

var (
	AppName    string
	AppVersion string
)

type Opts struct {
	Version func() `short:"v" long:"version" description:"show version"`

	GitHubOAuthKey    string   `long:"github-key" description:"GitHub OAuth Client ID"`
	GitHubOAuthSecret string   `long:"github-secret" description:"GitHub OAuth Client Secret"`
	SessionSecret     string   `long:"session-secret" description:"OAuth session secret"`
	BaseUrl           string   `long:"base-url" description:"this app's externally facing base URL. Required when configuring OAuth"`
	JwtSecret         string   `long:"jwt-secret" description:"jwt secret using to check whether user is logging in"`
	GitHubAllowUsers  []string `short:"u" long:"github-allow-user" description:"GitHub Username that you allowed to get private links"`

	// optional
	KubeConfigPath string `long:"kubeconfig" description:"filepath of KubeConfig"`
	BindAddr       string `long:"bind-addr" default:"0.0.0.0"`
	BindPort       int    `long:"bind-port" default:"8080"`
}

func Parse() (*Opts, error) {
	var opts Opts
	opts.Version = func() {
		fmt.Printf("%s %s\n", AppName, AppVersion)
		os.Exit(0)
	}
	if _, err := flags.Parse(&opts); err != nil {
		return nil, err
	}
	if opts.SessionSecret == "" {
		opts.SessionSecret = randstr.String(16)
	}
	if opts.JwtSecret == "" {
		opts.JwtSecret = randstr.String(16)
	}

	return &opts, nil
}
