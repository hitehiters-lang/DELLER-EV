package service

import (
	"context"
	pkgtrans "deller/internal/pkg/translator"
	"deller/internal/repository"
	"deller/internal/schema"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"
)

func LoadLocalisation(c *gin.Context) string {
	l := strings.ToLower(c.Query("l"))
	switch l {
	case "ru", "en", "kk", "zh":
		return l
	}
	return ""
}

func CacheCategory(cat schema.Categories) {
	pkgtrans.TranslateText("category_name_"+cat.Id, cat.Name, "zh")

	fmt.Println("[SERVER] - category \"" + cat.Id + "\" cached by creation")
}

func CacheManual(m schema.Manuals) {
	pkgtrans.TranslateText("manual_name_"+m.Slug, m.Name, "zh")

	fmt.Println("[SERVER] - manual \"" + m.Slug + "\" cached by creation")
}

func CacheProduct(p schema.Products) {
	s := p.Specification
	pkgtrans.TranslateText("spec_gear_"+s.Id.String(), s.Gear, "zh")
	pkgtrans.TranslateText("spec_pwm_"+s.Id.String(), s.PowerRatedMax, "zh")
	pkgtrans.TranslateText("spec_torq_"+s.Id.String(), s.TorqueMoment, "zh")
	pkgtrans.TranslateText("spec_mxspeed_"+s.Id.String(), s.MaxSpeed, "zh")
	pkgtrans.TranslateText("spec_acc_"+s.Id.String(), s.Battery, "zh")
	pkgtrans.TranslateText("spec_powres_"+s.Id.String(), s.PowerReserve, "zh")
	pkgtrans.TranslateText("spec_weight_"+s.Id.String(), s.Weight, "zh")
	pkgtrans.TranslateText("spec_cap_"+s.Id.String(), s.LoadCap, "zh")
	pkgtrans.TranslateText("spec_dim_"+s.Id.String(), s.DimensionsLWH, "zh")
	pkgtrans.TranslateText("spec_seat_"+s.Id.String(), s.SeatHeight, "zh")
	pkgtrans.TranslateText("spec_tire_"+s.Id.String(), s.Tires, "zh")
	pkgtrans.TranslateText("spec_br_"+s.Id.String(), s.Break, "zh")
	pkgtrans.TranslateText("spec_clear_"+s.Id.String(), s.Clearance, "zh")
	pkgtrans.TranslateText("spec_charg_"+s.Id.String(), s.Charging, "zh")
	pkgtrans.TranslateText("spec_mp_"+s.Id.String(), s.MoistureProtection, "zh")

	pkgtrans.TranslateText("product_name_"+p.Id, p.Name, "zh")
	pkgtrans.TranslateText("product_short_"+p.Id, p.Short, "zh")
	pkgtrans.TranslateText("product_desc_"+p.Id, p.Description, "zh")

	fmt.Println("[SERVER] - product \"" + p.Id + "\" cached by creation")
}
func CachePost(p schema.Posts) {
	pkgtrans.TranslateText("post_title_"+p.Slug, p.Title, "zh")
	pkgtrans.TranslateText("post_content_"+p.Slug, p.Content, "zh")

	fmt.Println("[SERVER] - post \"" + p.Slug + "\" cached by creation")
}

func CacheAllContent() {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequestWithContext(context.Background(), "GET", "/", nil)
	ctx.Request = req

	// caching categories
	err, _, categories := repository.GetCategories(ctx)
	if err != nil {
		fmt.Println("[SERVER] - " + err.Error())
		return
	}
	for _, category := range categories {
		CacheCategory(category)

		// caching products and specs
		err, _, products := repository.SelectProductsByCategory(ctx, category.Id)
		if err != nil {
			fmt.Println("[SERVER] - " + err.Error())
			return
		}
		for _, product := range products {
			err, _, full_product := repository.GetProductDetail(ctx, product.Id)
			if err != nil {
				fmt.Println("[SERVER] - " + err.Error())
				return
			}
			CacheProduct(full_product)
		}
	}
	fmt.Println("[SERVER] - PRODUCTS, SPECS AND CATEGORIES CACHED")

	// caching posts and manuals
	err, _, posts := repository.GetPosts(ctx, 200, 0)
	if err != nil {
		fmt.Println("[SERVER] - " + err.Error())
		return
	}
	for _, post := range posts {
		CachePost(post)
	}

	err, _, manuals := repository.GetManuals(ctx)
	if err != nil {
		fmt.Println("[SERVER] - " + err.Error())
		return
	}
	for _, manual := range manuals {
		CacheManual(manual)
	}
	fmt.Println("[SERVER] - BLOG AND MANUALS CACHED")
	fmt.Println("[SERVER] - CACHE ENDED SUCCESSFULLY")
}
