package main

import (
	"context"
	"fmt"
	"log"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/yaml"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	// Define the GroupVersionResource for Pods
	gvr := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "pods",
	}

	// List all Pods in all namespaces
	podList, err := dynClient.Resource(gvr).Namespace(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	// Print the list of Pods with their YAML manifests
	printResourcesWithYAML(podList, "pods", dynClient, gvr)
}

func printResourcesWithYAML(list *unstructured.UnstructuredList, resourceName string, dynClient dynamic.Interface, gvr schema.GroupVersionResource) {
	fmt.Printf("Found %d %s\n", len(list.Items), resourceName)
	for _, item := range list.Items {
		fmt.Printf("- %s\n", item.GetName())

		// Retrieve the full unstructured object for the pod
		pod, err := dynClient.Resource(gvr).Namespace(item.GetNamespace()).Get(context.TODO(), item.GetName(), metav1.GetOptions{})
		if err != nil {
			log.Printf("Error retrieving pod %s: %v", item.GetName(), err)
			continue
		}

		// Remove the managedFields section
		unstructured.RemoveNestedField(pod.Object, "metadata", "managedFields")

		// Convert the unstructured object to YAML
		podYAML, err := yaml.Marshal(pod.Object)
		if err != nil {
			log.Printf("Error converting pod %s to YAML: %v", item.GetName(), err)
			continue
		}

		// Print the YAML manifest
		fmt.Printf("---\n%s\n---\n", string(podYAML))
	}
}
