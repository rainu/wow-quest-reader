package system

import (
	"context"
	"github.com/sirupsen/logrus"
	"golang.design/x/clipboard"
)

type clipboardAccess struct{}

func NewClipboardAccess() (*clipboardAccess, error) {
	err := clipboard.Init()
	if err != nil {
		return nil, err
	}

	return &clipboardAccess{}, nil
}

func (c *clipboardAccess) Clear(ctx context.Context) {
	logrus.Debug("Clear clipboard content.")

	select {
	case <-ctx.Done():
	case <-clipboard.Write(clipboard.FmtText, []byte{}):
	}
}

func (c *clipboardAccess) Watch(ctx context.Context, contentChan chan string) {
	defer close(contentChan)

	logrus.Info("Start watching clipboard.")

	cbChan := clipboard.Watch(ctx, clipboard.FmtText)
	for {
		select {
		case content := <-cbChan:
			contentChan <- string(content)
		case <-ctx.Done():
			logrus.Info("Stop watching clipboard.")
			return
		}
	}
}
