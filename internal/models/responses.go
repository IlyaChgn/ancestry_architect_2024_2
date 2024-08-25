package models

type OkResponse struct {
	Code int `json:"code"`
	Body any `json:"body"`
}

type ErrResponse struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
}

type SeveralErrsResponse struct {
	Code   int      `json:"code"`
	Errors []string `json:"errors"`
}
