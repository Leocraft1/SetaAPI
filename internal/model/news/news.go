package news

type AllNews struct {
	Title string `json:"title"`
	Date  string `json:"date"`
	Link  string `json:"link"`
	Type  string `json:"type,omitempty"`
}

type AllNewsResponse struct {
	News []AllNews `json:"news"`
}

type News struct {
	Title   string `json:"title"`
	Date    string `json:"date"`
	Content string `json:"content"`
}
