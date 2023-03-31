package system

import (
	"golang.design/x/hotkey"
)

func _imm() map[string]hotkey.Modifier {
	m := map[string]hotkey.Modifier{}

	_mm(m, "Ctrl", hotkey.ModCtrl)
	_mm(m, "Shift", hotkey.ModShift)
	_mm(m, "Alt", hotkey.ModAlt)
	_mm(m, "Wind", hotkey.ModWin)

	return m
}
