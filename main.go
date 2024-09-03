package main

type PodDetails struct {
	PodName      string                       `json:"PodName"`
	PodNamespace string                       `json:"PodNamespace"`
	//PodYaml      string                       `json:"PodYaml"`
	PodLabels    map[string]map[string]string `json:"PodLabels"`
}


func main(


	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

)
