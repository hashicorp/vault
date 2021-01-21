package configutil

import "time"

type BarrierRotationConfig struct {
	Operations uint32
	Interval   time.Duration
}

type Barrier struct {
	BarrierRotationConfig
}
