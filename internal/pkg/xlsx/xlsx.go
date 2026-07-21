package xlsx

import (
	"context"
	"deller/internal/database"
	"deller/internal/schema"
	"fmt"
	"strconv"
	"strings"

	"github.com/uptrace/bun"
	"github.com/xuri/excelize/v2"
)

func normalizeWhitespace(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

type priceRow struct {
	RowNumber int
	Name      string
	PriceStr  string
	Price     int
}

type Processor struct {
	db *bun.DB
}

type ProcessResult struct {
	Total    int
	Matched  int
	Updated  int
	NotFound int
	Errors   int
	Details  []RowDetail
}

type RowDetail struct {
	Row      int
	Name     string
	Status   string
	OldPrice int
	NewPrice int
	Error    string
}

func NewProcessor() *Processor {
	return &Processor{
		db: database.GetDatabase(),
	}
}

func (p *Processor) ProcessFile(filePath string) (*ProcessResult, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть xlsx файл: %w", err)
	}
	defer f.Close()

	sheet := f.GetSheetName(0)
	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать лист: %w", err)
	}

	var priceRows []priceRow
	for i, row := range rows {
		if len(row) < 2 {
			continue
		}

		nameIdx := -1
		for j, cell := range row {
			trimmed := strings.TrimSpace(cell)
			if trimmed != "" {
				nameIdx = j
				break
			}
		}
		if nameIdx == -1 {
			continue
		}

		priceIdx := nameIdx + 1
		if priceIdx >= len(row) {
			continue
		}

		name := normalizeWhitespace(row[nameIdx])
		priceStr := strings.TrimSpace(row[priceIdx])
		if name == "" || priceStr == "" {
			continue
		}

		if strings.HasPrefix(name, "https://") || strings.HasPrefix(name, "http://") {
			continue
		}

		price, err := parsePrice(priceStr)
		if err != nil {
			priceRows = append(priceRows, priceRow{
				RowNumber: i + 1,
				Name:      name,
				PriceStr:  priceStr,
				Price:     0,
			})
			continue
		}

		priceRows = append(priceRows, priceRow{
			RowNumber: i + 1,
			Name:      name,
			PriceStr:  priceStr,
			Price:     price,
		})
	}

	return p.matchAndUpdate(priceRows)
}

func parsePrice(raw string) (int, error) {
	cleaned := strings.TrimSpace(raw)

	cleaned = strings.TrimSuffix(cleaned, "р.")
	cleaned = strings.TrimSuffix(cleaned, "р")
	cleaned = strings.TrimSuffix(cleaned, "руб.")
	cleaned = strings.TrimSuffix(cleaned, "руб")
	cleaned = strings.TrimSpace(cleaned)

	cleaned = strings.ReplaceAll(cleaned, ",", "")
	cleaned = strings.ReplaceAll(cleaned, " ", "")

	val, err := strconv.Atoi(cleaned)
	if err != nil {
		return 0, fmt.Errorf("не удалось распарсить цену '%s': %w", raw, err)
	}
	return val, nil
}

func (p *Processor) matchAndUpdate(priceRows []priceRow) (*ProcessResult, error) {
	ctx := context.Background()

	var products []schema.Products
	err := p.db.NewSelect().
		Model(&products).
		Column("id", "xlsx_name", "price").
		Where("xlsx_name IS NOT NULL").
		Where("xlsx_name != ''").
		Where("deleted_at IS NULL").
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить товары из БД: %w", err)
	}

	productMap := make(map[string]*schema.Products)
	for i := range products {
		key := normalizeWhitespace(products[i].XlsxName)
		if key != "" {
			productMap[key] = &products[i]
		}
	}

	type updateEntry struct {
		productId string
		oldPrice  int
		newPrice  int
	}

	result := &ProcessResult{}
	var updates []updateEntry

	for _, pr := range priceRows {
		trimmedName := normalizeWhitespace(pr.Name)

		if pr.Price == 0 && pr.PriceStr != "" {
			result.Errors++
			result.Details = append(result.Details, RowDetail{
				Row:    pr.RowNumber,
				Name:   trimmedName,
				Status: "parse_error",
				Error:  fmt.Sprintf("не удалось распарсить цену: '%s'", pr.PriceStr),
			})
			continue
		}

		if pr.Price == 0 && pr.PriceStr == "" {
			continue
		}

		prod, ok := productMap[trimmedName]
		if !ok {
			result.NotFound++
			result.Details = append(result.Details, RowDetail{
				Row:    pr.RowNumber,
				Name:   trimmedName,
				Status: "not_found",
			})
			continue
		}

		result.Matched++
		if prod.Price == pr.Price {
			result.Details = append(result.Details, RowDetail{
				Row:      pr.RowNumber,
				Name:     trimmedName,
				Status:   "unchanged",
				OldPrice: prod.Price,
				NewPrice: pr.Price,
			})
			continue
		}

		updates = append(updates, updateEntry{
			productId: prod.Id,
			oldPrice:  prod.Price,
			newPrice:  pr.Price,
		})
	}

	if len(updates) > 0 {
		tx, err := p.db.BeginTx(ctx, nil)
		if err != nil {
			return nil, fmt.Errorf("не удалось начать транзакцию: %w", err)
		}
		defer tx.Rollback()

		for _, u := range updates {
			_, err := tx.NewUpdate().
				Model((*schema.Products)(nil)).
				Set("price = ?", u.newPrice).
				Where("id = ?", u.productId).
				Exec(ctx)
			if err != nil {
				result.Errors++
				result.Details = append(result.Details, RowDetail{
					Row:      0,
					Name:     u.productId,
					Status:   "update_error",
					OldPrice: u.oldPrice,
					NewPrice: u.newPrice,
					Error:    err.Error(),
				})
				continue
			}
			result.Updated++
			result.Details = append(result.Details, RowDetail{
				Row:      0,
				Name:     u.productId,
				Status:   "updated",
				OldPrice: u.oldPrice,
				NewPrice: u.newPrice,
			})
		}

		if err := tx.Commit(); err != nil {
			return nil, fmt.Errorf("не удалось зафиксировать транзакцию: %w", err)
		}
	}

	result.Total = len(priceRows)

	return result, nil
}
