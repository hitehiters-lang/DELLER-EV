package schema

import "time"

type CachedLocalies struct {
	Id        int       `bun:"id,pk,autoincrement"`
	Slug      string    `bun:"slug"`
	LocFrom   string    `bun:"loc_from"`
	LocTo     string    `bun:"loc_to"`
	TransText string    `bun:"trans_text"`
	CreatedAt time.Time `bun:"created_at"`
}
