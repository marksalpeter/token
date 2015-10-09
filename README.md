#Token

This is a simple package for go that generates randomized base62 encoded tokens based on a single int64. It's ideal for shorturl services or for semi-secured randomized api primary keys.

## How it Works

**Token** is an alias for **int64**.  
Its implementation of the **fmt.Stringer** interface returns a **token.BASE62** encoded string based off of the number. 
Its implementation of the **json.Marshaler** interface encodes and decoded the **token.Token** to and from the same **token.BASE62** encoded string representation.  

Basically, the outside world will always address the token as its string equivolent and internally we can always be used as an int64 for fast, indexed, unique, lookups in various databases.

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
	bytes, _ := json.Marshal(&model)
	json.Unmarshal(bytes, &unmarshaled)

	fmt.Println(int64(model.ID))		// 2751173559858
	fmt.Println(model.ID)				// Mr1NSSu
	fmt.Println(string(bytes))			// {"ID":"Mr1NSSu"}
	fmt.Println(int64(unmarshaled.ID))	// 2751173559858
}
```

You can see it in action here:  
// TODO: add a link when it goes public

## Special Mentions

// TODO: get permission from ben to publish this
Special thanks to Ben Baron. The encode and decode functions are ported from a shorturl project of his and he graciously allowed me to publish them (i hope).
