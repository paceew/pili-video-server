package utils

import (
	"fmt"
	"testing"
)

func TestUuid(t *testing.T) {
	t.Run("...", testUuid)

}

func testUuid(t *testing.T) {
	id, _ := NewUUID()
	fmt.Printf("id:%v", id)
}
