package public_test

import (
	"testing"

	"git.sr.ht/~arielcostas/new.vigo360.es/public"
)

func TestFindMatchingTags(t *testing.T) {
	testCases := []struct {
		desc   string
		t1     []string
		t2     []string
		result int
	}{
		{
			desc:   "test case 1",
			t1:     []string{"1", "2", "3"},
			t2:     []string{"3", "2", "5"},
			result: 2,
		},
		{
			desc:   "test case 2",
			t1:     []string{"1", "2", "3"},
			t2:     []string{"3", "2", "1"},
			result: 3,
		},
		{
			desc:   "test case 3",
			t1:     []string{},
			t2:     []string{"3", "2", "5"},
			result: 0,
		},
		{
			desc:   "test case 4",
			t1:     []string{"1", "2", "3"},
			t2:     []string{},
			result: 0,
		},
		{
			desc:   "test case 5",
			t1:     []string{"1", "2", "3"},
			t2:     []string{"", "", "1"},
			result: 1,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			var got = public.FindMatchingTags(tC.t1, tC.t2)
			if got != tC.result {
				t.Fatalf("%s failed: expected %d got %d", tC.desc, tC.result, got)
			}
		})
	}
}
