package v1

// import (
// 	"context"
// 	"encoding/json"
// 	"log"
// 	"math/rand"
// 	"net/http"
// 	"net/smtp"
// 	"regexp"
// 	"strings"
// 	"time"

// 	models "DOCKER_TEST/api-gateway/api/models"
// 	pb "DOCKER_TEST/api-gateway/genproto/user-service"
// 	l "DOCKER_TEST/api-gateway/pkg/logger"
// 	"DOCKER_TEST/api-gateway/pkg/utils"

// 	"github.com/gin-gonic/gin"
// 	"google.golang.org/protobuf/encoding/protojson"
// )

// var Nums []byte

// func sendVerificationCode(email, code string) error {
// 	auth := smtp.PlainAuth("", "msrbobo@gmail.com", "jqpdwcxjrrrvwqip", "smtp.gmail.com")
// 	to := []string{email}
// 	from := "msrbobo@gmail.com"
// 	subject := "Verification Code"
// 	body := "Your verification code is: " + code
// 	msg := []byte("To: " + email + "\r\n" +
// 		"Subject: " + subject + "\r\n" +
// 		"\r\n" +
// 		body + "\r\n")

// 	err := smtp.SendMail("smtp.gmail.com:587", auth, from, to, msg)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func generateRandomCode() string {
// 	rand.Seed(time.Now().UnixNano())
// 	const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
// 	const length = 6
// 	code := make([]byte, length)
// 	for i := range code {
// 		code[i] = charset[rand.Intn(len(charset))]
// 	}
// 	log.Println(string(code))
// 	Nums = code
// 	return string(code)
// }

// func isValidEmail(email string) bool {
// 	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
// 	re := regexp.MustCompile(emailRegex)

// 	return re.MatchString(email)
// }

// func isValidPassword(password string) bool {
// 	if len(password) < 8 {
// 		return false
// 	}
// 	uppercaseRegex := regexp.MustCompile(`[A-Z]`)
// 	if !uppercaseRegex.MatchString(password) {
// 		return false
// 	}
// 	lowercaseRegex := regexp.MustCompile(`[a-z]`)
// 	if !lowercaseRegex.MatchString(password) {
// 		return false
// 	}
// 	digitRegex := regexp.MustCompile(`[0-9]`)
// 	if !digitRegex.MatchString(password) {
// 		return false
// 	}
// 	specialCharRegex := regexp.MustCompile(`[^A-Za-z0-9]`)
// 	if !specialCharRegex.MatchString(password) {
// 		return false
// 	}
// 	commonPatterns := []string{"password", "123456", "qwerty"}
// 	for _, pattern := range commonPatterns {
// 		if strings.Contains(password, pattern) {
// 			return false
// 		}
// 	}

// 	return true
// }

// CreateUser ...
// @Summary CreateUser
// @Description Api for creating a new user
// @Tags user
// @Accept json
// @Produce json
// @Param User body models.User true "createUserModel"
// @Success 200 {string} Success!
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/users/ [post]

// CreateUser generates a verification code and sends it to the user's email

// func (h *handlerV1) CreateUser(c *gin.Context) {
// 	var body models.User

// 	// Bind JSON data to the User struct
// 	if err := c.ShouldBindJSON(&body); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		log.Println(err)
// 		return
// 	}

// 	body.Email = strings.TrimSpace(body.Email)
// 	body.Email = strings.ToLower(body.Email)

// 	if !isValidEmail(body.Email) {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
// 		return
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
// 	defer cancel()

// 	exists, err := h.serviceManager.UserService().CheckUniquess(ctx, &pb.CheckUniqReq{
// 		Field: "email",
// 		Value: body.Email,
// 	})
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check email uniqueness"})
// 		h.log.Error("failed to check email uniqueness", l.Error(err))
// 		return
// 	}
// 	if exists.IsExist {
// 		c.JSON(http.StatusConflict, gin.H{"error": "This email is already in use, please use another email address"})
// 		return
// 	}

// 	if !isValidPassword(body.Password) {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Weak password"})
// 		return
// 	}

// 	verificationCode := generateRandomCode()

// 	bodyJSON, err := json.Marshal(body)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize user data"})
// 		return
// 	}

// 	// Store the verification code in Redis with a TTL of 2 minutes
// 	err = h.Redis.SetWithTTL(context.Background(), verificationCode, string(bodyJSON), int64(2*time.Minute/time.Second))

// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store verification code"})
// 		return
// 	}

// 	err = sendVerificationCode(body.Email, verificationCode)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send verification code"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Verification code sent successfully"})
// }

// VerifyCode post password
// @Summary VerifyCode
// @Description Api for getting password
// @Tags Check
// @Accept json
// @Produce json
// @Param code path string true "code"
// @Success 200 {object} models.Check
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/getpasword/{code} [GET]

// func (h *handlerV1) VerifyCode(c *gin.Context) {
// 	verificationCode := string(Nums)
// 	var body models.Check
// 	if err := c.ShouldBindJSON(&body); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		log.Println(err)
// 		return
// 	}
// 	userProvidedCode := body.Password

// 	if userProvidedCode != verificationCode {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid verification code"})
// 		return
// 	}

// 	h.CreateUser(c)
// }

// func (h *handlerV1) CreateUser(c *gin.Context) {
// 	var body models.User

// 	// Bind JSON data to the User struct
// 	if err := c.ShouldBindJSON(&body); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		log.Println(err)
// 		return
// 	}

// 	// Trim and normalize email
// 	body.Email = strings.TrimSpace(body.Email)
// 	body.Email = strings.ToLower(body.Email)

// 	// Check if the email is in a valid format
// 	if !isValidEmail(body.Email) {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
// 		return
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
// 	defer cancel()

// 	// Check if the email is already in use
// 	exists, err := h.serviceManager.UserService().CheckUniquess(ctx, &pb.CheckUniqReq{
// 		Field: "email",
// 		Value: body.Email,
// 	})
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check email uniqueness"})
// 		h.log.Error("failed to check email uniqueness", l.Error(err))
// 		return
// 	}
// 	if exists.IsExist {
// 		c.JSON(http.StatusConflict, gin.H{"error": "This email is already in use, please use another email address"})
// 		return
// 	}

// 	// Check if the password meets the required criteria
// 	if !isValidPassword(body.Password) {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Weak password"})
// 		return
// 	}

// 	// Generate a random verification code
// 	// code := generateRandomCode()
// 	rand.Seed(time.Now().UnixNano())

// 	// Define the characters to use in the code
// 	const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
// 	const length = 6 // Length of the verification code

// 	// Create a byte slice to store the code
// 	code := make([]byte, length)

// 	// Fill the byte slice with random characters from the charset
// 	for i := range code {
// 		code[i] = charset[rand.Intn(len(charset))]
// 	}

// 	// Convert the byte slice to a string and return it
// 	Nums = code

// 	// Send the verification code to the user's email
// 	err = sendVerificationCode(body.Email, string(Nums))
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send verification code"})
// 		return
// 	}

// 	// Return the verification code to the user for verification
// 	//c.JSON(http.StatusOK, gin.H{"verification_code": Nums})
// 	// If code verification is successful, proceed with user creation
// 	response, err := h.serviceManager.UserService().Create(ctx, &pb.User{
// 		Id:       body.Id,
// 		Name:     body.Name,
// 		LastName: body.Last_name,
// 		Email:    body.Email,
// 		Password: body.Password,
// 	})
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, response)
// }

// CreateUser ...
// @Summary CreateUser
// @Description Api for creating a new user
// @Tags user
// @Accept json
// @Produce json
// @Param User body models.User true "createUserModel"
// @Success 200 {string} Success!
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/users/ [post]

// CreateUser handles the creation of a new user
// func (h *handlerV1) CreateUser(c *gin.Context) {
// 	var body models.User

// 	// Bind JSON data to the User struct
// 	if err := c.ShouldBindJSON(&body); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		log.Println(err)
// 		return
// 	}

// 	// Trim and normalize email
// 	body.Email = strings.TrimSpace(body.Email)
// 	body.Email = strings.ToLower(body.Email)

// 	// Check if the email is in a valid format
// 	if !isValidEmail(body.Email) {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
// 		return
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
// 	defer cancel()

// 	// Check if the email is already in use
// 	exists, err := h.serviceManager.UserService().CheckUniquess(ctx, &pb.CheckUniqReq{
// 		Field: "email",
// 		Value: body.Email,
// 	})
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check email uniqueness"})
// 		h.log.Error("failed to check email uniqueness", l.Error(err))
// 		return
// 	}
// 	if exists.IsExist {
// 		c.JSON(http.StatusConflict, gin.H{"error": "This email is already in use, please use another email address"})
// 		return
// 	}

// 	// Check if the password meets the required criteria
// 	if !isValidPassword(body.Password) {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Weak password"})
// 		return
// 	}

// 	// TODO: Add code verification logic here
// 	code := generateRandomCode()

// 	// Send the code to the user's email
// 	err = sendVerificationCode(body.Email, code)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send verification code"})
// 		return
// 	}

// 	// Wait for the user to provide the code
// 	userProvidedCode := c.PostForm("code")

// 	// Compare the provided code with the generated one
// 	if userProvidedCode != code {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid verification code"})
// 		return
// 	}
// 	// If code verification is successful, proceed with user creation
// 	response, err := h.serviceManager.UserService().Create(ctx, &pb.User{
// 		Id:       body.Id,
// 		Name:     body.Name,
// 		LastName: body.Last_name,
// 		Email:    body.Email,
// 		Password: body.Password,
// 	})
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, response)
// }

// func (h *handlerV1) CreateUser(c *gin.Context) {
// 	var (
// 		body models.User
// 	)

// 	// Bind JSON data to the User struct
// 	if err := c.ShouldBindJSON(&body); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		log.Println(err)
// 		return
// 	}

// 	body.Email = strings.TrimSpace(body.Email)
// 	body.Email = strings.ToLower(body.Email)

// 	// Check if the email is a valid format
// 	if !isValidEmail(body.Email) {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
// 		return
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
// 	defer cancel()

// 	exists, err := h.serviceManager.UserService().CheckUniquess(ctx, &pb.CheckUniqReq{
// 		Field: "email",
// 		Value: body.Email,
// 	})

// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": err.Error(),
// 		})
// 		h.log.Error("failed to check email uniquess", l.Error(err))
// 		return
// 	}

// 	if exists.IsExist {
// 		c.JSON(http.StatusConflict, gin.H{
// 			"error": "This email already in use, please use another email address",
// 		})
// 		h.log.Error("failed to check email uniquess", l.Error(err))
// 		return
// 	}

// 	if !isValidPassword(body.Password) {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Weak password"})
// 		return
// 	}
// 	if CheckPassword(h *handlerV1){
// 	response, err := h.serviceManager.UserService().Create(ctx, &pb.User{
// 		Id:       body.Id,
// 		Name:     body.Name,
// 		LastName: body.Last_name,
// 		Email:    body.Email,
// 		Password: body.Password,
// 	})
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, response)
// }
// }

// CreateUser ...
// @Summary CreateUser
// @Description Api for creating a new user
// @Tags user
// @Accept json
// @Produce json
// @Param User body models.User true "createUserModel"
// @Success 200 {string} Success!
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/users/ [post]

// func (h *handlerV1) CreateUser(c *gin.Context) {
// 	var (
// 		body        models.User
// 		jspbMarshal protojson.MarshalOptions
// 	)
// 	jspbMarshal.UseProtoNames = true

// 	err := c.ShouldBindJSON(&body)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": err.Error(),
// 		})
// 		h.log.Error("failed to bind json", l.Error(err))
// 		return
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
// 	defer cancel()

// 	response, err := h.serviceManager.UserService().Create(ctx, &pb.User{
// 		Id:                   body.Id,
// 		Name:                 body.Name,
// 		LastName:             body.Last_name,
// 		XXX_NoUnkeyedLiteral: struct{}{},
// 		XXX_unrecognized:     []byte{},
// 		XXX_sizecache:        0,
// 	})
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": err.Error(),
// 		})
// 		h.log.Error("failed to create user", l.Error(err))
// 		return
// 	}
// 	c.JSON(http.StatusCreated, response)
// }

// GetUser gets user by id
// @Summary GetUser
// @Description Api for getting user
// @Tags user
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success 200 {object} models.User
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/users/{id} [GET]

// func (h *handlerV1) GetUser(c *gin.Context) {
// 	var jspbMarshal protojson.MarshalOptions
// 	jspbMarshal.UseProtoNames = true

// 	id := c.Param("id")

// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
// 	defer cancel()

// 	response, err := h.serviceManager.UserService().Get(
// 		ctx, &pb.UserRequest{
// 			UserId: id,
// 		})
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": err.Error(),
// 		})
// 		h.log.Error("failed to get user", l.Error(err))
// 		return
// 	}

// 	c.JSON(http.StatusOK, response)
// }

// // ListUsers gets user by id
// // @Summary ListUsers
// // @Description Api for getting user
// // @Tags user
// // @Accept json
// // @Produce json
// // @Param id path string true "ID"
// // @Success 200 {object} models.User
// // @Failure 400 {object} models.StandardErrorModel
// // @Failure 500 {object} models.StandardErrorModel
// // @Router /v1/users/{id} [post]

// func (h *handlerV1) ListUsers(c *gin.Context) {
// 	queryParams := c.Request.URL.Query()

// 	params, errStr := utils.ParseQueryParams(queryParams)
// 	if errStr != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": errStr[0],
// 		})
// 		h.log.Error("failed to parse query params json" + errStr[0])
// 		return
// 	}

// 	var jspbMarshal protojson.MarshalOptions
// 	jspbMarshal.UseProtoNames = true

// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
// 	defer cancel()

// 	response, err := h.serviceManager.UserService().GetAll(
// 		ctx, &pb.GetAllUsersRequest{
// 			Limit: params.Limit,
// 			Page:  params.Page,
// 		})
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": err.Error(),
// 		})
// 		h.log.Error("failed to list users", l.Error(err))
// 		return
// 	}

// 	c.JSON(http.StatusOK, response)
// }

// // Update user by id
// // @Summary UpdateUser
// // @Description Api for UpdateUser user by id
// // @Tags user
// // @Accept json
// // @Produce json
// // @Param User body models.User true "createUserModel"
// // @Success 200 {object} models.User
// // @Failure 400 {object} models.StandardErrorModel
// // @Failure 500 {object} models.StandardErrorModel
// // @Router /v1/usersput [PUT]
// func (h *handlerV1) UpdateUser(c *gin.Context) {
// 	var (
// 		body        pb.User
// 		jspbMarshal protojson.MarshalOptions
// 	)
// 	jspbMarshal.UseProtoNames = true

// 	err := c.ShouldBindJSON(&body)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": err.Error(),
// 		})
// 		h.log.Error("failed to bind json", l.Error(err))
// 		return
// 	}
// 	body.Id = c.Param("id")

// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
// 	defer cancel()

// 	response, err := h.serviceManager.UserService().Update(ctx, &body)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": err.Error(),
// 		})
// 		h.log.Error("failed to update user", l.Error(err))
// 		return
// 	}

// 	c.JSON(http.StatusOK, response)
// }

// // Delete user by id
// // @Summary DeleteUser
// // @Description Api for DeleteUser user by id
// // @Tags user
// // @Accept json
// // @Produce json
// // @Param id path string true "ID"
// // @Success 200 {object} models.User
// // @Failure 400 {object} models.StandardErrorModel
// // @Failure 500 {object} models.StandardErrorModel
// // @Router /v1/users/{id} [DELETE]
// func (h *handlerV1) DeleteUser(c *gin.Context) {
// 	var jspbMarshal protojson.MarshalOptions
// 	jspbMarshal.UseProtoNames = true

// 	guid := c.Param("id")
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
// 	defer cancel()

// 	response, err := h.serviceManager.UserService().Delete(
// 		ctx, &pb.UserRequest{
// 			UserId: guid,
// 		})
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": err.Error(),
// 		})
// 		h.log.Error("failed to delete user", l.Error(err))
// 		return
// 	}

// 	c.JSON(http.StatusOK, response)
// }
