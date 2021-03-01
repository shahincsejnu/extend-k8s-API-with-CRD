// followed this tuto: https://github.com/kubernetes/client-go/blob/master/examples/workqueue/main.go

package main

import (
	"flag"
	"fmt"
	"github.com/shahincsejnu/extend-k8s-API-with-CRD/pkg/apis/shahin.oka.com/v1alpha1"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
	"time"
	"k8s.io/klog/v2"
	ShahinV1alpha1 "github.com/shahincsejnu/extend-k8s-API-with-CRD/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"
)

// Controller demonstrates how to implement a controller with client-go.
type Controller struct {
	indexer cache.Indexer
	queue workqueue.RateLimitingInterface
	informer cache.Controller

	//crdClient c
	//kClient kub
}

// NewController creates a new Controller
func NewController(queue workqueue.RateLimitingInterface, indexer cache.Indexer, informer cache.Controller) *Controller {
	return &Controller{
		indexer:  indexer,
		queue:    queue,
		informer: informer,
	}
}

// Run begins watching and syncing
func (c *Controller) Run(threadiness int, stopCh chan struct{}) {
	defer runtime.HandleCrash()

	// let the workers stop when we are done
	defer c.queue.ShutDown()
	fmt.Println("Starting Teployment Controller")

	go c.informer.Run(stopCh)

	// wait for all involved caches to be synced, before processing items from the queue is started
	if !cache.WaitForCacheSync(stopCh, c.informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("Time out waiting for caches to sync"))
		return
	}

	for i:= 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	<-stopCh
	fmt.Println("Stopping Teployment Controller")

}

func (c *Controller) runWorker() {
	for c.processNextItem() {

	}
}

func (c *Controller) processNextItem() bool {
	// wait until there is a new item in the working queue
	key, quit := c.queue.Get()
	if quit {
		return false
	}

	// Tell the queue that we are done with processing this key. This unblocks the key for other workers
	// This allows safe parallel processing because two teployments with the same key are never processed in parallel
	defer c.queue.Done(key)

	// Invoke the method containing the business logic
	err := c.syncToStdout(key.(string))
	// Handle the error if something went wrong during the execution of the business logic
	c.handleErr(err, key)

	return true
}

// syncToStdout is the business logic of controller. In this controller it simple prints
// information about the teployment of stdout. In case an error happened, it has to simple return the error.
// The retry logic should not be part of the business logic
func (c *Controller) syncToStdout(key string) error {
	obj, exists, err := c.indexer.GetByKey(key)
	if err != nil {
		fmt.Errorf("Fetching object with key %s from store failed with %v", key, err)
		return err
	}

	if !exists {
		// below we will warm up our cache with a Teployment, so that we will see a delete for one teployment
		fmt.Printf("Teployment %s does not exist anymore\n", key)
	} else {
		// Note that you also have to check the uid if you have a local controlled resource, which
		// is dependent on the actual instance, to detect that a teployment was recreated with the same name
		fmt.Printf("Sync/Add/Update for Teployment %s\n", obj.(*v1alpha1.Teployment).GetName())

		//sdfds := obj.(v1alpha1.Teployment)
		//process()
	}

	return nil
}

//func (c *Controller) process(tep *v1alpha1.Teployment) error {
//	// Create depl
//	// Create service
//}

// handleErr checks if an error happened and makes sure we will retry later
func (c *Controller) handleErr(err error, key interface{}) {
	if err == nil {
		// Forget about the #AddRateLimited history of the key on every successful synchronization
		// This ensures that future processing of updates for this key is not delayed because of
		// an outdated error history
		c.queue.Forget(key)
		return
	}

	// This controller retries 5 times if something goes wrong. After that, it stops trying
	if c.queue.NumRequeues(key) < 5 {
		klog.Infof("Error syncing teployment %v: %v", key, err)

		// Re-enqueue the key rate limited. Based on the rate limiter on the
		// queue and the re-enqueue history, the key will be processed later again
		c.queue.AddRateLimited(key)
		return
	}

	c.queue.Forget(key)
	// Report to an external entity that, even after several retries, we could not successfully process this key
	runtime.HandleError(err)
	klog.Infof("Dropping teployment %q out of the queue: %v", key, err)
}

func main() {
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
			Name: "apiserver-teployment",
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
