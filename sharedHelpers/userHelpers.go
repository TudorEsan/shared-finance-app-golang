package sharedhelpers

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUserIdFromCtx(c *gin.Context) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(c.GetString("UserId"))
}
