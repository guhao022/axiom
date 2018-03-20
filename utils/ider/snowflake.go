package ider

import (
	"sync"
	"time"
)

type SnowFlake interface {
	Next() int64
}

const (
	CWorkerIdBits  = 10 // Num of WorkerId Bits
	CSenquenceBits = 12 // Num of Sequence Bits

	CWorkerIdShift  = 12
	CTimeStampShift = 22

	CSequenceMask = 0xfff // equal as getSequenceMask()
	CMaxWorker    = 0x3ff // equal as getMaxWorkerId()
)

var CEpoch = int64(time.Millisecond)

type IdGen struct {
	workerId      int64
	lastTimeStamp int64
	sequence      int64
	maxWorkerId   int64
	lock          *sync.Mutex
}

func NewID(workerid int64) SnowFlake {
	ig := new(IdGen)

	ig.maxWorkerId = getMaxWorkerId()

	if workerid > ig.maxWorkerId || workerid < 0 {
		panic("idgenerator worker not fit")
		return nil
	}
	ig.workerId = workerid
	ig.lastTimeStamp = -1
	ig.sequence = 0
	ig.lock = new(sync.Mutex)
	return ig
}

func getMaxWorkerId() int64 {
	return -1 ^ -1<<CWorkerIdBits
}

func getSequenceMask() int64 {
	return -1 ^ -1<<CSenquenceBits
}

// return in ms
func (iw *IdGen) timeGen() int64 {
	return time.Now().UnixNano() / 1000 / 1000
}

func (iw *IdGen) timeReGen(last int64) int64 {
	ts := time.Now().UnixNano()
	for {
		if ts < last {
			ts = iw.timeGen()
		} else {
			break
		}
	}
	return ts
}

func (ig *IdGen) Next() int64 {
	ig.lock.Lock()
	defer ig.lock.Unlock()
	ts := ig.timeGen()
	if ts == ig.lastTimeStamp {
		ig.sequence = (ig.sequence + 1) & CSequenceMask
		if ig.sequence == 0 {
			ts = ig.timeReGen(ts)
		}
	} else {
		ig.sequence = 0
	}

	if ts < ig.lastTimeStamp {
		panic("Clock moved backwards, Refuse gen id")
		return 0
	}
	ig.lastTimeStamp = ts
	ts = (ts-CEpoch)<<CTimeStampShift | ig.workerId<<CWorkerIdShift | ig.sequence
	return ts
}

// ParseId Func: reverse uid to timestamp, workid, seq
func ParseId(id int64) (t time.Time, ts int64, workerId int64, seq int64) {
	seq = id & CSequenceMask
	workerId = (id >> CWorkerIdShift) & CMaxWorker
	ts = (id >> CTimeStampShift) + CEpoch
	t = time.Unix(ts/1000, (ts%1000)*1000000)
	return
}
