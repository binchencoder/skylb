package util

import (
	"fmt"
	"path"
)

const (
	EndpointsKeyPrefix   = "/registry/services/endpoints"
	DefaultTargetRefKind = "Pod"
)

// CalculateKey returns the ETCD key for the given service.
func CalculateKey(namespace, serviceName string) string {
	return path.Join(EndpointsKeyPrefix, namespace, serviceName)
}

// CalculateEndpointKey returns the ETCD key for the given endpoint.
func CalculateEndpointKey(namespace, serviceName, host string, port int32) string {
	return path.Join(EndpointsKeyPrefix, namespace, serviceName, fmt.Sprintf("%s_%d", host, port))
}
