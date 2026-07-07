package service

import (
	pkgtrans "deller/internal/pkg/translator"
	"deller/internal/schema"
	"fmt"
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
