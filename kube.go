package main

import (
	"encoding/json"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Kube interface {
	Resources() []string
	Get(kind, namespace string) (map[string]interface{}, error)
}

type kube struct {
	clients map[string]rest.Interface
}

func NewKube(configPath string) (Kube, error) {
	config, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	clients := map[string]rest.Interface{
		"componentstatuses":      clientset.CoreV1().RESTClient(),
		"configmaps":             clientset.CoreV1().RESTClient(),
		"endpoints":              clientset.CoreV1().RESTClient(),
		"events":                 clientset.CoreV1().RESTClient(),
		"limitranges":            clientset.CoreV1().RESTClient(),
		"namespaces":             clientset.CoreV1().RESTClient(),
		"nodes":                  clientset.CoreV1().RESTClient(),
		"persistentvolumes":      clientset.CoreV1().RESTClient(),
		"persistentvolumeclaims": clientset.CoreV1().RESTClient(),
		"pods":                   clientset.CoreV1().RESTClient(),
		"podtemplates":           clientset.CoreV1().RESTClient(),
		"replicationcontrollers": clientset.CoreV1().RESTClient(),
		"resourcequotas":         clientset.CoreV1().RESTClient(),
		"secrets":                clientset.CoreV1().RESTClient(),
		"services":               clientset.CoreV1().RESTClient(),
		"serviceaccounts":        clientset.CoreV1().RESTClient(),
	}
	return &kube{clients}, nil
}

func (k *kube) Resources() []string {
	keys := make([]string, len(k.clients))
	i := 0
	for name := range k.clients {
		keys[i] = name
		i++
	}
	return keys
}

func (k *kube) Get(kind, namespace string) (map[string]interface{}, error) {
	req := k.clients[kind].Verb("GET").
		Resource(kind).
		Namespace(namespace).
		Timeout(10 * time.Second)

	var data struct {
		Items []map[string]interface{}
	}

	if raw, err := req.DoRaw(); err != nil {
		return nil, err
	} else if err := json.Unmarshal(raw, &data); err != nil {
		return nil, err
	}

	ret := make(map[string]interface{})

	for _, item := range data.Items {
		meta := item["metadata"].(map[string]interface{})
		name := meta["name"].(string)
		ret[name] = item
	}
	return ret, nil
}
