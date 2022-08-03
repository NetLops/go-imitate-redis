package main

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {

	b := []byte("*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n")
	fmt.Println(b[len(b)-2])
	fmt.Println('\r')
}
