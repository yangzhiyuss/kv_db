package logical

import (
	"encoding/json"
	"fmt"
	"kv-db/physical"
)

type BinaryNode struct {
	LeftRef  *BinaryNodeRef
	Key      string
	Value    *ValueRef
	RightRef *BinaryNodeRef
	Length   int64
}

func (bn *BinaryNode) storeRefs(valueStorage *physical.Storage, indexStorage *physical.Storage) {
	bn.Value.store(valueStorage)
	bn.LeftRef.store(valueStorage, indexStorage)
	bn.RightRef.store(valueStorage, indexStorage)
}

func (bn *BinaryNode) fromNode(param map[string]interface{}) {

}

// BinaryNodeRef /*------------------------------------------------------*/
/* 节点映射 */
type BinaryNodeRef struct {
	ValueRef
}

// 对值进行序列化
func (bnr *BinaryNodeRef) referentToBytes(referent interface{}) []byte {
	binaryNode := referent.(*BinaryNode)
	data := map[string]interface{}{
		"left":   binaryNode.LeftRef.address(),
		"key":    binaryNode.Key,
		"value":  binaryNode.Value.address(),
		"right":  binaryNode.RightRef.address(),
		"length": binaryNode.Length,
	}
	byteArr, _ := json.Marshal(data)
	return byteArr
}

// 对值反序列化
func (bnr *BinaryNodeRef) bytesToReferent(byteArr []byte) interface{} {
	var data map[string]interface{}
	err := json.Unmarshal(byteArr, &data)
	if err != nil {
		fmt.Println(err.Error())
	}
	leftRef := &BinaryNodeRef{ValueRef{_address: int64(data["left"].(float64))}}
	key := data["key"].(string)
	valueRef := &ValueRef{_address: int64(data["value"].(float64))}
	rightRef := &BinaryNodeRef{ValueRef{_address: int64(data["right"].(float64))}}
	length := int64(data["length"].(float64))
	return &BinaryNode{leftRef, key, valueRef, rightRef, length}
}

func (bnr *BinaryNodeRef) prepareToStore(valueStorage *physical.Storage, indexStorage *physical.Storage) {
	if bnr._referent != nil {
		bnr._referent.(*BinaryNode).storeRefs(valueStorage, indexStorage)
	}
}

func (bnr *BinaryNodeRef) get(indexStorage *physical.Storage) *BinaryNode {
	if bnr._referent == nil && bnr._address != 0 {
		data := indexStorage.Read(bnr._address)
		bnr._referent = bnr.bytesToReferent(data)
	}
	return bnr._referent.(*BinaryNode)
}

func (bnr *BinaryNodeRef) store(valueStorage *physical.Storage, indexStorage *physical.Storage) {
	if bnr._referent != nil && bnr._address == 0 {
		bnr.prepareToStore(valueStorage, indexStorage)
		data := bnr.referentToBytes(bnr._referent)
		bnr._address = indexStorage.Write(data)
	}
}
