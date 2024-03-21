package middleware

import (
	"fmt"
	//"log"
	"net/http"
	"strings"

	v1 "DOCKER_TEST/api-gateway/api/handlers/v1"
	"DOCKER_TEST/api-gateway/api/models"
	token "DOCKER_TEST/api-gateway/api/tokens"
	"DOCKER_TEST/api-gateway/config"

	"github.com/casbin/casbin/v2"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	//"github.com/spf13/cast"
)

// type JWTRoleAuth struct {
// 	enforcer   *casbin.Enforcer
// 	cfg        config.Config
// 	jwtHandler token.JWTHandler
// }

// func NewAuthorizer(e *casbin.Enforcer, jwtHandler token.JWTHandler, cfg config.Config) gin.HandlerFunc {
// 	a := &JWTRoleAuth{
// 		enforcer:   e,
// 		cfg:        cfg,
// 		jwtHandler: jwtHandler,
// 	}

// 	return func(c *gin.Context) {
// 		allow, err := a.CheckPermission(c.Request)
// 		if err != nil {
// 			v, _ := err.(*jwt.ValidationError)
// 			if v.Errors == jwt.ValidationErrorExpired {
// 				a.RequireRefresh(c)
// 			} else {
// 				// fmt.Println("xato tepada")
// 				a.RequirePermission(c)
// 			}
// 		} else if !allow {
// 			// fmt.Println("xato pasda")
// 			a.RequirePermission(c)
// 		}
// 	}
// }

// func (a *JWTRoleAuth) CheckPermission(r *http.Request) (bool, error) {
// 	user, err := a.GetRole(r)
// 	if err != nil {
// 		log.Println("error get role", err)
// 		return false, err
// 	}

// 	method := r.Method
// 	path := r.URL.Path

// 	// fmt.Println("userni metode", method)
// 	// fmt.Println("userni path", path)

// 	// fmt.Println("userni ", user)

// 	allowed, err := a.enforcer.Enforce(user, path, method)
// 	if err != nil {
// 		fmt.Println("failed to check permission: ", err)
// 		return false, err
// 	}

// 	return allowed, nil
// }

// func (a *JWTRoleAuth) GetRole(r *http.Request) (string, error) {
// 	var (
// 		role   string
// 		claims jwt.MapClaims
// 		err    error
// 	)

// 	jwtToken := r.Header.Get("Authorization")

// 	if jwtToken == "" {
// 		return "unauthorized", nil
// 	}

// 	a.jwtHandler.Token = jwtToken

// 	claims, err = a.jwtHandler.ExtractClaims()
// 	if err != nil {
// 		log.Println("error chack token", err)
// 		return "", err
// 	}

// 	if cast.ToString(claims["role"]) == "admin" {
// 		role = "admin"
// 	} else if cast.ToString(claims["role"]) == "user" {
// 		role = "user"
// 	} else if cast.ToString(claims["role"]) == "unauthorized" {
// 		role = "unauthorized"
// 	} else {
// 		role = "unknown"
// 	}

// 	fmt.Println(role)

// 	return role, nil
// }

// func (a *JWTRoleAuth) RequireRefresh(c *gin.Context) {
// 	c.JSON(http.StatusUnauthorized, gin.H{
// 		"error": "required refresh",
// 	})
// 	c.AbortWithStatus(401)
// }

// func (a *JWTRoleAuth) RequirePermission(c *gin.Context) {
// 	c.JSON(http.StatusForbidden, gin.H{
// 		"Error": "You have no access this page",
// 	})
// 	c.AbortWithStatus(403)
// }

type JwtRoleAuth struct {
	enforcer   *casbin.Enforcer
	cnf        config.Config
	jwtHandler token.JWTHandler
}

func NewAuthorizer(enforce *casbin.Enforcer, jwtHandler token.JWTHandler, cfg config.Config) gin.HandlerFunc {
	a := &JwtRoleAuth{
		enforcer:   enforce,
		cnf:        cfg,
		jwtHandler: jwtHandler,
	}
	return func(c *gin.Context) {
		allow, err := a.CheckPermission(c.Request)
		fmt.Println(allow)
		if err != nil {
			v, _ := err.(*jwt.ValidationError)
			if v.Errors == jwt.ValidationErrorExpired {
				a.RequireRefresh(c)
			} else {
				a.RequirePermission(c)
			}
		} else if !allow {
			a.RequirePermission(c)
		}
	}
}

// GetRole gets role from Authorization header if there is a token then it is
// parsed and in role got from role claim. If there is no token then role is
// unauthorized
func (a *JwtRoleAuth) GetRole(r *http.Request) (string, error) {
	var (
		claims jwt.MapClaims
		err    error
	)

	jwtToken := r.Header.Get("Authorization")
	if jwtToken == "" {
		return "unauthorized", nil
	} else if strings.Contains(jwtToken, "Basic") {
		return "unauthorized", nil
	}

	a.jwtHandler.Token = jwtToken
	claims, err = a.jwtHandler.ExtractClaims()
	if err != nil {
		return "unauthorized", err
	}

	return claims["role"].(string), nil
}

// CheckPermission checks whether user is allowed to use certain endpoint
func (a *JwtRoleAuth) CheckPermission(r *http.Request) (bool, error) {
	user, err := a.GetRole(r)
	if err != nil {
		return false, err
	}
	method := r.Method
	path := r.URL.Path
	fmt.Println(user, method, path)
	allowed, err := a.enforcer.Enforce(user, path, method)
	if err != nil {
		return false, err
	}

	return allowed, nil
}

// RequirePermission aborts request with 403 status
func (a *JwtRoleAuth) RequirePermission(c *gin.Context) {
	c.AbortWithStatusJSON(403, models.StandardResponse{
		Status:  v1.PermissionDenied,
		Message: "Permission denied",
	})
}

// RequireRefresh aborts request with 401 status
func (a *JwtRoleAuth) RequireRefresh(c *gin.Context) {
	c.AbortWithStatusJSON(401, models.StandardResponse{
		Status:  v1.AccessTokenExpired,
		Message: "Access token expired",
	})
}
