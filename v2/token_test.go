package token

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/marksalpeter/sugar"
)

type Test struct {
	Token    Token
	NilToken *Token
}

func TestToken(t *testing.T) {
	s := sugar.New(t)

	s.Assert("encode and Decode are consistant", func(log sugar.Log) bool {
		original := New()
		if transcoded, err := Decode(original.Encode()); err == nil {
			if original != transcoded {
				log("%s != %d", original, transcoded)
				return false
			}
		} else {
			log(err)
			return false
		}
		return true
	})

	s.Assert("Decode returns errors when tokens are not valid", func(log sugar.Log) bool {

		// invalid spaces
		invalidCharacters := `s p a c e`
		if _, err := Decode(invalidCharacters); err == nil {
			log("tokens with illegal characters in them don't throw an error")
			return false
		}

		// the token is larger than MAX_TOKEN_LENGTH
		tokenTooBig := `sfnalsdasdkasdnaerlaraksfnmaslrasdasadsadas`
		if _, err := Decode(tokenTooBig); err == nil {
			log("tokens larger that MAX_TOKEN_LENGTH don't throw an error")
			return false
		}

		// the token is smaller than MIN_TOKEN_LENGTH
		tokenTooSmall := ``
		if _, err := Decode(tokenTooSmall); err == nil {
			log("tokens larger that MAX_TOKEN_LENGTH don't throw an error")
			return false
		}

		return true
	})

	s.Assert("maxHashInt(tokenLength int) returns tokens of the correct length", func(log sugar.Log) bool {
		for i := MinTokenLength; i <= MaxTokenLength; i++ {
			min := Token(maxHashInt(i - 1)).Encode()
			max := Token(maxHashInt(i) - 1).Encode()
			if len(max) != i {
				log("failed on max -> %d != len(%s)", i, max)
				return false
			} else if len(min) != i {
				log("failed on min -> %d != len(%s)", i, min)
				return false
			}
		}
		return true
	})

	s.Assert("New(tokenLength int) panics when the tokenLength is out of range", func(log sugar.Log) bool {
		isPaniced := true
		testTokenLength := func(tokenLength int) {
			defer func() {
				if r := recover(); r == nil {
					log("did not panic when tokenLength == %d", tokenLength)
					isPaniced = isPaniced && false
				} else {
					isPaniced = isPaniced && true
				}
			}()
			New(tokenLength)
		}
		testTokenLength(MinTokenLength - 1)
		testTokenLength(MaxTokenLength + 1)
		return isPaniced
	})

	s.Assert("json.Marsheler is implemented to encode tokens as base62 strings and decode base62 strings back into tokens", func(log sugar.Log) bool {

		var unmarshaledTest Test
		test := Test{
			Token: New(),
			NilToken: func(token Token) *Token {
				return &token
			}(New()),
		}
		marshaledTest := fmt.Sprintf(
			`{"Token":"%s","NilToken":"%s"}`,
			test.Token.Encode(),
			test.NilToken.Encode(),
		)

		if bytes, err := json.Marshal(&test); err != nil {
			// failed to encode
			log(err)
			return false
		} else if !log.Compare(bytes, marshaledTest) {
			// faild to encode correcty
			return false
		} else if err := json.Unmarshal(bytes, &unmarshaledTest); err != nil {
			// failed to unmarshal
			log(err)
			return false
		} else {
			return log.Compare(unmarshaledTest, test)
		}

	})

	s.Assert("json.Marsheler.UnmarshalJson returns errors when tokens are not valid", func(log sugar.Log) bool {
		var test Test

		// invalid spaces
		invalidCharacters := `{"Token":"s p a c e"}`
		if err := json.Unmarshal([]byte(invalidCharacters), &test); err == nil {
			log("tokens with illegal characters in them don't throw an error")
			return false
		}

		// the token is larger than MaxTokenLength
		tokenTooBig := `{"Token":"sfnalsdasdkasdnaerlaraksfnmaslrasdasadsadas"}`
		if err := json.Unmarshal([]byte(tokenTooBig), &test); err == nil {
			log("tokens larger that MaxTokenLength don't throw an error")
			return false
		}

		// the token is smaller than MinTokenLength
		tokenTooSmall := `{"Token":""}`
		if err := json.Unmarshal([]byte(tokenTooSmall), &test); err == nil {
			log("tokens smaller that MinTokenLength don't throw an error")
			return false
		}

		// non-string tokens
		notAString := `{"Token":true}`
		if err := json.Unmarshal([]byte(notAString), &test); err == nil {
			log("Tokens that are not strings do not throw an error")
			return false
		}

		return true
	})

}
