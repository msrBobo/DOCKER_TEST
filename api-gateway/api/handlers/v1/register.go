package v1

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"DOCKER_TEST/api-gateway/api/helper/email"
	"DOCKER_TEST/api-gateway/api/models"
	token "DOCKER_TEST/api-gateway/api/tokens"
	pbu "DOCKER_TEST/api-gateway/genproto/user-service"
	"DOCKER_TEST/api-gateway/pkg/etc"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/cast"
	"golang.org/x/crypto/bcrypt"
)

// Register ...
// @Summary Register
// @Description Api for registration
// @tags Register
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param User body models.RegUser true "createUserModel"
// @Success 200 {string} string "Verification code sent to email"
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /register [POST]
func (h *handlerV1) Register(c *gin.Context) {
	var (
		body models.RegUser
		code string
	)

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to bind JSON"})
		return
	}

	body.Email = strings.TrimSpace(strings.ToLower(body.Email))

	// Trim and convert email to lowercase before validation
	e := body.Email

	// Validate email address
	if !models.IsValidEmail(e) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email address"})
		return
	}

	if err := body.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exists, err := h.servicemanager.UserService().CheckUniquess(c, &pbu.CheckUniqReq{
		Field: "email",
		Value: body.Email,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check email uniqueness"})
		return
	}

	if exists.IsExist {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	byteData, err := json.Marshal(body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal data"})
		return
	}

	code = etc.GenerateCode(6)
	err = h.redis.SetWithTTL(code+body.Email, string(byteData), int(time.Minute*100))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "UserCheck: redis.SetWithTTL()"})
		return
	}
	data := map[string]interface{}{
		"Code": code,
	}
	// send otp email
	err = email.SendEmail([]string{body.Email}, "Bobomurod\n", h.cfg, "./api/helper/email/emailotp.html", data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "UserCheck: email.SendEmail()"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Verification code sent to email"})
}

// VerifyCode ...
// @Summary Verify code
// @Description Api for verifying verification code
// @Tags Register
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param email query string true "Email"
// @Param code query string true "Verification code"
// @Success 200 {string} message "Verification successful"
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /verify [GET]
func (h *handlerV1) VerifyCode(c *gin.Context) {
	Email := c.Query("email")
	userCode := c.Query("code")
	code := userCode
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()
	Email = strings.TrimSpace(strings.ToLower(Email))
	e := Email
	// Validate email address
	if !models.IsValidEmail(e) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email address"})
		return
	}

	exists, err := h.servicemanager.UserService().CheckUniquess(ctx, &pbu.CheckUniqReq{
		Field: "email",
		Value: Email,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check email uniqueness"})
		return
	}

	if exists.IsExist {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	data, err := h.redis.Get(code + Email)
	if err != nil || code != userCode {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired verification code"})
		return
	}

	if cast.ToString(data) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "otp expired"})
		return
	}
	var user models.UserCreateReq
	if err := json.Unmarshal([]byte(cast.ToString(data)), &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal user data"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Id = uuid.New().String()

	// Create access and refresh tokens JWT
	h.jwthandler = token.JWTHandler{
		Sub:       user.Id,
		Role:      "user",
		SigninKey: h.cfg.SignInKey,
		Aud:       []string{"template-front"},
		Log:       h.log,
		Timout:    h.cfg.AccessTokenTimout,
	}

	access, refresh, err := h.jwthandler.GenerateAuthJWT()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "UserRegister: jwthandler.GenerateAuthJWT()"})
		return
	}

	// ctxWithCancel, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.ContextTimeout))
	// defer cancel()

	user.Password = string(hashedPassword)
	user.RefreshToken = refresh

	dataBytes, err := json.Marshal(&user)
	if err != nil {
		log.Println(err)
	}
	err = h.writer.ProducerMessage("test-topic", dataBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "While producing to rabbimq"})
		return
	}

	respUser := models.UserResponse{
		Id:           user.Id,
		First_Name:   user.First_Name,
		Last_Name:    user.Last_Name,
		Email:        user.Email,
		Password:     user.Password,
		RefreshToken: user.RefreshToken,
		AccessToken:  access,
	}

	// res, err := h.servicemanager.UserService().Create(ctxWithCancel, &pbu.UserCreateReq{
	// 	Id:           user.Id,
	// 	FirstName:    user.First_Name,
	// 	LastName:     user.Last_Name,
	// 	Email:        user.Email,
	// 	Password:     user.Password,
	// 	RefreshToken: user.RefreshToken,
	// })
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
	// 	return
	// }
	// res.AccessToken = access

	c.JSON(http.StatusOK, respUser)
}

// Login ...
// @Summary Login
// @Description Api for Login
// @Tags Register
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param email query string true "User Email"
// @Param password query string true "User Password"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /login [POST]
func (h *handlerV1) Login(c *gin.Context) {

	email := c.Query("email")
	password := c.Query("password")

	if email == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email and password are required"})
		return
	}
	email = strings.TrimSpace(strings.ToLower(email))
	// Trim and convert email to lowercase before validation
	e := email

	// Validate email address
	if !models.IsValidEmail(e) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email address"})
		return
	}
	ctxWithCancel, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.ContextTimeout))
	defer cancel()

	// Retrieve user from the database
	user, err := h.servicemanager.UserService().GetUserByEmail(ctxWithCancel, &pbu.Email{Email: email})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please check your gmail again"})
		return
	}

	// Compare hashed password with the provided password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	h.jwthandler = token.JWTHandler{
		Sub:       user.Id,
		Role:      "user",
		SigninKey: h.cfg.SignInKey,
		Aud:       []string{"template-front"},
		Log:       h.log,
		Timout:    h.cfg.AccessTokenTimout,
	}
	access, _, err := h.jwthandler.GenerateAuthJWT()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "LoginUser: jwthandler.GenerateAuthJWT()"})
		return
	}
	user.AccessToken = access

	c.JSON(http.StatusOK, user)
}

// GetUser gets user by id
// @Summary GetUser
// @Description Api for getting user
// @Tags user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id query string true "ID"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /getuser/:id [GET]
func (h *handlerV1) GetUser(c *gin.Context) {

	id := c.Query("id")
	isUUid := func(input string) bool {
		_, err := uuid.Parse(input)
		return err == nil
	}(id)

	if !isUUid {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error: is not uuid"})
		return
	}

	ctxWithCancel, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.ContextTimeout))
	defer cancel()

	response, err := h.servicemanager.UserService().Get(ctxWithCancel,
		&pbu.UserRequest{
			UserId: id,
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user"})
		h.log.Error("failed to get user")
		return
	}

	c.JSON(http.StatusOK, response)
}

// ListUsers gets user by id
// @Summary ListUsers
// @Description Api for getting all users
// @Tags user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query string true "User PAGES"
// @Param limit query string true "User LIMIT"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /getall [GET]
func (h *handlerV1) ListUsers(c *gin.Context) {
	page := c.Query("page")
	limit := c.Query("limit")
	ctxWithCancel, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.ContextTimeout))
	defer cancel()

	response, err := h.servicemanager.UserService().GetAll(
		ctxWithCancel, &pbu.GetAllUsersRequest{
			Page:  page,
			Limit: limit,
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to list users")
		return
	}

	c.JSON(http.StatusOK, response)
}

// Update user by id
// @Summary UpdateUser
// @Description Api for UpdateUser user by id
// @Tags user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param User body models.UpdUser true "UpdateUserModel"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /update [PUT]
func (h *handlerV1) UpdateUser(c *gin.Context) {

	var (
		body *models.UpdUser
	)

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to bind JSON"})
		return
	}

	if body.Email == "" || body.Password == "" {

		c.JSON(http.StatusBadRequest, gin.H{"error": "Email and password are required"})
		return
	}
	body.Email = strings.TrimSpace(strings.ToLower(body.Email))
	// Trim and convert email to lowercase before validation
	e := body.Email

	// Validate email address
	if !models.IsValidEmail(e) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email address"})
		return
	}

	if err := body.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctxWithCancel, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.ContextTimeout))
	defer cancel()

	response, err := h.servicemanager.UserService().Update(ctxWithCancel, &pbu.UserUpd{
		Id:        body.Id,
		FirstName: body.First_Name,
		LastName:  body.Last_Name,
		Email:     body.Email,
		Password:  body.Password,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to update user")
		return
	}

	c.JSON(http.StatusOK, response)
}

// Delete user by id
// @Summary DeleteUser
// @Description Api for DeleteUser user by id
// @Tags user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id query string true "ID"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /delete/:id [DELETE]
func (h *handlerV1) DeleteUser(c *gin.Context) {

	userId := c.Query("id")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.servicemanager.UserService().Delete(
		ctx, &pbu.UserRequest{
			UserId: userId,
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to delete user")
		return
	}

	c.JSON(http.StatusOK, response)
}
