package controllers

import (
	"flag"
	"path/filepath"

	"github.com/shahincsejnu/extend-k8s-API-with-CRD/pkg/apis/shahin.oka.com/v1alpha1"
	ShahinV1alpha1 "github.com/shahincsejnu/extend-k8s-API-with-CRD/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/client-go/util/workqueue"
)

func Start() {
	var kubeconfig *string

	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	// creates the connection
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}

	// creates the clientset
	clientset, err := ShahinV1alpha1.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// create the teployment watcher
	teploymentListWatcher := cache.NewListWatchFromClient(clientset.ShahinV1alpha1().RESTClient(), "teployments", "default", fields.Everything())

	// create the workqueue
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	// Bind the workqueue to a cache with the help of an informer. This way we make sure that
	// whenever the cache is updated, the teployment key is added to the workqueue.
	// Note that when we finally process the item from the workqueue, we might see a newer version
	// of the teployment than the version which was responsible for triggering the update.
	indexer, informer := cache.NewIndexerInformer(teploymentListWatcher, &v1alpha1.Teployment{}, 0, cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err == nil {
				queue.Add(key)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(newObj)
			if err == nil {
				queue.Add(key)
			}
		},
		DeleteFunc: func(obj interface{}) {
			// IndexerInformer uses a delta queue, therefore for deletes we have to use this key function
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err == nil {
				queue.Add(key)
			}
		},
	}, cache.Indexers{})

	controller := NewController(queue, indexer, informer)

	// we can now warm up the cache for initial synchronization
	// Let's suppose that we knew about a teployment "demo-teployment" on our last run, therefore add it to the cache
	// If this teployment is not there anymore, the controller will be notified about the removal after the cache has synchronized
	indexer.Add(&v1alpha1.Teployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "demo-teployment",
			Namespace: "default",
		},
	})

	// Now let's start the controller

	stop := make(chan struct{})
	defer close(stop)
	go controller.Run(1, stop)

	// wait forever, until user give ctrl+c
	select {}
}
