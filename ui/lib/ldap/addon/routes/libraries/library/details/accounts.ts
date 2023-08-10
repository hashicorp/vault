/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
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
