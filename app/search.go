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

	selectedSearchHistory int
}

func NewSearch() *Search {
	return &Search{
		SearchRes:             []SearchElement{},
		LastSearchText:        "",
		SrchHistory:           []string{},
		selectedSearchHistory: 0,
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

func (s *Search) SearchContent(text string, ins *Instance, saveSeachHyst bool, ignoreCase bool) {
	text = strings.Trim(text, " ")
	if text == "" {
		return
	}
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
	s.selectedSearchHistory = len(s.SrchHistory)

	s.LastSearchText = text

	if ignoreCase {
		text = strings.ToLower(text)
	}
	m := []SearchElement{}
	for i, fi := range ins.ShownContent {
		v := fi.Name
		vv := v
		if ignoreCase {
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
			if ignoreCase {
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

	if len(m) > 0 {
		ins.CurrentItem = m[0].OriginalIndex
	}
}

func (s *Search) SearchNext() string {
	histLen := len(s.SrchHistory)
	if histLen == 0 {
		return ""
	}
	s.selectedSearchHistory++
	t := ""
	if s.selectedSearchHistory > histLen {
		s.selectedSearchHistory = histLen
	}
	if s.selectedSearchHistory != histLen {
		t = s.SrchHistory[s.selectedSearchHistory]
	}
	return t
}

func (s *Search) SearchPrevious() string {
	if len(s.SrchHistory) == 0 {
		return ""
	}
	s.selectedSearchHistory--
	if s.selectedSearchHistory < 0 {
		s.selectedSearchHistory = 0
	}
	t := s.SrchHistory[s.selectedSearchHistory]
	return t
}
