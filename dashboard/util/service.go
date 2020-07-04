package util

import (
	"strings"

	vex "github.com/binchencoder/gateway-proto/data"
)

var (
	// ServiceIDsToNames contains the map from service IDs to names.
	ServiceIDsToNames = map[int32]string{}
	// ServiceNamesToIds contains the map from service names to IDs.
	ServiceNamesToIds = map[string]int32{}
)

func init() {
	for id, name := range vex.ServiceId_name {
		name = normalizedServiceName(name)
		ServiceIDsToNames[id] = name
		ServiceNamesToIds[name] = id
	}
}

func normalizedServiceName(name string) string {
	return strings.ToLower(strings.Replace(name, "_", "-", -1))
}
