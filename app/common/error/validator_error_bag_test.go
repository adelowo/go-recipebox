package error

import (
	"flag"
	"fmt"
	"os"
	"testing"
)

var validatorErrorBag *ValidatorErrorBag

func TestMain(m *testing.M) {

	validatorErrorBag = NewValidatorErrorBag()

	flag.Parse()

	os.Exit(m.Run())
}

func ExampleValidatorErrorBag_Add() {

	message := "Some message you want to tell the user based on his/her username"

	//Add a message with the key "username"
	validatorErrorBag.Add("username", message)

}

func ExampleValidatorErrorBag_Get() {

	message := "Some message you want to tell the user based on his/her username"

	//Add a message with the key "username"
	validatorErrorBag.Add("username", message)

	//Fetch a non existent key
	_, err := validatorErrorBag.Get("password")

	if err != nil {
		fmt.Println(err)
		//Output: Key, password does not exist on this bag
	}

}

func ExampleValidatorErrorBag_Count() {

	//We still have the username key on it

	fmt.Println(validatorErrorBag.Count())
	//Output: 1
}

func ExampleValidatorErrorBag_Has() {

	exists := validatorErrorBag.Has("password")
	fmt.Println(exists)
	//Output: false
}
