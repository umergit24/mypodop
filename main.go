package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/yaml"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	// Set up the HTTP server routes
	r := mux.NewRouter()
	r.HandleFunc("/", listPodsHandler)
	r.HandleFunc("/pod/{namespace}/{name}", podYAMLHandler)

	fmt.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func listPodsHandler(w http.ResponseWriter, r *http.Request) {
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

	// Generate HTML with the list of pods
	var builder strings.Builder
	builder.WriteString("<html><body>")
	builder.WriteString("<h1>Pod List</h1><ul>")

	for _, item := range podList.Items {
		podName := item.GetName()
		namespace := item.GetNamespace()
		builder.WriteString(fmt.Sprintf("<li><a href=\"/pod/%s/%s\">%s (%s)</a></li>", namespace, podName, podName, namespace))
	}

	builder.WriteString("</ul></body></html>")

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(builder.String()))
}

func podYAMLHandler(w http.ResponseWriter, r *http.Request) {
	// Get the pod name and namespace from the URL
	vars := mux.Vars(r)
	podName := vars["name"]
	namespace := vars["namespace"]

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

	// Retrieve the pod by name and namespace
	pod, err := dynClient.Resource(gvr).Namespace(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving pod %s: %v", podName, err), http.StatusInternalServerError)
		return
	}

	// Remove the managedFields section
	unstructured.RemoveNestedField(pod.Object, "metadata", "managedFields")

	// Convert the unstructured object to YAML
	podYAML, err := yaml.Marshal(pod.Object)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error converting pod %s to YAML: %v", podName, err), http.StatusInternalServerError)
		return
	}

	// Serve the YAML as plain text
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(podYAML))
}
