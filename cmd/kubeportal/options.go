package kubeportal

import "github.com/jessevdk/go-flags"

type Opts struct {
	GitHubOAuthKey    string   `long:"github-key" description:"GitHub OAuth Client ID"`
	GitHubOAuthSecret string   `long:"github-secret" description:"GitHub OAuth Client Secret"`
	SessionSecret     string   `long:"session-secret" description:"OAuth session secret"`
	GitHubAllowUsers  []string `short:"u" long:"github-allow-user" description:"GitHub Username that you allowed to get private links"`

	KubeConfigPath string `long:"kubeconfig" description:"filepath of KubeConfig"`
	BindAddr       string `long:"bind-addr" default:"0.0.0.0"`
	BindPort       int    `long:"bind-port" default:"8080"`
}

func Parse() (*Opts, error) {
	var opts Opts
	if _, err := flags.Parse(&opts); err != nil {
		return nil, err
	}
	return &opts, nil
}
