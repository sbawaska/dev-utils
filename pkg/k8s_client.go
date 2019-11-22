package devutil

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

type K8sClient struct {
	dc dynamic.Interface
}

func NewK8sClient() *K8sClient {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	// TODO remove above
	//config, err := rest.InClusterConfig()
	//if err != nil {
	//	panic(err.Error())
	//}
	// creates the clientset
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return &K8sClient{
		dc:dynamicClient,
	}
}

func (c *K8sClient) GetNestedString(streamName string, gvr schema.GroupVersionResource, fields... string) (string, error) {
	streams, err := c.dc.Resource(gvr).List(v1.ListOptions{})
	if err != nil {
		return "", err
	}

	for _, stream := range streams.Items {
		if stream.GetName() == streamName {
			topic, found, err := unstructured.NestedString(stream.UnstructuredContent(), fields...)
			if err != nil {
				return "", err
			}
			if !found {
				return "", errors.New("unexpected structure of status")
			}
			return topic, nil
		}
	}
	return "", errors.New(fmt.Sprintf("stream: %s not found", streamName))
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}