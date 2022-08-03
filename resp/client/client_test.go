package client

import (
	"fmt"
	"github.com/NetLops/go-imitate-redis/resp/reply"
	"strconv"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	ticker := time.NewTicker(1 * time.Second)
	//times := time.After(1 * time.Second)
	for range ticker.C {
		//for range times {
		fmt.Println("test")
	}
}

func TestClient(t *testing.T) {
	client, err := MakeClient("localhost:6379")
	if err != nil {
		t.Error(err)
	}
	client.Start()
	//result := client.Send([][]byte{
	//	[]byte("PING"),
	//})
	//if statusRet, ok := result.(*reply.PongReply); ok {
	//	if string(statusRet.ToBytes()) != "PONG" {
	//		t.Error("`ping` failed, result: " + string(statusRet.ToBytes()))
	//	}
	//}
	//fmt.Println(result)
	//result := client.Send([][]byte{
	//	[]byte("SET"),
	//	[]byte("a"),
	//	[]byte("a"),
	//})
	//fmt.Println(result)
	//if statusRet, ok := result.(*reply.StatusReply); ok {
	//	if statusRet.Status != "OK" {
	//		t.Error("`set` failed, result: " + statusRet.Status)
	//	}
	//}
	result := client.Send([][]byte{
		//[]byte("ping"),
		[]byte("GET"),
		[]byte("t3"),
	})

	fmt.Println(result)
	if bulkRet, ok := result.(*reply.BulkReply); ok {
		fmt.Println(string(bulkRet.ToBytes()))
		if string(bulkRet.Arg) != "34" {
			t.Error("`get` failed, result: " + string(bulkRet.Arg))
		}
	}

}

func TestClient_Send(t *testing.T) {
	client, err := MakeClient("localhost:6379")
	if err != nil {
		fmt.Println(err)
	}
	client.Start()
	for i := 0; i < 99999; i++ {

		temp := i
		go func() {

			//fmt.Println(temp)

			result := client.Send([][]byte{
				[]byte("PING"),
			})
			if statusRet, ok := result.(*reply.PongReply); ok {
				if string(statusRet.ToBytes()) != "PONG" {
					fmt.Println("`ping` failed, result: " + string(statusRet.ToBytes()))
				}
			}
			//fmt.Println(result)
			//result = client.Send([][]byte{
			//	[]byte("SET"),
			//	[]byte("a" + strconv.Itoa(temp)),
			//	[]byte(strconv.Itoa(temp)),
			//})
			////fmt.Println(result)
			//if statusRet, ok := result.(*reply.StatusReply); ok {
			//	if statusRet.Status != "OK" {
			//		fmt.Println("`set` failed, result: " + statusRet.Status)
			//	}
			//}

			result = client.Send([][]byte{
				[]byte("GET"),
				[]byte("a" + strconv.Itoa(temp)),
			})

			//fmt.Println(result)
			if bulkRet, ok := result.(*reply.BulkReply); ok {
				fmt.Println(string(bulkRet.Arg), "a"+strconv.Itoa(temp))
				if string(bulkRet.Arg) != strconv.Itoa(temp) {
					fmt.Println("`get` failed, result: "+string(bulkRet.Arg), "a"+strconv.Itoa(temp))
				}
			}
		}()
	}
	time.Sleep(10 * time.Second)
}

func BenchmarkName(b *testing.B) {
	for i := 0; i < b.N; i++ {

	}
}
