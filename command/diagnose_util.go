package command

type DiagnoseObserver interface {
	Exited(status int)
	Success(key string)
	Error(key string, err error)
	IsEnabled() bool
}

type NullDiagnoseObserver struct {
}

func (n *NullDiagnoseObserver) Exited(status int) {
}

func (n *NullDiagnoseObserver) Success(key string) {
}

func (n *NullDiagnoseObserver) Error(key string, err error) {
}

func (n *NullDiagnoseObserver) IsEnabled() bool {
	return false
}
