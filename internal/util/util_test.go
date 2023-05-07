// Package util contains helper functions
package util

import (
	"reflect"
	"testing"
	"time"
)

func TestGetClosestPastTimeRange(t *testing.T) {
	type args struct {
		fromStr string
		toStr   string
		rel     time.Time
	}
	tests := []struct {
		name     string
		args     args
		wantFrom time.Time
		wantTo   time.Time
		wantErr  bool
	}{
		{"invalid from", args{"ab:cd", "10:00", time.Time{}}, time.Time{}, time.Time{}, true},
		{"invalid to", args{"10:00", "ab:cd", time.Time{}}, time.Time{}, time.Time{}, true},
		{
			"all same day",
			args{"00:00", "01:00", time.Date(2023, time.January, 1, 23, 0, 0, 0, time.Local)},
			time.Date(2023, time.January, 1, 0, 0, 0, 0, time.Local),
			time.Date(2023, time.January, 1, 1, 0, 0, 0, time.Local),
			false,
		},
		{
			"from on prev day",
			args{"23:00", "10:00", time.Date(2023, time.January, 2, 12, 0, 0, 0, time.Local)},
			time.Date(2023, time.January, 1, 23, 0, 0, 0, time.Local),
			time.Date(2023, time.January, 2, 10, 0, 0, 0, time.Local),
			false,
		},
		{
			"both on prev day",
			args{"21:00", "22:00", time.Date(2023, time.January, 2, 12, 0, 0, 0, time.Local)},
			time.Date(2023, time.January, 1, 21, 0, 0, 0, time.Local),
			time.Date(2023, time.January, 1, 22, 0, 0, 0, time.Local),
			false,
		},
		{
			"partial overlap with rel",
			args{"12:00", "13:00", time.Date(2023, time.January, 2, 12, 0, 0, 0, time.Local)},
			time.Date(2023, time.January, 1, 12, 0, 0, 0, time.Local),
			time.Date(2023, time.January, 1, 13, 0, 0, 0, time.Local),
			false,
		},
		{
			"full 24h timeframe within rel",
			args{"12:00", "12:00", time.Date(2023, time.January, 2, 12, 0, 0, 0, time.Local)},
			time.Date(2023, time.January, 1, 12, 0, 0, 0, time.Local),
			time.Date(2023, time.January, 2, 12, 0, 0, 0, time.Local),
			false,
		},
		{
			"full 24h timeframe after rel",
			args{"13:00", "13:00", time.Date(2023, time.January, 3, 12, 0, 0, 0, time.Local)},
			time.Date(2023, time.January, 1, 13, 0, 0, 0, time.Local),
			time.Date(2023, time.January, 2, 13, 0, 0, 0, time.Local),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFrom, gotTo, err := GetClosestPastTimeRange(tt.args.fromStr, tt.args.toStr, tt.args.rel)
			if tt.wantErr {
				if err == nil {
					t.Errorf("GetClosestPastTimeRange() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if !reflect.DeepEqual(gotFrom, tt.wantFrom) {
				t.Errorf("GetClosestPastTimeRange() gotFrom = %v, want %v", gotFrom, tt.wantFrom)
			}
			if !reflect.DeepEqual(gotTo, tt.wantTo) {
				t.Errorf("GetClosestPastTimeRange() gotTo = %v, want %v", gotTo, tt.wantTo)
			}
		})
	}
}
