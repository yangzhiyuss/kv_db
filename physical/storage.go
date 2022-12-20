package physical

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

const (
	IntegerLength = 8
	INDEX_FILE    = 0
	VALUE_FILE    = 1
)

type Storage struct {
	_f        *os.File
	_fileType int
}

func (s *Storage) Init(f *os.File, fileType int) {
	s._f = f
	s._fileType = fileType
	s.initHead()
}

// 假如是索引文件的话就要初始化头节点。值节点直接写入和读取
func (s *Storage) initHead() {
	address := s.seekEnd()
	// 如果文件刚刚打开，进行初始化
	if address == 0 {
		byteArr := make([]byte, IntegerLength)
		_, _ = s._f.Write(byteArr)
		s.SyncStorage()
	}
}

// 获取文件末尾指针
func (s *Storage) seekEnd() int64 {
	endAddress, _ := s._f.Seek(0, io.SeekEnd)
	return endAddress
}

// 获取头
func (s *Storage) seekStartAddress() int64 {
	startAddress, _ := s._f.Seek(0, io.SeekStart)
	return startAddress
}

func (s *Storage) bytesToInt64(byteArr []byte) int64 {
	return int64(binary.LittleEndian.Uint64(byteArr))
}

func (s *Storage) int64ToBytes(integer int64) []byte {
	byteArr := make([]byte, IntegerLength)
	binary.LittleEndian.PutUint64(byteArr, uint64(integer))
	return byteArr
}

func (s *Storage) readInt64() int64 {
	byteArr := make([]byte, IntegerLength)
	_, _ = s._f.Read(byteArr)
	return s.bytesToInt64(byteArr)
}

func (s *Storage) writeInt64(integer int64) {
	_, err := s._f.Write(s.int64ToBytes(integer))
	if err != nil {
		fmt.Println(err.Error())
	}
}

// 追加写入
func (s *Storage) Write(data []byte) int64 {
	address := s.seekEnd()
	s.writeInt64(int64(len(data)))
	_, _ = s._f.Write(data)
	return address
}

func (s *Storage) CommitWrite(data []byte) int64 {
	address := s.seekEnd()
	s.writeInt64(int64(len(data)))
	_, _ = s._f.Write(data)
	return address
}

// 读取相应地址的数据
func (s *Storage) Read(address int64) []byte {
	_, _ = s._f.Seek(address, io.SeekStart)
	length := s.readInt64()
	data := make([]byte, length)
	_, _ = s._f.Read(data)
	return data
}

func (s *Storage) CommitRootAddress(rootAddress int64) {
	if s._fileType == INDEX_FILE {
		s.seekStartAddress()
		s.writeInt64(rootAddress)
	}

}

func (s *Storage) GetRootAddress() int64 {
	if s._fileType == INDEX_FILE {
		s.seekStartAddress()
		return s.readInt64()
	} else {
		return -1
	}

}

func (s *Storage) SyncStorage() {
	_ = s._f.Sync()
}

func (s *Storage) Close() {
	_ = s._f.Close()
}

func (s *Storage) Closed() bool {
	fileInfo, _ := s._f.Stat()
	return fileInfo == nil
}
