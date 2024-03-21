package api

import (
	_ "DOCKER_TEST/api-gateway/api/docs" // swag
	v1 "DOCKER_TEST/api-gateway/api/handlers/v1"
	"DOCKER_TEST/api-gateway/api/middleware"
	"DOCKER_TEST/api-gateway/config"
	"DOCKER_TEST/api-gateway/queue/rabbitmq/producer"
	"DOCKER_TEST/api-gateway/services"
	"database/sql"

	_ "github.com/lib/pq"

	t "DOCKER_TEST/api-gateway/api/tokens"

	"DOCKER_TEST/api-gateway/pkg/logger"
	"DOCKER_TEST/api-gateway/storage/redisrepo"

	"github.com/casbin/casbin/v2"

	"github.com/casbin/casbin/util"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//r := gin.Default()
// Initialize Casbin enforcer
// casbinEnforcer, err := casbin.NewEnforcer(option.Conf.AuthConfigPath, false)
// if err != nil {
// 	option.Logger.Error("casbin enforcer error", err)
// 	panic(err)
// }
// casbinEnforcer.GetRoleManager().AddMatchingFunc("keyMatch", util.KeyMatch)
// casbinEnforcer.GetRoleManager().AddMatchingFunc("keyMatch3", util.KeyMatch3)
// connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
// 	host, port, user, password, dbname)
// // Open a connection to the database
// db, err := sql.Open("postgres", connStr)
// if err != nil {
// 	panic(err)
// }
// defer db.Close()

// policies := [][]string{
// 	{"admin", "/v1/rbac/roles", "GET"},
// 	{"admin", "/v1/rbac/role/{role}", "DELETE"},
// 	{"admin", "/v1/rbac/user-roles", "GET"},
// 	{"admin", "/v1/rbac/add-user-role", "POST"},
// 	{"admin", "/v1/rbac/delete-user-role/{id}/{role}", "DELETE"},
// 	{"admin", "/v1/rbac/permissions", "GET"},
// 	{"admin", "/v1/rbac/add-policy", "POST"},
// 	{"admin", "/v1/rbac/delete-policy", "DELETE"},
// 	{"admin", "/v1/rbac/list-role-policies/{role}", "GET"},
// }

// for _, policy := range policies {
// 	_, err := casbinEnforcer.AddPolicy(policy)
// 	if err != nil {
// 		option.Logger.Error("error during investor enforcer add policies", zap.Error(err))
// 	}
// }

// // Query access control rules from the access_control table
// rows, err := db.Query("SELECT role, endpoint, method FROM access_control")
// if err != nil {
// 	panic(err)
// }
// defer rows.Close()

// // Load access control rules into Casbin enforcer
// for rows.Next() {
// 	var role, endpoint, method string
// 	if err := rows.Scan(&role, &endpoint, &method); err != nil {
// 		panic(err)
// 	}
// 	// Add access control rule to Casbin enforcer
// 	casbinEnforcer.AddPolicy(role, endpoint, method)
// }
// if err := rows.Err(); err != nil {
// 	panic(err)
// }

// // Pass the casbin enforcer and db object to the gin middleware or routes as needed
// // For example:
// r.Use(func(c *gin.Context) {
// 	c.Set("casbinEnforcer", casbinEnforcer)
// 	c.Set("db", db)
// 	c.Next()
// })

// Option ...
type Option struct {
	Conf           config.Config
	Logger         *logger.Logger
	JWTHandler     t.JWTHandler
	Redis          redisrepo.InMemoryStorageI
	ServiceManager services.IServiceManager
	Writer         producer.RabbitMQProducer
}

// @title Welcome to Bay Store
// @version 1.7
// @description     Here QA can test and frontend or mobile developers can get information of API endpoints.
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @Host localhost:5053
func New(option Option) *gin.Engine {

	casbinEnforcer, err := casbin.NewEnforcer(option.Conf.AuthConfigPath, option.Conf.CSVFilePath)
	if err != nil {
		option.Logger.Error("casbin enforcer error", err)
		panic(err)
	}
	err = casbinEnforcer.LoadPolicy()
	if err != nil {
		option.Logger.Error("casbin error load policy", err)
		panic(err)
	}

	casbinEnforcer.GetRoleManager().AddMatchingFunc("keyMatch", util.KeyMatch)
	casbinEnforcer.GetRoleManager().AddMatchingFunc("keyMatch3", util.KeyMatch3)

	// Open a connection to the database
	db, err := sql.Open("postgres", option.Conf.ConnStr)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	pool := &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", option.Conf.RedisAdres)
		},
	}

	jwtHandler := t.JWTHandler{
		SigninKey: option.Conf.SignInKey,
		Log:       option.Logger,
	}
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	h := v1.New(&v1.HandlerV1Config{
		ServiceManager: option.ServiceManager,
		Logger:         option.Logger,
		Cfg:            option.Conf,
		JWTHandler:     jwtHandler,
		Redis:          redisrepo.NewRedisRepo(pool),
		Writer:         option.Writer,
	})

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = []string{"*"}
	corsConfig.AllowBrowserExtensions = true
	corsConfig.AllowMethods = []string{"*"}
	router.Use(cors.New(corsConfig))

	router.Use(middleware.NewAuthorizer(casbinEnforcer, jwtHandler, option.Conf))

	api := router.Group("")
	// Serve Swagger UI

	api.POST("/login", h.Login)
	api.GET("/verify", h.VerifyCode)
	api.POST("/register", h.Register)

	api.GET("/getuser/:id", h.GetUser)
	api.GET("/getall", h.ListUsers)
	api.PUT("/update", h.UpdateUser)
	api.DELETE("/delete/:id", h.DeleteUser)

	url := ginSwagger.URL("swagger/doc.json")
	api.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	return router
}

// users
// api.POST("/comments", handlerV1.CreateCommnet)

// api.POST("/users", handlerV1.CreateUser)
// api.GET("/users/:id", handlerV1.GetUser)
// api.GET("/users", handlerV1.ListUsers)
// api.PUT("/users/upd", handlerV1.UpdateUser)
// api.DELETE("/users/:id", handlerV1.DeleteUser)
//api.POST("/refresh-token/", handlerV1.RefreshToken)

//api.GET("/getpasword/:code", handlerV1.VerifyCode)
//api.POST("/posts", handlerV1.CreatePost)
