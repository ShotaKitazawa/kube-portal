package flag

import (
	"fmt"
	"net/url"

	"github.com/urfave/cli/v3"
)

var (
	BindAddr = cli.StringFlag{
		Name:    "bind-addr",
		Usage:   "bind address",
		Sources: cli.EnvVars("BIND_ADDR"),
		Value:   "0.0.0.0:8080",
	}
	Development = cli.BoolFlag{
		Name:    "development",
		Usage:   "development mode",
		Sources: cli.EnvVars("DEVELOPMENT"),
		Value:   false,
	}
	KubeConfigPath = cli.StringFlag{
		Name:     "kubeconfig",
		Usage:    "filepath of KubeConfig",
		Sources:  cli.EnvVars("KUBECONFIG"),
		Required: false,
	}
	OIDCProviderURL = cli.StringFlag{
		Name:     "oidc-provider-url",
		Usage:    "The URL of OpenID Provider. (refer to https://openid.net/specs/openid-connect-core-1_0.html#Terminology)",
		Sources:  cli.EnvVars("OIDC_PROVIDER_URL"),
		Required: true,
		Validator: func(s string) error {
			if _, err := url.ParseRequestURI(s); err != nil {
				return fmt.Errorf("failed to parse URL: %w", err)
			}
			return nil
		},
	}
	OIDCClientID = cli.StringFlag{
		Name:     "oidc-client-id",
		Usage:    "OIDC Client Identifer. (refer to https://openid.net/specs/openid-connect-core-1_0.html#Terminology)",
		Sources:  cli.EnvVars("OIDC_CLIENT_ID"),
		Required: true,
	}
	OIDCClientSecret = cli.StringFlag{
		Name:     "oidc-client-secret",
		Usage:    "OIDC Client Secret. (refer to https://openid.net/specs/openid-connect-core-1_0.html#Terminology)",
		Sources:  cli.EnvVars("OIDC_CLIENT_SECRET"),
		Required: true,
	}

	RoleAttributePath = cli.StringFlag{
		Name:     "role-attribute-path",
		Usage:    "This value must be JMESPath format. Only entities that return true are allowed.",
		Sources:  cli.EnvVars("ROLE_ATTRIBUTE_PATH"),
		Required: true,
	}
	ServerURL = cli.StringFlag{
		Name:    "server-url",
		Usage:   "The URL of this server.",
		Sources: cli.EnvVars("SERVER_URL"),
		Value:   "http://localhost:8080",
		Validator: func(s string) error {
			if _, err := url.ParseRequestURI(s); err != nil {
				return fmt.Errorf("failed to parse URL: %w", err)
			}
			return nil
		},
	}
	ShowUntaggedLinks = cli.BoolFlag{
		Name:    "show-untagged-links",
		Usage:   "Show untagged links",
		Sources: cli.EnvVars("SHOW_UNTAGGED_LINKS"),
		Value:   false,
	}
)

func Flags(version, commit string) []cli.Flag {
	return []cli.Flag{
		&BindAddr,
		&Development,
		&KubeConfigPath,
		&OIDCProviderURL,
		&OIDCClientID,
		&OIDCClientSecret,
		&RoleAttributePath,
		&ServerURL,
		&ShowUntaggedLinks,
	}
}
