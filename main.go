package main

import (
	"bufio"
	"fmt"
	"github.com/quorzz/readygo/protocol"
	"net"
	"time"
)

// 初步的简单demo
func main() {

	conn := getConn()

	// fmt.Println("--- test normalizeArgs ---")
	// testNormalizeArgs()

	fmt.Println("--- test lpush ---")
	testLpush(conn)

	fmt.Println("--- test lrange ---")
	testLrange(conn)

	fmt.Println("--- test hmset ---")
	testHmset(conn)

	fmt.Println("--- test hgetall ---")
	testHgetAll(conn)
}

func testNormalizeArgs() {
	// request := protocol.Request{
	//  Command: "SET",
	//  Args: []interface{}{
	//      "a",
	//      []string{"k", "v", "po"},
	//      3,
	//      map[string]int{
	//          "u":  56,
	//          "tr": 9,
	//      },
	//      true,
	//      map[interface{}]interface{}{
	//          234: map[string]int{
	//              "u":  56,
	//              "tr": 9,
	//          },
	//          "aax": []string{"k", "v", "po"},
	//      },
	//  },
	// }

	r1 := map[string]int{
		"u":  56,
		"tr": 9,
	}

	fmt.Println(r1)
	res1 := protocol.NormalizeArgs(r1)
	fmt.Println(res1)
	fmt.Println("-----------------------")

	buf, _ := protocol.Pack("HMSET", res1...)
	fmt.Println(string(buf))
	fmt.Println("-----------------------")

	r2 := []interface{}{
		56, "tr", 9, "vb",
	}

	fmt.Println(r2)
	res2 := protocol.NormalizeArgs(r2)
	fmt.Println(res2)

	r3 := "977"
	res3 := protocol.NormalizeArgs(r3)
	fmt.Println(res3)
}

func testHmset(conn net.Conn) {

	r1 := map[string]string{
		"a":    "a",
		"b":    "b",
		"hao好": "好的",
	}

	res1 := protocol.NormalizeArgs("hash_test", r1)

	buf, _ := protocol.Pack("hmset", res1...)

	_, err := conn.Write(buf)

	if err != nil {
		fmt.Println("write error : ", err)
	}

	r := bufio.NewReader(conn)
	ret, err := protocol.ReadResponse(r)

	if err != nil {
		fmt.Println("response errors ", err)
	} else {
		fmt.Println(ret.Type)
		getData(ret)
	}
}

func testLpush(conn net.Conn) {

	args := protocol.NormalizeArgs(
		"list_test",
		[]interface{}{
			"韩梅梅",
			2,
			"list test 2",
			"list test 3",
		},
	)

	b, _ := protocol.Pack("rpush", args...)

	conn.Write(b)

	reader := bufio.NewReader(conn)

	ret, err := protocol.ReadResponse(reader)

	if err != nil {
		fmt.Println("response errors : ", err)
	} else {
		fmt.Println("response type : ", ret.Type)
		// fmt.Println("response data :", ret.Multi[1].Type)
		getData(ret)
	}
}

func testLrange(conn net.Conn) {

	args := protocol.NormalizeArgs(
		"list_test",
		0,
		-1,
	)

	b, _ := protocol.Pack("lrange", args...)

	conn.Write(b)

	reader := bufio.NewReader(conn)

	ret, err := protocol.ReadResponse(reader)

	if err != nil {
		fmt.Println("response errors : ", err)
	} else {
		fmt.Println("response type : ", ret.Type)
		getData(ret)
	}
}

func getData(r *protocol.Response) {
	fmt.Println(r)
	switch r.Type {
	case protocol.ResponseMutli:
		for _, rep := range r.Multi {
			fmt.Println("--- multi ---", string(rep.Bulk))
			getData(rep)
		}
	case protocol.ResponseStatus:
		fmt.Println("--- status response ---", r.Status)
	case protocol.ResponseBulk:
		fmt.Println("---bulk reaponse ---", string(r.Bulk))
	case protocol.ResponseInt:
		fmt.Println("--- number response ---", r.Integer)
	default:
		fmt.Println("not match")
	}
}

func getConn() net.Conn {
	conn, e := net.DialTimeout("tcp", "127.0.0.1:6379", time.Duration(10)*time.Second)

	if e != nil {
		fmt.Println("------- conn error -----", e)
	}
	return conn
}

func testHgetAll(conn net.Conn) {

	args := protocol.NormalizeArgs(
		"hash_test",
	)

	b, err := protocol.Pack("hgetall", args...)

	if err != nil {
		fmt.Println("write err : ", err)
	}
	conn.Write(b)

	reader := bufio.NewReader(conn)
	ret, err := protocol.ReadResponse(reader)

	if err != nil {
		fmt.Println("response errors : ", err)
	} else {
		fmt.Println("response type : ", ret.Type)

		m, err := ret.ToMap()
		if nil != err {
			fmt.Println("to map error", err)
		}
		fmt.Println(m)
	}

}
