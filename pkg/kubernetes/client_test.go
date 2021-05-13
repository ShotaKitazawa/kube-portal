package kubernetes

import (
	"context"
	"reflect"
	"testing"

	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/ShotaKitazawa/kube-portal/pkg/entities"
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
		want    []entities.IngressInfo
		wantErr bool
	}{
		{
			name: "testcase 01",
			fields: fields{
				fake.NewSimpleClientset(
					&networkingv1beta1.Ingress{
						ObjectMeta: metav1.ObjectMeta{
							Name:        "test01",
							Namespace:   "default",
							Annotations: map[string]string{},
						},
						Spec: networkingv1beta1.IngressSpec{
							Rules: []networkingv1beta1.IngressRule{{Host: "01.example.com"}},
						},
					},
				),
			},
			args: args{context.Background()},
			want: []entities.IngressInfo{
				{
					Fqdn:      "01.example.com",
					Proto:     "https",
					IconUrl:   entities.DefaultIconUrl,
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
				clientset: tt.fields.clientset,
			}
			got, err := c.ListIngressInfo(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ListIngressInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.ListIngressInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}
