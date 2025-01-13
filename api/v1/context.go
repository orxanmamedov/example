package apiv1

type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "context value " + k.name
}
