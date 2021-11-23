package impl

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/infraboard/mcube/exception"

	"github.com/infraboard/workflow/api/app/template"
)

func (s *impl) update(ctx context.Context, app *template.Template) error {
	_, err := s.col.UpdateOne(ctx, bson.M{"_id": app.Id}, bson.M{"$set": app})
	if err != nil {
		return exception.NewInternalServerError("update template(%s) error, %s", app.Id, err)
	}

	return nil
}
