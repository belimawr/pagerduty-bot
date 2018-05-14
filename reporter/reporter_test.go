package reporter

import (
	"testing"
	"time"
)

func Test_oncallReport_addDayForUser(t *testing.T) {
	r := New().(*oncallReport)

	user := "foo"
	day, _ := time.Parse("2006-01-02", "2018-05-14")
	key := day.Format("2006-01-02")

	r.AddDayForUser(user, day)

	if _, ok := r.m[user]; !ok {
		t.Fatalf("user %s not present", user)
	}

	if _, ok := r.m[user][key]; !ok {
		t.Errorf("Day %s must be present for user %q",
			key,
			user)
	}
}
