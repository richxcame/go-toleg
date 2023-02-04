package config

import (
	"encoding/json"
	"os"
)

var Clients []map[string]string

func init() {
	Clients = GetClients()
}

func GetClients() []map[string]string {
	var clients []map[string]string
	clientJSON, err := os.ReadFile("clients.json")
	if err != nil {
		return []map[string]string{}
	}
	err = json.Unmarshal(clientJSON, &clients)
	if err != nil {
		return []map[string]string{}
	}
	return clients
}
