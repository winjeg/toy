package main

import (
	"github.com/winjeg/toy/impl"
	"github.com/winjeg/toy/server"
)

func main() {
	server.Run(impl.StrStore, "123456", 6378)
}
