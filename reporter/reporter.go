package reporter

import "time"

type Reporter interface {
	AddDayForUser(user string, day time.Time)
	//	AddTimeForUser(user string, day time.Time, time time.Duration)
}

func New() Reporter {
	return &oncallReport{
		m: map[string]map[string]time.Duration{},
	}
}

// oncallReport - Holds the oncall days and missions of a user
type oncallReport struct {
	m map[string]map[string]time.Duration
}

// oncallTime - Describes oncall entries of a type
type oncallTime struct {
	Days    uint
	Minutes time.Duration
}

func (r *oncallReport) AddDayForUser(user string, day time.Time) {
	var userMap map[string]time.Duration
	var ok bool

	if userMap, ok = r.m[user]; !ok {
		userMap = map[string]time.Duration{}
	}

	key := day.Format("2006-01-02")
	userMap[key] = 0

	r.m[user] = userMap
}
