package service

import (
	"deller/internal/dto"
	pkgtrans "deller/internal/pkg/translator"
	"deller/internal/schema"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

func ConvertCategories(p []schema.Categories, code string) dto.CategoriesResponse {
	var cate []dto.Category
	var total int
	for _, e := range p {
		gories := dto.Category{
			Id:              e.Id,
			Name:            pkgtrans.TranslateText("category_name_"+e.Id, e.Name, code),
			ImagePath:       e.ImagePath,
			TotalInCategory: e.ProductsCount,
		}
		cate = append(cate, gories)
		total++
	}
	return dto.CategoriesResponse{
		Categories: cate,
		Total:      total,
	}
}

func ConvertCategory(p schema.Categories) dto.Category {
	cat := dto.Category{
		Id:              p.Id,
		Name:            p.Name,
		ImagePath:       p.ImagePath,
		TotalInCategory: 0,
	}
	return cat
}

func ConvertProducts(p []schema.Products, code string) dto.ProductsResponse {
	var prod []dto.Product
	var total int
	for _, e := range p {
		pwd := e.Price - int(float64(e.Price)*(float64(e.Discount)/100))
		uct := dto.Product{
			Id:                e.Id,
			CategoryId:        e.CategoryId,
			Spec:              e.Spec,
			Name:              pkgtrans.TranslateText("product_name_"+e.Id, e.Name, code),
			Short:             pkgtrans.TranslateText("product_short_"+e.Id, e.Short, code),
			Price:             e.Price,
			Discount:          e.Discount,
			PriceWithDiscount: pwd,
			MainImage:         e.MainImage,
			IsTrending:        e.IsTrending,
			XlsxName:          e.XlsxName,
			OnRequest:         e.OnRequest,
			OnStock:           e.OnStock,
			OnRequestPrice:    e.OnRequestPrice,
			OnStockPrice:      e.OnStockPrice,
		}
		prod = append(prod, uct)
		total++
	}
	return dto.ProductsResponse{
		Products: prod,
		Total:    total,
	}
}

func ConvertProductDetail(p schema.Products, v []dto.ProductImage, code string) dto.ProductDetail {
	pwd := p.Price - int(float64(p.Price)*(float64(p.Discount)/100))
	var spec dto.Specification
	if p.Specification != nil {
		s := p.Specification
		spec = dto.Specification{
			Id:                 s.Id,
			Gear:               pkgtrans.TranslateText("spec_gear_"+s.Id.String(), s.Gear, code),
			PowerRatedMax:      pkgtrans.TranslateText("spec_pwm_"+s.Id.String(), s.PowerRatedMax, code),
			TorqueMoment:       pkgtrans.TranslateText("spec_torq_"+s.Id.String(), s.TorqueMoment, code),
			MaxSpeed:           pkgtrans.TranslateText("spec_mxspeed_"+s.Id.String(), s.MaxSpeed, code),
			Battery:            pkgtrans.TranslateText("spec_acc_"+s.Id.String(), s.Battery, code),
			PowerReserve:       pkgtrans.TranslateText("spec_powres_"+s.Id.String(), s.PowerReserve, code),
			Weight:             pkgtrans.TranslateText("spec_weight_"+s.Id.String(), s.Weight, code),
			LoadCap:            pkgtrans.TranslateText("spec_cap_"+s.Id.String(), s.LoadCap, code),
			DimensionsLWH:      pkgtrans.TranslateText("spec_dim_"+s.Id.String(), s.DimensionsLWH, code),
			SeatHeight:         pkgtrans.TranslateText("spec_seat_"+s.Id.String(), s.SeatHeight, code),
			Tires:              pkgtrans.TranslateText("spec_tire_"+s.Id.String(), s.Tires, code),
			Break:              pkgtrans.TranslateText("spec_br_"+s.Id.String(), s.Break, code),
			Clearance:          pkgtrans.TranslateText("spec_clear_"+s.Id.String(), s.Clearance, code),
			Charging:           pkgtrans.TranslateText("spec_charg_"+s.Id.String(), s.Charging, code),
			MoistureProtection: pkgtrans.TranslateText("spec_mp_"+s.Id.String(), s.MoistureProtection, code),
		}
	}
	product := dto.ProductDetail{
		Id:                p.Id,
		CategoryId:        p.CategoryId,
		Spec:              p.Spec,
		Name:              pkgtrans.TranslateText("product_name_"+p.Id, p.Name, code),
		Short:             pkgtrans.TranslateText("product_short_"+p.Id, p.Short, code),
		Description:       pkgtrans.TranslateText("product_desc_"+p.Id, p.Description, code),
		Price:             p.Price,
		Discount:          p.Discount,
		Tags:              p.Tags,
		CreatedAt:         p.CreatedAt,
		UpdatedAt:         p.UpdatedAt,
		DeletedAt:         p.DeletedAt,
		PriceWithDiscount: pwd,
		MainImage:         p.MainImage,
		IsTrending:        p.IsTrending,
		XlsxName:          p.XlsxName,
		Specification:     spec,
		Images:            v,
		OnRequest:         p.OnRequest,
		OnStock:           p.OnStock,
		OnRequestPrice:    p.OnRequestPrice,
		OnStockPrice:      p.OnStockPrice,
	}
	return product
}

func GetImages(i []schema.ProductImages) []dto.ProductImage {
	var res []dto.ProductImage

	for _, e := range i {
		file_name := filepath.Base(e.ImagePath)
		url_path := "/images/" + file_name

		var img dto.ProductImage
		img = dto.ProductImage{
			Id:        e.Id,
			ProductId: e.ProductId,
			ImagePath: url_path,
		}
		res = append(res, img)
	}

	return res
}

func GenerateUniqueFileName(name, sharedDir string) (string, error) {

	cleanName := filepath.Clean(name)
	if strings.Contains(cleanName, "..") {
		return "", errors.New("invalid file name: path traversal detected")
	}

	absDir, err := filepath.Abs(sharedDir)
	if err != nil {
		return "", fmt.Errorf("failed to get abs dir: %w", err)
	}

	ext := filepath.Ext(cleanName)
	base := strings.TrimSuffix(cleanName, ext)

	checkPath := func(candidate string) (bool, error) {
		fullPath := filepath.Join(sharedDir, candidate)
		absPath, err := filepath.Abs(fullPath)
		if err != nil {
			return false, fmt.Errorf("failed to get abs path: %w", err)
		}
		if !strings.HasPrefix(absPath, absDir+string(os.PathSeparator)) && absPath != absDir {
			return false, errors.New("invalid file path: escaped shared directory")
		}

		info, err := os.Stat(fullPath)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return false, nil
			}
			return false, fmt.Errorf("failed to stat file: %w", err)
		}
		if info.IsDir() {
			return false, nil
		}
		return true, nil
	}

	exists, err := checkPath(cleanName)
	if err != nil {
		return "", err
	}
	if !exists {
		return cleanName, nil
	}

	const maxAttempts = 10000
	for i := 1; i <= maxAttempts; i++ {
		candidate := fmt.Sprintf("%s_%d%s", base, i, ext)
		exists, err := checkPath(candidate)
		if err != nil {
			return "", err
		}
		if !exists {
			return candidate, nil
		}
	}

	return "", errors.New("could not generate unique file name: too many attempts")
}

func ValidateRevBody(body dto.NewReview) error {
	var ret_err string
	var err error

	if body.Rate < 1 || body.Rate > 5 {
		ret_err = "Оценка невалидна"
	}

	if utf8.RuneCountInString(body.Content) > 990 {
		ret_err = "Слишком длинное содержание отзыва"
	}

	if ret_err != "" {
		err = errors.New(ret_err)
	}

	return err
}

func ConvertReviews(r []schema.Reviews) dto.ReviewResponse {
	var rev []dto.Review
	var total_rate int
	var total int = 0
	for _, e := range r {
		iews := dto.Review{
			Id:            e.Id,
			AuthorName:    e.Author.Name,
			AuthorIsAdmin: e.Author.IsAdmin,
			ProductId:     e.ProductId,
			ParentId:      e.ParentId,
			Rate:          e.Rate,
			Content:       e.Content,
			CreatedAt:     time.Now(),
		}
		rev = append(rev, iews)
		if iews.ParentId == uuid.Nil {
			total++
			total_rate += e.Rate
		}
	}
	var str_rate string
	var rate float64
	if total > 0 {
		str_rate = fmt.Sprintf("%.1f", float64(total_rate)/float64(total))
		rate, _ = strconv.ParseFloat(str_rate, 64)
	} else {
		rate = 0.0
	}

	return dto.ReviewResponse{
		Reviews: rev,
		Total:   total,
		Rate:    rate,
	}
}
