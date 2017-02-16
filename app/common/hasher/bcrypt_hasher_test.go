package hasher

func ExampleBcryptHasher_Hash() {

	hasher := BcryptHasher{}

	hasher.Hash("*72t723c(#fji3)@")
}
