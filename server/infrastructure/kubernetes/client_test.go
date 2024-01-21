package kubernetes

import (
	"context"
	"io"
	"log/slog"
	"testing"

	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/ShotaKitazawa/kube-portal/server/models"
	"github.com/google/go-cmp/cmp"
)

func TestClient_ListIngressInfo(t *testing.T) {
	type fields struct {
		clientset kubernetes.Interface
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    models.IngressInfoList
		wantErr bool
	}{
		{
			name: "testcase 01",
			fields: fields{
				fake.NewSimpleClientset(
					&networkingv1.Ingress{
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
				),
			},
			args: args{context.Background()},
			want: models.IngressInfoList{
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
			c := &Client{
				log:       slog.New(slog.NewJSONHandler(io.Discard, nil)),
				clientset: tt.fields.clientset,
			}
			got, err := c.ListIngress(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ListIngressInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("the result of Client.ListIngressInfo() is unexpected: %s", diff)
			}
		})
	}
}
