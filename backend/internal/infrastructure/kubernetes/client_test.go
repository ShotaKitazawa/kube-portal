package kubernetes

import (
	"context"
	"io"
	"log/slog"
	"testing"

	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/tools/cache"

	"github.com/ShotaKitazawa/kube-portal/internal/model"
	"github.com/google/go-cmp/cmp"
)

func newClient(ingressObjs, externalLinkObjs, httprouteObjs []any) *Client {
	ingressStore := cache.NewStore(cache.MetaNamespaceKeyFunc)
	for _, o := range ingressObjs {
		ingressStore.Add(o)
	}
	externalLinkStore := cache.NewStore(cache.MetaNamespaceKeyFunc)
	for _, o := range externalLinkObjs {
		externalLinkStore.Add(o)
	}
	c := &Client{
		log:               slog.New(slog.NewJSONHandler(io.Discard, nil)),
		ingressStore:      ingressStore,
		externalLinkStore: externalLinkStore,
	}
	if httprouteObjs != nil {
		httprouteStore := cache.NewStore(cache.MetaNamespaceKeyFunc)
		for _, o := range httprouteObjs {
			httprouteStore.Add(o)
		}
		c.httprouteStore = httprouteStore
	}
	return c
}

func TestClient_ListIngress(t *testing.T) {
	tests := []struct {
		name      string
		ingresses []*networkingv1.Ingress
		want      model.LinkList
		wantErr   bool
	}{
		{
			name: "testcase 01",
			ingresses: []*networkingv1.Ingress{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test01",
						Namespace: "default",
						Annotations: map[string]string{
							"kube-portal.kanatakita.com/rules.0.paths.0.enable": "true",
						},
					},
					Spec: networkingv1.IngressSpec{
						Rules: []networkingv1.IngressRule{{
							Host: "01.example.com",
							IngressRuleValue: networkingv1.IngressRuleValue{
								HTTP: &networkingv1.HTTPIngressRuleValue{
									Paths: []networkingv1.HTTPIngressPath{{
										Path: "/",
									}},
								},
							},
						}},
					},
				},
			},
			want: model.LinkList{
				{
					Name:      "test01",
					Hostname:  "01.example.com",
					Path:      "/",
					Proto:     "https",
					IconUrl:   "",
					Tags:      []string{},
					IsPrivate: false,
				},
			},
			wantErr: false,
		},
	}
	// TODO: add testcases

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			objs := make([]any, len(tt.ingresses))
			for i, ing := range tt.ingresses {
				objs[i] = ing
			}
			c := newClient(objs, nil, nil)
			got, err := c.ListIngress(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ListIngress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("the result of Client.ListIngress() is unexpected: %s", diff)
			}
		})
	}
}

func TestClient_ListExternalLink(t *testing.T) {
	tests := []struct {
		name          string
		externalLinks []*unstructured.Unstructured
		want          model.LinkList
		wantErr       bool
	}{
		{
			name: "basic external link",
			externalLinks: []*unstructured.Unstructured{
				{
					Object: map[string]any{
						"apiVersion": "kubeportal.kanatakita.com/v1alpha1",
						"kind":       "ExternalLink",
						"metadata": map[string]any{
							"name":      "test-link",
							"namespace": "default",
						},
						"spec": map[string]any{
							"title":     "Test Link",
							"href":      "https://example.com/app",
							"iconURL":   "https://example.com/icon.png",
							"tags":      []any{"prod"},
							"isPrivate": false,
						},
					},
				},
			},
			want: model.LinkList{
				{
					Name:      "Test Link",
					Hostname:  "example.com",
					Path:      "/app",
					Proto:     "https",
					IconUrl:   "https://example.com/icon.png",
					Tags:      []string{"prod"},
					IsPrivate: false,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			objs := make([]any, len(tt.externalLinks))
			for i, el := range tt.externalLinks {
				objs[i] = el
			}
			c := newClient(nil, objs, nil)
			got, err := c.ListExternalLink(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ListExternalLink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("the result of Client.ListExternalLink() is unexpected: %s", diff)
			}
		})
	}
}

func TestClient_ListHTTPRoute(t *testing.T) {
	tests := []struct {
		name       string
		httproutes []*unstructured.Unstructured
		wantNil    bool // nil httprouteStore (Gateway API absent)
		want       model.LinkList
		wantErr    bool
	}{
		{
			name:    "returns nil when Gateway API is not available",
			wantNil: true,
			want:    nil,
			wantErr: false,
		},
		{
			name: "basic httproute",
			httproutes: []*unstructured.Unstructured{
				{
					Object: map[string]any{
						"apiVersion": "gateway.networking.k8s.io/v1",
						"kind":       "HTTPRoute",
						"metadata": map[string]any{
							"name":      "test-route",
							"namespace": "default",
							"annotations": map[string]any{
								"kube-portal.kanatakita.com/rules.0.paths.0.enable": "true",
								"kube-portal.kanatakita.com/rules.0.paths.0.name":   "My App",
							},
						},
						"spec": map[string]any{
							"hostnames": []any{"app.example.com"},
							"rules": []any{
								map[string]any{
									"matches": []any{
										map[string]any{
											"path": map[string]any{
												"type":  "PathPrefix",
												"value": "/app",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			want: model.LinkList{
				{
					Name:      "My App",
					Hostname:  "app.example.com",
					Path:      "/app",
					Proto:     "https",
					IconUrl:   "",
					Tags:      []string{},
					IsPrivate: false,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var httprouteObjs []any
			if !tt.wantNil {
				httprouteObjs = make([]any, len(tt.httproutes))
				for i, r := range tt.httproutes {
					httprouteObjs[i] = r
				}
			}
			c := newClient(nil, nil, httprouteObjs)
			got, err := c.ListHTTPRoute(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ListHTTPRoute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("the result of Client.ListHTTPRoute() is unexpected: %s", diff)
			}
		})
	}
}
