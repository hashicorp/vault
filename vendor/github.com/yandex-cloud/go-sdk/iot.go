// Copyright (c) 2019 Yandex LLC. All rights reserved.
// Author: Andrey Khaliullin <avhaliullin@yandex-team.ru>

package ycsdk

import (
	"github.com/yandex-cloud/go-sdk/gen/iot/devices"
)

const (
	IoTDevicesServiceID Endpoint = "iot-devices"
)

func (sdk *SDK) IoT() *IoT {
	return &IoT{sdk: sdk}
}

type IoT struct {
	sdk *SDK
}

func (m *IoT) Devices() *devices.Devices {
	return devices.NewDevices(m.sdk.getConn(IoTDevicesServiceID))
}
