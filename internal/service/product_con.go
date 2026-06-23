package service

import (
	"deller/internal/dto"
	"deller/internal/schema"
)

func ConvertProducts(p []schema.Products) dto.ProductResponse {
	var prod []dto.Product
	var total int
	for _, e := range p {
		pwd := e.Price - int(float64(e.Price)*(float64(e.Discount)/100))
		uct := dto.Product{
			Id:                e.Id,
			CategoryId:        e.CategoryId,
			Stack:             e.Stack,
			Name:              e.Name,
			Short:             e.Short,
			Price:             e.Price,
			Discount:          e.Discount,
			PriceWithDiscount: pwd,
		}
		prod = append(prod, uct)
		total++
	}
	return dto.ProductResponse{
		Products: prod,
		Total:    total,
	}
}

/*
func ConvertProductMore(p []schema.Products) dto.ProductResponse {
	var prod []dto.Product
	for _, e := range p {
		pwd := e.Price - (e.Price * (e.Discount / 100))
		uct := dto.Product{
			Id:                e.Id,
			CategoryId:        e.CategoryId,
			Stack:             e.Stack,
			Name:              e.Name,
			Short:             e.Short,
			Description:       e.Description,
			Price:             e.Price,
			Discount:          e.Discount,
			Tags:              e.Tags,
			CreatedAt:         e.CreatedAt,
			UpdatedAt:         e.UpdatedAt,
			DeletedAt:         e.DeletedAt,
			PriceWithDiscount: pwd,
		}
		prod = append(prod, uct)
	}

}*/
