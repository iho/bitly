package handlers

type UrlCreationRequest struct {
	Url string `json:"url"`
}

type UrlCreationResponse struct {
	Code string `json:"code"`
}

type UrlGetRequest struct {
	Code string `uri:"code" binding:"required"`
}
