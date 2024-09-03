package main

import (
	"context"
	"fmt"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
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

	// Create a Dynamic Client and a Kubernetes Clientset
	dynClient, err := dynamic.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating dynamic client: %v", err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating Kubernetes clientset: %v", err)
	}

	// Create a Discovery Client
	discoveryClient := clientset.Discovery()

	// Get the list of all API resources available
	serverResources, err := discoveryClient.ServerPreferredResources()
	if err != nil {
		log.Fatalf("Error retrieving server preferred resources: %v", err)
	}

	// Iterate over all resources and list them
	for _, resourceList := range serverResources {
		for _, resource := range resourceList.APIResources {
			// Skip subresources (like pod/logs, pod/status) and non-listable resources
			if containsSlash(resource.Name) || !resource.Namespaced {
				continue
			}

			gvr := schema.GroupVersionResource{
				Group:    resourceList.GroupVersion,
				Version:  resource.Version,
				Resource: resource.Name,
			}

			// Adjust for core group
			if gvr.Group == "v1" {
				gvr.Version = gvr.Group
				gvr.Group = ""
			}

			// List the resources
			list, err := dynClient.Resource(gvr).Namespace(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				fmt.Printf("Error listing %s: %v\n", gvr.Resource, err)
				continue
			}

			// Print the resources
			printResources(list, gvr.Resource)
		}
	}
}

// containsSlash checks if a string contains a slash, indicating it's a subresource
func containsSlash(s string) bool {
	return len(s) > 0 && s[0] == '/'
}

// printResources prints the count and names of resources
func printResources(list *unstructured.UnstructuredList, resourceName string) {
	fmt.Printf("Found %d %s\n", len(list.Items), resourceName)
	for _, item := range list.Items {
		fmt.Printf("- %s\n", item.GetName())
	}
}
