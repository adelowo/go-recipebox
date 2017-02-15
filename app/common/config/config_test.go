package config

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func ExampleReadConfig() {

	config, _ := ReadConfig("../../../config/config.json")

	fmt.Println(config)
	//Output: {:9091 localhost}
}

func TestReadConfig(t *testing.T) {

	config, err := ReadConfig("../../../config/config.json")

	if err == nil {
		t.Log("Passing")
	} else {
		t.Error("Failing")
	}

	c := &Configuration{":9091", "localhost"}

	assert.ObjectsAreEqual(c, config)
}

func TestReadConfigFailsBecauseOfInvalidFile(t *testing.T) {

	_, err := ReadConfig("../../../config/config.jso")

	if err != nil {
		t.Log("Test passed: An invalid file cannot be read")
	} else {
		t.Error("Test is failing... Why is an invalid file being read ?")
	}

}
