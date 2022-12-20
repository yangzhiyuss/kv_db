package inter

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestClient(t *testing.T) {
	student := &Student{Person{
		name:    "sssss",
		address: 11,
	}}
	PrintMsg(student)
	person := &Person{
		name:    "aaaa",
		address: 111,
	}
	PrintMsg(person)
}

func TestJson(t *testing.T) {
	data := map[string]interface{}{
		"age":  11,
		"name": "zhangsan",
	}

	byteArr, _ := json.Marshal(data)

	var newData map[string]interface{}
	_ = json.Unmarshal(byteArr, &newData)
	fmt.Println(int64(newData["age"].(float64)))
	fmt.Println(newData["name"].(string))
}
