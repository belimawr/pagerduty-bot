package storeReporter

import (
	"io"
	"time"
)

// Reporter - Generates the report
type Reporter interface {
	Report(io.Writer)
}

// Store - Interface to store the data
type Store interface {
	AddDayForUser(user string, day time.Time)
	AddTimeForUser(user string, day time.Time, time time.Duration)
}

// StoreReporter - a Store that also generates Reports
type StoreReporter interface {
	Store
	Reporter
}

// DayTyper - indentifies the type of the day
type DayTyper interface {
	Type(time.Time) string
}

// DayTyperFunc - functions can be of type DayTyper
type DayTyperFunc func(time.Time) string

// Type - wraps a function to implement DayTyper
func (f DayTyperFunc) Type(t time.Time) string {
	return f(t)
}
