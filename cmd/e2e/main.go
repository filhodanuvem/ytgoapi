package main

import "github.com/filhodanuvem/ytgoapi/e2e"

func main() {
	ts := e2e.NewPostSuccessfulSuite()
	ts.Run()
}
