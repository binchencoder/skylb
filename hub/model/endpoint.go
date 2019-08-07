package model

import "fmt"

// ServiceEndpoint represents a service endpoint.
// (A simplified version of pb.InstanceEndpoint)
type ServiceEndpoint struct {
	IP   string
	Port int32
}

func (se ServiceEndpoint) String() string {
	return fmt.Sprintf("%s:%d", se.IP, se.Port)
}

// ServiceEndpoints holds service endpoints.
type ServiceEndpoints map[string]ServiceEndpoint
