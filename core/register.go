package core

import (
	"github.com/WALL-EEEEEEE/Axiom/common"
	"github.com/bobg/go-generics/v3/maps"
)

type ServType struct {
	common.Enum[string]
	Value string
}

func (s ServType) String() string {

	return s.Value
}
func (s ServType) Index() string {
	return ""
}
func (s ServType) Values() []ServType {
	return servTypes
}

type ServHub map[string]map[string]interface{}

var (
	SPIDER ServType = ServType{Value: "SPIDER"}
	PIPE   ServType = ServType{Value: "PIPE"}
	TASK   ServType = ServType{Value: "TASK"}
)

var servTypes = []ServType{
	SPIDER,
	PIPE,
	TASK,
}

var Reg Registry

type Registry struct {
	hub ServHub
}

type Serv interface {
	GetName() string
	GetType() []ServType
}

func GetSupportedServTypes() []ServType {
	return servTypes
}

func NewRegistry() Registry {
	hub := make(ServHub)
	for _, r_type := range servTypes {
		hub[r_type.Value] = map[string]interface{}{}
	}
	return Registry{hub: hub}
}

func (reg *Registry) GetByType(servTypes ...ServType) map[string]interface{} {
	var servs = make(map[string]interface{})
	if len(servTypes) < 1 {
		maps.Each(reg.hub, func(K string, V map[string]interface{}) {
			maps.Copy(servs, V)
		})
		return servs
	}
	for _, servType := range servTypes {
		maps.Copy(servs, reg.hub[servType.Value])
	}
	return servs
}

func (reg *Registry) Register(serv Serv) {
	for _, typ := range serv.GetType() {
		reg.hub[typ.Value][serv.GetName()] = serv
	}
}
func (reg *Registry) UnRegister(serv Serv) {
	for _, typ := range serv.GetType() {
		interfaces := reg.GetByType(typ)
		delete(interfaces, serv.GetName())
	}
}

func init() {
	Reg = NewRegistry()
}
