package signals

// NilSignal is a special signal that is blank or "nil"
type NilSignal int

func (s *NilSignal) String() string { return "SIGNIL" }
func (s *NilSignal) Signal()        {}
