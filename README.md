#Token

This is a simple package for go that generates randomized base62 encoded tokens based on an integer. It's ideal for short url services or for any short, unique, randomized tokens you need to use throughout your app.

## How it Works

`Token` is an alias for `uint64`.

The `Token.Encode()` method returns a base62 encoded string based off of the uint64.

`Token` implements the `encoding.TextMarshaler` and `encoding.TextUnmarshaler` interfaces to encode and decode to and from the base62 string representation of the `uint64`

Basically, the outside world will always see the token as a base62 encoded string, but in your app you will always be able to use the token as a `uint64` for fast, indexed, unique, lookups in various databases.

**IMPORTANT:** Remember to always check for collisions when adding randomized tokens to a database

## Example

```go  
package main

import (
	"fmt"
	"github.com/marksalpeter/token"
)

type Model struct {
    ID	token.Token `json:"id"`
}

func main() {
	// create a new model
	model := Model {
		ID:	token.New(), // creates a new, random uint64 token
	}
	fmt.Println(model.ID)          // 2751173559858
	fmt.Println(model.ID.Encode()) // Mr1NSSu

	// encode the model as json
	marshaled, err := json.Marshal(&model)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(marshaled)) // {"id":"Mr1NSSu"}

	// decode the model
	var unmarshaled Model
	if err := json.Unmarshal(marshaled, &unmarshaled); err != nil {
		panic(err)
	}
	fmt.Println(unmarshaled.ID)    // 2751173559858

}
```

## Special Mentions

Special thanks to [@einsteinx2](https://github.com/einsteinx2). The encode and decode functions are ported from a short url project of his and he graciously allowed me to publish them.
