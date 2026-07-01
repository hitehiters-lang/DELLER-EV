package service

import (
	"deller/internal/dto"
	"deller/internal/schema"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ConvertCategories(p []schema.Categories) dto.CategoriesResponse {
	var cate []dto.Category
	var total int
	for _, e := range p {
		gories := dto.Category{
			Id:              e.Id,
			Name:            e.Name,
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

func ConvertProducts(p []schema.Products) dto.ProductsResponse {
	var prod []dto.Product
	var total int
	for _, e := range p {
		pwd := e.Price - int(float64(e.Price)*(float64(e.Discount)/100))
		uct := dto.Product{
			Id:                e.Id,
			CategoryId:        e.CategoryId,
			Spec:              e.Spec,
			Name:              e.Name,
			Short:             e.Short,
			Price:             e.Price,
			Discount:          e.Discount,
			PriceWithDiscount: pwd,
			MainImage:         e.MainImage,
			IsTrending:        e.IsTrending,
		}
		prod = append(prod, uct)
		total++
	}
	return dto.ProductsResponse{
		Products: prod,
		Total:    total,
	}
}

func ConvertProductDetail(p schema.Products, v []dto.ProductImage) dto.ProductDetail {
	pwd := p.Price - int(float64(p.Price)*(float64(p.Discount)/100))
	var spec dto.Specification
	if p.Specification != nil {
		s := p.Specification
		spec = dto.Specification{
			Id:                 s.Id,
			Gear:               s.Gear,
			PowerRatedMax:      s.PowerRatedMax,
			TorqueMoment:       s.TorqueMoment,
			MaxSpeed:           s.MaxSpeed,
			Battery:            s.Battery,
			PowerReserve:       s.PowerReserve,
			Weight:             s.Weight,
			LoadCap:            s.LoadCap,
			DimensionsLWH:      s.DimensionsLWH,
			SeatHeight:         s.SeatHeight,
			Tires:              s.Tires,
			Break:              s.Break,
			Clearance:          s.Clearance,
			Charging:           s.Charging,
			MoistureProtection: s.MoistureProtection,
		}
	}
	product := dto.ProductDetail{
		Id:                p.Id,
		CategoryId:        p.CategoryId,
		Spec:              p.Spec,
		Name:              p.Name,
		Short:             p.Short,
		Description:       p.Description,
		Price:             p.Price,
		Discount:          p.Discount,
		Tags:              p.Tags,
		CreatedAt:         p.CreatedAt,
		UpdatedAt:         p.UpdatedAt,
		DeletedAt:         p.DeletedAt,
		PriceWithDiscount: pwd,
		MainImage:         p.MainImage,
		Specification:     spec,
		Images:            v,
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

func GenerateUniqueFileName(name string) (string, error) {
	sharedDir := "./shared"

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
