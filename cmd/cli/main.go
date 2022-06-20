package main

import (
	config "github.com/jeaguil/piefai/_con"

	polygon "github.com/polygon-io/client-go/rest"
)

func main() {
	c := polygon.New(config.GetAPIKey())
	_ = c
}
