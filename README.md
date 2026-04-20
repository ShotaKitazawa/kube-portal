# kube-portal

A Kubernetes dashboard that aggregates links from Ingress objects and `ExternalLink` CRDs into a single portal page. Access control is based on OIDC Bearer tokens.

## How it works

kube-portal reads the cluster and exposes a single page at `/`. Links are sourced from:

- **Ingress objects** — opt-in via `kube-portal.kanatakita.com/` annotations
- **ExternalLink CRD** — arbitrary external URLs registered as Kubernetes custom resources

Links marked `is-private: true` are hidden unless the request carries a valid OIDC token.

## Deploy

```sh
kubectl apply -f https://github.com/ShotaKitazawa/kube-portal/releases/latest/download/install.yaml
```

The `install.yaml` release asset is built with kustomize and includes the CRD, RBAC, Deployment, and Service.

## Configuration

All flags accept equivalent environment variables (e.g. `--bind-addr` → `BIND_ADDR`).

| Flag | Default | Description |
|---|---|---|
| `--bind-addr` | `0.0.0.0:8080` | Listen address |
| `--kubeconfig` | — | Path to kubeconfig (omit for in-cluster) |
| `--oidc-provider-url` | — | OIDC issuer URL |
| `--oidc-client-id` | — | OIDC client ID |
| `--oidc-audience` | — | Expected `aud` claim in access tokens |
| `--role-attribute-path` | — | JMESPath expression evaluated against the userinfo payload; must return `true` to allow access |
| `--show-untagged-links` | `false` | Show links that have no tags |

## Ingress annotations

Annotation keys follow the pattern `kube-portal.kanatakita.com/rules.<ruleIdx>.paths.<pathIdx>.<key>` to support multi-rule/multi-path Ingresses.

| Key | Description |
|---|---|
| `enable` | Must be present to show this path |
| `name` | Display name (default: `.metadata.name`) |
| `proto` | Protocol (default: `https`) |
| `icon-url` | Icon URL |
| `tags` | Comma-separated tags |
| `is-private` | Hide unless authenticated (default: `false`) |

## ExternalLink CRD

```yaml
apiVersion: kube-portal.kanatakita.com/v1alpha1
kind: ExternalLink
metadata:
  name: example
spec:
  title: "Example"
  href: "https://example.com"
  iconURL: "https://example.com/favicon.ico"
  tags: ["infra"]
  isPrivate: false
```
