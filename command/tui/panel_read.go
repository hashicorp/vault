package tui

import (
	"fmt"
	"sort"
	
	"github.com/hashicorp/vault/api"
	"github.com/nsf/termbox-go"
)

type ReadPanel struct {
	path       string
	kv         map[string]interface{}
	sortedKeys []string
}

func NewReadPanel(client *api.Client, path string) (*ReadPanel, error) {

	secret, err := client.Logical().ReadWithData(path, map[string][]string{})
	if err != nil {
		return nil, err
	}
	if secret == nil {
		return nil, fmt.Errorf("no value found at %s", path)
	}

	var keys []string
	for k := range secret.Data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	l := &ReadPanel{
		kv:         secret.Data,
		path:       path,
		sortedKeys: keys,
	}

	return l, nil
}

func (r *ReadPanel) handleAction(curserY int) (TerminalPanel, error) {
	return nil, nil
}

func (r *ReadPanel) draw(curserY int) error {

	label := fmt.Sprintf("Vault path: %s", r.path)

	for j, ch := range label {
		termbox.SetCell(j, 0, ch, termbox.ColorGreen|termbox.AttrBold, defaultBgColor)
	}

	for idx, key := range r.sortedKeys {

		entry := fmt.Sprintf("%s: %s", key, r.kv[key])

		for j, ch := range entry {
			termbox.SetCell(j+1, idx+2, ch, defaultFgColor, defaultBgColor)
		}
	}

	return nil
}

func (r *ReadPanel) maxY() int {
	return len(r.kv) - 1
}
