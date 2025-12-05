package clientservice

import (
	"sync"

	"github.com/ayayaakasvin/web-socket-test/internal/models"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type ClientService struct {
	clients map[uint]*models.Client
	mutex   *sync.RWMutex
	logger  *logrus.Logger
}

func NewClientService(logger *logrus.Logger) *ClientService {
	return &ClientService{
		clients: make(map[uint]*models.Client),
		mutex:   new(sync.RWMutex),
		logger:  logger,
	}
}

func (cm *ClientService) Register(c *models.Client) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	if oldClient, exists := cm.clients[c.UserInfo.ID]; exists {
		oldClient.Conn.Close()
	}

	cm.clients[c.UserInfo.ID] = &models.Client{
		UserInfo:     c.UserInfo,
		Conn:         c.Conn,
		ConnectionID: uuid.NewString(),
	}
	cm.logger.WithField("user id", c.UserInfo.ID).Info("Client registered")
}

func (cm *ClientService) Unregister(userId uint) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	if client, exists := cm.clients[userId]; exists {
		client.Conn.Close()
		delete(cm.clients, userId)
		cm.logger.WithField("id", userId).Info("Client unregistered")
	}
}

func (cm *ClientService) Snapshot() []*models.Client {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	list := make([]*models.Client, 0, len(cm.clients))
	for _, c := range cm.clients {
		list = append(list, c)
	}
	return list
}

func (cm *ClientService) Close() int {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	counter := 0
	for id, client := range cm.clients {
		client.Conn.Close()
		delete(cm.clients, id)
		counter++
	}

	cm.logger.Infof("All %d clients closed", counter)

	return counter
}
