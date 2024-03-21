package models

import (
	"errors"
	"net/mail"
	"unicode"
)

// StandardErrorModel represents a standard error response.
type StandardErrorModel struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
	// Include any other fields you need for error responses
}

type ListRolesResponse struct {
}

type ResponseError struct {
}
type Empty struct {
}
type Endpoints struct {

}

type CreateUserRoleResponse struct {

}

type CreateUserRoleRequest struct {

} 


type ListUserRolesResponse struct {

} 

type ListRolePolicyResponse struct {

	} 

type UserResponse struct {
	Id           string `json:"id"`
	First_Name   string `json:"first_name"`
	Last_Name    string `json:"last_name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type RegUser struct {
	First_Name string `json:"first_name"`
	Last_Name  string `json:"last_name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}

type UpdUser struct {
	Id         string `json:"id"`
	First_Name string `json:"first_name"`
	Last_Name  string `json:"last_name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}

type UserCreateReq struct {
	Id           string `json:"id"`
	First_Name   string `json:"first_name"`
	Last_Name    string `json:"last_name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	RefreshToken string `json:"refresh_token"`
}

type StandardResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func (rm *RegUser) Validate() error {
	if rm.Email == "" {
		return errors.New("email is required")
	}
	if !isValidPassword(rm.Password) {
		return errors.New("password must be at least 8 characters long, contain uppercase and lowercase letters, and symbols")
	}
	return nil
}

func (rm *UpdUser) Validate() error {
	if rm.Email == "" {
		return errors.New("email is required")
	}
	if !isValidPassword(rm.Password) {
		return errors.New("password must be at least 8 characters long, contain uppercase and lowercase letters, and symbols")
	}
	return nil
}

// isValidPassword checks if the password meets the required criteria
func isValidPassword(password string) bool {
	hasUpperCase := false
	hasLowerCase := false
	hasSymbol := false
	if len(password) < 8 {
		return false
	}
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpperCase = true
		case unicode.IsLower(char):
			hasLowerCase = true
		case unicode.IsSymbol(char) || unicode.IsPunct(char):
			hasSymbol = true
		}
	}
	return hasUpperCase && hasLowerCase && hasSymbol
}

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
