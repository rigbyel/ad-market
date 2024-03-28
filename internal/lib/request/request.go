package request

type UserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type AdvertRequest struct {
	Header   string `json:"header"`
	Body     string `json:"body,omitempty"`
	Price    int    `json:"price"`
	ImageURL string `json:"image_url"`
}
