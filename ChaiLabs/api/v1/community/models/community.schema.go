package communityModel

import (
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Community struct {
	ID              primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	SnapshotSpaceId string               `json:"snapshotSpaceId" bson:"snapshotSpaceId,omitempty" validate:"required,unique"`
	Name            string               `json:"name" bson:"name,omitempty" validate:"required,unique"`
	Creator         primitive.ObjectID   `json:"creator" bson:"creator,omitempty" validate:"required"`
	Blocks          []primitive.ObjectID `json:"blocks" bson:"blocks,omitempty"`
	Description     string               `json:"description" bson:"description,omitempty"`
	Category        string               `json:"category" bson:"category,omitempty"`
	Image           string               `json:"image" bson:"image,omitempty"`
	CreatedAt       time.Time            `json:"createdAt" bson:"createdAt,omitempty"`
}

// ValidateCommunity validates the Community struct.
//
// No parameters.
// Returns an error.
func (c *Community) ValidateCommunity() error {
	validate := validator.New()
	return validate.Struct(c)
}
