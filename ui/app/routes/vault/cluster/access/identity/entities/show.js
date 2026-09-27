/**
 * Copyright IBM Corp. 2016, 2025
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
    const tabs = TABS['entity'];

    if (!tabs.includes(section)) {
      const error = new Error(`Invalid section: ${section}`);
      error.httpStatus = 404;
      throw error;
    }

    const { data } = await this.api.identity.entityReadById(params.item_id);
    const canAddAlias = (await this.capabilities.for('groupAlias').canCreate) || false;
    const aliases =
      section === 'aliases'
        ? await attachAliasCapabilities({
            aliases: data.aliases,
            identityType: 'entity',
            capabilities: this.capabilities,
          })
        : data.aliases;

    return hash({
      model: { ...data, aliases, identityType: 'entity', canAddAlias },
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
