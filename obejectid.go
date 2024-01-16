package tinyid

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"sync/atomic"
	"time"
)

// 改良Mongodb ObjectId生成算法
// 时间戳部分由mongodb的4字节提升为5字节, 避免2038的问题
// 同时pid部分由mongodb的5字节改为4字节, 整体保持12字节不变

var (
	processUnique   = processUniqueBytes()
	objectIdCounter = readRandomUint32()
)

type ObjectID [12]byte

func (id ObjectID) Hex() string {
	return hex.EncodeToString(id[:])
}

func NewObjectID() ObjectID {
	return NewObjectIDFromTimestamp(time.Now())
}

func NewObjectIDFromTimestamp(timestamp time.Time) ObjectID {
	var b [12]byte

	binary.BigEndian.PutUint32(b[0:5], uint32(timestamp.Unix()))
	copy(b[5:9], processUnique[:])
	putUint24(b[9:12], atomic.AddUint32(&objectIdCounter, 1))

	return b
}

func putUint24(b []byte, v uint32) {
	b[0] = byte(v >> 16)
	b[1] = byte(v >> 8)
	b[2] = byte(v)
}

func processUniqueBytes() [4]byte {
	var b [4]byte
	_, err := io.ReadFull(rand.Reader, b[:])
	if err != nil {
		panic(fmt.Errorf("cannot initialize objectid package with crypto.rand.Reader: %w", err))
	}

	return b
}

func readRandomUint32() uint32 {
	var b [4]byte
	_, err := io.ReadFull(rand.Reader, b[:])
	if err != nil {
		panic(fmt.Errorf("cannot initialize objectid package with crypto.rand.Reader: %w", err))
	}

	return (uint32(b[0]) << 0) | (uint32(b[1]) << 8) | (uint32(b[2]) << 16) | (uint32(b[3]) << 24)
}
