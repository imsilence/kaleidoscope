package proxy

import (
	"github.com/sirupsen/logrus"

	"github.com/imsilence/kaleidoscope/proxy/config"
	"github.com/imsilence/kaleidoscope/proxy/service"
	_ "github.com/imsilence/kaleidoscope/proxy/service/init"
)

type Manager struct {
	config   *config.ProxyConfig
	services map[string]service.Service
}

func (m *Manager) Init(c *config.ProxyConfig) {
	m.config = c
}

func (m *Manager) startService() {
	for _, serviceConfig := range m.config.Services {
		if provider, ok := service.DefaultManager.Get(serviceConfig.Type); ok {
			provider.Init(m.config.Addr, serviceConfig)
			go func(provider service.Service) {
				if err := provider.ListenAndServe(); err != nil {
					logrus.WithFields(logrus.Fields{
						"provider": provider.Name(),
						"error":    err,
					}).Debug("error service listen")
				} else {
					logrus.WithFields(logrus.Fields{
						"provider": provider.Name(),
					}).Error("service listen...")
				}
			}(provider)
		} else {
			logrus.WithFields(logrus.Fields{
				"name":   serviceConfig.Type,
				"config": serviceConfig,
			}).Error("service provider not found")
		}
	}
}

func (m *Manager) stopService() {

}

func (m *Manager) Start() {
	m.startService()
}

func (m *Manager) Stop() {
	m.stopService()
}

func NewManager() *Manager {
	return &Manager{}
}

var DefaultManager = NewManager()
