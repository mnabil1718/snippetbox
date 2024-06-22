package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHumanDate(t *testing.T) {
	t.Parallel() // run this concurrently just an experiment

	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2020, 12, 17, 10, 0, 0, 0, time.UTC),
			want: "17 Dec 2020 at 10:00",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm:   time.Date(2020, 12, 17, 10, 0, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "17 Dec 2020 at 09:00",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			str := humanDate(test.tm)
			assert.Equal(t, test.want, str)
		})
	}

}
