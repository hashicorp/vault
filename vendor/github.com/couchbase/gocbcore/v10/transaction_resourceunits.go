package gocbcore

import (
	"errors"
)

type resourceUnitCallback func(result *ResourceUnitResult)

func noopResourceUnitCallback(*ResourceUnitResult) {}

func (t *transactionAttempt) ReportResourceUnits(units *ResourceUnitResult) {
	if units == nil {
		return
	}

	t.recordResourceUnit(units)
}

func (t *transactionAttempt) ReportResourceUnitsError(err error) {
	if err == nil {
		return
	}

	var kerr *KeyValueError
	if errors.As(err, &kerr) {
		if kerr.Internal.ResourceUnits != nil {
			t.recordResourceUnit(kerr.Internal.ResourceUnits)
		}
	}
}
