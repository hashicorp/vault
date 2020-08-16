// Copyright (c) 2019 Yandex LLC. All rights reserved.
// Author: Dmitry Novikov <novikoff@yandex-team.ru>

package ycsdk

import (
	"github.com/yandex-cloud/go-sdk/gen/marketplace/metering"
)

const (
	MarketplaceMeteringServiceID Endpoint = "marketplace"
)

type Marketplace struct {
	sdk *SDK
}

func (m *Marketplace) Metering() *metering.Metering {
	return metering.NewMetering(m.sdk.getConn(MarketplaceMeteringServiceID))
}
