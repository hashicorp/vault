/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { hash } from 'rsvp';

import type LdapLibraryModel from 'vault/models/ldap/library';

export default class LdapLibraryRoute extends Route {
  model() {
    const model = this.modelFor('libraries.library') as LdapLibraryModel;
    return hash({
      library: model,
      statuses: model.fetchStatus(),
    });
  }
}
