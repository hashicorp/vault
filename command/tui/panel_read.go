package tui

import (
	"fmt"
	"github.com/hashicorp/vault/api"
	"github.com/nsf/termbox-go"
)

type ReadPanel struct {
	path    string
	kv map[string]interface{}
}

func NewReadPanel(client *api.Client, path string) (*ReadPanel, error) {

	secret, err := client.Logical().ReadWithData(path, map[string][]string{})
	if err != nil {
		return nil, err
	}
	if secret == nil {
		return nil, fmt.Errorf("no value found at %s", path)
	}


	l := &ReadPanel{
		kv: secret.Data,
		path:    path,
	}

	return l, nil
}

func (r *ReadPanel) handleAction(curserY int) (TerminalPanel, error) {
	return nil, nil
}

func (r *ReadPanel) draw(curserY int) error {

	label := fmt.Sprintf("Vault path: %s", r.path)

	for j, ch := range label {
		termbox.SetCell(j, 0, ch, termbox.ColorGreen | termbox.AttrBold, defaultBgColor)
	}

	i := 0
	for k, v := range r.kv { // this loop is not stable it's just a POC ;)

		entry := fmt.Sprintf("%s: %s", k, v)

		for j, ch := range entry {
			termbox.SetCell(j+1, i+2, ch, defaultFgColor, defaultBgColor)
		}

		i++
	}

	return nil
}

func (r *ReadPanel) maxY() int {
	return len(r.kv) - 1
}