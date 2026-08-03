package dto

type UploadMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Link    string `json:"link"`
}
type Message struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
type AddCategoryBody struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	ImageName string `json:"image_name"`
}
type AddProductBody struct {
	Id            string         `json:"id"`
	CategoryId    string         `json:"category_id"`
	Name          string         `json:"name"`
	Short         string         `json:"short"`
	Description   string         `json:"description"`
	Price         int            `json:"price"`
	Discount      int            `json:"discount"`
	Tags          []string       `json:"tags"`
	MainImage     string         `json:"main_image" example:"/images/example.png"`
	XlsxName      string         `json:"xlsx_name"`
	Specification AdditionalSpec `json:"specification"`
	IsTrending    bool           `json:"is_trending"`
	OnRequest     bool           `json:"on_request"`
	OnStock       bool           `json:"on_stock"`
}
type AdditionalSpec struct {
	Gear               string `json:"gear"`
	PowerRatedMax      string `json:"power_rated_max"`
	TorqueMoment       string `json:"torque_moment"`
	MaxSpeed           string `json:"max_speed"`
	Battery            string `json:"battery"`
	PowerReserve       string `json:"power_reserve"`
	Weight             string `json:"weight"`
	LoadCap            string `json:"load_cap"`
	DimensionsLWH      string `json:"dimensions_l_w_h"`
	SeatHeight         string `json:"seat_height"`
	Tires              string `json:"tires"`
	Break              string `json:"break"`
	Clearance          string `json:"clearance"`
	Charging           string `json:"charging"`
	MoistureProtection string `json:"moisture_protection"`
}
type PutProductBody struct {
	CategoryId    string         `json:"category_id"`
	Name          string         `json:"name"`
	Short         string         `json:"short"`
	Description   string         `json:"description"`
	Price         int            `json:"price"`
	Discount      int            `json:"discount"`
	Tags          []string       `json:"tags"`
	MainImage     string         `json:"main_image" example:"/images/example.png"`
	XlsxName      string         `json:"xlsx_name"`
	Specification AdditionalSpec `json:"specification"`
	IsTrending    bool           `json:"is_trending"`
	OnRequest     bool           `json:"on_request"`
	OnStock       bool           `json:"on_stock"`
}
type PutCategoryBody struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	ImageName string `json:"image_name"`
}
type NewPost struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	ImageName string `json:"image_name"`
}
type NewComment struct {
	Content string `json:"content"`
}
type PutPostBody struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	ImageName string `json:"image_name"`
}
