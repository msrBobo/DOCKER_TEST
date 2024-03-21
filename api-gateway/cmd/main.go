package main

import (
	"DOCKER_TEST/api-gateway/api"
	"DOCKER_TEST/api-gateway/config"
	"DOCKER_TEST/api-gateway/pkg/logger"
	"DOCKER_TEST/api-gateway/queue/rabbitmq/producer"
	"DOCKER_TEST/api-gateway/services"
)

func main() {

	cfg := config.Load()
	log := logger.New("api_gateway")

	serviceManager, err := services.NewServiceManager(&cfg)
	if err != nil {
		log.Error("gRPC dial error")
	}

	rabbitMQAddr := "amqp://guest:guest@rabbitmq:5672/"

	writer, err := producer.NewRabbitMQProducer(rabbitMQAddr)
	if err != nil {
		log.Fatal("error creating RabbitMQ producer", err)
	}
	
	
	defer writer.Close()

	server := api.New(api.Option{
		Conf:           cfg,
		Logger:         log,
		ServiceManager: serviceManager,
		Writer:         *writer,
	})

	if err := server.Run(cfg.HTTPPort); err != nil {
		log.Fatal("failed to run http server")
		panic(err)
	}
}
