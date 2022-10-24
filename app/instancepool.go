package app

import (
	"atfm/app/config"
	"atfm/app/models"
	"net/rpc"
)

type InstancePool struct {
	instances []*Instance

	rpcClient *rpc.Client

	config *config.Config
}

func NewInstancePool(client *rpc.Client, config *config.Config) *InstancePool {
	return &InstancePool{
		instances: []*Instance{},
		rpcClient: client,
		config:    config,
	}
}

func (p InstancePool) AddInstanceDefault() (*Instance, int, error) {
	return p.AddInstance(p.config.StartDir, p.config.StartBasepath, models.LOCALFM)
}

func (p *InstancePool) AddInstance(openPath, basePath string, mod models.FsMod) (*Instance, int, error) {
	id := len(p.instances)
	ins := NewInstance(mod, openPath, basePath, id, p.config.Display.ShowOpenParent)
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

func (p *InstancePool) RefreshInstances(path string, content []models.FileInfo, selfDelete bool) {
	for _, v := range p.instances {
		if v.DirPath != path {
			continue
		}
		if selfDelete {
			err := v.OpenDir(p.config.StartDir, p.config.StartBasepath, models.LOCALFM)
			if err != nil {
				return
			}
			continue
		}
		citem := v.CurrentItem
		v.Content = content
		v.ShownContent = v.GetShownContent(content)
		if citem > len(v.ShownContent) {
			citem = len(v.ShownContent)
		}
		v.CurrentItem = citem
	}
}
