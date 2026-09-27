/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { syncDestinations } from 'core/helpers/sync-destinations';
import { resolveGlyph } from 'core/utils/all-engines-metadata';

import type { SyncDestination } from 'vault/helpers/sync-destinations';

export default class SelectTypeComponent extends Component {
  /**
   * Returns the full list of sync destinations with their icon resolved for the
   * current theme. Icons ending in `-color` are stripped to their monochrome
   * variant in dark mode (e.g. `aws-color` → `aws`).
   */
  get destinations(): Array<SyncDestination> {
    return syncDestinations().map((d) => ({ ...d, icon: resolveGlyph(d.icon) ?? d.icon }));
  }
}
