## Kubernetes Resource Explorer

This application retrieves and displays information about Kubernetes resources running on a cluster. It provides a user-friendly web interface to explore Pods, Services and other resources.

### Inputs

This application reads Kubernetes configuration from the default location (`~/.kube/config`) and connects to the cluster accordingly.

### Outputs

This application serves a React web application on `http://localhost:8080`.  The web application displays the following information:

- List of all available Kubernetes resource types.  This list is retrieved dynamically from the cluster, ensuring it reflects the available APIs.
- For each resource type, a list of all resources of that type running on the cluster. This information includes resources across all namespaces.
- Detailed information (manifest file) for each individual resource. The manifest can be viewed as either JSON or YAML. The `managedFields` section is removed from the YAML output for clarity. 

The application uses CORS to allow any domain to access the data, making it easy to integrate the frontend with a backend running on a different port or domain.
