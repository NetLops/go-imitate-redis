package wildcard

import (
	"fmt"
	"testing"
)

func TestWildcard(t *testing.T) {
	pattern := CompilePattern("[a-z]*")
	match := pattern.IsMatch("test1")
	fmt.Println(match)
}
