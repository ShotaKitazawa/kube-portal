# kube-portal

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
![test](https://github.com/ShotaKitazawa/kube-portal/workflows/test/badge.svg)

## What is this?

kube-portal works as portal-site for links published by Kubernetes Ingress.
You can also ["customize the link's icon"](./docs/ingress-annotations.md) and ["set up private links"](./docs/private-links.md) by setting annotations to Ingress object.

## Quick Start

* install

```bash
kubectl create namespace kube-portal
kubectl apply -f https://raw.githubusercontent.com/ShotaKitazawa/kube-portal/master/manifests/install.yaml
```

* port-forward

```bash
kubectl -n kube-portal port-forward svc/kube-portal 8080:8080
```

Then, display like the below page when access to <http://localhost:8080>.

![page-top](https://raw.githubusercontent.com/ShotaKitazawa/kube-portal/image/page-top.png)

## More informations

* [/docs](/docs)
