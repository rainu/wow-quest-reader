package model

func (i Info) IsValid() bool {
	if i.L == "" && i.Player.Name == "" && i.Player.Sex == 0 {
		return false
	}

	return true
}
