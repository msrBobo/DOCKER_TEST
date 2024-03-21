package main

import (
	"DOCKER_TEST/user-service/config"
	pb "DOCKER_TEST/user-service/genproto/user-service"
	"DOCKER_TEST/user-service/pkg/db"
	"DOCKER_TEST/user-service/pkg/logger"
	"DOCKER_TEST/user-service/queue/rabbitmq/consumer"
	"DOCKER_TEST/user-service/service"
	"DOCKER_TEST/user-service/storage/postgres"
	"fmt"
	"net"

	"google.golang.org/grpc"
)

func main() {
	cfg := config.Load()

	log := logger.New(cfg.LogLevel, "DOCKER_TEST/user-service")
	defer logger.Cleanup(log)

	log.Info("main: sqlxConfig",
		logger.String("host", cfg.PostgresHost),
		logger.Int("port", cfg.PostgresPort),
		logger.String("database", cfg.PostgresDatabase))

	connDB, err, _ := db.ConnectToDB(cfg)
	if err != nil {
		log.Fatal("sqlx connection to postgres error", logger.Error(err))
	}
	rabbitMQAddr := "amqp://guest:guest@my-rabbitmq:5672/"
	rabbitConsumer, err := consumer.NewRabbitMQConsumer(rabbitMQAddr, "test-topic")
	if err != nil {
		fmt.Printf("Error while initializing consumer: %v", err)
	}

	
	defer rabbitConsumer.Close()
	// Start consuming messages
	err = rabbitConsumer.ConsumeMessage(postgres.ConsumerHandler)
	if err != nil {
		fmt.Printf("Error while consuming messages: %v", err)
	}

	userService := service.NewUserService(connDB, log)

	lis, err := net.Listen("tcp", cfg.RPCPort)
	if err != nil {
		log.Fatal("Error while listening: %v", logger.Error(err))
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, userService)
	log.Info("main: server running",
		logger.String("port", cfg.RPCPort))

	if err := s.Serve(lis); err != nil {
		log.Fatal("Error while listening: %v", logger.Error(err))

	}

}

// concumer, err := concumer.NewKafkaConsumerInit([]string{"localhost:9092"}, "test-topic", "1")
// if err != nil {
// 	log.Fatal("Error while init concumer: %v", logger.Error(err))
// }
// defer concumer.Close()

// err = concumer.ConcumerMessage(ConsumerHandler)
// if err != nil {
// 	log.Fatal("ConcumerMessage: %v", logger.Error(err))
// }
