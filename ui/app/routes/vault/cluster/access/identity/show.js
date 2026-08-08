/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { hash } from 'rsvp';
import Route from '@ember/routing/route';
import { TABS } from 'vault/helpers/tabs-for-identity-show';
import { service } from '@ember/service';

export default class IdentityShowRoute extends Route {
  @service router;
  @service api;
  @service capabilities;

  async model(params) {
    const { section } = params;
    const itemType = this.modelFor('vault.cluster.access.identity');
    const tabs = TABS[itemType];

    if (!tabs.includes(section)) {
      const error = new Error(`Invalid section: ${section}`);
      error.httpStatus = 404;
      throw error;
    }

    const methodType = itemType === 'entity' ? 'entityReadById' : 'groupReadById';
    const { data } = await this.api.identity[methodType](params.item_id);
    const canAddAlias = (await this.capabilities.for('groupAlias').canCreate) || false;

    return hash({
      model: { ...data, identityType: itemType, canAddAlias },
      section,
    });
  }

  afterModel(resolvedModel) {
    const { section, model } = resolvedModel;

    if (model.identityType === 'group' && model?.type === 'internal' && section === 'aliases') {
      return this.router.transitionTo('vault.cluster.access.identity.show', model.id, 'details');
    }
  }

  setupController(controller, resolvedModel) {
    const { model, section } = resolvedModel;
    controller.setProperties({
      model,
      section,
    });
  }
}
