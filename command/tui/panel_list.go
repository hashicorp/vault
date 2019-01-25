package tui

import (
	"fmt"
	"path"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/nsf/termbox-go"
)

type ListPanel struct {
	path    string
	entries []string
	actions []panelAction
}

func NewListPanel(client *api.Client, path string) (*ListPanel, error) {

	secret, err := client.Logical().List(path)
	if err != nil {
		return nil, err
	}

	listData, ok := extractListData(secret)
	if !ok {
		return nil, fmt.Errorf("could not extract list data")
	}

	entries := make([]string, len(listData))
	actions := make([]panelAction, len(listData))

	for i := 0; i < len(listData); i++ {
		entry := fmt.Sprintf("%s", listData[i])

		entries[i] = entry
		actions[i] = listActionGenerator(client, path, entry)
	}

	l := &ListPanel{
		actions: actions,
		entries: entries,
		path:    path,
	}

	return l, nil
}

func (l *ListPanel) handleAction(curserY int) (TerminalPanel, error) {
	return l.actions[curserY]()
}

func (l *ListPanel) draw(curserY int) error {

	label := fmt.Sprintf("Vault path: %s", l.path)

	for j, ch := range label {
		termbox.SetCell(j, 0, ch, termbox.ColorGreen | termbox.AttrBold, defaultBgColor)
	}

	for i := 0; i < len(l.entries); i++ {
		entry := fmt.Sprintf("%s", l.entries[i])

		bgColor := defaultBgColor
		if i == curserY {
			bgColor |= termbox.AttrReverse
		}

		fgColor := defaultFgColor
		if !strings.HasSuffix(entry, "/") {
			fgColor = termbox.AttrBold
		}

		for j, ch := range entry {
			termbox.SetCell(j+1, i+2, ch, fgColor, bgColor)
		}
	}

	return nil
}

func (l *ListPanel) maxY() int {
	return len(l.entries) - 1
}

func listActionGenerator(client *api.Client, base, entry string) panelAction {

	if strings.HasSuffix(entry, "/") {

		return func() (TerminalPanel, error) {
			l, err := NewListPanel(client, path.Join(base, entry))
			if err != nil {
				return nil, err
			}
			return l, nil
		}

	} else {

		return func() (TerminalPanel, error) {

			r, err := NewReadPanel(client, path.Join(base, entry))
			if err != nil {
				return nil, err
			}
			return r, nil
		}
	}
}
