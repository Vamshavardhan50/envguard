// testdata/go-project/main.go
// Sample Go usage for envguard tests.

package main

import "os"

func main() {
	port := os.Getenv("SERVER_PORT")
	secret := os.Getenv("DB_PASSWORD")
	_ = port
	_ = secret
}
