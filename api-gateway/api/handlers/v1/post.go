package v1

// import (
// 	"context"
// 	"net/http"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"google.golang.org/protobuf/encoding/protojson"

// 	models "DOCKER_TEST/api-gateway/api/models"
// 	pbp "DOCKER_TEST/api-gateway/genproto/post-service"
// 	l "DOCKER_TEST/api-gateway/pkg/logger"
// )

// // CreateUser ...
// // @Summary CreateUser
// // @Description Api for creating a new user
// // @Tags user
// // @Accept json
// // @Produce json
// // @Param User body models.User true "createUserModel"
// // @Success 200 {object} models.User
// // @Failure 400 {object} models.StandardErrorModel
// // @Failure 500 {object} models.StandardErrorModel
// // @Router /v1/users/ [post]

// func (h *handlerV1) CreatePost(c *gin.Context) {
// 	var (
// 		body        models.Post
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

// 	response, err := h.serviceManager.PostService().Create(ctx, &pbp.Post{
// 		Id:       body.Id,
// 		Title:    body.Title,
// 		ImageUrl: body.Image_Url,
// 		OwnerId:  body.OwnerId,
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

// // GetPost ...
// // @Summary GetPost
// // @Description Api for creating a new user
// // @Tags user
// // @Accept json
// // @Produce json
// // @Param User body models.User true "createUserModel"
// // @Success 200 {object} models.User
// // @Failure 400 {object} models.StandardErrorModel
// // @Failure 500 {object} models.StandardErrorModel
// // @Router /v1/users/ [get]

// func (h *handlerV1) GetPost(c *gin.Context) {
// 	var jspbMarshal protojson.MarshalOptions
// 	jspbMarshal.UseProtoNames = true

// 	id := c.Param("id")

// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
// 	defer cancel()

// 	response, err := h.serviceManager.PostService().Get(
// 		ctx, &pbp.GetRequest{
// 			PostId: id,
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
