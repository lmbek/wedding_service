package tracer

// ContextKey is an interface representing a key suitable for use with context.Context.
// The concrete implementation is unexported to avoid leaking concrete types.
type ContextKey interface {
	// Value returns the actual key object to be used with context.WithValue.
	// The returned value is comparable.
	Value() any
}

type contextKey struct{}

// NewContextKey constructs a new ContextKey instance.
func NewContextKey() ContextKey {
	return &contextKey{}
}

func (k *contextKey) Value() any {
	return k
}
