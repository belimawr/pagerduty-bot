package main

import (
	"time"

	pagerduty "github.com/PagerDuty/go-pagerduty"
	"github.com/pkg/errors"
)

func wasOnCall(o pagerduty.OnCall) (bool, error) {
	end, err := time.Parse("2006-01-02T15:04:05-07:00", o.End)
	if err != nil {
		return false, errors.WithStack(err)
	}

	return end.Hour() == 23, nil
}
