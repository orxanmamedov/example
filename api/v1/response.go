package apiv1

type responseFoo struct {
	ID string `validate:"required" json:"id"`
}
