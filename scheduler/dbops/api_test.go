package dbops

import (
	"fmt"
	"testing"
)

func TestMain(m *testing.M) {
	m.Run()
}

func TestWriteLike(t *testing.T) {
	t.Run("", testWriteLike)
}

type test struct {
	a string
	b string
	c int
}

func testWriteLike(t *testing.T) {
	t1 := &test{a: "testing", b: "1"}
	t2 := &test{a: "testing2", b: "2"}

	var tclie []*test

	tclie = append(tclie, t1)
	tclie = append(tclie, t2)

	for _, row := range tclie {
		likestr := "like_" + row.b
		row.a = likestr
		row.c = 5
	}

	for _, row := range tclie {
		fmt.Printf("a:%v,b:%v,c:%v\n", row.a, row.b, row.c)
	}
}
