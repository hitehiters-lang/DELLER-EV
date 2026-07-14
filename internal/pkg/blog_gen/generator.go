package blog_gen

import (
	"context"
	"deller/internal/config"
	"deller/internal/database"
	"deller/internal/schema"
	"deller/internal/service"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type themesResponse struct {
	Themes []string `json:"themes"`
}

type postResponse struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type imagePromptResponse struct {
	ImagePrompt string `json:"image_prompt"`
}

type commentResponse struct {
	Content    string `json:"content"`
	AuthorName string `json:"author_name"`
}

func generatePost() error {
	ctx := context.Background()
	db := database.GetDatabase()

	var recentPosts []schema.Posts
	err := db.NewSelect().Model(&recentPosts).
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		Limit(10).
		Scan(ctx)
	if err != nil {
		return fmt.Errorf("get recent posts: %w", err)
	}

	recentTitles := make([]string, 0, len(recentPosts))
	for _, p := range recentPosts {
		recentTitles = append(recentTitles, p.Title)
	}

	theme, err := generateTheme(recentTitles)
	if err != nil {
		return fmt.Errorf("generate theme: %w", err)
	}
	fmt.Println("[BLOG_GEN] selected theme:", theme)

	post, err := generatePostContent(theme)
	if err != nil {
		return fmt.Errorf("generate post: %w", err)
	}
	fmt.Println("[BLOG_GEN] generated post:", post.Title)

	imagePrompt, err := generateImagePrompt(post.Content)
	if err != nil {
		return fmt.Errorf("generate image prompt: %w", err)
	}
	fmt.Println("[BLOG_GEN] image prompt:", imagePrompt)

	imgBytes, err := fetchImageBytes(imagePrompt)
	if err != nil {
		return fmt.Errorf("fetch image: %w", err)
	}

	imagePath, err := saveImage(imgBytes)
	if err != nil {
		return fmt.Errorf("save image: %w", err)
	}

	slug := service.CreateSlugFromTitle(post.Title)

	aiAuthorID := config.GetAIAuthorID()
	newPost := schema.Posts{
		Slug:       slug,
		Title:      post.Title,
		Content:    post.Content,
		AuthorID:   aiAuthorID,
		IsApproved: true,
		CreatedAt:  time.Now(),
		ImagePath:  imagePath,
	}

	_, err = db.NewInsert().Model(&newPost).Exec(ctx)
	if err != nil {
		return fmt.Errorf("insert post: %w", err)
	}
	fmt.Println("[BLOG_GEN] post saved with slug:", slug)

	go service.CachePost(newPost)

	err = generateComments(ctx, db, slug, post.Title, post.Content)
	if err != nil {
		return fmt.Errorf("generate comments: %w", err)
	}

	return nil
}

func generateTheme(recentTitles []string) (string, error) {
	titlesStr := strings.Join(recentTitles, "\n- ")
	if titlesStr == "" {
		titlesStr = "(пока нет постов)"
	}
	userPrompt := strings.Replace(promptPostThemes, "{{RECENT_POSTS}}", titlesStr, 1)

	raw, err := callLLM("Ты - генератор тем для постов об электромотоциклах, скутерах и байках.", userPrompt)
	if err != nil {
		return "", err
	}

	jsonStr, err := extractJSON(raw)
	if err != nil {
		return "", err
	}

	var resp themesResponse
	if err := json.Unmarshal([]byte(jsonStr), &resp); err != nil {
		return "", fmt.Errorf("unmarshal themes: %w", err)
	}

	if len(resp.Themes) == 0 {
		return "", fmt.Errorf("no themes generated")
	}

	return resp.Themes[rand.Intn(len(resp.Themes))], nil
}

func generatePostContent(theme string) (postResponse, error) {
	userPrompt := strings.Replace(promptPostGen, "{{THEME}}", theme, 1)

	raw, err := callLLM("Ты - автор постов для блога об электромотоциклах, скутерах и байках.", userPrompt)
	if err != nil {
		return postResponse{}, err
	}

	jsonStr, err := extractJSON(raw)
	if err != nil {
		return postResponse{}, err
	}

	var resp postResponse
	if err := json.Unmarshal([]byte(jsonStr), &resp); err != nil {
		return postResponse{}, fmt.Errorf("unmarshal post: %w", err)
	}

	return resp, nil
}

func generateImagePrompt(content string) (string, error) {
	userPrompt := strings.Replace(promptImageGen, "{{POST_CONTENT}}", content, 1)

	raw, err := callLLM("Ты - генератор промптов для изображений электромотоциклов и скутеров.", userPrompt)
	if err != nil {
		return "", err
	}

	jsonStr, err := extractJSON(raw)
	if err != nil {
		return "", err
	}

	var resp imagePromptResponse
	if err := json.Unmarshal([]byte(jsonStr), &resp); err != nil {
		return "", fmt.Errorf("unmarshal image prompt: %w", err)
	}

	return resp.ImagePrompt, nil
}

func generateComments(ctx context.Context, db *bun.DB, slug, title, content string) error {
	var existingComments []schema.Comments
	err := db.NewSelect().Model(&existingComments).
		Where("post_slug = ?", slug).
		Where("deleted_at IS NULL").
		Scan(ctx)
	if err != nil {
		return fmt.Errorf("get existing comments: %w", err)
	}

	promptTemplate := promptComment

	var commentIDs []string
	var generatedComments []string

	for i := 0; i < 2; i++ {
		allExisting := make([]string, 0, len(existingComments)+len(generatedComments))
		for _, c := range existingComments {
			allExisting = append(allExisting, c.Author.Name+": "+c.Content)
		}
		allExisting = append(allExisting, generatedComments...)

		existingForPrompt := ""
		if len(allExisting) > 0 {
			existingForPrompt = "Уже существующие комментарии (не повторяй их имена и содержание):\n" + strings.Join(allExisting, "\n")
		}

		userPrompt := strings.Replace(promptTemplate, "{{POST_TITLE}}", title, 1)
		userPrompt = strings.Replace(userPrompt, "{{POST_CONTENT}}", content, 1)
		userPrompt = strings.Replace(userPrompt, "{{EXISTING_COMMENTS}}", existingForPrompt, 1)

		raw, err := callLLM("Ты - комментатор блога об электромотоциклах, скутерах и байках.", userPrompt)
		if err != nil {
			return fmt.Errorf("generate comment %d: %w", i+1, err)
		}

		jsonStr, err := extractJSON(raw)
		if err != nil {
			return fmt.Errorf("extract comment %d json: %w", i+1, err)
		}

		var resp commentResponse
		if err := json.Unmarshal([]byte(jsonStr), &resp); err != nil {
			return fmt.Errorf("unmarshal comment %d: %w", i+1, err)
		}

		generatedComments = append(generatedComments, resp.AuthorName+": "+resp.Content)

		authorID, err := findOrCreateUser(resp.AuthorName)
		if err != nil {
			return fmt.Errorf("find/create user for comment %d: %w", i+1, err)
		}

		commentID, err := generateUUID()
		if err != nil {
			return err
		}

		newComment := schema.Comments{
			Id:        commentID,
			AuthorId:  authorID,
			PostSlug:  slug,
			Content:   resp.Content,
			CreatedAt: time.Now(),
		}

		_, err = db.NewInsert().Model(&newComment).Exec(ctx)
		if err != nil {
			return fmt.Errorf("insert comment %d: %w", i+1, err)
		}

		commentIDs = append(commentIDs, commentID.String())
		fmt.Printf("[BLOG_GEN] comment %d by %s saved\n", i+1, resp.AuthorName)
	}

	replyPromptTemplate := promptComResponse

	randomIdx := rand.Intn(len(commentIDs))
	var randomComment schema.Comments
	err = db.NewSelect().Model(&randomComment).
		Where("id = ?", commentIDs[randomIdx]).
		Scan(ctx)
	if err != nil {
		return fmt.Errorf("get random comment: %w", err)
	}

	replyUserPrompt := strings.Replace(replyPromptTemplate, "{{COMMENT_TO_REPLY}}", randomComment.Content, 1)

	raw, err := callLLM("Ты - комментатор блога об электромотоциклах, скутерах и байках.", replyUserPrompt)
	if err != nil {
		return fmt.Errorf("generate reply: %w", err)
	}

	jsonStr, err := extractJSON(raw)
	if err != nil {
		return fmt.Errorf("extract reply json: %w", err)
	}

	var replyResp commentResponse
	if err := json.Unmarshal([]byte(jsonStr), &replyResp); err != nil {
		return fmt.Errorf("unmarshal reply: %w", err)
	}

	replyAuthorID, err := findOrCreateUser(replyResp.AuthorName)
	if err != nil {
		return fmt.Errorf("find/create user for reply: %w", err)
	}

	replyID, err := generateUUID()
	if err != nil {
		return err
	}

	replyComment := schema.Comments{
		Id:        replyID,
		AuthorId:  replyAuthorID,
		PostSlug:  slug,
		Content:   replyResp.Content,
		ParentId:  randomComment.Id,
		CreatedAt: time.Now(),
	}

	_, err = db.NewInsert().Model(&replyComment).Exec(ctx)
	if err != nil {
		return fmt.Errorf("insert reply: %w", err)
	}

	fmt.Printf("[BLOG_GEN] reply by %s saved\n", replyResp.AuthorName)
	return nil
}

func generateUUID() (uuid.UUID, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return uuid.Nil, fmt.Errorf("generate uuid: %w", err)
	}
	return id, nil
}
