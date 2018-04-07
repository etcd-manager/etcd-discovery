// Code generated by go-enum
// DO NOT EDIT!

package etcd

import (
	"fmt"
	"strings"
)

const (
	// ProcessTypeDirect is a ProcessType of type Direct
	ProcessTypeDirect ProcessType = iota
	// ProcessTypeStaticPod is a ProcessType of type StaticPod
	ProcessTypeStaticPod
)

const _ProcessTypeName = "DirectStaticPod"

var _ProcessTypeMap = map[ProcessType]string{
	0: _ProcessTypeName[0:6],
	1: _ProcessTypeName[6:15],
}

func (i ProcessType) String() string {
	if str, ok := _ProcessTypeMap[i]; ok {
		return str
	}
	return fmt.Sprintf("ProcessType(%d)", i)
}

var _ProcessTypeValue = map[string]ProcessType{
	_ProcessTypeName[0:6]:                   0,
	strings.ToLower(_ProcessTypeName[0:6]):  0,
	_ProcessTypeName[6:15]:                  1,
	strings.ToLower(_ProcessTypeName[6:15]): 1,
}

// ParseProcessType attempts to convert a string to a ProcessType
func ParseProcessType(name string) (ProcessType, error) {
	if x, ok := _ProcessTypeValue[name]; ok {
		return ProcessType(x), nil
	}
	return ProcessType(0), fmt.Errorf("%s is not a valid ProcessType", name)
}
