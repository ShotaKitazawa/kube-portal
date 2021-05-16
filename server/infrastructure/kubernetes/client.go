package kubernetes

import (
	"context"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/ShotaKitazawa/kube-portal/server/entities"
)

type Client struct {
	clientset kubernetes.Interface
}

func NewClient(kubeconfigPath string) (*Client, error) {
	kubeConfig, err := buildConfig(kubeconfigPath)
	if err != nil {
		return nil, err
	}
	client, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return nil, err
	}
	return &Client{client}, nil
}

func buildConfig(kubeconfig string) (cfg *rest.Config, err error) {
	var errBuildConfig, errInClusterConfig error
	if kubeconfig != "" {
		cfg, errBuildConfig = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if errBuildConfig == nil {
			return cfg, nil
		}
	} else {
		cfg, err = clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
		if err == nil {
			return cfg, nil
		}
	}
	cfg, errInClusterConfig = rest.InClusterConfig()
	if errInClusterConfig != nil {
		if errBuildConfig != nil {
			return nil, errBuildConfig
		}
		return nil, errInClusterConfig
	}
	return cfg, nil
}

func (c *Client) ListIngressInfo(ctx context.Context) ([]entities.IngressInfo, error) {
	// TODO: support networking.k8s.io/v1
	ings, err := c.clientset.NetworkingV1beta1().Ingresses("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var result []entities.IngressInfo
	for _, ing := range ings.Items {

		/* check ignore flag */
		ignoreStr := ing.Annotations["kube-portal.kanatakita.com/ignore"]
		if strings.ToLower(ignoreStr) == "true" {
			continue
		}

		/* get & validate values */
		// name
		// - using metadata.name if annotation is empty
		name, ok := ing.Annotations["kube-portal.kanatakita.com/name"]
		if !ok {
			name = ing.Name
		}
		// fqdn
		// - TODO: ref ing.Spec.Rules[1:].Host
		fqdn := ing.Spec.Rules[0].Host
		// proto
		// - allow only "http" or "https"
		proto, ok := ing.Annotations["kube-portal.kanatakita.com/proto"]
		if !ok || !(strings.ToLower(proto) == "http" || strings.ToLower(proto) == "https") {
			proto = "https"
		}
		// icon_url
		// - allow empty
		iconUrl := ing.Annotations["kube-portal.kanatakita.com/icon-url"]
		// is_private
		var isPrivate bool
		isPrivateStr, ok := ing.Annotations["kube-portal.kanatakita.com/is-private"]
		if ok && strings.ToLower(isPrivateStr) == "true" {
			isPrivate = true
		}

		/* set values to result */
		result = append(result, entities.IngressInfo{
			Name:      name,
			Fqdn:      fqdn,
			Proto:     proto,
			IconUrl:   iconUrl,
			IsPrivate: isPrivate,
		})
	}

	return result, nil
}
