package tests

import "github.com/shamcode/simd/record"

type StatusEnum uint8

const (
	StatusActive StatusEnum = iota + 1
	StatusDisabled
)

func (s StatusEnum) String() string {
	switch s {
	case StatusActive:
		return "ACTIVE"
	case StatusDisabled:
		return "DISABLED"
	default:
		return ""
	}
}

type Tag uint16

const (
	TagTester Tag = iota + 1
	TagConfirmed
	TagFree
)

func (t Tag) String() string {
	switch t {
	case TagTester:
		return "tester"
	case TagConfirmed:
		return "confirmed"
	case TagFree:
		return "free"
	default:
		return ""
	}
}

type Tags map[Tag]struct{}

func (t Tags) Has(item Tag) bool {
	_, ok := t[item]
	return ok
}

type CounterKey uint16

const (
	CounterKeyUnreadMessages CounterKey = iota + 1
	CounterKeyPendingTasks
)

type Counters map[CounterKey]uint32

func (c Counters) HasKey(key CounterKey) bool {
	_, ok := c[key]
	return ok
}
func (c Counters) HasValue(check record.MapValueComparator[uint32]) (bool, error) {
	for _, item := range c {
		res, err := check.Compare(item)
		if nil != err {
			return false, err
		}

		if res {
			return true, nil
		}
	}

	return false, nil
}

type HasCounterValueEqual uint32

func (c HasCounterValueEqual) Compare(item uint32) (bool, error) {
	return item == uint32(c), nil
}

type User struct {
	ID       int64
	Name     string
	Status   StatusEnum
	Score    int
	IsOnline bool
	Tags     Tags
	Counters Counters
}

func (u *User) GetID() int64 { return u.ID }

var userFields = record.NewFields()

var userID = record.NewIDGetter[*User]()

var userName = record.ComparableGetter[*User, string]{
	Field: userFields.New("name"),
	Get:   func(item *User) string { return item.Name },
}

var userStatus = record.ComparableGetter[*User, StatusEnum]{
	Field: userFields.New("status"),
	Get:   func(item *User) StatusEnum { return item.Status },
}

var userScore = record.ComparableGetter[*User, int]{
	Field: userFields.New("score"),
	Get:   func(item *User) int { return item.Score },
}

var userIsOnline = record.BoolGetter[*User]{
	Field: userFields.New("is_online"),
	Get:   func(item *User) bool { return item.IsOnline },
}

var userTags = record.SetGetter[*User, Tag]{
	Field: userFields.New("tags"),
	Get:   func(item *User) record.Set[Tag] { return item.Tags },
}

var userCounters = record.MapGetter[*User, CounterKey, uint32]{
	Field: userFields.New("counters"),
	Get:   func(item *User) record.Map[CounterKey, uint32] { return item.Counters },
}

type byOnline struct {
	onlineToUp bool
}

func (sorting byOnline) Calc(item *User) int64 {
	if sorting.onlineToUp == item.IsOnline {
		return 0
	}

	return 1
}
