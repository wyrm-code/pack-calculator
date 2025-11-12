package solver

import (
	"errors"
	"sort"
)

// State used for DP reconstruction.
type state struct {
	OK    bool
	Packs int
	Prev  int
	Size  int
}

// Solve returns a pack breakdown that ships at least `items` using the
// smallest total items; within that, using the fewest number of packs.
func Solve(items int, sizes []int) (map[int]int, int, error) {
	if items <= 0 {
		return nil, 0, errors.New("items must be > 0")
	}
	if len(sizes) == 0 {
		return nil, 0, errors.New("at least one pack size is required")
	}

	// sanitize: positive, remove dups
	sizeSet := map[int]struct{}{}
	var sizeBucket []int
	for _, s := range sizes {
		if s <= 0 {
			return nil, 0, errors.New("pack sizes must be positive")
		}
		if _, ok := sizeSet[s]; !ok {
			sizeSet[s] = struct{}{}
			sizeBucket = append(sizeBucket, s)
		}
	}
	sort.Ints(sizeBucket) // ascending for DP convenience

	maxS := sizeBucket[len(sizeBucket)-1]
	limit := items + maxS - 1

	dp := make([]state, limit+1)
	dp[0] = state{OK: true, Packs: 0, Prev: -1, Size: 0}

	for _, s := range sizeBucket {
		for t := s; t <= limit; t++ {
			if dp[t-s].OK {
				cand := dp[t-s].Packs + 1
				if !dp[t].OK || cand < dp[t].Packs {
					dp[t] = state{OK: true, Packs: cand, Prev: t - s, Size: s}

				}
			}
		}
	}

	// choose best S >= items with min total; then min packs
	bestS := -1
	bestP := 0
	for S := items; S <= limit; S++ {
		if !dp[S].OK {
			continue
		}
		if bestS == -1 || S < bestS || (S == bestS && dp[S].Packs < bestP) {
			bestS = S
			bestP = dp[S].Packs
		}
	}
	if bestS == -1 {
		return nil, 0, errors.New("no solution found (unexpected with positive sizes)")
	}
	// reconstruct
	counts := map[int]int{}
	for cur := bestS; cur > 0; {
		st := dp[cur]
		counts[st.Size]++
		cur = st.Prev
	}
	return counts, bestS, nil
}

// SortedDescKeys returns the keys descending (helper for printing stable output).
func SortedDescKeys(m map[int]int) []int {
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(keys)))
	return keys
}
