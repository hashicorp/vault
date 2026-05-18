/**
 * Copyright IBM Corp. 2026, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationSerializer from '../application';

export default class OidcAssignmentSerializer extends ApplicationSerializer {
  primaryKey = 'name';
}
