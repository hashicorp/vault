/**
 * Copyright IBM Corp. 2016, 2026
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
    const identityType = this.modelFor('vault.cluster.access.identity');
    const itemType = identityType + '-alias';
    const tabs = TABS[itemType];

    if (!tabs.includes(section)) {
      const error = new Error(`Invalid section: ${section}`);
      error.httpStatus = 404;
      throw error;
    }
    const methodType = identityType === 'group' ? 'groupReadAliasById' : 'entityReadAliasById';
    const { data } = await this.api.identity[methodType](params.item_alias_id);

    return hash({
      model: { ...data, itemType, identityType },
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
