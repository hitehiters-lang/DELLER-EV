package schema

import "time"

type Mails struct {
	Id          int       `bun:"id,pk,autoincrement"`
	Name        string    `bun:"name"`
	PhoneNumber string    `bun:"phone_number"`
	Email       string    `bun:"email"`
	Comment     string    `bun:"comment"`
	IsBanned    bool      `bun:"is_banned"`
	CreatedAt   time.Time `bun:"created_at"`
}
