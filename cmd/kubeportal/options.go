package kubeportal

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/thanhpk/randstr"
)

var (
	AppName    string
	AppVersion string
)

type Opts struct {
	Version           bool
	GitHubOAuthKey    string
	GitHubOAuthSecret string
	BaseUrl           string
	JwtSecret         string
	ExpiredHour       int
	GitHubAllowUsers  string
	KubeConfigPath    string
	BindAddr          string
}

func Parse() (*Opts, error) {
	var opts Opts
	flag.BoolVar(&opts.Version, "version", false, "show version")
	flag.StringVar(&opts.GitHubOAuthKey, "github-client-id", "", "GitHub OAuth Client ID (optional, required if enable Login feature)")
	flag.StringVar(&opts.GitHubOAuthSecret, "github-client-secret", "", "GitHub OAuth Client Secret (optional, required if enable Login feature)")
	flag.StringVar(&opts.BaseUrl, "base-url", "", "app's externally facing base URL (optional, required if enable Login feature)")
	flag.StringVar(&opts.JwtSecret, "jwt-secret", randstr.String(16), "jwt secret using to check whether user is logging in (optional)")
	flag.IntVar(&opts.ExpiredHour, "expired-hour", 12, "jwt expired time (hour)")
	flag.StringVar(&opts.GitHubAllowUsers, "github-allow-users", "", "specified GitHub Usernames by comma-separated that you allowed to get private links")
	flag.StringVar(&opts.KubeConfigPath, "kubeconfig", "", "filepath of KubeConfig")
	flag.StringVar(&opts.BindAddr, "bind-addr", "0.0.0.0:8080", "bind address")
	flag.VisitAll(func(f *flag.Flag) {
		if s := os.Getenv(strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))); s != "" {
			f.Value.Set(s)
		}
	})
	flag.Parse()

	if opts.Version {
		fmt.Printf("%s %s\n", AppName, AppVersion)
		os.Exit(0)
	}

	return &opts, nil
}
