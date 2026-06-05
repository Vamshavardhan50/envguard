// main.go
// Entry point for the envguard CLI.

package main

import "github.com/Vamshavardhan50/envguard/cmd"

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cmd.Execute(version, commit, date)
}
