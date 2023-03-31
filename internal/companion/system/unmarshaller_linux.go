package system

import (
	"golang.design/x/hotkey"
)

func _imm() map[string]hotkey.Modifier {
	m := map[string]hotkey.Modifier{}

	_mm(m, "Ctrl", hotkey.ModCtrl)
	_mm(m, "Shift", hotkey.ModShift)
	_mm(m, "Mod1", hotkey.Mod1)
	_mm(m, "Mod2", hotkey.Mod2)
	_mm(m, "Mod3", hotkey.Mod3)
	_mm(m, "Mod4", hotkey.Mod4)
	_mm(m, "Mod5", hotkey.Mod5)

	return m
}
