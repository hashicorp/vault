package ycsdk

import (
	"github.com/yandex-cloud/go-sdk/gen/functions"
	"github.com/yandex-cloud/go-sdk/gen/triggers"
)

type Serverless struct {
	sdk *SDK
}

const (
	FunctionServiceID Endpoint = "serverless-functions"
	TriggerServiceID  Endpoint = "serverless-triggers"
)

func (s *Serverless) Functions() *functions.Function {
	return functions.NewFunction(s.sdk.getConn(FunctionServiceID))
}

func (s *Serverless) Triggers() *triggers.Trigger {
	return triggers.NewTrigger(s.sdk.getConn(TriggerServiceID))
}
