package processor

import (
	"github.com/rainu/wow-quest-reader/internal/companion/system"
	"github.com/stretchr/testify/assert"
	"golang.design/x/hotkey"
	"testing"
)

func TestIsEq(t *testing.T) {
	assert.True(t, isEq(system.Hotkey{Modifier: []hotkey.Modifier{}, Key: 0x2E}, system.Hotkey{Modifier: []hotkey.Modifier{}, Key: 0x2E}))
	assert.False(t, isEq(system.Hotkey{Modifier: []hotkey.Modifier{}, Key: 0x2E}, system.Hotkey{Modifier: []hotkey.Modifier{}, Key: 0x20}))
	assert.False(t, isEq(system.Hotkey{Modifier: []hotkey.Modifier{}, Key: 0x2E}, system.Hotkey{Modifier: []hotkey.Modifier{hotkey.Mod1}, Key: 0x2E}))
}
