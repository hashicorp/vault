package command

type panelAction func() (TerminalPanel, error)

type TerminalPanel interface {
	draw(curserY int) error
	handleAction(curserY int) (TerminalPanel, error)
	maxY() int
}
