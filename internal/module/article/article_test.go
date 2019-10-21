package article

import (
	"dyc/internal/db"
	"testing"
)

func TestNewIndex(t *testing.T) {
	db.NewElasticsearchDefault()
	total, _, err := NewIndex(1)
	if err != nil {
		t.Fatal(err)
	}

	if total <= 0 {
		t.Errorf("exepted: > 0\ngot: %d", total)
	}
}
