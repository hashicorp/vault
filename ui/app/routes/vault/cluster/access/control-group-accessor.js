/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Route from '@ember/routing/route';

export default class ControlGroupAccessorRoute extends Route {
  @service api;
  @service capabilities;
  @service version;

  async beforeModel() {
    await this.version.fetchFeatures();
    return super.beforeModel(...arguments);
  }

  async model(params) {
    if (!this.version.hasControlGroups) return null;
    const accessor = params.accessor;
    const [{ data }, { canUpdate: canAuthorize }, { canUpdate: canConfigure }] = await Promise.all([
      this.api.sys.systemWriteControlGroupRequest({ accessor }),
      this.capabilities.for('controlGroupAuthorize'),
      this.capabilities.for('controlGroupConfig'),
    ]);
    const requestEntity = data?.request_entity;
    return {
      id: accessor,
      approved: Boolean(data?.approved),
      requestPath: data?.request_path,
      // an empty request_entity means the requester is unknown
      requestEntity: requestEntity && Object.keys(requestEntity).length ? requestEntity : null,
      authorizations: (data?.authorizations || []).map((authz) => ({
        id: authz.entity_id,
        name: authz.entity_name,
      })),
      canAuthorize,
      canConfigure,
    };
  }

  // The requesting entity and approvers link to the entity details page only when the user can read them
  async afterModel(model) {
    const ids = new Set((model?.authorizations || []).map((authz) => authz.id));
    if (model?.requestEntity?.id) ids.add(model.requestEntity.id);

    const pathById = {};
    for (const id of ids) {
      pathById[id] = this.capabilities.pathFor('identityCapabilities', { identityType: 'entity', id });
    }
    const paths = Object.values(pathById);
    const capabilities = paths.length ? await this.capabilities.fetch(paths) : {};
    this.canReadEntity = {};
    for (const id of ids) {
      this.canReadEntity[id] = Boolean(capabilities[pathById[id]]?.canRead);
    }
  }

  setupController(controller) {
    super.setupController(...arguments);
    controller.set('canReadEntity', this.canReadEntity);
  }
}
