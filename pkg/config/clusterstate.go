//go:generate go-enum -f=clusterstate.go --lower --flag
package config

// ClusterState x ENUM(
// New,
// Existing
// )
type ClusterState int32
