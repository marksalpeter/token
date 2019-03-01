package token

// Error is a token error
type Error string

// Error implements the `errors.Error` interface
func (e Error) Error() string {
	return string(e)
}

const (
	// ErrTokenTooSmall is the error returned or panic'd when a base62 token is smaller than `MinTokenLength`
	ErrTokenTooSmall = Error("the base62 token is smaller than MinTokenLength")

	// ErrTokenTooBig is the error returned or panic'd when a base62 token is larger than `MaxTokenLength`
	ErrTokenTooBig = Error("the base62 token is larger than MaxTokenLength")

	// ErrInvalidCharacter is the error returned or panic'd when a non `Base62` string is being parsed
	ErrInvalidCharacter = Error("there was a non base62 character in the token")
)
