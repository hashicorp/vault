package tui

type (
	PanelStack struct {
		top    *node
		length int
	}
	node struct {
		value TerminalPanel
		prev  *node
	}
)

// Create a new stack
func NewPanelStack() *PanelStack {
	return &PanelStack{nil, 0}
}

// Return the number of items in the stack
func (s *PanelStack) Len() int {
	return s.length
}

// View the top item on the stack
func (s *PanelStack) Peek() TerminalPanel {
	if s.length == 0 {
		return nil
	}
	return s.top.value
}

// Pop the top item of the stack and return it
func (s *PanelStack) Pop() TerminalPanel {
	if s.length == 0 {
		return nil
	}

	n := s.top
	s.top = n.prev
	s.length--
	return n.value
}

// Push a value onto the top of the stack
func (s *PanelStack) Push(value TerminalPanel) {
	n := &node{value, s.top}
	s.top = n
	s.length++
}
