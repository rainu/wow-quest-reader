package processor

import (
	"context"
	"github.com/rainu/wow-quest-client/internal/quest/model"
	"github.com/sirupsen/logrus"
)

func (p *processor) runPersistence(
	ctx context.Context,
	questChan chan model.Quest,
	npcChan chan model.NonPlayerCharacter,
	itemChan chan model.Item,
	objectChan chan model.Object,
) {
	logrus.Debug("Start persistence worker.")
	defer func() {
		logrus.Debug("Stop persistence worker.")
	}()

	questClosed := false
	npcClosed := false
	itemClosed := false
	objectClosed := false

	for !questClosed || !npcClosed || !itemClosed || !objectClosed {
		select {
		case <-ctx.Done():
			//context closed
			return
		case q, ok := <-questChan:
			if !ok {
				questClosed = true
				continue
			}

			log := logrus.WithField("quest_id", q.Id).WithField("locale", q.Locale)

			if vErr := q.IsValid(); vErr != nil {
				log.
					WithError(vErr).
					Warning("Quest is invalid.")
				continue
			}

			err := p.store.SaveQuest(ctx, q)
			if err != nil {
				log.WithError(err).Error("Error while persisting quest!")
			}
		case npc, ok := <-npcChan:
			if !ok {
				npcClosed = true
				continue
			}

			log := logrus.WithField("npc_id", npc.Id)

			err := p.store.SaveNpc(ctx, npc)
			if err != nil {
				log.WithError(err).Error("Error while persisting npc!")
			}
		case item, ok := <-itemChan:
			if !ok {
				itemClosed = true
				continue
			}

			log := logrus.WithField("item_id", item.Id)

			err := p.store.SaveItem(ctx, item)
			if err != nil {
				log.WithError(err).Error("Error while persisting item!")
			}
		case object, ok := <-objectChan:
			if !ok {
				objectClosed = true
				continue
			}

			log := logrus.WithField("object_id", object.Id)

			err := p.store.SaveObject(ctx, object)
			if err != nil {
				log.WithError(err).Error("Error while persisting object!")
			}
		}
	}
}
