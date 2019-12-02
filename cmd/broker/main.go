package main

import broker "github.com/ingtranet/knaing-broker"

func main() {
	app := broker.NewApp()
	app.Run()
}
