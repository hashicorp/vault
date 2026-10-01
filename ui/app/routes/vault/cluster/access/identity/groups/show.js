/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { hash } from 'rsvp';
import Route from '@ember/routing/route';
import { TABS } from 'vault/helpers/tabs-for-identity-show';
import { service } from '@ember/service';
import { attachAliasCapabilities } from 'vault/utils/identity-helpers';

export default class IdentityShowRoute extends Route {
  @service router;
  @service api;
  @service capabilities;

  async model(params) {
    const { section } = params;
    const tabs = TABS['group'];

    if (!tabs.includes(section)) {
      const error = new Error(`Invalid section: ${section}`);
      error.httpStatus = 404;
      throw error;
    }

    const { data } = await this.api.identity.groupReadById(params.item_id);
    // Adding an alias is an update operation on identity/group-alias (no existence check).
    const { canUpdate: canAddAlias } = await this.capabilities.for('groupAlias');
    const { canUpdate: canEditGroup } = await this.capabilities.for('identityCapabilities', {
      identityType: 'group',
      id: params.item_id,
    });
    const alias =
      section === 'aliases'
        ? await attachAliasCapabilities({
            aliases: data.alias,
            identityType: 'group',
            capabilities: this.capabilities,
          })
        : data.alias;

    return hash({
      model: { ...data, alias, identityType: 'group', canAddAlias, canEditGroup },
      section,
    });
  }

  afterModel(resolvedModel) {
    const { section, model } = resolvedModel;

    if (model?.type === 'internal' && section === 'aliases') {
      return this.router.transitionTo('vault.cluster.access.identity.groups.show', model.id, 'details');
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
