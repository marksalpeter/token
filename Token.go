package token

// This is a simple package for go that generates randomized base62 encoded tokens based on a single integer.
// It's ideal for shorturl services or for semi-secured randomized api primary keys.
//
// How it Works
//
// `Token` is an alias for `uint64`.
// Its `Token.Encode()` method interface returns a `Base62` encoded string based off of the number.
// Its implementation of the `json.Marshaler` interface encodes and decoded the `Token` to and from the same
// `Base62` encoded string representation.
//
// Basically, the outside world will always address the token as its string equivolent and internally we can
// always be used as an `uint64` for fast, indexed, unique, lookups in various databases.
//
// **IMPORTANT:** Remember to always check for collisions when adding randomized tokens to a database

import (
	"bytes"
	"fmt"
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

// Token is an alias of an int64 that is json marshalled into a base62 encoded token
type Token uint64

// Encode encodes the token into a base62 string
func (t Token) Encode() string {
	number := uint64(t)
	if number == 0 {
		return ""
	}

	var chars []byte
	for number > 0 {
		result := number / base62Len
		remainder := number % base62Len
		chars = append(chars, Base62[remainder])
		number = result
	}

	for i, j := 0, len(chars)-1; i < j; i, j = i+1, j-1 {
		chars[i], chars[j] = chars[j], chars[i]
	}

	return string(chars)
}

// UnmarshalText implements the `encoding.TextMarshaler` interface
func (t *Token) UnmarshalText(data []byte) error {
	str := string(data)
	strLen := len(data)

	// decode the token
	decoded, err := Decode(str[1 : strLen-1])
	if err != nil {
		return err
	}

	// the token was successfully decoded
	*t = decoded
	return nil
}

// MarshalText implements the `encoding.TextMarsheler` interface
func (t Token) MarshalText() ([]byte, error) {
	return []byte(t.Encode()), nil
}

// New returns a `Base62` encoded `Token` of *up to* `DefaultTokenLength`
// if you pass in a `tokenLength` between `MinTokenLength` and `MaxTokenLength` this will return
// a `Token` of *up to* that length instead if you pass in a `tokenLength` that is out of range it will panic
func New(tokenLength ...int) Token {

	// calculate the max hash int based on the token length
	var max uint64
	if tokenLength != nil {
		isInRange := tokenLength[0] >= MinTokenLength && tokenLength[0] <= MaxTokenLength
		if isInRange {
			max = maxHashInt(tokenLength[0])
		} else {
			panic(fmt.Errorf("tokenLength ∉ [%d,%d]", MinTokenLength, MaxTokenLength))
		}
	} else {
		max = maxHashInt(DefaultTokenLength)
	}

	// generate a psuedo random token
	rand.Seed(time.Now().UTC().UnixNano())
	number := uint64(rand.Int63n(int64(max & math.MaxInt64)))

	return Token(number)
}

// Decode returns a token from a 1-12 character base62 encoded string
func Decode(token string) (Token, error) {

	number := uint64(0)
	idx := 0.0
	chars := []byte(Base62)

	charsLength := float64(len(chars))
	tokenLength := float64(len(token))

	if tokenLength > MaxTokenLength {
		return Token(0), fmt.Errorf("%d > MaxTokenLength (%d)", int(tokenLength), MaxTokenLength)
	} else if tokenLength < MinTokenLength {
		return Token(0), fmt.Errorf("%d < MinTokenLength (%d)", int(tokenLength), MinTokenLength)
	}

	for _, c := range []byte(token) {
		power := tokenLength - (idx + 1)
		index := bytes.IndexByte(chars, c)
		if index < 0 {
			return Token(0), fmt.Errorf("%q is not present in %s", c, Base62)
		}
		number += uint64(index) * uint64(math.Pow(charsLength, power))
		idx++
	}

	return Token(number), nil
}

// maxHashInt returns the largest possible int that will yeild a base62 encoded token of the specified length
func maxHashInt(length int) uint64 {
	return uint64(math.Max(0, math.Min(math.MaxUint64, math.Pow(float64(base62Len), float64(length)))))
}
