/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { hash } from 'rsvp';
import Route from '@ember/routing/route';
import { TABS } from 'vault/helpers/tabs-for-identity-show';
import { service } from '@ember/service';
import { fetchIdentityItems } from 'vault/utils/identity-helpers';

const RELATED_IDENTITIES = {
  entity: {
    groups: [
      {
        identityType: 'group',
        modelKey: 'groups',
        idKeys: ['group_ids', 'direct_group_ids', 'inherited_group_ids'],
      },
    ],
  },
  group: {
    members: [
      { identityType: 'group', modelKey: 'groups', idKeys: ['member_group_ids'] },
      { identityType: 'entity', modelKey: 'entities', idKeys: ['member_entity_ids'] },
    ],
    'parent-groups': [{ identityType: 'group', modelKey: 'groups', idKeys: ['parent_group_ids'] }],
  },
};

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
    const [response, canAddAlias] = await Promise.all([
      this.api.identity[methodType](params.item_id),
      this.capabilities.for('groupAlias').canCreate,
    ]);
    const relatedIdentities = await this.fetchRelatedIdentities(itemType, section, response.data);

    return hash({
      model: {
        ...response.data,
        ...relatedIdentities,
        identityType: itemType,
        canAddAlias: canAddAlias || false,
      },
      section,
    });
  }

  async fetchRelatedIdentities(itemType, section, model) {
    const identities = (RELATED_IDENTITIES[itemType]?.[section] || []).filter(({ idKeys }) =>
      idKeys.some((key) => model[key]?.length)
    );
    const results = await Promise.allSettled(
      identities.map(({ identityType }) => fetchIdentityItems({ identityType, api: this.api }))
    );

    return identities.reduce((related, { modelKey }, index) => {
      const result = results[index];
      related[modelKey] = result.status === 'fulfilled' ? result.value : [];
      return related;
    }, {});
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
