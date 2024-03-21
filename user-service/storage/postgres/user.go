package postgres

import (
	pb "DOCKER_TEST/user-service/genproto/user-service"
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
)

type userRepo struct {
	db *sqlx.DB
}

// NewUserRepo ...
func NewUserRepo(db *sqlx.DB) *userRepo {
	return &userRepo{db: db}
}

// CheckUniquenessByEmail checks if the given email is unique.
func (r *userRepo) CheckUniquenessByEmail(req *pb.CheckUniqReq) (bool, error) {
	var exists int
	err := r.db.QueryRow("SELECT count(1) FROM users WHERE email = $1", req.Value).Scan(&exists)

	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

func (r *userRepo) GetUserByEmail(email string) (*pb.UserResponse, error) {
	query := `SELECT id, first_name, last_name, email, password,refresh_token, created_at, updated_at FROM users WHERE email = $1`
	var user pb.UserResponse
	err := r.db.QueryRow(query, email).Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.RefreshToken, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) Create(ctx context.Context, user *pb.UserCreateReq) (*pb.UserResponse, error) {

	respUser := &pb.UserResponse{}
	query := `INSERT INTO users (
		id,
		first_name,
		last_name,
		email,
		password,
		refresh_token
	)
	VALUES ($1, $2, $3, $4, $5,$6)
	RETURNING
		id,
		first_name,
		last_name,
		email,
		password,
		refresh_token,
		created_at,
		updated_at
	`
	err := r.db.QueryRow(
		query,
		user.Id,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Password,
		user.RefreshToken,
	).Scan(
		&respUser.Id,
		&respUser.FirstName,
		&respUser.LastName,
		&respUser.Email,
		&respUser.Password,
		&respUser.RefreshToken,
		&respUser.CreatedAt,
		&respUser.UpdatedAt,
	)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var (
		CreatedAt time.Time
		UpdatedAt time.Time
	)
	respUser.CreatedAt = CreatedAt.Format(time.RFC1123)
	respUser.UpdatedAt = UpdatedAt.Format(time.RFC1123)

	return respUser, nil
}

func (r *userRepo) Get(user *pb.UserRequest) (*pb.UserResponse, error) {
	query := `SELECT id, first_name, last_name,email,password,refresh_token,created_at,updated_at FROM users WHERE id = $1`
	var respUser pb.UserResponse
	err := r.db.QueryRow(query, user.UserId).Scan(&respUser.Id, &respUser.FirstName, &respUser.LastName, &respUser.Email, &respUser.Password, &respUser.RefreshToken, &respUser.CreatedAt, &respUser.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &respUser, nil
}

func (r *userRepo) GetAll(req *pb.GetAllUsersRequest) (*pb.GetAllUsersResponse, error) {

	limit, _ := strconv.Atoi(req.Limit)
	page, _ := strconv.Atoi(req.Page)
	offset := limit * (page - 1)
	query := `SELECT id, first_name, last_name, email,password,refresh_token, created_at, updated_at FROM users LIMIT $1 OFFSET $2`
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	var allUsers pb.GetAllUsersResponse
	for rows.Next() {
		var user pb.UserResponse
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.RefreshToken, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		allUsers.AllUsers = append(allUsers.AllUsers, &user)
	}
	return &allUsers, nil
}

func (r *userRepo) Update(user *pb.UserUpd) (*pb.UserResponse, error) {
	query := `
        UPDATE
            users
        SET
            first_name = $2,
            last_name = $3,
            email = $4,
            password = $5
        WHERE
            id = $1
        RETURNING
            id,
            first_name,
            last_name,
            email,
            password,
            refresh_token,
            created_at,
            updated_at
    `
	var respUser pb.UserResponse
	err := r.db.QueryRow(
		query,
		user.Id,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Password,
	).Scan(
		&respUser.Id,
		&respUser.FirstName,
		&respUser.LastName,
		&respUser.Email,
		&respUser.Password,
		&respUser.RefreshToken,
		&respUser.CreatedAt,
		&respUser.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no rows were affected: %v", err)
		}
		return nil, err
	}
	return &respUser, nil
}

func (r *userRepo) Delete(user *pb.UserRequest) (*pb.UserResponse, error) {
	log.Println(user.UserId)
	query := `DELETE FROM users WHERE id = $1 RETURNING id, first_name, last_name, email, password, refresh_token, created_at, updated_at`
	var respUser pb.UserResponse
	err := r.db.QueryRow(query, user.UserId).Scan(&respUser.Id, &respUser.FirstName, &respUser.LastName, &respUser.Email, &respUser.Password, &respUser.RefreshToken, &respUser.CreatedAt, &respUser.UpdatedAt)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &respUser, nil
}
