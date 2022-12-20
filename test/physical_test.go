package test

import (
	"kv-db/physical"
	"os"
)

func getStorage() *physical.Storage {
	file, _ := os.OpenFile("index", os.O_CREATE, 0666)
	storage := new(physical.Storage)
	storage.Init(file, physical.VALUE_FILE)
	return storage
}

func write() {
	select {}
}
