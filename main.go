package main

import (
	"os"
	"regexp"
	"time"

	pagerduty "github.com/PagerDuty/go-pagerduty"
	"github.com/belimawr/pagerduty-bot/config"
	reporter "github.com/belimawr/pagerduty-bot/storeReporter"
	"github.com/caarlos0/env"
)

var r = regexp.MustCompile(`mission\.(?P<duration>(\d+(m|h)){1,3})$`)

func main() {
	cfg := config.Config{}

	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}

	logger := cfg.Logger()

	client := pagerduty.NewClient(cfg.AccessToken)

	usersMap := map[string]string{}

	storage := reporter.New(reporter.DayTyperFunc(dayTyperFn))

	//========================== List Users
	users, err := client.ListUsers(pagerduty.ListUsersOptions{
		APIListObject: pagerduty.APIListObject{
			Limit: 100,
		},
	})

	if err != nil {
		logger.Fatal().Msg(err.Error())
	}

	for _, u := range users.Users {
		usersMap[u.ID] = u.Name
	}

	//========================= List Oncall
	onCall, err := client.ListOnCalls(pagerduty.ListOnCallOptions{
		APIListObject: pagerduty.APIListObject{
			Limit: 100,
		},
		TimeZone:            cfg.TimeZone,
		EscalationPolicyIDs: cfg.EscalationPlocies,
		Since:               "2018-05-01",
		Until:               "2018-05-10",
	})

	for _, o := range onCall.OnCalls {
		if o.EscalationLevel == 2 {
			oncall, err := wasOnCall(o)
			if err != nil {
				logger.Error().Err(err).Msg("could not verify oncall")
			}

			if oncall {
				day, err := time.Parse(time.RFC3339, o.End)
				if err != nil {
					logger.Error().Err(err).Msg("parsing oncall day")
					continue
				}

				storage.AddDayForUser(o.User.Summary, day)
				logger.Debug().Msgf("Oncall: %s - %s - %s, %s",
					o.Start, o.End, o.User.Summary, o.User.ID)
			}
		}
	}

	//========================== List log entries
	apiOpt := pagerduty.APIListObject{
		Limit:  100,
		Offset: 0,
	}

	for {
		apiOpt.Offset = apiOpt.Offset + apiOpt.Limit
		logEntries, err := client.ListLogEntries(pagerduty.ListLogEntriesOptions{
			APIListObject: apiOpt,
			Since:         "2018-05-01",
			Until:         "2018-05-10",
			TimeZone:      cfg.TimeZone,
		})
		if err != nil {
			logger.Panic().Msgf("listing log entries: %+v", err)
		}
		logger.Debug().Msgf("LogEntries offset: %d, limit: %d, more: %t",
			logEntries.Offset,
			logEntries.Limit,
			logEntries.More)

		for _, l := range logEntries.LogEntries {
			if l.Type == "annotate_log_entry" &&
				l.Channel.Type == "note" &&
				l.Agent.Type == "user_reference" {
				user := usersMap[l.Agent.ID]

				lst := r.FindStringSubmatch(l.Channel.Summary)
				var duration time.Duration
				var err error
				if len(lst) >= 2 {
					token := lst[1]
					duration, err = time.ParseDuration(token)
					if err == nil {
						storage.AddTimeForUser(user, time.Now(), duration)
					} else {
						logger.Error().Msgf("error parsing: %q, err: %s", token, err)
					}
					logger.
						Debug().
						Msgf("Agent-type: %s, agent-ID: %s, createdAt: %s, "+
							"summary: %#v",
							l.Agent.Type,
							user,
							l.CreatedAt,
							l.Channel.Summary)
				} else {
					logger.Warn().Msgf("could not parse: %q", l.Channel.Summary)
				}
			}
		}

		if !logEntries.More {
			break
		}
	}

	storage.Report(os.Stdout)
}

func dayTyperFn(t time.Time) string {
	if t.Weekday() == time.Saturday ||
		t.Weekday() == time.Sunday {
		return "non-business-day"
	}
	return "business-day"
}
