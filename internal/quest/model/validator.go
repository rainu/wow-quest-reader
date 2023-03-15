package model

import "fmt"

func (q *Quest) IsValid() error {
	if q.Id < 0 {
		return fmt.Errorf("missing id")
	}
	if q.Locale == "" {
		return fmt.Errorf("missing locale")
	}
	if q.Title == "" {
		return fmt.Errorf("missing title")
	}
	if q.Description == "" && q.Progress == "" && q.Completion == "" && !q.Obsolete {
		return fmt.Errorf("missing description/progress/completion")
	}
	if q.StartNPC == nil && q.StartObject == nil && q.StartItem == nil && !q.Obsolete {
		return fmt.Errorf("no start npc/object/item")
	}
	if q.EndNPC == nil && q.EndObject == nil && q.EndItem == nil && !q.Obsolete {
		return fmt.Errorf("no end npc/object/item")
	}

	if q.StartNPC != nil && q.StartNPC.Id < 0 {
		return fmt.Errorf("missing npc id")
	}
	if q.EndNPC != nil && q.EndNPC.Id < 0 {
		return fmt.Errorf("missing npc id")
	}

	if q.StartObject != nil && q.StartObject.Id < 0 {
		return fmt.Errorf("missing object id")
	}
	if q.EndObject != nil && q.EndObject.Id < 0 {
		return fmt.Errorf("missing object id")
	}

	if q.StartItem != nil && q.StartItem.Id < 0 {
		return fmt.Errorf("missing item id")
	}
	if q.EndItem != nil && q.EndItem.Id < 0 {
		return fmt.Errorf("missing item id")
	}

	return nil
}
