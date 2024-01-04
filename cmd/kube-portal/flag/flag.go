package flag

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/thanhpk/randstr"
)

type Opts struct {
	BaseUrl           string
	BindAddr          string
	ExpiredHour       int
	GitHubAllowUsers  string
	GitHubOAuthKey    string
	GitHubOAuthSecret string
	JwtSecret         string
	KubeConfigPath    string
	ShowUntaggedLinks bool
}

func Parse(version, commit string) (*Opts, error) {
	var showVersion bool
	flag.BoolVar(&showVersion, "version", false,
		"show version")

	var opts Opts
	flag.StringVar(&opts.BaseUrl, "base-url", "", "app's externally facing base URL (optional, required if enable Login feature)")
	flag.StringVar(&opts.BindAddr, "bind-addr", "0.0.0.0:8080", "bind address")
	flag.IntVar(&opts.ExpiredHour, "expired-hour", 12, "jwt expired time (hour)")
	flag.StringVar(&opts.GitHubAllowUsers, "github-allow-users", "", "specified GitHub Usernames by comma-separated that you allowed to get private links")
	flag.StringVar(&opts.GitHubOAuthKey, "github-client-id", "", "GitHub OAuth Client ID (optional, required if enable Login feature)")
	flag.StringVar(&opts.GitHubOAuthSecret, "github-client-secret", "", "GitHub OAuth Client Secret (optional, required if enable Login feature)")
	flag.StringVar(&opts.JwtSecret, "jwt-secret", "", "jwt secret using to check whether user is logging in (optional)")
	flag.StringVar(&opts.KubeConfigPath, "kubeconfig", "", "filepath of KubeConfig")
	flag.BoolVar(&opts.ShowUntaggedLinks, "show-untagged-links", true, "show untagged links")
	flag.VisitAll(func(f *flag.Flag) {
		if s := os.Getenv(strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))); s != "" {
			_ = f.Value.Set(s)
		}
	})
	flag.Parse()

	if showVersion {
		fmt.Printf("Version: %s (Commit: %s)\n", version, commit)
		return nil, fmt.Errorf("")
	}

	if opts.JwtSecret == "" {
		opts.JwtSecret = randstr.String(32)
	}

	return &opts, nil
}
