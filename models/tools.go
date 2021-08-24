package models

import "sort"

func RankByPoint(m map[int]int) PosList {
	pl := make(PosList, len(m))
	i := 0
	for k, v := range m {
		pl[i] = Pos{k, v}
		i++
	}
	sort.Sort(pl)
	return pl
}

type Pos struct {
	Key   int
	Value int
}

type PosList []Pos

func (p PosList) Len() int           { return len(p) }
func (p PosList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PosList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
