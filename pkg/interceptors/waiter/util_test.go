package waiter

import (
	"testing"

	"github.com/rebuy-de/kubernetes-deployment/pkg/kubeutil"
)

func TestIsDeployment(t *testing.T) {
	cases := []struct {
		name     string
		want     bool
		manifest string
	}{
		{
			name: "apps_v1",
			want: true,
			manifest: `apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.7.9
        ports:
        - containerPort: 80`,
		},
		{
			name: "extensions_v1beta1",
			want: true,
			manifest: `apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.7.9
        ports:
        - containerPort: 80`,
		},
		{
			name: "service",
			want: false,
			manifest: `apiVersion: v1
kind: Service
metadata:
  name: my-service
spec:
  selector:
    app: MyApp
  ports:
    - protocol: TCP
      port: 80
      targetPort: 9376`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			obj, err := kubeutil.Decode([]byte(tc.manifest))
			if err != nil {
				t.Fatal(err)
			}

			have := isDeployment(obj)
			if tc.want != have {
				t.Fatalf("Assertion failed. Want: %t. Have: %t.", tc.want, have)
			}
		})
	}

}
