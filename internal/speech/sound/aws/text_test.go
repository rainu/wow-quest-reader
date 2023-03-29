package aws

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_transformText(t *testing.T) {
	in := `Hello DOOM SLAYER. You are welcome here! <You are thinking about this.>`
	out := `<speak>Hello <prosody volume="loud">DOOM SLAYER</prosody>. You are welcome here! You are thinking about this.</speak>`

	assert.Equal(t, out, transformText(in))
}
