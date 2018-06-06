package storeReporter

import (
	"reflect"
	"testing"
	"time"
)

func Test_addToSet(t *testing.T) {
	set := []string{}
	toAdd := []string{"foo", "bar", "foo", "bar"}
	expected := []string{"foo", "bar"}

	for _, el := range toAdd {
		set = addToSet(set, el)
	}

	if !reflect.DeepEqual(expected, set) {
		t.Errorf("expecting set: %v, got: %v",
			expected, set)
	}
}

func Test_New_not_nil_map(t *testing.T) {
	a := DayTyperFunc(func(_ time.Time) string {
		return "foo"
	})

	inM := New(a).(*inMemory)

	if inM.m == nil {
		t.Errorf("inM.m cannot be nil")
	}

	if inM.typer == nil {
		t.Error("inM.typer cannot be nil")
	}
}

func Test_inMemory_addDayForUser(t *testing.T) {
	expected := "2018-05-14"
	s := &inMemory{
		m: map[string]onCallReport{},
		typer: DayTyperFunc(func(_ time.Time) string {
			return "foo"
		}),
	}

	user := "foo"
	day, _ := time.Parse("2006-01-02", "2018-05-14")
	key := day.Format("2006-01-02")
	kind := s.typer.Type(day)

	s.AddDayForUser(user, day)

	if _, ok := s.m[user]; !ok {
		t.Fatalf("user %s not present", user)
	}

	if len(s.m[user].days[kind]) == 0 {
		t.Fatalf("%q for user %q cannot be empty",
			kind,
			user)
	}

	if s.m[user].days[kind][0] != expected {
		t.Errorf("Day %s must be present for user %q",
			key,
			user)
	}
}

func Test_inMemory_AddTimerForUser(t *testing.T) {
	val := 1 * time.Second
	expectedMission := 3 * val
	user := "foo"
	r := &inMemory{
		m: map[string]onCallReport{},
	}

	for i := 0; i < 3; i++ {

		r.AddTimeForUser(user, time.Now(), val)
	}

	if r.m[user].mission != expectedMission {
		t.Errorf("expecting r.m[%q].mission: %v, got: %v",
			user, r.m[user].mission, expectedMission)
	}
}
