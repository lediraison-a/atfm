package app

import (
	"atfm/app/models"
	"net/rpc"
)

type InstancePool struct {
	instances []*Instance
	rpcClient *rpc.Client
}

func NewInstancePool(client *rpc.Client) *InstancePool {
	return &InstancePool{
		instances: []*Instance{},
		rpcClient: client,
	}
}

func (p *InstancePool) AddInstance(openPath, basePath string, mod models.FsMod) (*Instance, int, error) {
	id := len(p.instances)
	ins := NewInstance(mod, openPath, basePath, id)
	ins.rpcClient = p.rpcClient
	err := ins.OpenDir(openPath, basePath, mod)
	if err != nil {
		return nil, -1, err
	}
	p.instances = append(p.instances, ins)
	return ins, id, nil
}

func (p *InstancePool) GetInstance(index int) *Instance {
	return p.instances[index]
}
