package assert_test

import (
	"testing"

	"github.com/LaBatata101/goinsta/assert"
)

type Test struct {
	field1 string
	field2 int
	field3 bool
	field4 float32
}

func TestSnapshotStruct(t *testing.T) {
	t1 := Test{
		"hello",
		10,
		false,
		3.14,
	}
	assert.Snapshot(t, t1)
}

func TestSnapshotString(t *testing.T) {
	assert.Snapshot(t, "This is a string")
}

func TestSnapshotBigOutput(t *testing.T) {
	assert.Snapshot(t, `
Lorem ipsum dolor sit amet, consectetur adipiscing elit. Ut ac interdum ex. Fusce iaculis ex nunc, ac interdum ex
tempus a. Donec efficitur accumsan cursus. Mauris efficitur sem quis est dictum posuere. Integer faucibus facilisis
finibus. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Phasellus sit amet tortor tincidunt, aliquam dui at,
ultrices sem. Integer eu quam a nunc placerat faucibus vitae vitae lacus. Class aptent taciti sociosqu ad litora
torquent per conubia nostra, per inceptos himenaeos. Vivamus eu venenatis ex.

Donec tellus turpis, sagittis at sapien ac, elementum volutpat arcu. Fusce maximus leo sit amet est dictum pharetra.
Nam in sem erat. Proin bibendum sem dignissim tortor condimentum faucibus. Praesent quis metus sit amet magna euismod
rhoncus. Aliquam porttitor est consequat finibus auctor. Etiam quis lectus a sem congue tempus non eget purus. Sed
tristique tortor in auctor auctor. Praesent vel tellus ut metus dignissim porttitor. Nulla et lobortis justo. Mauris
nibh nisi, mollis id quam vitae, vulputate pretium nibh. Morbi lobortis efficitur purus vitae tincidunt. Nam feugiat
maximus feugiat. Etiam lacinia blandit erat ut pulvinar. Phasellus hendrerit, eros iaculis convallis commodo, tortor
ligula gravida nulla, vitae luctus mauris diam a nulla.

Morbi tristique justo a massa gravida sodales. Integer non eros efficitur, iaculis mi sed, lobortis massa. Donec ut
ullamcorper urna, sed consectetur turpis. Curabitur vitae ante sodales lacus interdum faucibus. Duis purus nibh,
aliquam ut erat in, iaculis porta diam. Phasellus cursus feugiat ultrices. Suspendisse non mauris quis
lectus lobortis sollicitudin nec sed ante. Ut sed velit vehicula, tincidunt eros eu, aliquam dolor. Nulla diam lacus,
feugiat at turpis in, elementum vestibulum eros.`)
}
