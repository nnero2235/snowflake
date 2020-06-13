package snowflake

import (
	"fmt"
	"sync"
	"testing"
)

func TestNewSnowflake(t *testing.T) {
	snow := NewSnowflake(0,0)
	id := snow.NextId()
	t.Log(len(fmt.Sprintf("%b",id)))
	s := snow.ParseId(id)
	t.Log(s)
}

func TestNewSnowflakeMulti(t *testing.T){
	snow := NewSnowflake(29,15)
	for i:=0;i<4097;i++{
		id := snow.NextId()
		t.Log(id)
		s := snow.ParseId(id)
		t.Log(s)
	}
}

func TestNewSnowflakeMultiThread(t *testing.T) {
	snow := NewSnowflake(29,15)
	var wg sync.WaitGroup
	for i:=0;i<4097;i++{
		wg.Add(1)
		go func() {
			id := snow.NextId()
			t.Log(id)
			s := snow.ParseId(id)
			t.Log(s)
			wg.Done()
		}()
	}
	wg.Wait()
}