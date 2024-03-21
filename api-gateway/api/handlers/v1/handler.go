package v1

import (
	t "DOCKER_TEST/api-gateway/api/tokens"
	"DOCKER_TEST/api-gateway/config"
	"DOCKER_TEST/api-gateway/pkg/logger"
	"DOCKER_TEST/api-gateway/queue/rabbitmq/producer"
	"DOCKER_TEST/api-gateway/services"
	"DOCKER_TEST/api-gateway/storage/redisrepo"
)

type handlerV1 struct {
	log            *logger.Logger
	cfg            config.Config
	jwthandler     t.JWTHandler
	redis          redisrepo.InMemoryStorageI
	servicemanager services.IServiceManager
	writer         producer.RabbitMQProducer
}

type HandlerV1Config struct {
	Logger         *logger.Logger
	Cfg            config.Config
	JWTHandler     t.JWTHandler
	Redis          redisrepo.InMemoryStorageI
	ServiceManager services.IServiceManager
	Writer         producer.RabbitMQProducer
}

// New ...
func New(c *HandlerV1Config) *handlerV1 {
	return &handlerV1{
		log:            c.Logger,
		cfg:            c.Cfg,
		jwthandler:     c.JWTHandler,
		redis:          c.Redis,
		servicemanager: c.ServiceManager,
		writer:         c.Writer,
	}
}
