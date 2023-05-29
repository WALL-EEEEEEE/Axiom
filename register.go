package core

type ServType int64

type ServHub map[ServType]map[string]interface{}

const (
	SPIDER ServType = iota
	PIPE
	TASK
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
		hub[r_type] = map[string]interface{}{}
	}
	return Registry{hub: hub}
}

func (reg *Registry) GetByType(servType ServType) map[string]interface{} {
	return reg.hub[servType]
}

func (reg *Registry) Register(serv Serv) {
	for _, typ := range serv.GetType() {
		interfaces := reg.GetByType(typ)
		interfaces[serv.GetName()] = serv
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
