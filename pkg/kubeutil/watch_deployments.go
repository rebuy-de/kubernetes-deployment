package kubeutil

import (
	"context"
	"time"

	apps "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

func WatchDeployments(ctx context.Context, client kubernetes.Interface, selector fields.Selector) chan *apps.Deployment {
	lw := cache.NewListWatchFromClient(
		client.AppsV1().RESTClient(),
		"deployments",
		v1.NamespaceAll,
		selector)

	stop := make(chan struct{}, 1)
	results := make(chan *apps.Deployment)

	store, controller := cache.NewInformer(
		lw,
		&apps.Deployment{},
		60*time.Second,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				results <- obj.(*apps.Deployment)
			},
			UpdateFunc: func(old, obj interface{}) {
				results <- obj.(*apps.Deployment)
			},
			DeleteFunc: func(obj interface{}) {
				results <- obj.(*apps.Deployment)
			},
		})

	for _, obj := range store.List() {
		results <- obj.(*apps.Deployment)
	}

	go controller.Run(stop)

	go func() {
		<-ctx.Done()
		close(stop)
		close(results)
	}()

	return results
}
