package token

import (
	"testing"
	"github.com/marksalpeter/sugar"
	"encoding/json"
	"regexp"
	"fmt"
)

type Test struct {
	Token		Token
	NilToken 	*Token
}

func TestToken(t * testing.T) {
	sugar.New(t).
	
	Assert("encode and Decode are consistant", func (log sugar.Log) bool {
		original := New()
		if transcoded, err := Decode(encode(original)); err == nil {
			if original != transcoded {
				log("%s != %d", original, transcoded)
				return false
			}	
		} else {
			log(err)
			return false
		}
		return true
	}).
	
	Assert("Decode returns errors when tokens are not valid", func (log sugar.Log) bool {
		
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
	}).
	
	Assert("fmt.Stringer is implemented to return the encoded token", func (log sugar.Log) bool {
		token := New()
		return encode(token) == token.String()
	}).
	
	Assert("maxHashInt(tokenLength int) returns tokens of the correct length", func (log sugar.Log) bool {
		for i := MIN_TOKEN_LENGTH; i <= MAX_TOKEN_LENGTH; i++ {
			min := encode(Token(maxHashInt(i - 1)))
			max := encode(Token(maxHashInt(i) - 1))
			if len(max) != i {
				log("failed on max -> %d != len(%s)", i, max)
				return false
			} else if len(min) != i {
				log("failed on min -> %d != len(%s)", i, min)
				return false
			} 
		}		
		return true		
	}).
	
	Assert("New(tokenLength int) panics when the tokenLength is out of range", func (log sugar.Log) bool {
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
			New(tokenLength - 1)
		}
		testTokenLength(MIN_TOKEN_LENGTH - 1)
		testTokenLength(MAX_TOKEN_LENGTH + 1)
		return isPaniced
	}).
	
	Assert("json.Marsheler is implemented to encode tokens as base62 strings and decode base62 strings back into tokens", func (log sugar.Log) bool {
		
		var token = New()
		var test = Test {
			Token: New(),
			NilToken: &token,
		}
		testRegexp := regexp.MustCompile(fmt.Sprintf("\"Token\":\"%s\",\"NilToken\":\"%s\"", test.Token, test.NilToken))
		
		// encoding
		if bytes, err := json.Marshal(&test); err == nil {
			if !testRegexp.Match(bytes) {
				log("%+v != %s", test, bytes)
				return false
			}
			
			// decoding
			var unmarshal Test
			if err := json.Unmarshal(bytes, &unmarshal); err == nil {
				if test.Token != unmarshal.Token || *test.NilToken != *unmarshal.NilToken {
					log("%+v != %+v (expected)", unmarshal, test)
					return false
				}
			} else {
				log("decoding failed: %s", err)
				return false
			}
			
		} else {
			log("encoding failed: %s", err)
			return false
		}
		
		return true
	}).
	
	Assert("json.Marsheler.UnmarshalJson returns errors when tokens are not valid", func (log sugar.Log) bool {
		var test Test
		
		// invalid spaces
		invalidCharacters := `{"Token":"s p a c e"}`
		if err := json.Unmarshal([]byte(invalidCharacters), &test); err == nil {
			log("tokens with illegal characters in them don't throw an error")
			return false
		}
		
		// the token is larger than MAX_TOKEN_LENGTH
		tokenTooBig := `{"Token":"sfnalsdasdkasdnaerlaraksfnmaslrasdasadsadas"}`
		if err := json.Unmarshal([]byte(tokenTooBig), &test); err == nil {
			log("tokens larger that MAX_TOKEN_LENGTH don't throw an error")
			return false
		}
		
		// the token is smaller than MIN_TOKEN_LENGTH
		tokenTooSmall := `{"Token":""}`
		if err := json.Unmarshal([]byte(tokenTooSmall), &test); err == nil {
			log("tokens larger that MAX_TOKEN_LENGTH don't throw an error")
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