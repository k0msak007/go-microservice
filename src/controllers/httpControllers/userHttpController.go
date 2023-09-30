package httpcontrollers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/k0msak007/go-microservice/src/config"
	"github.com/k0msak007/go-microservice/src/models"
	pbItem "github.com/k0msak007/go-microservice/src/proto/item"
	"github.com/k0msak007/go-microservice/src/repositories"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserHttpController struct {
	Cfg            *config.Config
	UserRepository *repositories.UserRepository
}

func (h *UserHttpController) FindOneUser(c echo.Context) error {
	ctx := c.Request().Context()

	userId := strings.Trim(c.Param("user_id"), " ")

	user, err := h.UserRepository.FindOneUser(ctx, userId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &models.Error{
			Message: err.Error(),
		})
	}

	connItem, err := grpc.Dial(h.Cfg.Grpc.ItemAppUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Error{
			Message: err.Error(),
		})
	}
	defer connItem.Close()

	clientItem := pbItem.NewItemServiceClient(connItem)

	streamItems, err := clientItem.FindItems(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Error{
			Message: err.Error(),
		})
	}

	for _, item := range user.Item {
		if err := streamItems.Send(
			&pbItem.ItemReq{
				Id: item.ObjectId.Hex(),
			},
		); err != nil {
			log.Printf("%v.Send(%v) = %v", streamItems, &pbItem.ItemReq{Id: item.ObjectId.Hex()}, err)
			return c.JSON(http.StatusInternalServerError, &models.Error{
				Message: err.Error(),
			})
		}
		fmt.Printf("order: %v has been streamed\n", &pbItem.ItemReq{Id: item.ObjectId.Hex()})
	}
	results, err := streamItems.CloseAndRecv()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Error{
			Message: err.Error(),
		})
	}
	if len(results.Data) == 0 {
		return c.JSON(http.StatusInternalServerError, &models.Error{
			Message: err.Error(),
		})
	}

	user.Item = make([]models.Item, 0)
	for _, item := range results.Data {
		user.Item = append(user.Item, models.Item{
			ObjectId: func() primitive.ObjectID {
				result, _ := primitive.ObjectIDFromHex(item.Id)
				return result
			}(),
			Title:       item.Title,
			Description: item.Description,
			Damage:      item.Damage,
		})
	}

	return c.JSON(http.StatusOK, user)
}
