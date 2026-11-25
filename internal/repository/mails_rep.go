package repository

import (
	"deller/internal/database"
	"deller/internal/dto"
	"deller/internal/schema"
	"time"

	"github.com/gin-gonic/gin"
)

func GetMails(c *gin.Context, f string) (error, int, []schema.Mails) {
	var mails []schema.Mails
	ctx := c.Request.Context()

	db := database.GetDatabase()
	q := db.NewSelect().Model(&mails)
	switch f {
	case "n":
		q = q.Where("NOT is_banned")
	case "b":
		q = q.Where("is_banned")
	}
	ct, err := q.ScanAndCount(ctx)
	if err != nil {
		return err, 0, nil
	}
	return nil, ct, mails
}

func AddMail(c *gin.Context, body dto.SendContactMail, ban bool) error {
	new_mail := schema.Mails{
		Name:        body.Name,
		PhoneNumber: body.PhoneNumber,
		Email:       body.Email,
		Comment:     body.Comment,
		IsBanned:    ban,
		CreatedAt:   time.Now(),
	}
	ctx := c.Request.Context()

	db := database.GetDatabase()
	err := db.NewInsert().Model(&new_mail).Returning("*").Scan(ctx)
	return err
}
