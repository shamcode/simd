package btree

import (
	"fmt"
	"github.com/shamcode/simd/asserts"
	"github.com/shamcode/simd/indexes/compute"
	"github.com/shamcode/simd/storage"
	"sort"
	"testing"
)

func concatIDs(ids []storage.IDIterator) []int {
	var result []int
	for _, idsStorage := range ids {
		idsStorage.Iterate(func(id int64) {
			result = append(result, int(id))
		})
	}
	sort.Ints(result)
	return result
}

func TestBTree(t *testing.T) {
	tree := NewTree(3, true)
	n := 10
	for i := 1; i <= n; i++ {
		idStorage := storage.CreateUniqueIDStorage()
		idStorage.Add(int64(i))
		tree.Set(compute.IntKey(i), idStorage)
	}

	t.Run("GetForKey", func(t *testing.T) {
		for i := 1; i <= n; i++ {
			asserts.Equals(t, 1, tree.Get(compute.IntKey(i)).Count(), "check count")
		}
	})

	t.Run("LessThan", func(t *testing.T) {
		testCases := []struct {
			key           compute.IntKey
			expectedCount int
			expectedIDS   []int
		}{
			{key: 0, expectedCount: 0, expectedIDS: nil},
			{key: 1, expectedCount: 0, expectedIDS: nil},
			{key: 2, expectedCount: 1, expectedIDS: []int{1}},
			{key: 3, expectedCount: 2, expectedIDS: []int{1, 2}},
			{key: 10, expectedCount: 9, expectedIDS: []int{1, 2, 3, 4, 5, 6, 7, 8, 9}},
			{key: 11, expectedCount: 10, expectedIDS: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
			{key: 12, expectedCount: 10, expectedIDS: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
		}
		for _, testCase := range testCases {
			count, ids := tree.LessThan(testCase.key)
			asserts.Equals(t, testCase.expectedCount, count, fmt.Sprintf("check count for %d", testCase.key))
			asserts.Equals(t, testCase.expectedIDS, concatIDs(ids), fmt.Sprintf("check _int64 for %d", testCase.key))
		}
	})

	t.Run("LessOrEqual", func(t *testing.T) {
		testCases := []struct {
			key           compute.IntKey
			expectedCount int
			expectedIDS   []int
		}{
			{key: 0, expectedCount: 0},
			{key: 1, expectedCount: 1, expectedIDS: []int{1}},
			{key: 2, expectedCount: 2, expectedIDS: []int{1, 2}},
			{key: 3, expectedCount: 3, expectedIDS: []int{1, 2, 3}},
			{key: 10, expectedCount: 10, expectedIDS: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
			{key: 11, expectedCount: 10, expectedIDS: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
			{key: 12, expectedCount: 10, expectedIDS: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
		}
		for _, testCase := range testCases {
			count, ids := tree.LessOrEqual(testCase.key)
			asserts.Equals(t, testCase.expectedCount, count, fmt.Sprintf("check count for %d", testCase.key))
			asserts.Equals(t, testCase.expectedIDS, concatIDs(ids), fmt.Sprintf("check _int64 for %d", testCase.key))
		}
	})

	t.Run("GreaterThan", func(t *testing.T) {
		testCases := []struct {
			key           compute.IntKey
			expectedCount int
			expectedIDS   []int
		}{
			{key: -1, expectedCount: 10, expectedIDS: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
			{key: 0, expectedCount: 10, expectedIDS: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
			{key: 1, expectedCount: 9, expectedIDS: []int{2, 3, 4, 5, 6, 7, 8, 9, 10}},
			{key: 2, expectedCount: 8, expectedIDS: []int{3, 4, 5, 6, 7, 8, 9, 10}},
			{key: 3, expectedCount: 7, expectedIDS: []int{4, 5, 6, 7, 8, 9, 10}},
			{key: 10, expectedCount: 0, expectedIDS: nil},
			{key: 11, expectedCount: 0, expectedIDS: nil},
		}
		for _, testCase := range testCases {
			count, ids := tree.GreaterThan(testCase.key)
			asserts.Equals(t, testCase.expectedCount, count, fmt.Sprintf("check count for %d", testCase.key))
			asserts.Equals(t, testCase.expectedIDS, concatIDs(ids), fmt.Sprintf("check _int64 for %d", testCase.key))
		}
	})

	t.Run("GreaterOrEqual", func(t *testing.T) {
		testCases := []struct {
			key           compute.IntKey
			expectedCount int
			expectedIDS   []int
		}{
			{key: -1, expectedCount: 10, expectedIDS: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
			{key: 0, expectedCount: 10, expectedIDS: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
			{key: 1, expectedCount: 10, expectedIDS: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
			{key: 2, expectedCount: 9, expectedIDS: []int{2, 3, 4, 5, 6, 7, 8, 9, 10}},
			{key: 3, expectedCount: 8, expectedIDS: []int{3, 4, 5, 6, 7, 8, 9, 10}},
			{key: 10, expectedCount: 1, expectedIDS: []int{10}},
			{key: 11, expectedCount: 0, expectedIDS: nil},
		}
		for _, testCase := range testCases {
			count, ids := tree.GreaterOrEqual(testCase.key)
			asserts.Equals(t, testCase.expectedCount, count, fmt.Sprintf("check count for %d", testCase.key))
			asserts.Equals(t, testCase.expectedIDS, concatIDs(ids), fmt.Sprintf("check _int64 for %d", testCase.key))
		}
	})
}
