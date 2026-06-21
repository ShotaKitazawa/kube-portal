package kubernetes

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/elliotchance/orderedmap/v3"
	"github.com/mattn/go-jsonpointer"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"

	kubeportalv1alpha1 "github.com/ShotaKitazawa/kube-portal/internal/infrastructure/kubernetes/api/v1alpha1"
	"github.com/ShotaKitazawa/kube-portal/internal/model"
	"github.com/ShotaKitazawa/kube-portal/internal/model/port"
)

const ingressAnnotationPrefix = "kube-portal.kanatakita.com/"

var ingressAnnotationMatcher *regexp.Regexp

var externalLinkGVR = schema.GroupVersionResource{
	Group:    kubeportalv1alpha1.GroupVersion.Group,
	Version:  kubeportalv1alpha1.GroupVersion.Version,
	Resource: "externallinks",
}

func init() {
	ingressAnnotationMatcher = regexp.MustCompile(
		`^rules\.(\d+)\.paths\.(\d+)\.(.+)$`)
}

type Client struct {
	log               *slog.Logger
	ingressStore      cache.Store
	externalLinkStore cache.Store
}

var _ port.Kubernetes = (*Client)(nil)

func NewClient(ctx context.Context, logger *slog.Logger, kubeconfigPath string) (*Client, error) {
	kubeConfig, err := buildConfig(kubeconfigPath)
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return nil, err
	}
	dynClient, err := dynamic.NewForConfig(kubeConfig)
	if err != nil {
		return nil, err
	}

	stopCh := ctx.Done()

	factory := informers.NewSharedInformerFactory(clientset, 0)
	ingressInformer := factory.Networking().V1().Ingresses().Informer()

	dynFactory := dynamicinformer.NewDynamicSharedInformerFactory(dynClient, 0)
	externalLinkInformer := dynFactory.ForResource(externalLinkGVR).Informer()

	factory.Start(stopCh)
	dynFactory.Start(stopCh)

	syncCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	if !cache.WaitForCacheSync(syncCtx.Done(), ingressInformer.HasSynced, externalLinkInformer.HasSynced) {
		return nil, fmt.Errorf("timed out waiting for caches to sync")
	}

	return &Client{
		log:               logger,
		ingressStore:      ingressInformer.GetStore(),
		externalLinkStore: externalLinkInformer.GetStore(),
	}, nil
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

func (c *Client) ListIngress(_ context.Context) (model.LinkList, error) {
	var result model.LinkList
	for _, item := range c.ingressStore.List() {
		ing, ok := item.(*networkingv1.Ingress)
		if !ok {
			c.log.Warn("unexpected object type in ingress store")
			continue
		}

		m := orderedmap.NewOrderedMap[string, model.Link]()
		logger := c.log.With(
			slog.String("name", ing.Name),
			slog.String("namespace", ing.Namespace),
		)
		for key, val := range ing.Annotations {
			logger = logger.With(
				slog.String("annotation", key),
			)
			// Skip if it's not related annotation
			if !strings.HasPrefix(key, ingressAnnotationPrefix) {
				continue
			}
			key = strings.TrimPrefix(key, ingressAnnotationPrefix)

			// Extract value from annotations
			l := ingressAnnotationMatcher.FindStringSubmatch(key)
			if len(l) != 4 {
				logger.Warn("skip: annotation didn't match the matcher")
				continue
			}
			ruleIdx, _ := strconv.Atoi(l[1])
			pathIdx, _ := strconv.Atoi(l[2])
			action := l[3]

			//
			// Fill model.IngressInfo
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
				logger.Warn("skip: annotation is not supported")
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
					logger.Warn(fmt.Sprintf(
						"skip: Ingress.spec.rules[%d].host is empty", ruleIdx))
					continue
				}
				rule, ok := rvRule.(string)
				if !ok {
					logger.Warn(fmt.Sprintf(
						"skip: Ingress.spec.rules[%d].host is empty", ruleIdx))
					continue
				}
				tmpIngressInfo.Hostname = rule
				// Path
				rvPath, err := jsonpointer.Get(obj,
					fmt.Sprintf("/rules/%d/http/paths/%d/path", ruleIdx, pathIdx))
				if err != nil {
					logger.Debug(fmt.Sprintf(
						"Ingress.spec.rules[%d].http.paths[%d].path is empty: use / as path",
						ruleIdx, pathIdx))
					tmpIngressInfo.Path = "/"
				}
				path, ok := rvPath.(string)
				if !ok {
					logger.Debug(fmt.Sprintf(
						"Ingress.spec.rules[%d].http.paths[%d].path is empty: use / as path",
						ruleIdx, pathIdx))
					tmpIngressInfo.Path = "/"
				}
				tmpIngressInfo.Path = path
			}
			m.Set(fmt.Sprintf("%d-%d", ruleIdx, pathIdx), tmpIngressInfo)
		}

		// append result from orderedMap
		for key := range m.Keys() {
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

func (c *Client) ListExternalLink(_ context.Context) (model.LinkList, error) {
	if c.externalLinkStore == nil {
		return nil, nil
	}
	var result []model.Link
	for _, item := range c.externalLinkStore.List() {
		u, ok := item.(*unstructured.Unstructured)
		if !ok {
			c.log.Warn("unexpected object type in external link store")
			continue
		}
		exLink := kubeportalv1alpha1.ExternalLink{}
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, &exLink); err != nil {
			return nil, err
		}
		/* get & validate values */
		href, err := url.Parse(exLink.Spec.Href)
		if err != nil {
			return nil, err
		}
		/* set values to result */
		result = append(result, model.Link{
			Name:      exLink.Spec.Title,
			Hostname:  href.Host,
			Path:      href.Path,
			Proto:     href.Scheme,
			IconUrl:   exLink.Spec.IconURL,
			Tags:      exLink.Spec.Tags,
			IsPrivate: exLink.Spec.IsPrivate,
		})
	}
	return result, nil
}
