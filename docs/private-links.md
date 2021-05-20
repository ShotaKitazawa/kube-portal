# Privete Links

If specifying `kube-portal.kanatakita.com/is-private` annotation, The link associated with that Ingress Object will become a private link, that will not be visible unless you do an OAuth authentication & authorization on GitHub.

### Expose kube-portal for OAuth callback

You must expose kube-portal for OAuth callback.
In the following steps, assume that you are exposed you app at `https://portal.example.com`

### Get GitHub OAuth ID & Secret

You access to https://github.com/settings/applications/new and create OAuth application

![create-oauth-app](https://raw.githubusercontent.com/ShotaKitazawa/kube-portal/image/create-oauth-app.png)

### Demo

* create below files

```yaml
# kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
- https://raw.githubusercontent.com/ShotaKitazawa/kube-portal/master/manifests/install.yaml
patchesStrategicMerge:
- deployment.yaml
```

```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: kube-portal
  namespace: kube-portal
spec:
  template:
    spec:
      containers:
        - name: kube-portal
          env:
            - name: BASE_URL
              value: https://portal.kanatakita.com
            - name: GITHUB_ALLOW_USERS
              value: ShotaKitazawa
            - name: GITHUB_CLIENT_ID
              value: <input GitHub OAuth Client ID>
            - name: GITHUB_CLIENT_SECRET
              value: <input GitHub OAuth Client Secret>
            - name: JWT_SECRET
              value: this-is-jwt-secret
```

* apply above kube-portal manifest

```
kustomize build . | kubectl apply -f-
```

* apply below Ingress manifest

```yaml
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  name: portfolio
  annotations:
    kube-portal.kanatakita.com/name: "Portfolio"
    kube-portal.kanatakita.com/icon-url: "https://portfolio.kanatakita.com/images/shoshinsha-shirokuro.jpg"
    kube-portal.kanatakita.com/is-private: "true"
spec:
  rules:
  - host: portfolio.kanatakita.com
    http:
      paths:
      - backend:
          serviceName: portfolio
          servicePort: 8080
```

* access portal site

![demo-with-arrow](https://raw.githubusercontent.com/ShotaKitazawa/kube-portal/image/demo-with-arrow.png)

* login by GitHub OAuth

![demo-loggedin](https://raw.githubusercontent.com/ShotaKitazawa/kube-portal/image/demo-loggedin.png)
