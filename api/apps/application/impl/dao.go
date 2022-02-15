package impl

import (
	"context"

	"github.com/infraboard/mcube/exception"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/infraboard/workflow/api/apps/application"
)

func (s *service) update(ctx context.Context, app *application.Application) error {
	_, err := s.col.UpdateOne(ctx, bson.M{"_id": app.Id}, bson.M{"$set": app})
	if err != nil {
		return exception.NewInternalServerError("update application(%s) error, %s", app.Id, err)
	}

	return nil
}
