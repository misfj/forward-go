package core

import (
	"context"
	apiServer2 "forward-go/core/apiServer"
	"forward-go/log"
	"sync"
)

type Server interface {
	Init() error
	Name() string
	Startup(ctx context.Context) error
	Close() error
}
type Core struct {
	servers map[string]Server
	wg      sync.WaitGroup
}

func New() *Core {
	c := Core{
		servers: make(map[string]Server),
	}
	apiServer := apiServer2.NewApiServer()
	c.servers[apiServer.Name()] = apiServer
	return &c
}
func (c *Core) Init() (err error) {
	defer func() {
		if err == nil {
			return
		}
		for _, server := range c.servers {
			server.Close()
		}
	}()
	for _, server := range c.servers {
		err = server.Init()
		if err != nil {
			log.Error(err)
			return err
		}
	}
	return nil
}
func (c *Core) Run(ctx context.Context) {
	for _, serv := range c.servers {
		c.wg.Add(1)
		go func(serv Server, wg *sync.WaitGroup) {
			err := serv.Startup(ctx)
			if err != nil {
				log.Error(err)
				return
			}
			//serv.Close()
			wg.Done()
		}(serv, &c.wg)
	}
}
func (c *Core) Close() error {
	for _, serv := range c.servers {
		serv.Close()
	}
	return nil
}

func (c *Core) Wait() {
	c.wg.Wait()
}
