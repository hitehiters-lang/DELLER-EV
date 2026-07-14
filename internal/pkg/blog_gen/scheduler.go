package blog_gen

import (
	"context"
	"deller/internal/config"
	"deller/internal/database"
	"deller/internal/schema"
	"fmt"
	"time"
)

func isGenerationDay() bool {
	today := time.Now().Weekday()
	return today == time.Monday || today == time.Wednesday || today == time.Friday
}

func isWithinGenerationWindow() bool {
	now := time.Now().Hour()
	return now >= 12 && now <= 21
}

func alreadyGeneratedToday() bool {
	ctx := context.Background()
	db := database.GetDatabase()

	aiAuthorID := config.GetAIAuthorID()
	today := time.Now().Format("2006-01-02")

	count, err := db.NewSelect().Model((*schema.Posts)(nil)).
		Where("author_id = ?", aiAuthorID).
		Where("created_at::date = ?", today).
		Where("deleted_at IS NULL").
		Count(ctx)
	if err != nil {
		fmt.Println("[BLOG_GEN] error checking today's posts:", err)
		return false
	}
	return count > 0
}

func canGenerate() bool {
	return isGenerationDay() && isWithinGenerationWindow() && !alreadyGeneratedToday()
}
