package dto

type XlsxImportResponse struct {
	Total    int                `json:"total"`
	Matched  int                `json:"matched"`
	Updated  int                `json:"updated"`
	NotFound int                `json:"not_found"`
	Errors   int                `json:"errors"`
	Details  []XlsxRowDetail    `json:"details"`
}

type XlsxRowDetail struct {
	Row      int    `json:"row"`
	Name     string `json:"name"`
	Status   string `json:"status"`
	OldPrice int    `json:"old_price"`
	NewPrice int    `json:"new_price"`
	Error    string `json:"error,omitempty"`
}
