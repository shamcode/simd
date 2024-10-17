package btree

import (
	"sort"
	"testing"

	asserts "github.com/shamcode/assert"
	"github.com/shamcode/simd/indexes"
	"github.com/shamcode/simd/record"
	"github.com/shamcode/simd/where"
	"github.com/shamcode/simd/where/comparators"
)

type _int64 []int64

func (s _int64) Len() int           { return len(s) }
func (s _int64) Less(i, j int) bool { return s[i] < s[j] }
func (s _int64) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func TestIndex(t *testing.T) {
	index := NewInt64BTreeIndex(record.ID, 8, true)
	var id int64
	for id = 1; id <= 10; id++ {
		key := index.Compute().ForValue(id)
		index.ConcurrentStorage().GetOrCreate(key).Add(id)
	}

	t.Run("weight", func(t *testing.T) {
		testCases := []struct {
			condition        where.Condition
			expectedCanApply bool
			expectedWeight   indexes.IndexWeight
		}{
			{
				condition: where.Condition{
					Cmp: comparators.Int64FieldComparator{
						Cmp:    where.LT,
						Getter: record.ID,
						Value:  []int64{5},
					},
				},
				expectedCanApply: true,
				expectedWeight:   indexes.IndexWeightLow,
			},
			{
				condition: where.Condition{
					Cmp: comparators.Int64FieldComparator{
						Cmp:    where.EQ,
						Getter: record.ID,
						Value:  []int64{1},
					},
				},
				expectedCanApply: true,
				expectedWeight:   indexes.IndexWeightMedium,
			},
			{
				condition: where.Condition{
					Cmp: comparators.Int64FieldComparator{
						Cmp:    where.InArray,
						Getter: record.ID,
						Value:  []int64{1, 5},
					},
				},
				expectedCanApply: true,
				expectedWeight:   indexes.IndexWeightMedium,
			},
			{
				condition: where.Condition{
					WithNot: true,
					Cmp: comparators.Int64FieldComparator{
						Cmp:    where.EQ,
						Getter: record.ID,
						Value:  []int64{1},
					},
				},
				expectedCanApply: true,
				expectedWeight:   indexes.IndexWeightHigh,
			},
		}

		for _, test := range testCases {
			test := test
			t.Run(test.condition.String(), func(t *testing.T) {
				t.Parallel()
				canApply, weight := index.Weight(test.condition)
				asserts.Equals(t, test.expectedCanApply, canApply, "can apply")
				asserts.Equals(t, test.expectedWeight, weight, "weight")
			})
		}
	})

	t.Run("select", func(t *testing.T) {
		testCases := []struct {
			condition     where.Condition
			expectedCount int
			expectedIDs   []int64
		}{
			{
				condition: where.Condition{
					Cmp: comparators.Int64FieldComparator{
						Cmp:    where.LT,
						Getter: record.ID,
						Value:  []int64{5},
					},
				},
				expectedCount: 4,
				expectedIDs:   []int64{1, 2, 3, 4},
			},
			{
				condition: where.Condition{
					Cmp: comparators.Int64FieldComparator{
						Cmp:    where.LE,
						Getter: record.ID,
						Value:  []int64{5},
					},
				},
				expectedCount: 5,
				expectedIDs:   []int64{1, 2, 3, 4, 5},
			},
			{
				condition: where.Condition{
					Cmp: comparators.Int64FieldComparator{
						Cmp:    where.EQ,
						Getter: record.ID,
						Value:  []int64{1},
					},
				},
				expectedCount: 1,
				expectedIDs:   []int64{1},
			},
			{
				condition: where.Condition{
					Cmp: comparators.Int64FieldComparator{
						Cmp:    where.GT,
						Getter: record.ID,
						Value:  []int64{5},
					},
				},
				expectedCount: 5,
				expectedIDs:   []int64{6, 7, 8, 9, 10},
			},
			{
				condition: where.Condition{
					Cmp: comparators.Int64FieldComparator{
						Cmp:    where.GE,
						Getter: record.ID,
						Value:  []int64{5},
					},
				},
				expectedCount: 6,
				expectedIDs:   []int64{5, 6, 7, 8, 9, 10},
			},
			{
				condition: where.Condition{
					Cmp: comparators.Int64FieldComparator{
						Cmp:    where.InArray,
						Getter: record.ID,
						Value:  []int64{1, 5},
					},
				},
				expectedCount: 2,
				expectedIDs:   []int64{1, 5},
			},
			{
				condition: where.Condition{
					WithNot: true,
					Cmp: comparators.Int64FieldComparator{
						Cmp:    where.EQ,
						Getter: record.ID,
						Value:  []int64{5},
					},
				},
				expectedCount: 9,
				expectedIDs:   []int64{1, 2, 3, 4, 6, 7, 8, 9, 10},
			},
		}

		for _, test := range testCases {
			test := test
			t.Run(test.condition.String(), func(t *testing.T) {
				t.Parallel()
				count, idsStorage, err := index.Select(test.condition)
				asserts.Success(t, err)
				var ids []int64
				for _, store := range idsStorage {
					store.Iterate(func(id int64) {
						ids = append(ids, id)
					})
				}
				sort.Sort(_int64(ids))
				asserts.Equals(t, test.expectedIDs, ids, "ids")
				asserts.Equals(t, test.expectedCount, count, "count")
			})
		}
	})
}
