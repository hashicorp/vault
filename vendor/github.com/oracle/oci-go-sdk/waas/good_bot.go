// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Web Application Acceleration and Security Services API
//
// OCI Web Application Acceleration and Security Services
//

package waas

import (
	"github.com/oracle/oci-go-sdk/common"
)

// GoodBot The good bot settings. Good bots provides a list of bots which are managed by known providers.
type GoodBot struct {

	// The unique key for the bot.
	Key *string `mandatory:"true" json:"key"`

	// Enables or disables the bot.
	IsEnabled *bool `mandatory:"true" json:"isEnabled"`

	// The bot name.
	Name *string `mandatory:"false" json:"name"`

	// The description of the bot.
	Description *string `mandatory:"false" json:"description"`
}

func (m GoodBot) String() string {
	return common.PointerString(m)
}
