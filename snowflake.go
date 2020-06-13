package snowflake

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

const (
	initEpoch        int64 = 1420041600000
	workIdBits             = 5
	dataCenterIdBits       = 5
	sequenceIdBits         = 12
	maxDataCenterId        = 31 //~(-1 << dataCenterIdBits)
	maxWorkId              = 31 //~(-1 << workIdBits)

	timestampBitsMove  = dataCenterIdBits + workIdBits + sequenceIdBits
	dataCenterBitsMove = workIdBits + sequenceIdBits
	workIdBitsMove     = sequenceIdBits
	sequenceMask       = 4095 //~(-1 << sequenceBits)
)

type Snowflake struct {
	lastTimestamp int64
	dataCenterId  int64
	workId        int64
	sequence      int64
	mutex         sync.Mutex
}

//can parse id into the part info built with string.
func (snow *Snowflake) ParseId(id int64) string {
	seq := id & 0xfff
	workId := (id & 0x1f000) >> workIdBitsMove
	dataCenterId := (id & 0x3e0000) >> dataCenterBitsMove
	timestamp := (id >> timestampBitsMove) + initEpoch
	dateTime := time.Unix(timestamp/1000, 0).Format("2006-01-02 15:04:05.999")
	return fmt.Sprintf("DateTime: %s -> dataCenter: %d -> worker: %d -> sequence: %d",
		dateTime, dataCenterId, workId, seq)
}

func (snow *Snowflake) NextId() int64 {
	//must be locked to ensure the unique id
	snow.mutex.Lock()
	defer snow.mutex.Unlock()

	now := time.Now().Unix() * 1000
	if now == snow.lastTimestamp { //need use sequence incr
		snow.sequence = (snow.sequence + 1) & sequenceMask
		if snow.sequence == 0 { //overflow wait to next timestamp
			now = waitNextTimestamp(now)
		}
	} else if now < snow.lastTimestamp { //clock is back, error
		panic("clock is backed. Can't generate id.")
	} else {
		snow.sequence = 0
	}

	snow.lastTimestamp = now

	return (now-initEpoch)<<timestampBitsMove | snow.dataCenterId<<dataCenterBitsMove | snow.workId<<workIdBitsMove | snow.sequence
}

func waitNextTimestamp(now int64) int64 {
	next := time.Now().Unix() * 1000
	for now >= next {
		next = time.Now().Unix() * 1000
	}
	return next
}

func NewSnowflake(dataCenterId, workId int64) *Snowflake {
	if dataCenterId > maxDataCenterId || dataCenterId < 0 {
		panic("dataCenterId's range in 0 ~ " + strconv.Itoa(maxDataCenterId))
	}
	if workId > maxWorkId || workId < 0 {
		panic("workId's range in 0 ~ " + strconv.Itoa(maxWorkId))
	}
	return &Snowflake{
		lastTimestamp: -1,
		dataCenterId:  dataCenterId,
		workId:        workId,
		sequence:      0,
	}
}
