package sanity

import (
	"strconv"
	"testing"
)

func TestIf(t *testing.T) {
	AssertEqual(t, If(true, 1, 2), 1)
	AssertEqual(t, If(false, 1, 2), 2)
	AssertEqual(t, If(true, "1", "2"), "1")
	AssertEqual(t, If(false, "1", "2"), "2")
}

func TestMap(t *testing.T) {
	AssertSliceEqual(t,
		Map([]string{"1", "3", "5"},
			func(s string) float64 {
				if v, err := strconv.ParseFloat(s, 32); err == nil {
					return v
				} else {
					panic(err)
				}
			}),
		[]float64{1, 3, 5})
}

func TestFilter(t *testing.T) {
	AssertSliceEqual(t,
		Filter([]uint64{0, 2, 4, 6, 8, 10},
			func(i int, v uint64) bool {
				return i == 0 || v > 5
			}),
		[]uint64{0, 6, 8, 10})
}

func TestChannelFns(t *testing.T) {
	makeTestChan := func() <-chan string {
		c := make(chan string)
		go func() {
			c <- "56"
			c <- "78"
			c <- "34"
			c <- "12"
			close(c)
		}()
		return c
	}
	{
		// FilterChan
		cf := FilterChan(makeTestChan(), func(v string) bool { return v[0] >= '5' })
		AssertEqual(t, <-cf, "56")
		AssertEqual(t, <-cf, "78")
		if _, ok := <-cf; ok {
			t.Errorf("Filtered channel isn't closed")
		}
	}
	{
		s := ChanToSlice(makeTestChan())
		AssertSliceEqual(t, s, []string{"56", "78", "34", "12"})
	}
}

func TestSaneSort(t *testing.T) {
	lessfn := func(a, b string) bool { return a < b }
	values := []string{"56", "78", "34", "12"}
	SaneSortSlice(values, lessfn)
	AssertSliceEqual(t,
		values,
		[]string{"12", "34", "56", "78"})
}
