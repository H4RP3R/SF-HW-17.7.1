package counter

import (
	"fmt"
	"sync"
	"testing"
)

func TestCounterIncrement(t *testing.T) {
	type args struct {
		maxVal int
		step   int
	}
	tests := []struct {
		count int
		args
		want int
	}{
		{42, args{maxVal: 42, step: 1}, 42},
		{112, args{maxVal: 88, step: 1}, 88},
		{2211, args{maxVal: 2211, step: 2}, 2210},
		{99999, args{maxVal: 66666, step: 5}, 66665},
	}

	for _, tt := range tests {
		name := fmt.Sprintf("count times: %d, args: %+v", tt.count, tt.args)
		t.Run(name, func(t *testing.T) {
			c := New(tt.maxVal, tt.step)
			for i := 0; i < tt.count; i++ {
				c.Increment()
			}
			if c.Val() != tt.want {
				t.Errorf("got: %d, want: %d", c.Val(), tt.want)
			}
			c.Close()
		})
	}
}

func TestCounterIncrementMultiThreading(t *testing.T) {
	type args struct {
		maxVal int
		step   int
	}
	tests := []struct {
		count int
		args
		want int
		gNum int
	}{
		{42, args{maxVal: 42, step: 1}, 42, 2},
		{112, args{maxVal: 88, step: 1}, 88, 4},
		{2100, args{maxVal: 2000, step: 1}, 2000, 10},
		{2211, args{maxVal: 2211, step: 2}, 2210, 15},
		{99999, args{maxVal: 66666, step: 5}, 66665, 121},
	}

	for _, tt := range tests {
		name := fmt.Sprintf("count times: %d, args: %+v", tt.count, tt.args)
		t.Run(name, func(t *testing.T) {
			var wg sync.WaitGroup
			c := New(tt.maxVal, tt.step)
			wg.Add(tt.gNum)
			for i := 0; i < tt.gNum; i++ {
				go func() {
					defer wg.Done()
					for i := 0; i < tt.count/tt.gNum; i++ {
						c.Increment()
					}
				}()
			}
			wg.Wait()
			if c.Val() != tt.want {
				t.Errorf("got: %d, want: %d", c.Val(), tt.want)
			}
			c.Close()
		})
	}
}
