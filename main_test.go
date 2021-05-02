package main

import (
	"fmt"
	"testing"
)

func Test_isInInterval(t *testing.T) {
	tests := []struct {
		page     int
		interval string
		want     bool
	}{
		{page: 5, interval: "5", want: true},
		{page: 10, interval: "3", want: false},
		{page: 100, interval: "100,101", want: true},
		{page: 29, interval: "30-35", want: false},
		{page: 30, interval: "30-35", want: true},
		{page: 33, interval: "30-35", want: true},
		{page: 35, interval: "30-35", want: true},
		{page: 36, interval: "30-35", want: false},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("%d in %s", tt.page, tt.interval)
		t.Run(testname, func(t *testing.T) {
			if got := isInInterval(tt.page, tt.interval); got != tt.want {
				t.Errorf("isInInterval() = %v, want %v", got, tt.want)
			}
		})
	}
}
