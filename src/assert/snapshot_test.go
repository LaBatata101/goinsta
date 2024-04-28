package assert_test

import (
	"goinsta/assert"
	"testing"
)

type Test struct {
	field1 string
	field2 int
}

func TestHello(t *testing.T) {
	t1 := Test{
		"hello",
		10,
	}
	assert.Snapshot(t, t1)
}

func TestHello2(t *testing.T) {
	t1 := Test{
		"world",
		10,
	}
	assert.Snapshot(t, t1)
}
