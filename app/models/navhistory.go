package models

import "atfm/generics"

type NavHistoryRec struct {
	Path     string
	Index    int
	Mod      FsMod
	BasePath string
}

type NavHistory struct {
	historyBack    *generics.Stack[NavHistoryRec]
	historyForward *generics.Stack[NavHistoryRec]
}

func (n *NavHistory) GetHistoryStack(mod NavHistoryMod) *generics.Stack[NavHistoryRec] {
	if mod {
		return n.historyBack
	} else {
		return n.historyForward
	}
}

func NewNavHistory() *NavHistory {
	return &NavHistory{
		historyBack:    generics.NewStack[NavHistoryRec](),
		historyForward: generics.NewStack[NavHistoryRec](),
	}
}
