package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/yaml"
)

var (
	resourceData = make(map[string]map[string]*unstructured.Unstructured)
	dynClient    dynamic.Interface
)

func main() {
	// Load Kubernetes configuration
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

	config, err := kubeconfig.ClientConfig()
	if err != nil {
		log.Fatalf("Error getting Kubernetes config: %v", err)
	}

	// Create a Dynamic Client and a Kubernetes Clientset
	dynClient, err = dynamic.NewForConfig(config)
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

	// Initialize resource data
	initializeResourceData(serverResources)

	// Set up HTTP routes
	http.Handle("/", enableCors(http.HandlerFunc(mainPageHandler)))
	http.Handle("/list/", enableCors(http.HandlerFunc(resourceListHandler)))
	http.Handle("/details/", enableCors(http.HandlerFunc(resourceDetailHandler)))


	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func initializeResourceData(serverResources []*metav1.APIResourceList) {
	for _, resourceList := range serverResources {
		for _, resource := range resourceList.APIResources {
			if containsSlash(resource.Name) || !resource.Namespaced {
				continue
			}

			gvr := schema.GroupVersionResource{
				Group:    resourceList.GroupVersion,
				Version:  resource.Version,
				Resource: resource.Name,
			}

			if gvr.Group == "v1" {
				gvr.Version = gvr.Group
				gvr.Group = ""
			}

			// List the resources
			list, err := dynClient.Resource(gvr).Namespace(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				log.Printf("Error listing %s: %v\n", gvr.Resource, err)
				continue
			}

			// Store resources in memory
			if resourceData[gvr.Resource] == nil {
				resourceData[gvr.Resource] = make(map[string]*unstructured.Unstructured)
			}
			for _, item := range list.Items {
				resourceData[gvr.Resource][item.GetName()] = &item
			}
		}
	}
}

// mainPageHandler serves the list of all resource types in JSON
func mainPageHandler(w http.ResponseWriter, r *http.Request) {
	resourceTypes := make([]string, 0, len(resourceData))
	for resourceType := range resourceData {
		resourceTypes = append(resourceTypes, resourceType)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resourceTypes)
}

// resourceListHandler serves the list of resources for a specific type in JSON
func resourceListHandler(w http.ResponseWriter, r *http.Request) {
	resourceType := r.URL.Path[len("/list/"):]

	resources, exists := resourceData[resourceType]
	if !exists {
		http.NotFound(w, r)
		return
	}

	resourceNames := make([]string, 0, len(resources))
	for name := range resources {
		resourceNames = append(resourceNames, name)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resourceNames)
}


// resourceDetailHandler serves the details of a specific resource in YAML
func resourceDetailHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.SplitN(r.URL.Path[len("/details/"):], "/", 2)
	if len(parts) != 2 {
		http.NotFound(w, r)
		return
	}
	resourceType, name := parts[0], parts[1]

	resources, exists := resourceData[resourceType]
	if !exists {
		http.NotFound(w, r)
		return
	}

	resource, exists := resources[name]
	if !exists {
		http.NotFound(w, r)
		return
	}

	// Remove managedFields section before converting to YAML
	unstructured.RemoveNestedField(resource.Object, "metadata", "managedFields")

	// Convert to YAML
	resourceYAML, err := yaml.Marshal(resource.Object)
	if err != nil {
		http.Error(w, "Error converting to YAML", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/yaml")
	w.Write(resourceYAML)
}


// containsSlash checks if a string contains a slash, indicating it's a subresource
func containsSlash(s string) bool {
	return len(s) > 0 && s[0] == '/'
}


func enableCors(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")  // Allow any domain, replace "*" with specific domain for production
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        // Handle preflight requests (OPTIONS method)
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        next.ServeHTTP(w, r)
    })
}
