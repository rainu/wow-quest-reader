package model

func (q QuestInformation) IsValid() bool {
	if q.Id == "" {
		return false
	}
	if q.Description == "" && q.Completion == "" && q.Progress == "" {
		return false
	}

	return true
}
