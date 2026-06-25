package pkgauth

import (
	"context"
	"deller/internal/database"
	"deller/internal/schema"
	"net/http"
)

func GetCurrentUser(r *http.Request) (schema.Users, error) {
	var db = database.GetDatabase()
	var user schema.Users
	var ctx = context.Background()

	if id, ok := r.Context().Value(UserIDKey).(int); ok {
		err := db.NewSelect().Model(&user).Where("id = ?", id).Scan(ctx)
		if err != nil {
			return schema.Users{}, err
		}
	}
	return user, nil
}
