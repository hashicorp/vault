/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { hash } from 'rsvp';
import Route from '@ember/routing/route';
import { TABS } from 'vault/helpers/tabs-for-identity-show';
import { service } from '@ember/service';

export default class IdentityAliasesShowRoute extends Route {
  @service api;

  async model(params) {
    const { section } = params;
    const tabs = TABS['entity-alias'];

    if (!tabs.includes(section)) {
      const error = new Error(`Invalid section: ${section}`);
      error.httpStatus = 404;
      throw error;
    }

    const { data } = await this.api.identity.entityReadAliasById(params.item_alias_id);

    return hash({
      model: { ...data, itemType: 'entity-alias', identityType: 'entity' },
      section,
    });
  }

  setupController(controller, resolvedModel) {
    const { model, section } = resolvedModel;
    controller.setProperties({
      model,
      section,
    });
  }
}
