package main

import "encoding/json"

func UnmarshalShadow(data []byte) (Shadow, error) {
	var r Shadow
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Shadow) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Shadow struct {
	Device    string   `json:"device"`
	Name      string   `json:"name"`
	Metadata  Metadata `json:"metadata"`
	State     State    `json:"state"`
	Timestamp string   `json:"timestamp"`
	Version   int64    `json:"version"`
}

type Metadata struct {
	Reported map[string]interface{} `json:"reported"`
	Desired  map[string]interface{} `json:"desired"`
}

type State struct {
	Reported map[string]interface{} `json:"reported"`
	Desired  map[string]interface{} `json:"desired"`
}
