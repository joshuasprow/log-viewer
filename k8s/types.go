package k8s

func Ptr[V any](v V) *V {
	return &v
}

type Result[V any] struct {
	V   V
	Err error
}
