package main

import (
	"encoding/json"
	"flag"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Kube interface {
	Get(kind, namespace string) (map[string]interface{}, error)
}

type kube struct {
	rest *rest.RESTClient
}

func NewKube(configPath string) Kube {
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		panic(err.Error())
	}

	config.NegotiatedSerializer = serializer.DirectCodecFactory{
		CodecFactory: scheme.Codecs,
	}

	config.GroupVersion = &v1.SchemeGroupVersion
	config.APIPath = "/api"
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}
	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}
	client, err := rest.RESTClientFor(config)
	if err != nil {
		panic(err.Error())
	}
	return &kube{client}
}

func (k *kube) Get(kind, namespace string) (map[string]interface{}, error) {
	req := k.rest.Verb("GET").
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
