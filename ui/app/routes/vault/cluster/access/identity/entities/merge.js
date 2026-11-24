/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import MergeEntitiesForm from 'vault/forms/identity/merge-entities';

export default class VaultClusterAccessIdentityMergeRoute extends Route {
  model() {
    return { form: new MergeEntitiesForm({}, { isNew: true }) };
  }
}
