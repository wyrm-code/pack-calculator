package solver

import "testing"

func TestExamples(t *testing.T) {
	sizes := []int{250, 500, 1000, 2000, 5000}

	check := func(items int, wantTotal int, want map[int]int) {
		got, total, err := Solve(items, sizes)
		if err != nil { t.Fatalf("Solve(%d) error: %v", items, err) }
		if total != wantTotal {
			t.Fatalf("Solve(%d) total=%d want=%d", items, total, wantTotal)
		}
		for s, q := range want {
			if got[s] != q {
				t.Fatalf("Solve(%d) packs[%d]=%d want=%d", items, s, got[s], q)
			}
		}
		// verify no extra unexpected packs
		for s := range got {
			if want[s] == 0 && got[s] != 0 {
				t.Fatalf("Solve(%d) unexpected size %d qty %d", items, s, got[s])
			}
		}
	}

	check(1, 250, map[int]int{250:1})
	check(250, 250, map[int]int{250:1})
	check(251, 500, map[int]int{500:1})
	check(501, 750, map[int]int{500:1, 250:1})
	check(12001, 12250, map[int]int{5000:2, 2000:1, 250:1})
}

func TestCustomSizes(t *testing.T) {
	// money analogy from the prompt
	sizes := []int{25, 10, 5} // quarters, dimes, nickels
	got, total, err := Solve(100, sizes) // $1
	if err != nil { t.Fatal(err) }
	if total != 100 { t.Fatalf("total=%d want=100", total) }
	// Among optimal totals, fewest packs is 4 quarters
	if got[25] != 4 || len(got) != 1 {
		t.Fatalf("want 4x25, got %#v", got)
	}
}

func TestInvalid(t *testing.T) {
	if _, _, err := Solve(0, []int{1}); err == nil {
		t.Fatal("expected error for items=0")
	}
	if _, _, err := Solve(1, []int{}); err == nil {
		t.Fatal("expected error for empty sizes")
	}
	if _, _, err := Solve(1, []int{-1, 5}); err == nil {
		t.Fatal("expected error for negative size")
	}
}
