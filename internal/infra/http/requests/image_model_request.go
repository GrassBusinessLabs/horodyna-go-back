package requests

type RegisterIamgeModelRequest struct {
	Image ImageRequest `json:"title" validate:"required,gte=1,max=40"`
}

type UpdateIamgeModelRequest struct {
	Name    string `json:"name" validate:"required,gte=1,max=40"`
	From    string `json:"from" validate:"required,email"`
	To      string `json:"to" validate:"required,email"`
	GoTrans bool   `json:"go_trans" validate:"required"`
}
