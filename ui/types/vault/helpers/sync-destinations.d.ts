/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { DestinationType } from 'sync/utils/constants';
import { DestinationName, DestinationRoleTypeOption } from 'vault/sync';

export interface SyncDestination {
  name: DestinationName;
  type: DestinationType;
  // Widened to string so resolved (dark-mode) icon names are assignable.
  // The static data always uses the `-color` variants; resolveGlyph may strip
  // that suffix at runtime, producing a plain icon name like `aws`.
  icon: string;
  category: 'cloud' | 'dev-tools';
  maskedParams: Array<string>;
  readonlyParams: Array<string>;
  defaultValues: object;
  roleTypeOptions?: Array<DestinationRoleTypeOption>;
}
