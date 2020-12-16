package generator

import (
	"testing"
)

func TestNew(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Error(r)
			return
		}
	}()
	t.Logf("%+v", New())
}
