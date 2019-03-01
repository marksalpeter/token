package token_test

import (
	"encoding/json"
	"fmt"

	"github.com/marksalpeter/token"
)

type Model struct {
	ID token.Token
}

func ExampleToken() {

	model := Model{
		ID: token.New(),
	}
	var unmarshaled Model
	marshaled, _ := json.Marshal(&model)
	json.Unmarshal(marshaled, &unmarshaled)

	fmt.Println(model.ID)
	fmt.Println(model.ID.Encode())
	fmt.Println(string(marshaled))
	fmt.Println(uint64(unmarshaled.ID))

	// Output:
	// 2751173559858
	// Mr1NSSu
	// {"ID":"Mr1NSSu"}
	// 2751173559858

}
