package logical

import (
	"encoding/json"
	"fmt"
	"kv-db/physical"
)

type ValueRef struct {
	_referent interface{}
	_address  int64
}

func (vr *ValueRef) Address() int64 {
	return vr._address
}

// 对值进行序列化
func (vr *ValueRef) referentToBytes(referent interface{}) []byte {
	byteArr, err := json.Marshal(referent)
	if err != nil {
		fmt.Println(err.Error())
	}
	return byteArr
}

// 对值反序列化
func (vr *ValueRef) bytesToReferent(byteArr []byte) interface{} {
	var referent interface{}
	err := json.Unmarshal(byteArr, referent)
	if err != nil {
		fmt.Println(err.Error())
	}
	return referent
}

func (vr *ValueRef) Get(valueStorage *physical.Storage) interface{} {
	//referent为空，
	//地址不为空，还没有进行读取
	if vr._referent == nil && vr._address != 0 {
		data := valueStorage.Read(vr._address)
		referent := vr.bytesToReferent(data)
		vr._referent = referent
	}
	return vr._referent
}

func (vr *ValueRef) Store(valueStorage *physical.Storage) {
	// referent不为空，
	//地址为0，还没有进行存储
	if vr._referent != nil && vr._address == 0 {
		data := vr.referentToBytes(vr._referent)
		vr._address = valueStorage.Write(data)
	}
}
