package service

import (
	"github.com/imsilence/kaleidoscope/proxy/config"
	"github.com/sirupsen/logrus"
)

type Service interface {
	Name() string
	Init(string, config.ServiceConfig)
	ListenAndServe() error
	Shutdown() error
}

type Manager struct {
	services map[string]Service
}

func NewManager() *Manager {
	return &Manager{
		services: make(map[string]Service),
	}
}

func (m *Manager) Register(s Service) {
	name := s.Name()
	if _, ok := m.services[name]; ok {
		logrus.Fatal("Service Provider Name is Exists")
	}
	m.services[name] = s
}

func (m *Manager) Get(name string) (Service, bool) {
	service, ok := m.services[name]
	return service, ok
}

var DefaultManager *Manager = NewManager()
