//go:generate go-enum -f=processtype.go --lower
package etcd

// ProcessType x ENUM(
// Direct,
// StaticPod
// )
type ProcessType int32
