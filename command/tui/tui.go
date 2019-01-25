package tui

import (
	"github.com/hashicorp/vault/api"
	"github.com/nsf/termbox-go"
)

const foregroundColor = termbox.ColorDefault
const backgroundColor = termbox.ColorBlack

type TerminalUI struct {
	client     *api.Client
	panels     *PanelStack
	selectionY int
}

// New initializes a new terminal user interface by listing the
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

	return tui, err
}

// Starts will initialize termbox and start the event loop for handling key events
func (t *TerminalUI) Start() error {

	err := termbox.Init()
	if err != nil {
		return err
	}
	defer termbox.Close()

	termbox.SetInputMode(termbox.InputEsc)
	reallocBackBuffer(termbox.Size())

	t.draw()

mainloop:
	for {

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

		if err != nil {
			return err
		}

		t.draw()
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

func (t *TerminalUI) draw() error {
	termbox.Clear(foregroundColor, backgroundColor)

	t.panels.Peek().draw(t.selectionY)

	for i := 0; i < bbw; i++ {
		termbox.SetCell(i, bbh - 1, 0, 0, termbox.ColorWhite)
	}

	for i, c := range t.client.Address() {
		termbox.SetCell(bbw - len(t.client.Address()) - 1 + i, bbh - 1, c, termbox.ColorBlack | termbox.AttrBold, termbox.ColorWhite)
	}

	termbox.Flush()

	return nil
}
