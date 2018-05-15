package reporter

import (
	"testing"
	"time"
)

func Test_oncallReport_addDayForUser(t *testing.T) {
	expected := "2018-05-14"
	r := New().(*inMemory)

	user := "foo"
	day, _ := time.Parse("2006-01-02", "2018-05-14")
	key := day.Format("2006-01-02")
	kind := r.typer.Type(day)

	r.AddDayForUser(user, day)

	if _, ok := r.m[user]; !ok {
		t.Fatalf("user %s not present", user)
	}

	if len(r.m[user].days[kind]) == 0 {
		t.Fatalf("%q for user %q cannot be empty",
			kind,
			user)
	}

	if r.m[user].days[kind][0] != expected {
		t.Errorf("Day %s must be present for user %q",
			key,
			user)
	}
}
