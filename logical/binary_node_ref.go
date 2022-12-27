package logical

import (
	"google.golang.org/protobuf/proto"
	"kv-db/physical"
	nodebytes "kv-db/protoc"
)

type BinaryNode struct {
	LeftRef  *BinaryNodeRef
	Key      string
	ValueRef    *ValueRef
	RightRef *BinaryNodeRef
	Length   int64
}

func (bn *BinaryNode) storeRefs(valueStorage *physical.Storage, indexStorage *physical.Storage) {
	bn.ValueRef.Store(valueStorage)
	bn.LeftRef.Store(valueStorage, indexStorage)
	bn.RightRef.Store(valueStorage, indexStorage)
}

func (bn *BinaryNode) fromNode(param map[string]interface{}) *BinaryNode {
	newNode := &BinaryNode{
		LeftRef: bn.LeftRef,
		Key: bn.Key,
		ValueRef: bn.ValueRef,
		RightRef: bn.RightRef,
		Length: bn.Length,
	}

	if param["left_ref"] != nil {
		nodeRef := param["left_ref"].(*BinaryNodeRef)
		newNode.Length += nodeRef.Length() - bn.LeftRef.Length()
		newNode.LeftRef = nodeRef
	} else if param["right_ref"] != nil {
		nodeRef := param["left_ref"].(*BinaryNodeRef)
		newNode.Length += nodeRef.Length() - bn.RightRef.Length()
		newNode.LeftRef = nodeRef
	} else if param["value_ref"] != nil {
		newValue := param["value_ref"].(*ValueRef)
		newNode.ValueRef = newValue
	}
	return newNode

}

// BinaryNodeRef /*------------------------------------------------------*/
/* 节点映射 */
type BinaryNodeRef struct {
	ValueRef
}

func (bnr *BinaryNodeRef) Length() int64 {
	if bnr._referent == nil && bnr._address != 0 {
		return -1
	}
	if bnr._referent != nil {
		node:= bnr._referent.(*BinaryNode)
		return node.Length
	} else {
		return 0
	}
}

// 对值进行序列化
func (bnr *BinaryNodeRef) referentToBytes(referent interface{}) []byte {
	binaryNode := referent.(*BinaryNode)
	data := &nodebytes.BinaryNodeRefBytes{
		LeftRef: binaryNode.LeftRef.Address(),
		Key: binaryNode.Key,
		ValueRef: binaryNode.ValueRef.Address(),
		RightRef: binaryNode.RightRef.Address(),
		Length: binaryNode.Length,
	}
	byteArr, _ := proto.Marshal(data)
	return byteArr
}

// 对值反序列化
func (bnr *BinaryNodeRef) bytesToReferent(byteArr []byte) interface{} {
	data := new(nodebytes.BinaryNodeRefBytes)
	_ =proto.Unmarshal(byteArr, data)
	leftRef := &BinaryNodeRef{ValueRef{_address: data.LeftRef}}
	key := data.Key
	valueRef := &ValueRef{_address: data.ValueRef}
	rightRef := &BinaryNodeRef{ValueRef{_address: data.RightRef}}
	length := data.Length
	return &BinaryNode{leftRef, key, valueRef, rightRef, length}
}

func (bnr *BinaryNodeRef) prepareToStore(valueStorage *physical.Storage, indexStorage *physical.Storage) {
	if bnr._referent != nil {
		bnr._referent.(*BinaryNode).storeRefs(valueStorage, indexStorage)
	}
}

func (bnr *BinaryNodeRef) Get(indexStorage *physical.Storage) *BinaryNode {
	if bnr._referent == nil && bnr._address != 0 {
		data := indexStorage.Read(bnr._address)
		bnr._referent = bnr.bytesToReferent(data)
	}
	return bnr._referent.(*BinaryNode)
}

func (bnr *BinaryNodeRef) Store(valueStorage *physical.Storage, indexStorage *physical.Storage) {
	if bnr._referent != nil && bnr._address == 0 {
		bnr.prepareToStore(valueStorage, indexStorage)
		data := bnr.referentToBytes(bnr._referent)
		bnr._address = indexStorage.Write(data)
	}
}
