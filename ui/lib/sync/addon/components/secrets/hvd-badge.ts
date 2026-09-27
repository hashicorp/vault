/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';

import type FlagsService from 'vault/services/flags';

/**
 * @module HvdBadge
 * Renders the "Standard feature" badge on HVD managed clusters, and nothing otherwise.
 * Sync page headers all surface this badge, so it lives here rather than being repeated
 * alongside every `Page::Header` invocation.
 *
 * @example
 * <Page::Header @title="Secrets Sync">
 *   <:badges><Secrets::HvdBadge /></:badges>
 * </Page::Header>
 */
export default class SyncHvdBadgeComponent extends Component {
  @service declare readonly flags: FlagsService;
}
