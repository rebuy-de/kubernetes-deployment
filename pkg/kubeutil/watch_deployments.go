package kubeutil

import (
	"context"
	"time"

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
	"k8s.io/client-go/tools/cache"
)

func WatchDeployments(ctx context.Context, client kubernetes.Interface, selector fields.Selector) chan *v1beta1.Deployment {
	lw := cache.NewListWatchFromClient(
		client.ExtensionsV1beta1().RESTClient(),
		"deployments",
		api.NamespaceAll,
		selector)

	stop := make(chan struct{}, 1)
	results := make(chan *v1beta1.Deployment)

	store, controller := cache.NewInformer(
		lw,
		&v1beta1.Deployment{},
		60*time.Second,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				results <- obj.(*v1beta1.Deployment)
			},
			UpdateFunc: func(old, obj interface{}) {
				results <- obj.(*v1beta1.Deployment)
			},
			DeleteFunc: func(obj interface{}) {
				results <- obj.(*v1beta1.Deployment)
			},
		})

	for _, obj := range store.List() {
		results <- obj.(*v1beta1.Deployment)
	}

	go controller.Run(stop)

	go func() {
		<-ctx.Done()
		close(stop)
		close(results)
	}()

	return results
}
