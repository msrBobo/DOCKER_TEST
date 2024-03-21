package v1

// import (
// 	"context"
// 	"net/http"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"google.golang.org/protobuf/encoding/protojson"

// 	models "DOCKER_TEST/api-gateway/api/models"
// 	pbc "DOCKER_TEST/api-gateway/genproto/comment-service"
// 	l "DOCKER_TEST/api-gateway/pkg/logger"
// )

// // CreateCommnet ...
// // @Summary CreateCommnet
// // @Description Api for creating a new Commnets
// // @Tags Commnets
// // @Accept json
// // @Produce json
// // @Param Commnet body models.Comment true "CreateCommnet"
// // @Success 200 {object} models.Comment
// // @Failure 400 {object} models.StandardErrorModel
// // @Failure 500 {object} models.StandardErrorModel
// // @Router /v1/users/ [post]

// func (h *handlerV1) CreateCommnet(c *gin.Context) {
// 	var (
// 		body        models.Comment
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

// 	response, err := h.serviceManager.CommentService().CreateComment(ctx, &pbc.CreateComment{
// 		Id:      body.Id,
// 		Content: body.Content,
// 		PostId:  body.PostId,
// 		OwnerId: body.OwnerId,
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

// func (h *handlerV1) GetCommnet(c *gin.Context) {
// 	var jspbMarshal protojson.MarshalOptions
// 	jspbMarshal.UseProtoNames = true

// 	id := c.Param("id")

// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
// 	defer cancel()

// 	response, err := h.serviceManager.PostService().GetPost(
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
