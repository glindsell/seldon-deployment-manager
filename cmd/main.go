package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	"encoding/json"

	v1alpha2client "github.com/SeldonIO/seldon-operator/pkg/client/clientset/versioned/typed/machinelearning/v1alpha2"
	v1alpha2api "github.com/seldonio/seldon-operator/pkg/apis/machinelearning/v1alpha2"

	"github.com/glindsell/seldon/internal/seldonwrapper"
)

func main() {
	var kubeconfig *string
	model := os.Args[1]

	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}

	seldonClientset, err := v1alpha2client.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	file, err := os.Open(model)
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err.Error())
	}

	var dep v1alpha2api.SeldonDeployment

	err = json.Unmarshal(byteValue, &dep)
	if err != nil {
		panic(err.Error())
	}

	sc, err := seldonwrapper.NewSeldonWrapper(*seldonClientset)
	if err != nil {
		panic(err)
	}

	//TODO: Fix timestamp workaround
	dep.Spec.Predictors[0].ComponentSpecs[0].Metadata.CreationTimestamp = metav1.Now()

	log.Printf("Creating resource...\n")

	err = sc.CreateDeployment(&dep)
	if err != nil {
		panic(err)
	}

	c1 := make(chan string)
	c2 := make(chan string)

	go func() {
		err = sc.WaitForDeploymentState(dep.Name, "Available")
		if err != nil {
			panic(err)
		}
		c1 <- "Resource available."
	}()

	msg1 := <-c1
	log.Printf("%v\n", msg1)

	err = sc.DeleteDeployment(dep.Name)
	if err != nil {
		panic(err)
	}

	go func() {
		time.Sleep(120 * time.Second)
		c2 <- "Timed out waiting for deletion to finish."
	}()

	go func() {
		err = sc.WaitForDeploymentNotFound(dep.Name)
		if err != nil {
			panic(err)
		}
		c2 <- "Resource deleted."
	}()

	msg2 := <-c2
	log.Printf("%v\n", msg2)
	log.Printf("Exiting.\n")
}
