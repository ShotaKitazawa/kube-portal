package kubernetes

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/elliotchance/orderedmap/v2"
	"github.com/mattn/go-jsonpointer"
	"github.com/sirupsen/logrus"
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

const ingressAnnotationPrefix = "kube-portal.kanatakita.com/"

var ingressAnnotationMatcher *regexp.Regexp

func init() {
	ingressAnnotationMatcher = regexp.MustCompile(
		`^rules\.(\d+)\.paths\.(\d+)\.(.+)$`)
}

type Client struct {
	clientset kubernetes.Interface
	dynamic   dynamic.Interface
	log       *logrus.Logger
}

var _ ports.Kubernetes = (*Client)(nil)

func NewClient(l *logrus.Logger, kubeconfigPath string) (*Client, error) {
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
	return &Client{client, dynamic, l}, nil
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

	var result models.IngressInfoList
	for _, ing := range ings.Items {
		m := orderedmap.NewOrderedMap[string, models.IngressInfo]()
		for key, val := range ing.Annotations {
			// Skip if it's not related annotation
			if !strings.HasPrefix(key, ingressAnnotationPrefix) {
				continue
			}
			key = strings.TrimPrefix(key, ingressAnnotationPrefix)

			// Extract value from annotations
			l := ingressAnnotationMatcher.FindStringSubmatch(key)
			if len(l) != 4 {
				c.log.Warn("TODO_1")
				continue
			}
			ruleIdx, _ := strconv.Atoi(l[1])
			pathIdx, _ := strconv.Atoi(l[2])
			action := l[3]

			//
			// Fill models.IngressInfo
			//
			tmpIngressInfo, _ := m.Get(fmt.Sprintf("%d-%d", ruleIdx, pathIdx))
			// Fill from annotations
			switch action {
			case "enable":
				// pass
			case "name":
				tmpIngressInfo.Name = val
			case "proto":
				tmpIngressInfo.Proto = val
			case "icon-url":
				tmpIngressInfo.IconUrl = val
			case "tags":
				tmpIngressInfo.Tags = strings.Split(val, ",")
			case "is-private":
				if strings.ToLower(val) == "true" {
					tmpIngressInfo.IsPrivate = true
				}
			default:
				c.log.Warn("TODO")
				continue
			}
			// Fill Fqdn & Path if first time
			if tmpIngressInfo.Hostname == "" || tmpIngressInfo.Path == "" {
				b, err := json.Marshal(&ing.Spec)
				if err != nil {
					return nil, err
				}
				var obj interface{}
				if err := json.Unmarshal(b, &obj); err != nil {
					return nil, err
				}
				// Hostname
				rvRule, err := jsonpointer.Get(obj,
					fmt.Sprintf("/rules/%d/host", ruleIdx))
				if err != nil {
					c.log.Warn("TODO_2")
					return nil, err
				} else if rule, ok := rvRule.(string); !ok {
					c.log.Warn("TODO_2")
					continue
				} else {
					tmpIngressInfo.Hostname = rule
				}
				// Path
				rvPath, err := jsonpointer.Get(obj,
					fmt.Sprintf("/rules/%d/http/paths/%d/path", ruleIdx, pathIdx))
				if err != nil {
					c.log.Warn("TODO_3")
					tmpIngressInfo.Path = "/"
				} else if path, ok := rvPath.(string); !ok {
					c.log.Warn("TODO_4")
					tmpIngressInfo.Path = "/"
				} else {
					tmpIngressInfo.Path = path
				}
			}
			m.Set(fmt.Sprintf("%d-%d", ruleIdx, pathIdx), tmpIngressInfo)
		}

		// append result from orderedMap
		for _, key := range m.Keys() {
			val, _ := m.Get(key)
			// defaulting
			if val.Name == "" {
				val.Name = ing.GetName()
			}
			if val.Proto == "" {
				val.Proto = "https"
			}
			if val.Tags == nil {
				val.Tags = []string{}
			}
			result = append(result, val)
		}
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
			Hostname:  u.Host,
			Path:      u.Path,
			Proto:     u.Scheme,
			IconUrl:   exLink.Spec.IconURL,
			Tags:      exLink.Spec.Tags,
			IsPrivate: exLink.Spec.IsPrivate,
		})
	}
	return result, nil
}
