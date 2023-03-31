package system

import (
	"fmt"
	"golang.design/x/hotkey"
	"strconv"
	"strings"
)

func (h *Hotkey) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var sValue string
	err := unmarshal(&sValue)
	if err != nil {
		return err
	}

	sValue = strings.TrimSpace(sValue)
	sValue = strings.ToLower(sValue)
	sValue = strings.ReplaceAll(sValue, " ", "")

	keys := strings.Split(sValue, "+")

	hk := Hotkey{}
	for _, key := range keys {
		if strings.HasPrefix(key, "[") && strings.HasSuffix(key, "]") {
			key = strings.ReplaceAll(key, "[", "")
			key = strings.ReplaceAll(key, "]", "")

			m, exists := modMap[key]
			if !exists {
				return fmt.Errorf("invalid modifier: %s", key)
			}
			hk.Modifier = append(hk.Modifier, m)
		} else if strings.HasPrefix(key, "<") && strings.HasSuffix(key, ">") {
			key = strings.ReplaceAll(key, "<", "")
			key = strings.ReplaceAll(key, ">", "")

			if strings.HasPrefix(key, "0x") {
				key = strings.ReplaceAll(key, "0x", "")
				r, err := strconv.ParseUint(key, 16, 16)
				if err != nil {
					return fmt.Errorf("invalid hex value for key (%s): %e", key, err)
				}
				hk.Key = hotkey.Key(uint16(r))
			} else {
				var exists bool
				hk.Key, exists = keyMap[key]
				if !exists {
					return fmt.Errorf("unknown key: %s", key)
				}
			}
		}
	}

	*h = hk
	return nil
}

var keyMap = _ikm()
var modMap = _imm()

func _ikm() map[string]hotkey.Key {
	m := map[string]hotkey.Key{}

	_km(m, "Space", hotkey.KeySpace)
	_km(m, "1", hotkey.Key1)
	_km(m, "2", hotkey.Key2)
	_km(m, "3", hotkey.Key3)
	_km(m, "4", hotkey.Key4)
	_km(m, "5", hotkey.Key5)
	_km(m, "6", hotkey.Key6)
	_km(m, "7", hotkey.Key7)
	_km(m, "8", hotkey.Key8)
	_km(m, "9", hotkey.Key9)
	_km(m, "0", hotkey.Key0)
	_km(m, "A", hotkey.KeyA)
	_km(m, "B", hotkey.KeyB)
	_km(m, "C", hotkey.KeyC)
	_km(m, "D", hotkey.KeyD)
	_km(m, "E", hotkey.KeyE)
	_km(m, "F", hotkey.KeyF)
	_km(m, "G", hotkey.KeyG)
	_km(m, "H", hotkey.KeyH)
	_km(m, "I", hotkey.KeyI)
	_km(m, "J", hotkey.KeyJ)
	_km(m, "K", hotkey.KeyK)
	_km(m, "L", hotkey.KeyL)
	_km(m, "M", hotkey.KeyM)
	_km(m, "N", hotkey.KeyN)
	_km(m, "O", hotkey.KeyO)
	_km(m, "P", hotkey.KeyP)
	_km(m, "Q", hotkey.KeyQ)
	_km(m, "R", hotkey.KeyR)
	_km(m, "S", hotkey.KeyS)
	_km(m, "T", hotkey.KeyT)
	_km(m, "U", hotkey.KeyU)
	_km(m, "V", hotkey.KeyV)
	_km(m, "W", hotkey.KeyW)
	_km(m, "X", hotkey.KeyX)
	_km(m, "Y", hotkey.KeyY)
	_km(m, "Z", hotkey.KeyZ)
	_km(m, "Return", hotkey.KeyReturn)
	_km(m, "Escape", hotkey.KeyEscape)
	_km(m, "Delete", hotkey.KeyDelete)
	_km(m, "Tab", hotkey.KeyTab)
	_km(m, "Left", hotkey.KeyLeft)
	_km(m, "Right", hotkey.KeyRight)
	_km(m, "Up", hotkey.KeyUp)
	_km(m, "Down", hotkey.KeyDown)
	_km(m, "F1", hotkey.KeyF1)
	_km(m, "F2", hotkey.KeyF2)
	_km(m, "F3", hotkey.KeyF3)
	_km(m, "F4", hotkey.KeyF4)
	_km(m, "F5", hotkey.KeyF5)
	_km(m, "F6", hotkey.KeyF6)
	_km(m, "F7", hotkey.KeyF7)
	_km(m, "F8", hotkey.KeyF8)
	_km(m, "F9", hotkey.KeyF9)
	_km(m, "F10", hotkey.KeyF10)
	_km(m, "F11", hotkey.KeyF11)
	_km(m, "F12", hotkey.KeyF12)
	_km(m, "F13", hotkey.KeyF13)
	_km(m, "F14", hotkey.KeyF14)
	_km(m, "F15", hotkey.KeyF15)
	_km(m, "F16", hotkey.KeyF16)
	_km(m, "F17", hotkey.KeyF17)
	_km(m, "F18", hotkey.KeyF18)
	_km(m, "F19", hotkey.KeyF19)
	_km(m, "F20", hotkey.KeyF20)

	return m
}

func _km(m map[string]hotkey.Key, s string, k hotkey.Key) {
	m[strings.ToLower(s)] = k
}

func _mm(m map[string]hotkey.Modifier, s string, k hotkey.Modifier) {
	m[strings.ToLower(s)] = k
}
