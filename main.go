package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"

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
	http.HandleFunc("/", mainPageHandler)
	http.HandleFunc("/resources/", resourceListHandler)
	http.HandleFunc("/resource/", resourceDetailHandler)

	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// initializeResourceData populates resourceData with resource information
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

// mainPageHandler serves the main page with resource headings
func mainPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := `<html><body>
	<h1>Kubernetes Resources</h1>
	<ul>{{range $resourceType, $resources := .}}
		<li><a href="/resources/{{$resourceType}}">{{$resourceType}}</a></li>
	{{end}}</ul>
	</body></html>`
	t := template.Must(template.New("main").Parse(tmpl))
	t.Execute(w, resourceData)
}

// resourceListHandler serves the list of resources of a specific type
func resourceListHandler(w http.ResponseWriter, r *http.Request) {
	resourceType := r.URL.Path[len("/resources/"):]

	resources, exists := resourceData[resourceType]
	if !exists {
		http.NotFound(w, r)
		return
	}

	tmpl := `<html><body>
	<h1>{{.ResourceType}}</h1>
	<ul>{{range $name, $resource := .Resources}}
		<li><a href="/resource/{{$name}}">{{$name}}</a></li>
	{{end}}</ul>
	</body></html>`
	t := template.Must(template.New("resourceList").Parse(tmpl))
	t.Execute(w, map[string]interface{}{
		"ResourceType": resourceType,
		"Resources":    resources,
	})
}

// resourceDetailHandler serves the YAML and labels of a specific resource
func resourceDetailHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path[len("/resource/"):]

	for resourceType, resources := range resourceData {
		if resource, exists := resources[name]; exists {
			// Remove managedFields section
			unstructured.RemoveNestedField(resource.Object, "metadata", "managedFields")

			// Extract labels
			labels, found, err := unstructured.NestedStringMap(resource.Object, "metadata", "labels")
			if err != nil {
				http.Error(w, "Error extracting labels", http.StatusInternalServerError)
				return
			}

			// Convert to YAML
			resourceYAML, err := yaml.Marshal(resource.Object)
			if err != nil {
				http.Error(w, "Error converting to YAML", http.StatusInternalServerError)
				return
			}

			// Prepare the labels and YAML for output
			labelLines := ""
			if found {
				for key, value := range labels {
					labelLines += fmt.Sprintf("%s: %s\n", key, value)
				}
			} else {
				labelLines = "No labels found."
			}

			w.Header().Set("Content-Type", "text/plain")
			fmt.Fprintf(w, "Resource Type: %s\n\nLabels:\n%s\n\nYAML:\n%s\n", resourceType, labelLines, resourceYAML)
			return
		}
	}

	http.NotFound(w, r)
}

// containsSlash checks if a string contains a slash, indicating it's a subresource
func containsSlash(s string) bool {
	return len(s) > 0 && s[0] == '/'
}
