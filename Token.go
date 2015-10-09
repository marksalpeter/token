package token

import (
	"math"
	"math/rand"
	"fmt"
	"bytes"
	"time"
	"errors"
)

const (
	BASE62					= "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	DEFAULT_TOKEN_LENGTH	= 7
	MAX_TOKEN_LENGTH 		= 11
	MIN_TOKEN_LENGTH 		= 1
)

var (
	base62Len				= len(BASE62)
)

// this is an alias of an int64 that is json marshalled and stringified into a base62 encoded token
type Token uint64

// implements the **fmt.Stringer** interface to encode the token into a base62 string
func (t Token) String () string {
	return encode(t)
}

// implements the **json.Marsheler** interface to decode the token from a base62 string back into an int64
func (t *Token) UnmarshalJSON(data []byte) error {
	length		:= len(data)
	str			:= string(data)
	isString	:= (str[0] == '"' || str[0] == '\'') && (str[length - 1] == '"' || str[length - 1] == '\'')
	if isString {
		if decoded, err := Decode(str[1:length - 1]); err == nil {
			*t = decoded
			return nil
		} else {
			return err
		}
	} else {
		return errors.New("attempted to parse a non-string token") 
	}
}

// implements the **json.Marsheler** interface to encode the token into a base62 string
func (t Token) MarshalJSON() ([]byte, error) {
	token := encode(t)
	if token == "" {
		return []byte{}, nil
	} else {
	    return []byte("\"" + encode(t) + "\""), nil		
	}
}

// returns a **BASE62** encoded **Token** of *up to* **DEFAULT_TOKEN_LENGTH**  
// if you pass in a **tokenLength** between **MIN_TOKEN_LENGTH** and **MAX_TOKEN_LENGTH** this will return 
// a **Token** of *up to* that length instead if you pass in a **tokenLength** that is out of range it will panic
func New(tokenLength ...int) Token {
	
	// calculate the max hash int based on the token length
	var max uint64
	if tokenLength != nil {
		isInRange := tokenLength[0] >= MIN_TOKEN_LENGTH && tokenLength[0] <= MAX_TOKEN_LENGTH 
		if isInRange {
			max = maxHashInt(tokenLength[0])
		} else {
			panic(fmt.Sprintf("tokenLength ∉ [%d,%d]", MIN_TOKEN_LENGTH, MAX_TOKEN_LENGTH))
			return Token(0)
		}	
	} else {
		max = maxHashInt(DEFAULT_TOKEN_LENGTH)
	}
	
	// generate a psuedo random token
	rand.Seed(time.Now().UTC().UnixNano())
	number := uint64(rand.Int63n(int64(max & math.MaxInt64)))
	
	return Token(number)
}

// returns a token from a 1-12 character base62 encoded string
func Decode(token string) (Token, error) {
	
    number := uint64(0)
    idx    := 0.0
    chars  := []byte(BASE62)

	charsLength := float64(len(chars))
	tokenLength := float64(len(token))
	
	if tokenLength > MAX_TOKEN_LENGTH {
		return Token(0), errors.New(fmt.Sprintf("%d > MAX_TOKEN_LENGTH (%d)", int(tokenLength), MAX_TOKEN_LENGTH))
	} else if tokenLength < MIN_TOKEN_LENGTH {
		return Token(0), errors.New(fmt.Sprintf("%d < MIN_TOKEN_LENGTH (%d)", int(tokenLength), MIN_TOKEN_LENGTH))
	}

	for _, c := range []byte(token) {
		power := tokenLength - (idx + 1)
		index := bytes.IndexByte(chars, c)
		if index < 0 {
			return Token(0), errors.New(fmt.Sprintf("%q is not present in %s", c, BASE62))
		}
		number += uint64(index) * uint64(math.Pow(charsLength, power))
		idx++
	}

	return Token(number), nil
}

func encode(token Token) string {
	
	number	:= uint64(token)
	if number == 0 {
		return ""
    }
	
    chars 	:= make([]byte, 0)
    length	:= uint64(len(BASE62))

	for number > 0 {
		result    := number / length
		remainder := number % length
		chars   = append(chars, BASE62[remainder])		
		number  = result
	}

	for i, j := 0, len(chars) - 1; i < j; i, j = i + 1, j - 1 {
		chars[i], chars[j] = chars[j], chars[i]
	}

    return string(chars)
}

// returns the largest possible int that will yeild a base62 encoded token of the specified length
func maxHashInt(length int) uint64 {
	return uint64(math.Max(0, math.Min(math.MaxUint64, math.Pow(float64(base62Len), float64(length)))))
}