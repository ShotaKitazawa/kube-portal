# Ingress Annotations

kube-portal create links refers to `.spec.rules[0].host`.

You can control the behavior of kube-portal by specifying the following Annotations to Ingress Object.

| Annotations                             | Description                                                                          |
|:---------------------------------------:|:-------------------------------------------------------------------------------------|
| `kube-portal.kanatakita.com/name`       | displayed link name (default: using `.metadata.name`)                                |
| `kube-portal.kanatakita.com/icon-url`   | displayed link's icon url                                                            |
| `kube-portal.kanatakita.com/proto`      | link's protocol (default: `https`)                                                   |
| `kube-portal.kanatakita.com/is-private` | if this annotation is "true", don't display link unless logged in (default: `false`) |
| `kube-portal.kanatakita.com/ignore`     | if this annotations is "true", don't display link (default: `false`)                 |

### Demo

* install kube-portal to k8s

```shell
kubectl apply -f https://raw.githubusercontent.com/ShotaKitazawa/kube-portal/master/manifests/install.yaml
```

* apply below manifest

```yaml
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  name: portal
  annotations:
    kube-portal.kanatakita.com/name: "Portal Site"
    kube-portal.kanatakita.com/icon-url: "https://raw.githubusercontent.com/ShotaKitazawa/kube-portal/image/kube-portal.png"
spec:
  rules:
  - host: portal.kanatakita.com
    http:
      paths:
      - backend:
          serviceName: kube-portal
          servicePort: 8080
```

* access portal site

![demo](https://raw.githubusercontent.com/ShotaKitazawa/kube-portal/image/demo.png)
