package kubernetes

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	kubeportalv1alpha1 "github.com/ShotaKitazawa/kube-portal/server/infrastructure/kubernetes/api/v1alpha1"
	"github.com/ShotaKitazawa/kube-portal/server/models"
	"github.com/ShotaKitazawa/kube-portal/server/models/ports"
)

type Client struct {
	clientset kubernetes.Interface
	dynamic   dynamic.Interface
}

var _ ports.Kubernetes = (*Client)(nil)

func NewClient(kubeconfigPath string) (*Client, error) {
	kubeConfig, err := buildConfig(kubeconfigPath)
	if err != nil {
		return nil, err
	}
	client, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return nil, err
	}
	dynamic, err := dynamic.NewForConfig(kubeConfig)
	if err != nil {
		return nil, err
	}
	return &Client{client, dynamic}, nil
}

func buildConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	// else
	cfg, err := rest.InClusterConfig()
	if err != nil {
		cfg, err = clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize rest.Config. failed both InClusterConfig() and BuildConfigFromFlags() with RecommendedHomeFile")
		}
	}
	return cfg, nil
}

func (c *Client) ListIngress(ctx context.Context) (models.IngressInfoList, error) {
	ings, err := c.clientset.NetworkingV1().Ingresses("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var result []models.IngressInfo
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
		// path
		// - TODO: ref ing.Spec.Rules[1:].http.Paths[1:].path
		path := ing.Spec.Rules[0].HTTP.Paths[0].Path
		if path == "" {
			path = "/"
		}

		// proto
		// - allow only "http" or "https"
		proto, ok := ing.Annotations["kube-portal.kanatakita.com/proto"]
		if !ok || !(strings.ToLower(proto) == "http" || strings.ToLower(proto) == "https") {
			proto = "https"
		}
		// icon_url
		iconUrl := ing.Annotations["kube-portal.kanatakita.com/icon-url"]
		// tags
		tags := []string{}
		tagsStr, ok := ing.Annotations["kube-portal.kanatakita.com/tags"]
		if ok {
			tags = strings.Split(tagsStr, ",")
		}
		// is_private
		var isPrivate bool
		isPrivateStr, ok := ing.Annotations["kube-portal.kanatakita.com/is-private"]
		if ok && strings.ToLower(isPrivateStr) == "true" {
			isPrivate = true
		}

		/* set values to result */
		result = append(result, models.IngressInfo{
			Name:      name,
			Fqdn:      fqdn,
			Path:      path,
			Proto:     proto,
			IconUrl:   iconUrl,
			Tags:      tags,
			IsPrivate: isPrivate,
		})
	}

	return result, nil
}

func (c *Client) ListExternalLink(ctx context.Context) (models.IngressInfoList, error) {
	gvr := schema.GroupVersionResource{
		Group:    kubeportalv1alpha1.GroupVersion.Group,
		Version:  kubeportalv1alpha1.GroupVersion.Version,
		Resource: "externallinks",
	}
	uList, err := c.dynamic.Resource(gvr).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	var result []models.IngressInfo
	for _, u := range uList.Items {
		exLink := kubeportalv1alpha1.ExternalLink{}
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, &exLink); err != nil {
			return nil, err
		}
		/* get & validate values */
		u, err := url.Parse(exLink.Spec.Href)
		if err != nil {
			return nil, err
		}
		/* set values to result */
		result = append(result, models.IngressInfo{
			Name:      exLink.Spec.Title,
			Fqdn:      u.Host,
			Path:      u.Path,
			Proto:     u.Scheme,
			IconUrl:   exLink.Spec.IconURL,
			Tags:      exLink.Spec.Tags,
			IsPrivate: exLink.Spec.IsPrivate,
		})
	}
	return result, nil
}
