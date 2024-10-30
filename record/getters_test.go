package record

import (
	"sort"
	"testing"

	asserts "github.com/shamcode/assert"
)

type user struct {
	id     int64
	bool   bool
	int    int
	enum8  enum8
	enum16 enum16
	int32  int32
	string string
}

func (u user) GetID() int64 { return u.id }

type enum8 uint8

func (e enum8) Value() uint8 { return uint8(e) }

type enum16 uint16

func (e enum16) Value() uint16 { return uint16(e) }

func TestGetters(t *testing.T) {
	fields := NewFields()

	users := []user{
		{
			id:     1,
			bool:   true,
			int:    3,
			enum8:  1,
			enum16: 20,
			int32:  100,
			string: "cccc",
		},
		{
			id:     2,
			bool:   false,
			int:    2,
			enum8:  3,
			enum16: 10,
			int32:  150,
			string: "aaaa",
		},
		{
			id:     3,
			bool:   true,
			int:    1,
			enum8:  2,
			enum16: 15,
			int32:  120,
			string: "bbbb",
		},
	}

	t.Run("less", func(t *testing.T) {
		testCases := []struct {
			getter interface {
				Less(a, b user) bool
				String() string
			}
			expectedOrder []int64
		}{
			{
				getter:        NewIDGetter[user](),
				expectedOrder: []int64{1, 2, 3},
			},
			{
				getter: BoolGetter[user]{
					Field: fields.New("bool"),
					Get:   func(item user) bool { return item.bool },
				},
				expectedOrder: []int64{2, 1, 3},
			},
			{
				getter: ComparableGetter[user, int]{
					Field: fields.New("int"),
					Get:   func(item user) int { return item.int },
				},
				expectedOrder: []int64{3, 2, 1},
			},
			{
				getter: EnumGetter[user, uint8]{
					Field: fields.New("enum8"),
					Get:   func(item user) Enum[uint8] { return item.enum8 },
				},
				expectedOrder: []int64{1, 3, 2},
			},
			{
				getter: EnumGetter[user, uint16]{
					Field: fields.New("enum16"),
					Get:   func(item user) Enum[uint16] { return item.enum16 },
				},
				expectedOrder: []int64{2, 3, 1},
			},
			{
				getter: ComparableGetter[user, int32]{
					Field: fields.New("int32"),
					Get:   func(item user) int32 { return item.int32 },
				},
				expectedOrder: []int64{1, 3, 2},
			},
			{
				getter: ComparableGetter[user, string]{
					Field: fields.New("string"),
					Get:   func(item user) string { return item.string },
				},
				expectedOrder: []int64{2, 3, 1},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.getter.String(), func(t *testing.T) {
				sort.SliceStable(users, func(i, j int) bool {
					return tc.getter.Less(users[i], users[j])
				})
				ids := make([]int64, len(users))
				for i := range users {
					ids[i] = users[i].GetID()
				}
				asserts.Equals(t, tc.expectedOrder, ids, tc.getter.String())
			})
		}
	})
}
