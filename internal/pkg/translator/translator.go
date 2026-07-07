package pkgtrans

import (
	"context"
	"database/sql"
	"deller/internal/config"
	"deller/internal/database"
	"deller/internal/schema"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var LibreKey string
var client = &http.Client{}

func TranslatorCheck() error {
	resp, err := http.Get("http://127.0.0.1:5000/health")
	if err != nil {
		fmt.Println("[SERVER] - Переводчик не отвечает: ", err)
		return err
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("[SERVER] - Translator status: ", string(body))
	LibreKey = config.GetLibre()

	return nil
}

func TranslateText(slug, text, code string) string {

	if code == "kk" || code == "ru" || text == "" {
		return text
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	db := database.GetDatabase()
	var cached schema.CachedLocalies

	err := db.NewSelect().Model(&cached).Where("slug = ?", slug).Where("loc_to = ?", code).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			if code == "zh" {
				err = db.NewSelect().Model(&cached).Where("slug = ?", slug).Where("loc_to = ?", "en").Scan(ctx)
				if err != nil {
					translation(slug, text, "en", "ru")
				}
				if cached.TransText == "" {
					cached.TransText = text
				}
				return translation(slug, cached.TransText, code, "en")
			}
		}
		return text
	}
	return cached.TransText
}

func translation(slug, text, code, from string) string {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()
	db := database.GetDatabase()
	var cached schema.CachedLocalies

	data := url.Values{}
	data.Set("q", text)
	data.Set("source", from)
	data.Set("target", code)
	data.Set("format", "html")
	data.Set("alternatives", "0")
	data.Set("api_key", LibreKey)

	req, err := http.NewRequestWithContext(ctx, "POST", "http://127.0.0.1:5000/translate", strings.NewReader(data.Encode()))
	if err != nil {
		return err.Error()
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return err.Error()
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err.Error()
	}

	var result struct {
		TranslatedText string `json:"translatedText"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return err.Error()
	}
	cached = schema.CachedLocalies{
		Slug:      slug,
		LocFrom:   "ru",
		LocTo:     code,
		TransText: result.TranslatedText,
		CreatedAt: time.Now(),
	}
	_, err = db.NewInsert().Model(&cached).Exec(ctx)
	if err != nil {
		fmt.Println("[SERVER] - ", err)
	}
	return result.TranslatedText
}
