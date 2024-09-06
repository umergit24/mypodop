package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	http.HandleFunc("/", servePods)
	fmt.Println("Starting server at http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func servePods(w http.ResponseWriter, r *http.Request) {
	// Load the Kubernetes configuration
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

	config, err := kubeconfig.ClientConfig()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting Kubernetes config: %v", err), http.StatusInternalServerError)
		return
	}

	// Create a Dynamic Client
	dynClient, err := dynamic.NewForConfig(config)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating dynamic client: %v", err), http.StatusInternalServerError)
		return
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
		http.Error(w, fmt.Sprintf("Error listing pods: %v", err), http.StatusInternalServerError)
		return
	}

	// Serve the list of Pods as HTML
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<html><body>")
	fmt.Fprintf(w, "<h1>Pods List</h1>")
	fmt.Fprintf(w, "<p>Found %d pods:</p>", len(podList.Items))
	fmt.Fprintf(w, "<ul>")
	for _, item := range podList.Items {
		fmt.Fprintf(w, "<li>%s</li>", item.GetName())
	}
	fmt.Fprintf(w, "</ul>")
	fmt.Fprintf(w, "</body></html>")
}
