package lungo

import (
	"fmt"
	"testing"
)

func TestAssertNil(t *testing.T) {
	assertNil(t, nil)

	i := (*any)(nil)
	assertNil(t, i)

	var s []int
	assertNil(t, s)

	var m map[int]int
	assertNil(t, m)

	var f func()
	assertNil(t, f)
}

func TestAssertNotNil(t *testing.T) {
	assertNotNil(t, 1)
	assertNotNil(t, "1")

	i := struct{ foo string }{foo: "bar"}
	assertNotNil(t, i)

	s := []int{1}
	assertNotNil(t, s)

	m := map[int]int{1: 2}
	assertNotNil(t, m)

	f := func() { fmt.Println("Foo") }
	assertNotNil(t, f)
}

func TestAssertEqualNil(t *testing.T) {
	assertEqual(t, nil, nil)

	i := (*any)(nil)
	assertEqual(t, i, i)

	var s []int
	assertEqual(t, s, s)

	var m map[int]int
	assertEqual(t, m, m)

	var f func()
	assertEqual(t, f, f)
}

func TestAssertEqual(t *testing.T) {
	assertEqual(t, "a", "a")
	assertEqual(t, 1, 1)
	assertEqual(t, []int{0}, []int{0})
	assertEqual(t, map[int]int{0: 1}, map[int]int{0: 1})
}

func TestAssertPanic(t *testing.T) {
	assertPanic(t, "Foo", func() {
		panic("Foo")
	})
}
