# Options

* kube-portal options
    * specifiable alse environment variables (`--base-url` -> `BASE_URL`)

```
Usage of kube-portal
  -base-url string
        app's externally facing base URL (optional, required if enable Login feature)
  -bind-addr string
        bind address (default "0.0.0.0:8080")
  -expired-hour int
        jwt expired time (hour) (default 12)
  -github-allow-users string
        specified GitHub Usernames by comma-separated that you allowed to get private links
  -github-client-id string
        GitHub OAuth Client ID (optional, required if enable Login feature)
  -github-client-secret string
        GitHub OAuth Client Secret (optional, required if enable Login feature)
  -jwt-secret string
        jwt secret using to check whether user is logging in (optional)
  -kubeconfig string
        filepath of KubeConfig
  -version
        show version
```
