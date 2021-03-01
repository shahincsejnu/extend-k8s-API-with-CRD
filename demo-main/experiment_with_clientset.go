package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"path/filepath"

	"github.com/davecgh/go-spew/spew"

	ShahinV1alpha1 "github.com/shahincsejnu/extend-k8s-API-with-CRD/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Fatalf("Could not get Kubernetes config: %s", err)
	}

	skc := ShahinV1alpha1.NewForConfigOrDie(config)

	tpmnt, err := skc.ShahinV1alpha1().Teployments("default").Get(context.TODO(), "apiserver-teployment5", metav1.GetOptions{})
	if err != nil {
		fmt.Printf("%v", err.Error())
	} else {
		spew.Dump(tpmnt)
	}
}
