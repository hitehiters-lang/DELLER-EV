package xlsx

import (
	"deller/internal/schema"
	"fmt"

	"github.com/xuri/excelize/v2"
)

const (
	colA = "A"
	colB = "B"
	colC = "C"

	headerCompany  = "Электроприводы   DellerQS"
	headerAddress1 = "Краснодарский край, г. Кореновск, ул. Пурыхина, 1к"
	headerAddress2 = "Алматы, с. Кыргауылды, ул. Садовая 71"
	headerSite     = "https://DellerEV.com/                                                +7(918)430-18-86"
	headerFiz      = "ФизЛицо"
)

func writeHeader(f *excelize.File, sheet string) error {
	styleBold, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Size: 12, Family: "Calibri"},
	})
	if err != nil {
		return fmt.Errorf("не удалось создать стиль жирного шрифта: %w", err)
	}

	styleNormal, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Size: 12, Family: "Calibri"},
	})
	if err != nil {
		return fmt.Errorf("не удалось создать стиль обычного шрифта: %w", err)
	}

	f.SetCellValue(sheet, "B1", headerCompany)
	f.SetCellStyle(sheet, "B1", "B1", styleBold)

	f.SetCellValue(sheet, "B2", headerAddress1)
	f.SetCellStyle(sheet, "B2", "B2", styleNormal)

	f.SetCellValue(sheet, "B3", headerAddress2)
	f.SetCellStyle(sheet, "B3", "B3", styleNormal)

	f.SetCellValue(sheet, "B4", headerSite)
	f.SetCellStyle(sheet, "B4", "B4", styleNormal)

	f.SetCellValue(sheet, "C4", headerFiz)
	f.SetCellStyle(sheet, "C4", "C4", styleBold)

	return nil
}

func ExportPriceList(products []schema.Products, catMap map[string]string, groupByCategory bool) (string, error) {
	f := excelize.NewFile()
	sheet := f.GetSheetName(0)

	f.SetColWidth(sheet, colA, colA, 5.14)
	f.SetColWidth(sheet, colB, colB, 92.14)
	f.SetColWidth(sheet, colC, colC, 17.43)

	if err := writeHeader(f, sheet); err != nil {
		return "", err
	}

	styleNormal, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Size: 12, Family: "Calibri"},
	})

	priceFmt := "#,##0 \"р.\""
	stylePrice, _ := f.NewStyle(&excelize.Style{
		Font:         &excelize.Font{Size: 12, Family: "Calibri"},
		CustomNumFmt: &priceFmt,
	})

	row := 6
	var lastCategory string

	for _, p := range products {
		catName := catMap[p.CategoryId]

		if groupByCategory && lastCategory != "" && catName != lastCategory {
			row++
		}
		lastCategory = catName

		name := p.XlsxName
		if name == "" {
			name = p.Name
		}

		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), name)
		f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), styleNormal)

		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), p.Price)
		f.SetCellStyle(sheet, fmt.Sprintf("C%d", row), fmt.Sprintf("C%d", row), stylePrice)

		row++
	}

	tmpPath := "./shared/tmp_export.xlsx"
	if err := f.SaveAs(tmpPath); err != nil {
		return "", fmt.Errorf("не удалось сохранить xlsx файл: %w", err)
	}

	return tmpPath, nil
}
