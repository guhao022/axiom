package ider

import (
	"encoding/binary"
	"encoding/hex"
	"hash/crc32"
	"os"
	"time"
)

type IDGener struct {
	machineCode [3]byte
	pidCode     [2]byte
	second      int64
	idChan      chan string
	inc         uint32
}

func ID() <-chan string {
	ig := &IDGener{
		second: time.Now().Unix(),
		inc:    0,
		idChan: make(chan string),
	}
	ig.machineCode, ig.pidCode = ig.machinePidEncode()

	go ig.generUniqueId()

	return ig.idChan
}

func (ig *IDGener) machinePidEncode() ([3]byte, [2]byte) {
	hostname, _ := os.Hostname()

	buf := [4]byte{}
	binary.BigEndian.PutUint32(buf[:], crc32.ChecksumIEEE([]byte(hostname)))
	machineCode := [3]byte{buf[1], buf[2], buf[3]}

	pid := os.Getpid()
	pidCode := [2]byte{}
	binary.BigEndian.PutUint16(pidCode[:], uint16(pid))

	return machineCode, pidCode
}

func (ig *IDGener) timeIncEncode() ([4]byte, [3]byte) {
	now := time.Now().Unix()
	timeEncode := [4]byte{}
	incEncode := [3]byte{}
	buf := [8]byte{}

	if ig.second != now {
		ig.second = now
		ig.inc = 0
	}

	binary.BigEndian.PutUint64(buf[:], uint64(ig.second))
	timeEncode = [4]byte{buf[4], buf[5], buf[6], buf[7]}

	binary.BigEndian.PutUint32(buf[:], uint32(ig.inc))
	incEncode = [3]byte{buf[1], buf[2], buf[3]}
	ig.inc++

	return timeEncode, incEncode
}

func (ig *IDGener) generUniqueId() {
	var uniqueId [12]byte
	for {
		timeEncode, incEncode := ig.timeIncEncode()
		var i int = 0
		for _, b := range timeEncode {
			uniqueId[i] = b
			i++
		}

		for _, b := range ig.machineCode {
			uniqueId[i] = b
			i++
		}

		for _, b := range ig.pidCode {
			uniqueId[i] = b
			i++
		}

		for _, b := range incEncode {
			uniqueId[i] = b
			i++
		}

		ig.idChan <- hex.EncodeToString(uniqueId[:])

	}
}
