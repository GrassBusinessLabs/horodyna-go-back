package requests

type RegisterIamgeModelRequest struct {
	Image ImageModelRequest `json:"image" validate:"required,gte=1,max=40"`
}

type ImageModelRequest struct {
	Name     string `json:"name" validate:"required"`
	Data     string `json:"data" validate:"required"`
	Entity   string `json:"entity" validate:"required"`
	EntityId uint64 `json:"entityId" validate:"required"`
}

type UpdateIamgeModelRequest struct {
	Name    string `json:"name" validate:"required,gte=1,max=40"`
	From    string `json:"from" validate:"required,email"`
	To      string `json:"to" validate:"required,email"`
	GoTrans bool   `json:"go_trans" validate:"required"`
}
