package kubeutil

import (
	"context"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

func WatchPods(ctx context.Context, client kubernetes.Interface, selector fields.Selector) chan *v1.Pod {
	lw := cache.NewListWatchFromClient(
		client.CoreV1().RESTClient(),
		"pods",
		v1.NamespaceAll,
		selector)

	stop := make(chan struct{}, 1)
	results := make(chan *v1.Pod)

	store, controller := cache.NewInformer(
		lw,
		&v1.Pod{},
		60*time.Second,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				results <- obj.(*v1.Pod)
			},
			UpdateFunc: func(old, obj interface{}) {
				results <- obj.(*v1.Pod)
			},
			DeleteFunc: func(obj interface{}) {
				results <- obj.(*v1.Pod)
			},
		})

	for _, obj := range store.List() {
		results <- obj.(*v1.Pod)
	}

	go controller.Run(stop)

	go func() {
		<-ctx.Done()
		close(stop)
		close(results)
	}()

	return results
}
