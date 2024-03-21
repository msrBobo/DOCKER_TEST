package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	pb "DOCKER_TEST/user-service/genproto/user-service"

	"github.com/jackc/pgx/v4"
)

type RefreshToken struct {
	RefreshToken string `json:"refresh_token"`
}

type UserCreateReq struct {
	Id           string `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	RefreshToken string `json:"refresh_token"`
}

func ConsumerHandler(message []byte) {
	var user UserCreateReq

	if err := json.Unmarshal(message, &user); err != nil {
		log.Printf("Error unmarshaling message: %v", err)
		return
	}

	var refreshToken RefreshToken
	if err := json.Unmarshal(message, &refreshToken); err != nil {
		log.Printf("Error unmarshaling refresh token: %v", err)
		return
	}
	user.RefreshToken = refreshToken.RefreshToken

	respUser, err := Create(context.Background(), &user)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(respUser)
}

func Create(ctx context.Context, user *UserCreateReq) (*pb.UserResponse, error) {
	connString := "user=postgres dbname=userdb password=1234 host=db port=5432 sslmode=disable"
	conn, err := pgx.Connect(ctx, connString)
	// Example connection string
	if err != nil {
		return nil, err
	}
	defer conn.Close(ctx)

	query := `
		INSERT INTO users (
			id,
			first_name,
			last_name,
			email,
			password,
			refresh_token
		) VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, first_name, last_name, email, password, refresh_token
	`

	var respUser pb.UserResponse
	err = conn.QueryRow(ctx, query, user.Id, user.FirstName, user.LastName, user.Email, user.Password, user.RefreshToken).
		Scan(&respUser.Id, &respUser.FirstName, &respUser.LastName, &respUser.Email, &respUser.Password, &respUser.RefreshToken)
	if err != nil {
		return nil, err
	}

	return &respUser, nil
}
