// +build !enterprise

package logical

type entReq struct {
	ControlGroup interface{}
}

func (r *Request) EntReq() *entReq {
	return &entReq{}
}

func (r *Request) SetEntReq(*entReq) {
}
