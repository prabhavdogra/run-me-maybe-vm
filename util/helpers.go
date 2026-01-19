package util

// OneOf checks if a value exists in a list of values
func OneOf[T comparable](value T, options ...T) bool {
	for _, option := range options {
		if value == option {
			return true
		}
	}
	return false
}

// OneOf checks if a value doesn't exist in a list of values
func NotOneOf[T comparable](value T, options ...T) bool {
	for _, option := range options {
		if value == option {
			return false
		}
	}
	return true
}

// GetPointer returns a pointer to the given value
func GetPointer[T any](value T) *T {
	return &value
}
