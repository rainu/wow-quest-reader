package model

func (q QuestInformation) IsValid() bool {
	if q.Id == "" {
		return false
	}
	if q.Text == "" && q.Completion == "" && q.Progress == "" {
		return false
	}

	return true
}
