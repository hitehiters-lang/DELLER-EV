package service

import (
	"deller/internal/dto"
	pkgtrans "deller/internal/pkg/translator"
	"deller/internal/schema"
)

func ConvertManuals(m []schema.Manuals, t int, code string) dto.ManualResponse {
	var man []dto.Manual
	for _, e := range m {
		uals := dto.Manual{
			Id:        e.Id,
			Slug:      e.Slug,
			Name:      pkgtrans.TranslateText("manual_name_"+e.Slug, e.Name, code),
			ImagePath: e.ImagePath,
			FilePath:  e.FilePath,
			CreatedAt: e.CreatedAt,
		}
		man = append(man, uals)
	}
	return dto.ManualResponse{
		Manuals: man,
		Total:   t,
	}
}
