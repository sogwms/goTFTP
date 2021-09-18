package main

import (
	"tftp/tftp"
)

func main() {
	println("Running...")
	tftp.Server(".")
}
