package service

import (
	"deller/internal/dto"
	"deller/internal/schema"
	"path/filepath"
)

func ConvertCategories(p []schema.Categories) dto.CategoriesResponse {
	var cate []dto.Category
	var total int
	for _, e := range p {
		gories := dto.Category{
			Id:              e.Id,
			Name:            e.Name,
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
	pwd := p.Price - (p.Price * (p.Discount / 100))
	var spec dto.Specification
	if p.Specification != nil {
		spec = dto.Specification{
			Id:                 p.Specification.Id,
			Gear:               p.Specification.Gear,
			PowerRatedMax:      p.Specification.PowerRatedMax,
			TorqueMoment:       p.Specification.TorqueMoment,
			MaxSpeed:           p.Specification.MaxSpeed,
			Battery:            p.Specification.Battery,
			PowerReserve:       p.Specification.PowerReserve,
			Weight:             p.Specification.Weight,
			LoadCap:            p.Specification.LoadCap,
			DimensionsLWH:      p.Specification.DimensionsLWH,
			SeatHeight:         p.Specification.SeatHeight,
			Tires:              p.Specification.Tires,
			Break:              p.Specification.Break,
			Clearance:          p.Specification.Clearance,
			Charging:           p.Specification.Charging,
			MoistureProtection: p.Specification.MoistureProtection,
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
		Specification:     spec,
		Images:            v,
	}
	return product
}

func GetImages(i []schema.ProductImages) []dto.ProductImage {
	var res []dto.ProductImage

	for _, e := range i {
		file_name := filepath.Base(e.ImagePath)
		url_path := "/api/v1/images/" + file_name

		var img dto.ProductImage
		img = dto.ProductImage{
			Id:        e.Id,
			ProductId: e.ProductId,
			IsMain:    e.IsMain,
			ImagePath: url_path,
		}
		res = append(res, img)
	}

	return res
}

func GetMainImage(i schema.ProductImages) dto.ProductImage {

	file_name := filepath.Base(i.ImagePath)
	url_path := "/api/v1/images/" + file_name

	var img dto.ProductImage
	img = dto.ProductImage{
		Id:        i.Id,
		ProductId: i.ProductId,
		IsMain:    i.IsMain,
		ImagePath: url_path,
	}

	return img
}
