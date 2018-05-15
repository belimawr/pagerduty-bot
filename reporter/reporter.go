package reporter

import (
	"fmt"
	"io"
	"time"
)

type Reporter interface {
	AddDayForUser(user string, day time.Time)
	//	AddTimeForUser(user string, day time.Time, time time.Duration)
	Report(io.Writer)
}

type DayTyper interface {
	Type(time.Time) string
}

type DayTyperFunc func(time.Time) string

func (f DayTyperFunc) Type(t time.Time) string {
	return f(t)
}

func New() Reporter {
	return &inMemory{
		m:     map[string]onCallReport{},
		typer: DayTyperFunc(dayTyperFn),
	}
}

func dayTyperFn(t time.Time) string {
	if t.Weekday() == time.Saturday ||
		t.Weekday() == time.Sunday {
		return "non-business-day"
	}
	return "business-day"
}

// inMemory - Holds the oncall days and missions of a user
type inMemory struct {
	m     map[string]onCallReport
	typer DayTyper
}

type onCallReport struct {
	days    map[string][]string
	mission time.Duration
}

func (r *inMemory) AddDayForUser(user string, day time.Time) {
	var userData onCallReport
	var ok bool

	key := day.Format("2006-01-02")
	kind := r.typer.Type(day)

	if userData, ok = r.m[user]; !ok {
		userData = onCallReport{
			days: map[string][]string{},
		}
		r.m[user] = userData
	}

	r.m[user].days[kind] = addToSet(r.m[user].days[kind], key)
}

func (r inMemory) Report(w io.Writer) {
	for user := range r.m {
		fmt.Fprintf(w, "%s:", user)
		for kind := range r.m[user].days {
			fmt.Fprintf(w, "\n\t%s: ", kind)
			for _, day := range r.m[user].days[kind] {
				fmt.Fprintf(w, "%s, ", day)
			}
		}
		fmt.Fprintf(w, "\n")
	}
}

func addToSet(set []string, el string) []string {
	for _, e := range set {
		if e == el {
			return set
		}
	}

	return append(set, el)
}
