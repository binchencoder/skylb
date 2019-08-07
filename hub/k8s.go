package hub

import (
	"fmt"
	"time"

	"github.com/golang/glog"
	api "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

var (
	k8sv1 corev1.CoreV1Interface
)

func getClient() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}

func (eh *endpointsHub) startK8sWatcher() {
	k8sCliset, err := getClient()
	if err != nil {
		glog.Fatalln(err)
	}
	glog.Infoln("Established connection to Kubernetes API server.")

	k8sv1 = k8sCliset.CoreV1()
	wlist := cache.NewListWatchFromClient(k8sCliset.CoreV1().RESTClient(), "endpoints", metav1.NamespaceAll, fields.Everything())
	go func() {
		_, controller := cache.NewInformer(wlist, &api.Endpoints{}, 60*time.Second,
			cache.ResourceEventHandlerFuncs{
				AddFunc: func(obj interface{}) {},
				UpdateFunc: func(oldObj, newObj interface{}) {
					if newEps, ok := newObj.(*api.Endpoints); ok {
						key := eh.calculateKey(newEps.ObjectMeta.Namespace, newEps.ObjectMeta.Name)
						newEps, _ := newObj.(*api.Endpoints)
						eh.applyK8sEndpoints(key, newEps)
					}
				},
				DeleteFunc: func(obj interface{}) {},
			})

		stop := make(chan struct{})
		defer close(stop)
		controller.Run(stop)
	}()
}

func (eh *endpointsHub) applyK8sEndpoints(key string, eps *api.Endpoints) {
	var so *serviceObject
	eh.WithRLock(func() error {
		so, _ = eh.services[key]
		return nil
	})
	if so == nil {
		// It's normal to get nil service objects because not all services
		// use SkyLB.
		glog.V(6).Infof("serviceObject is nil for key %#v", key)
		return
	}
	eh.applyEndpoints(so, eps)
}

func (eh *endpointsHub) fetchK8sEndpoints(namespace, serviceName string) (*api.Endpoints, error) {
	epsList, err := k8sv1.Endpoints(namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, eps := range epsList.Items {
		if eps.Name == serviceName {
			return &eps, nil
		}
	}
	return nil, fmt.Errorf("endpoints for service %s.%s was not found", namespace, serviceName)
}
