package services

// AppsListResponseItemModel ...
type AppsListResponseItemModel struct {
	Title string `json:"title"`
	Slug  string `json:"slug"`
	Owner struct {
		Name string `json:"name"`
	} `json:"owner"`
}

// AppsListResponseModel ...
type AppsListResponseModel struct {
	Data []AppsListResponseItemModel `json:"data"`
}
