#Token

This is a simple package for go that generates randomized base62 encoded tokens based on a single integer. It's ideal for shorturl services or for semi-secured randomized api primary keys.

## How it Works

`Token` is an alias for `uint64`.  
Its `Token.Encode()` method interface returns a `BASE62` encoded string based off of the number.  
Its implementation of the `json.Marshaler` interface encodes and decoded the `Token` to and from the same 
`BASE62` encoded string representation.
Its implementation of the `json.Marshaler` interface encodes and decoded the `Token` to and from the same `BASE62` encoded string representation.  

Basically, the outside world will always address the token as its string equivolent and internally we can always be used as an `uint64` for fast, indexed, unique, lookups in various databases.

**IMPORTANT:** Remember to always check for collisions when adding randomized tokens to a database

## Example

```go  
package main

import (
	"fmt"
	"github.com/marksalpeter/token"	
)

type Model struct {
    ID	token.Token
}

func main() {
	model := Model {
		ID:	token.New(),
	}
	var unmarshaled Model
	marshaled, _ := json.Marshal(&model)
	json.Unmarshal(marshaled, &unmarshaled)

	fmt.Println(uint64(model.ID))		// 2751173559858
	fmt.Println(model.ID)				// Mr1NSSu
	fmt.Println(string(marshaled))		// {"ID":"Mr1NSSu"}
	fmt.Println(uint64(unmarshaled.ID))	// 2751173559858
}
```

You can see it in action here:  
// TODO: add a link when it goes public

## Special Mentions

Special thanks to [@einsteinx2](https://github.com/einsteinx2). The encode and decode functions are ported from a shorturl project of his and he graciously allowed me to publish them.
// TODO: get permission from ben to publish this