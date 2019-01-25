package tui

import (
	"fmt"
	"github.com/hashicorp/vault/api"
	"github.com/nsf/termbox-go"
)

const defaultFgColor = termbox.ColorDefault
const defaultBgColor = termbox.ColorBlack

type TerminalUI struct {
	client     *api.Client
	panels     *PanelStack
	selectionY int
}

// New initializes a new terminal user interface by listing
// Vault's data at the given path.
func New(client *api.Client, path string) (*TerminalUI, error) {

	tui := &TerminalUI{
		client: client,
		panels: NewPanelStack(),
	}

	l, err := NewListPanel(client, path)
	if err != nil {
		return nil, err
	}

	tui.panels.Push(l)

	return tui, nil
}

// Start will initialize termbox and start the event loop for handling key events
func (t *TerminalUI) Start() error {

	err := termbox.Init()
	if err != nil {
		return err
	}
	defer termbox.Close()

	termbox.SetInputMode(termbox.InputEsc)
	reallocBackBuffer(termbox.Size())

	t.draw(nil)

mainloop:
	for {

		err = nil

		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {

			case termbox.KeyEsc:
				break mainloop

			case termbox.KeyArrowUp:
				t.navigateUp()

			case termbox.KeyArrowDown:
				t.navigateDown()

			case termbox.KeyArrowRight:
				err = t.navigateRight()
			case termbox.KeyEnter:
				err = t.navigateRight()

			case termbox.KeyArrowLeft:
				t.navigateLeft()
			case termbox.KeyBackspace:
				t.navigateLeft()

			default:
				switch ev.Ch {
				case rune('q'):
					break mainloop
				case rune('h'):
					t.navigateLeft()
				case rune('j'):
					t.navigateDown()
				case rune('k'):
					t.navigateUp()
				case rune('l'):
					err = t.navigateRight()
				}
			}

		case termbox.EventResize:
			reallocBackBuffer(ev.Width, ev.Height)
		}

		t.draw(err)
	}

	return nil
}

func (t *TerminalUI) navigateUp()  {
	if t.selectionY > 0 {
		t.selectionY--
	}
}

func (t *TerminalUI) navigateDown() {
	if t.selectionY < t.panels.Peek().maxY() {
		t.selectionY++
	}
}

func (t *TerminalUI) navigateLeft() {
	if t.panels.Len() <= 1 {
		return
	}
	t.panels.Pop()
	t.selectionY = 0
}

func (t *TerminalUI) navigateRight() error {

	panel, err := t.panels.Peek().handleAction(t.selectionY)
	if err != nil {
		return err
	}

	if panel != nil {
		t.panels.Push(panel)
		t.selectionY = 0
	}

	return nil
}

func (t *TerminalUI) draw(err error) {
	termbox.Clear(defaultFgColor, defaultBgColor)

	t.panels.Peek().draw(t.selectionY)

	if err != nil {
		footerFgColor := termbox.ColorWhite
		footerBgColor := termbox.ColorRed

		for i := 0; i < bbw; i++ {
			termbox.SetCell(i, bbh - 1, 0, 0,  footerBgColor)
		}

		for i, c := range err.Error() {
			termbox.SetCell(1 + i, bbh - 1, c, footerFgColor, footerBgColor)
		}

	} else{
		footerFgColor := termbox.ColorWhite
		footerBgColor := termbox.ColorBlue

		for i := 0; i < bbw; i++ {
			termbox.SetCell(i, bbh - 1, 0, 0,  footerBgColor)
		}

		addressLabel := fmt.Sprintf("Address: %s", t.client.Address())
		for i, c := range addressLabel {
			termbox.SetCell(bbw - len(addressLabel) - 1 + i, bbh - 1, c, footerFgColor | termbox.AttrBold, footerBgColor)
		}

		hintLabel := "Press 'q' or 'ESC' to quit"
		for i, c := range hintLabel {
			termbox.SetCell(1 + i, bbh - 1, c, footerFgColor, footerBgColor)
		}
	}

	termbox.Flush()
}
