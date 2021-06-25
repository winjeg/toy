package main

import (
	"fmt"
	"github.com/winjeg/toy/store"
)

func main() {
	defer store.CleanUp()
	store.Set("k2", "v2")
	d, _ := store.Get("k2")
	fmt.Println("current value of k2: " + string(d))
	store.Del("k2")
	d2, _ := store.Get("k2")
	fmt.Println("current value after del: " + string(d2))
}