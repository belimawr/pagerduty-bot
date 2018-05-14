package main

import (
	pagerduty "github.com/PagerDuty/go-pagerduty"
	"github.com/belimawr/pagerduty-bot/config"
	"github.com/caarlos0/env"
)

func main() {
	cfg := config.Config{}

	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}

	logger := cfg.Logger()

	client := pagerduty.NewClient(cfg.AccessToken)

	usersMap := map[string]string{}

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
				logger.Debug().Msgf("Oncall: %s - %s - %s",
					o.Start, o.End, o.User.Summary)
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
				logger.
					Debug().
					Msgf("Agent-type: %s, agent-ID: %s, createdAt: %s, "+
						"summary: %#v",
						l.Agent.Type,
						usersMap[l.Agent.ID],
						l.CreatedAt,
						l.Channel.Summary)
			}
		}

		if !logEntries.More {
			break
		}
	}
}
