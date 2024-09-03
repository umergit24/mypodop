package main

import (
	"context"
	"fmt"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// Load the Kubernetes configuration
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

	config, err := kubeconfig.ClientConfig()
	if err != nil {
		log.Fatalf("Error getting Kubernetes config: %v", err)
	}

	// Create a Dynamic Client
	dynClient, err := dynamic.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating dynamic client: %v", err)
	}

	// Define the GroupVersionResource for deployments
	gvr := schema.GroupVersionResource{
		Group:    "apps",
		Version:  "v1",
		Resource: "deployments",
	}

	// List all deployments in all namespaces
	podList, err := dynClient.Resource(gvr).Namespace(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	// Print the list of deployments
	printResources(podList, "deployments")

}

func printResources(list *unstructured.UnstructuredList, resourceName string) {
	fmt.Printf("Found %d %s\n", len(list.Items), resourceName)
	for _, item := range list.Items {
		fmt.Printf("- %s\n", item.GetName())
	}
}
