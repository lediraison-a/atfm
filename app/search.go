package app

import (
	"sort"
	"strings"
)

type SearchElement struct {
	Source, Target   string
	OriginalIndex    int
	InTargetStrIndex []int
}

type Search struct {
	SearchRes      []SearchElement
	LastSearchText string
	SrchHistory    []string

	ignoreCase bool
	incSearch  bool
}

func NewSearch() *Search {
	return &Search{
		SearchRes:      []SearchElement{},
		LastSearchText: "",
		SrchHistory:    []string{},
		ignoreCase:     true,
		incSearch:      true,
	}
}

func (s *Search) ResetSearch() {
	s.SearchRes = []SearchElement{}
}

func (s *Search) GetLastSrch() string {
	if len(s.SrchHistory) == 0 {
		return ""
	} else {
		return s.SrchHistory[len(s.SrchHistory)-1]
	}
}

func (s *Search) SearchJumpForward(ins *Instance) {
	if len(s.SearchRes) == 0 {
		return
	}

	for _, v := range s.SearchRes {
		if v.OriginalIndex > ins.CurrentItem {
			ins.CurrentItem = v.OriginalIndex
			return
		}
	}
	ins.CurrentItem = s.SearchRes[0].OriginalIndex
}

func (s *Search) SearchJumpBackward(ins *Instance) {
	if len(s.SearchRes) == 0 {
		return
	}

	for i := len(s.SearchRes) - 1; i >= 0; i-- {
		v := s.SearchRes[i]
		if v.OriginalIndex < ins.CurrentItem {
			ins.CurrentItem = v.OriginalIndex
			return
		}
	}
	ins.CurrentItem = s.SearchRes[len(s.SearchRes)-1].OriginalIndex
}

func (s *Search) SearchContent(text string, ins *Instance) {
	Mystrutil := func(s, subs string) []int {
		ns := s
		r := []int{}
		for strings.Contains(ns, subs) {
			i := strings.Index(ns, subs)
			r = append(r, i, len(subs))
			ns = ns[i+len(subs):]
		}
		return r
	}

	s.SrchHistory = append(s.SrchHistory, text)

	s.LastSearchText = text

	if s.ignoreCase {
		text = strings.ToLower(text)
	}
	m := []SearchElement{}
	for i, fi := range ins.ShownContent {
		v := fi.Name
		vv := v
		if s.ignoreCase {
			vv = strings.ToLower(v)
		}
		if strings.HasPrefix(vv, text) {
			m = append(m, SearchElement{
				Source:           v,
				Target:           text,
				OriginalIndex:    i,
				InTargetStrIndex: []int{0, len(text)},
			})
		}
	}
	if len(m) == 0 {
		for i, fi := range ins.ShownContent {
			v := fi.Name
			vv := v
			if s.ignoreCase {
				vv = strings.ToLower(v)
			}
			if strings.Contains(vv, text) {
				m = append(m, SearchElement{
					Source:           v,
					Target:           text,
					OriginalIndex:    i,
					InTargetStrIndex: Mystrutil(vv, text),
				})
			}
		}
	}

	sort.Slice(s.SearchRes, func(i, j int) bool {
		return s.SearchRes[i].OriginalIndex < s.SearchRes[j].OriginalIndex
	})

	s.SearchRes = m

	if s.incSearch && len(m) > 0 {
		ins.CurrentItem = m[0].OriginalIndex
	}
}
