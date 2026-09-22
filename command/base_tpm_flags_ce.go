// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: MPL-2.0

//go:build !enterprise

package command

func (c *BaseCommand) setupTPMFlags(f *FlagSet) {}
