package service

import (
	"github.com/nightlord189/ulp/internal/config"
	"github.com/nightlord189/ulp/internal/db"
	"sync"
)

// Service - структура со ссылками на зависимости
type Service struct {
	Config      *config.Config
	DB          *db.Manager
	portManager *portManager
}

// NewService - конструктор Service
func NewService(cfg *config.Config, db *db.Manager) *Service {
	service := Service{
		Config:      cfg,
		DB:          db,
		portManager: initPortManager(cfg.TestPortStart),
	}
	return &service
}

type portManager struct {
	ports     map[int]bool
	lock      *sync.Mutex
	startPort int
}

func initPortManager(startPort int) *portManager {
	return &portManager{
		ports:     make(map[int]bool, 0),
		lock:      &sync.Mutex{},
		startPort: startPort,
	}
}

func (p *portManager) getAndLockPort() int {
	p.lock.Lock()
	defer p.lock.Unlock()
	port := p.startPort
	for true {
		ok, _ := p.ports[port]
		if !ok {
			p.ports[port] = true
			return port
		}
		port++
	}
	return 0
}

func (p *portManager) freePort(port int) {
	p.lock.Lock()
	delete(p.ports, port)
	defer p.lock.Unlock()
}
