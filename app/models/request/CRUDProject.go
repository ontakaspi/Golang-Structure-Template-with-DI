package request

type CreateExampleProject struct {
	Name string `json:"name" validate:"required"`
}
