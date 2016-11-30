package token

// This is a simple package for go that generates randomized base62 encoded tokens based on a single integer.
// It's ideal for shorturl services or for semi-secured randomized api primary keys.
//
// How it Works
//
// `Token` is an alias for `uint64`.
// Its `Token.Encode()` method interface returns a `Base62` encoded string based off of the number.
// Its implementation of the `encoding.TextMarshaler` and `encoding.TextUnmarshaler` interfaces encodes and
// decodes the `Token` when its being marshalled or unmarshalled as json or xml.
//
// Basically, the outside world will always address the token as its string equivolent and internally we can
// always be used as an `uint64` for fast, indexed, unique, lookups in various databases.
//
// **IMPORTANT:** Remember to always check for collisions when adding randomized tokens to a database

import (
	"bytes"
	"math"
	"math/rand"
	"time"
)

const (
	// Base62 is a string respresentation of every possible base62 character
	Base62 = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// MaxTokenLength is the largest possible character length of a token
	MaxTokenLength = 10

	// MinTokenLength is the smallest possible character length of a token
	MinTokenLength = 1

	// DefaultTokenLength is the default size of a token
	DefaultTokenLength = 9
)

var (
	base62Len = uint64(len(Base62))
)

// init initializes the random number generator
func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// Token is an alias of an uint64 that is marshalled into a base62 encoded token
type Token uint64

// Encode encodes the token into a base62 string
func (t Token) Encode() string {
	bs, _ := t.MarshalText()
	return string(bs)
}

// UnmarshalText implements the `encoding.TextUnmarshaler` interface
func (t *Token) UnmarshalText(data []byte) error {

	number := uint64(0)
	idx := 0.0
	chars := []byte(Base62)

	charsLength := float64(len(chars))
	tokenLength := float64(len(data))

	if tokenLength > MaxTokenLength {
		return ErrTokenTooBig
	} else if tokenLength < MinTokenLength {
		return ErrTokenTooSmall
	}

	for _, c := range data {
		power := tokenLength - (idx + 1)
		index := bytes.IndexByte(chars, c)
		if index < 0 {
			return ErrInvalidCharacter
		}
		number += uint64(index) * uint64(math.Pow(charsLength, power))
		idx++
	}

	// the token was successfully decoded
	*t = Token(number)
	return nil
}

// MarshalText implements the `encoding.TextMarsheler` interface
func (t Token) MarshalText() ([]byte, error) {
	number := uint64(t)
	var chars []byte

	if number == 0 {
		return chars, nil
	}

	for number > 0 {
		result := number / base62Len
		remainder := number % base62Len
		chars = append(chars, Base62[remainder])
		number = result
	}

	for i, j := 0, len(chars)-1; i < j; i, j = i+1, j-1 {
		chars[i], chars[j] = chars[j], chars[i]
	}

	return chars, nil
}

// New returns a `Base62` encoded `Token` of *up to* `DefaultTokenLength`
// if you pass in a `tokenLength` between `MinTokenLength` and `MaxTokenLength` this will return
// a `Token` of *up to* that length instead if you pass in a `tokenLength` that is out of range it will panic
func New(tokenLength ...int) Token {

	// calculate the max hash int based on the token length
	var max uint64
	if tokenLength == nil {
		max = maxHashInt(DefaultTokenLength)
	} else if tl := tokenLength[0]; tl < MinTokenLength {
		panic(ErrTokenTooSmall)
	} else if tl > MaxTokenLength {
		panic(ErrTokenTooBig)
	} else {
		max = maxHashInt(tl)
	}

	// generate a psuedo random token
	number := uint64(rand.Int63n(int64(max & math.MaxInt64)))

	return Token(number)
}

// Decode returns a token from a 1-12 character base62 encoded string
func Decode(token string) (Token, error) {
	var t Token
	err := (&t).UnmarshalText([]byte(token))
	return t, err
}

// maxHashInt returns the largest possible int that will yeild a base62 encoded token of the specified length
func maxHashInt(length int) uint64 {
	return uint64(math.Max(0, math.Min(math.MaxUint64, math.Pow(float64(base62Len), float64(length)))))
}
