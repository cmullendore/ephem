package main

import (
	"net/http"
	"strings"
	"testing"
	"time"
)

func Test_main(t *testing.T) {
	testGet()
}

func testGet() {
	reader := strings.NewReader("String Test 1")

	go main()

	time.Sleep(1000)

	resp, err := http.Get("http://localhost:8080/a/b/c")

	if err != nil {
		println(err.Error())
	}
	println(resp.Status)

	resp2, err2 := http.Post("http://localhost:8080/a/b/c", "application/json", reader)
	if err2 != nil {
		println(err.Error())
	}
	println(resp2.Status)

	println("Got here")
	//println(resp.Status)

}
