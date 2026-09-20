/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { hash } from 'rsvp';
import Route from '@ember/routing/route';
import { TABS } from 'vault/helpers/tabs-for-identity-show';
import { service } from '@ember/service';
import { fetchRelatedIdentityItems } from 'vault/utils/identity-helpers';

const RELATED_IDENTITIES = {
  groups: [
    {
      identityType: 'group',
      modelKey: 'groups',
      idKeys: ['group_ids', 'direct_group_ids', 'inherited_group_ids'],
    },
  ],
};

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

    const [response, canAddAlias] = await Promise.all([
      this.api.identity.entityReadById(params.item_id),
      this.capabilities.for('groupAlias').canCreate,
    ]);
    const relatedIdentities = await fetchRelatedIdentityItems({
      api: this.api,
      model: response.data,
      relations: RELATED_IDENTITIES[section] || [],
    });

    return hash({
      model: {
        ...response.data,
        ...relatedIdentities,
        identityType: 'entity',
        canAddAlias: canAddAlias || false,
      },
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
