package test

 import (
	"kv-db/physical"
	"os"
	"testing"
)

func getStorage() *physical.Storage {
	file, _ := os.OpenFile("index", os.O_CREATE, 0666)
	storage := new(physical.Storage)
	storage.Init(file, physical.VALUE_FILE)
	return storage
}

func TestWrite(t *testing.T) {
	storage := getStorage()
	data01 := "ssssssss"
	b1 := []byte(data01)
	offset := storage.Write(b1)
	b2 := storage.Read(offset)
	data02 := string(b2)
	println(data02)
}