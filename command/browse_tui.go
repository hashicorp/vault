package command

import (
	"github.com/hashicorp/vault/api"
	"github.com/nsf/termbox-go"
)

type TerminalUI struct {
	client *api.Client
	panels *PanelStack
	posY   int
}

func NewTerminalUI(client *api.Client, path string) (*TerminalUI, error) {

	tui := &TerminalUI{
		client: client,
		panels: NewPanelStack(),
	}

	l, err := NewListPanel(client, path)
	if err != nil {
		return nil, err
	}

	tui.panels.Push(l)

	return tui, err
}

func (t *TerminalUI) handleInput(input termbox.Key) error {

	switch input {
	case termbox.KeyArrowUp:
		if t.posY > 0 {
			t.posY--
		}

	case termbox.KeyArrowDown:
		if t.posY < t.panels.Peek().maxY() {
			t.posY++
		}

	case termbox.KeyArrowRight:
		fallthrough
	case termbox.KeyEnter:

		panel, err := t.panels.Peek().handleAction(t.posY)
		if err != nil {
			return err
		}

		t.panels.Push(panel)
		t.posY = 0

	case termbox.KeyArrowLeft:
		fallthrough
	case termbox.KeyBackspace:

		if t.panels.Len() <= 1 {
			return nil
		}

		t.panels.Pop()
		t.posY = 0
	}

	return nil
}

func (t *TerminalUI) draw() error {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	t.panels.Peek().draw(t.posY)

	termbox.Flush()

	return nil
}
