package config

import (
	"fmt"
	"strings"
	"testing"
)

func TestName(t *testing.T) {
	str := "   index test test"
	any := strings.IndexAny(str, "test")
	fmt.Println(strings.IndexAny("chicken", "aeiouy"))
	fmt.Println(strings.IndexAny("crwth", "aeiouy"))
	fmt.Println(any)

}
