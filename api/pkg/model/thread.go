package model

import (
	"net/mail"
	"sort"
	"strings"
	"time"
)

type Thread struct {
	ID           string
	Messages     []Message
	Participants []mail.Address
	Subject      string
	Updated      time.Time
}

func (m *Model) ThreadList(to string) ([]Thread, error) {
	ms, err := m.MessageList(to)
	if err != nil {
		return nil, err
	}

	tsh := map[string][]Message{}

	for _, m := range ms {
		tsh[m.Thread()] = append(tsh[m.Thread()], m)
	}

	ts := []Thread{}

	for id, ms := range tsh {
		sort.Slice(ms, func(i, j int) bool {
			return ms[i].Received.Before(ms[j].Received)
		})

		ash := map[mail.Address]bool{}

		for _, m := range ms {
			a, err := m.From()
			if err != nil {
				return nil, err
			}

			ash[*a] = true
		}

		ps := []mail.Address{}

		for a := range ash {
			ps = append(ps, a)
		}

		s, err := ms[0].Subject()
		if err != nil {
			return nil, err
		}

		t := Thread{
			ID:           id,
			Messages:     ms,
			Participants: ps,
			Subject:      strings.TrimPrefix(s, "Re: "),
			Updated:      ms[len(ms)-1].Received,
		}

		ts = append(ts, t)
	}

	sort.Slice(ts, func(i, j int) bool {
		return ts[i].Updated.After(ts[j].Updated)
	})

	return ts, nil
}
