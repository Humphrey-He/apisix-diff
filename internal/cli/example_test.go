package cli_test

import (
	"fmt"
)

func Example_planExitCodes() {
	fmt.Println("0: no diff, validation passed")
	fmt.Println("1: diff detected")
	fmt.Println("2: validation failed")
	// Output:
	// 0: no diff, validation passed
	// 1: diff detected
	// 2: validation failed
}

func Example_planCommand() {
	fmt.Println("apidiff plan -f ./apisix.yaml --admin-url http://127.0.0.1:9180 --token <X-API-KEY>")
	// Output:
	// apidiff plan -f ./apisix.yaml --admin-url http://127.0.0.1:9180 --token <X-API-KEY>
}

func Example_validateCommand() {
	fmt.Println("apidiff validate -f ./apisix.yaml --rules ./rules.yaml")
	// Output:
	// apidiff validate -f ./apisix.yaml --rules ./rules.yaml
}
