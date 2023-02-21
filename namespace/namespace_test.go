package namespace

import (
	"github.com/shamcode/assert"
	"testing"
)

type user struct {
	id              int64
	computedCounter int
}

func (u *user) GetID() int64 {
	return u.id
}

func (u *user) ComputeFields() {
	u.computedCounter++
}

func TestComputeFields(t *testing.T) {
	item := &user{id: 1}
	store := CreateNamespace()

	asserts.Success(t, store.Insert(item))
	asserts.Equals(t, 1, item.computedCounter, "compute on insert")

	asserts.Success(t, store.Upsert(item))
	asserts.Equals(t, 2, item.computedCounter, "compute on upsert (update)")

	newItem := &user{id: 2}
	asserts.Success(t, store.Upsert(newItem))
	asserts.Equals(t, 1, newItem.computedCounter, "compute on upsert (insert)")

	updatedItem := &user{id: 1}
	asserts.Success(t, store.Upsert(updatedItem))
	asserts.Equals(t, 1, updatedItem.computedCounter, "compute on upsert (update)")
}
