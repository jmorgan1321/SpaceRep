package core

import (
	"fmt"
	"testing"

	"github.com/jmorgan1321/SpaceRep/internal/test"
)

func TestBucketNextPrevBucket(t *testing.T) {
	tests := []struct {
		in, exp Bucket
		method  string
	}{
		{method: "next", in: Daily, exp: Weekly},
		{method: "next", in: Weekly, exp: Monthly},
		{method: "next", in: Monthly, exp: Yearly},
		{method: "next", in: Yearly, exp: Yearly},

		{method: "prev", in: Yearly, exp: Monthly},
		{method: "prev", in: Monthly, exp: Weekly},
		{method: "prev", in: Weekly, exp: Daily},
		{method: "prev", in: Daily, exp: Daily},
	}
	for i, tt := range tests {
		if tt.method == "prev" {
			test.ExpectEQ(t, tt.exp, tt.in.PrevBucket(), fmt.Sprintf("test %v", i))
		} else {
			test.ExpectEQ(t, tt.exp, tt.in.NextBucket(), fmt.Sprintf("test %v", i))
		}
	}
}

func TestBucketGetMaxCount(t *testing.T) {
	tests := []struct {
		in  Bucket
		exp int
	}{
		{in: Daily, exp: 8},
		{in: Weekly, exp: 4},
		{in: Monthly, exp: 2},
		{in: Yearly, exp: 1},
	}
	for i, tt := range tests {
		test.ExpectEQ(t, tt.exp, tt.in.GetMaxCount(), fmt.Sprintf("test %v", i))
	}
}

func TestInfoUpdateBucket(t *testing.T) {
	tests := []struct {
		in, exp Info
	}{
		{in: Info{Count: Daily.GetMaxCount() - 1, B: Daily}, exp: Info{Count: Daily.GetMaxCount() - 1, B: Daily}},
		{in: Info{Count: Daily.GetMaxCount(), B: Daily}, exp: Info{Count: 0, B: Weekly}},

		{in: Info{Count: Weekly.GetMaxCount() - 1, B: Weekly}, exp: Info{Count: Weekly.GetMaxCount() - 1, B: Weekly}},
		{in: Info{Count: Weekly.GetMaxCount(), B: Weekly}, exp: Info{Count: 0, B: Monthly}},

		{in: Info{Count: Monthly.GetMaxCount() - 1, B: Monthly}, exp: Info{Count: Monthly.GetMaxCount() - 1, B: Monthly}},
		{in: Info{Count: Monthly.GetMaxCount(), B: Monthly}, exp: Info{Count: 0, B: Yearly}},

		{in: Info{Count: Yearly.GetMaxCount() - 1, B: Yearly}, exp: Info{Count: Yearly.GetMaxCount() - 1, B: Yearly}},
		{in: Info{Count: Yearly.GetMaxCount(), B: Yearly}, exp: Info{Count: 0, B: Yearly}},

		{in: Info{Count: 0, B: Yearly}, exp: Info{Count: 0, B: Yearly}},
		{in: Info{Count: -1, B: Yearly}, exp: Info{Count: 0, B: Monthly}},

		{in: Info{Count: 0, B: Monthly}, exp: Info{Count: 0, B: Monthly}},
		{in: Info{Count: -1, B: Monthly}, exp: Info{Count: 0, B: Weekly}},

		{in: Info{Count: 0, B: Weekly}, exp: Info{Count: 0, B: Weekly}},
		{in: Info{Count: -1, B: Weekly}, exp: Info{Count: 0, B: Daily}},

		{in: Info{Count: 0, B: Daily}, exp: Info{Count: 0, B: Daily}},
		{in: Info{Count: -1, B: Daily}, exp: Info{Count: 0, B: Daily}},
	}
	for i, tt := range tests {
		tt.in.UpdateBucket()
		test.ExpectEQ(t, tt.exp.B, tt.in.B, fmt.Sprintf("test %v: Bucket", i))
		test.ExpectEQ(t, tt.exp.Count, tt.in.Count, fmt.Sprintf("test %v: Count", i))
	}
}

type mockCard struct {
	Info
}

// Display
func (*mockCard) Type() string    { return "" }
func (*mockCard) Name() string    { return "" }
func (*mockCard) Clone(Info) Card { return nil }
func (*mockCard) Tmpl() string    { return "" }
