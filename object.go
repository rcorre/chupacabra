package main

import "encoding/json"

// Object represents a generic kubernetes object
type Object struct {
	Name string
	Body string
}

func ParseObjectList(data []byte) ([]Object, error) {
	var list struct {
		Items []map[string]interface{}
	}

	if err := json.Unmarshal(data, &list); err != nil {
		return nil, err
	}

	objects := make([]Object, len(list.Items))

	for i, item := range list.Items {
		meta := item["metadata"].(map[string]interface{})
		name := meta["name"].(string)

		body, err := json.MarshalIndent(item, "", "    ")
		if err != nil {
			return nil, err
		}

		objects[i] = Object{
			Name: name,
			Body: string(body),
		}
	}
	return objects, nil
}
