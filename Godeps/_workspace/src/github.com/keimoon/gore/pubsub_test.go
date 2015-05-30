package gore

import (
	"strconv"
	"testing"
)

func TestPubsub(t *testing.T) {
	tp := &testPubsub{t: t}
	tp.setup()
	tp.test()
}

type testPubsub struct {
	t        *testing.T
	pubReady chan bool
	subReady chan bool
	pubEnd   chan bool
	subEnd   chan bool
}

func (tp *testPubsub) setup() {
	tp.pubReady = make(chan bool, 1)
	tp.subReady = make(chan bool, 1)
	tp.pubEnd = make(chan bool, 1)
	tp.subEnd = make(chan bool, 1)
}

func (tp *testPubsub) test() {
	go tp.publisher()
	go tp.subscriber()
	<-tp.pubEnd
	<-tp.subEnd
}

func (tp *testPubsub) publisher() {
	t := tp.t
	defer func() {
		tp.pubEnd <- true
	}()
	conn, err := Dial("localhost:6379")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	tp.pubReady <- true
	<-tp.subReady

	for i := 0; i < 250; i++ {
		err := Publish(conn, "test", i)
		if err != nil {
			t.Fatal(err)
		}
	}
	for i := 250; i < 500; i++ {
		err := Publish(conn, "touhou", i)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func (tp *testPubsub) subscriber() {
	t := tp.t
	defer func() {
		tp.subEnd <- true
	}()
	conn, err := Dial("localhost:6379")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	subs := NewSubscriptions(conn)
	subs.Subscribe("test")
	subs.PSubscribe("tou*")
	tp.subReady <- true
	<-tp.pubReady
	check := 0
	for message := range subs.Message() {
		if message == nil {
			break
		}
		if strconv.FormatInt(int64(check), 10) != string(message.Message) {
			t.Fatal(check, string(message.Message))
		}
		check++
		if check == 500 {
			break
		}
	}
}
