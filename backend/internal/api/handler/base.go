package handler

type AuthParam struct {
	Authorization string `required:"true" header:"Authorization" example:"BEARER <token>" doc:"BEARER <token>"`
}
