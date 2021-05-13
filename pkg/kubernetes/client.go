package kubernetes

import (
	"context"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/ShotaKitazawa/kube-portal/pkg/entities"
)

type Client struct {
	clientset kubernetes.Interface
}

func NewClient(kubeconfigPath string) (*Client, error) {
	if kubeconfigPath == "" {
		kubeconfigPath = clientcmd.RecommendedHomeFile
	}
	kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, err
	}
	client, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return nil, err
	}
	return &Client{client}, nil
}

func (c *Client) ListIngressInfo(ctx context.Context) ([]entities.IngressInfo, error) {
	ings, err := c.clientset.NetworkingV1().Ingresses("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var result []entities.IngressInfo
	for _, ing := range ings.Items {

		// check ignore flag
		ignoreStr := ing.Annotations["kube-portal.kanatakita.com/ignore"]
		if strings.ToLower(ignoreStr) == "true" {
			continue
		}

		// correct values
		fqdn := ing.Spec.Rules[0].Host
		iconUrl, ok := ing.Annotations["kube-portal.kanatakita.com/icon-url"]
		if !ok {
			iconUrl = entities.DefaultIconUrl
		}
		proto, ok := ing.Annotations["kube-portal.kanatakita.com/proto"]
		if !ok { // TODO: validation
			proto = "https"
		}
		var isPrivate bool
		isPrivateStr, ok := ing.Annotations["kube-portal.kanatakita.com/is-private"]
		if ok && strings.ToLower(isPrivateStr) == "true" {
			isPrivate = true
		}
		result = append(result, entities.IngressInfo{
			Fqdn:      fqdn,
			Proto:     proto,
			IconUrl:   iconUrl,
			IsPrivate: isPrivate,
		})
	}

	return result, nil
}
