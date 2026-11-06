/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Route from '@ember/routing/route';
import ListRoute from 'vault/mixins/list-route';
import { paginate } from 'core/utils/paginate-list';

export default Route.extend(ListRoute, {
  api: service(),

  async model() {
    const { urls } = this.modelFor('vault.cluster.access.method.item');
    const { page, pageFilter } = this.paramsFor(this.routeName);
    try {
      const response = await this.api.request.get(`${urls.list}?list=true`);
      const { data } = await response.json();
      const ids = paginate(data?.keys ?? [], { page: Number(page) || 1, filter: pageFilter });
      const items = ids.map((id) => ({ id }));
      items.meta = ids.meta;
      return items;
    } catch (error) {
      const { status, response } = await this.api.parseError(error);
      // The application route sends Control Group errors to the approval page.
      if (response?.isControlGroupError) {
        throw response;
      }
      // Vault returns 404 when a mount has no items yet.
      if (status === 404) {
        return [];
      }
      throw error;
    }
  },

  actions: {
    willTransition() {
      window.scrollTo(0, 0);
      return true;
    },
    reload() {
      this.refresh();
    },
  },

  setupController(controller) {
    this._super(...arguments);
    const { itemType, pathInfo, urls } = this.modelFor('vault.cluster.access.method.item');
    controller.set('itemType', itemType);
    controller.set('methodModel', this.modelFor('vault.cluster.access.method'));
    controller.set('urls', urls);
    controller.set(
      'paths',
      pathInfo.paths.filter((path) => path.navigation && path.itemType.includes(itemType))
    );
  },
});
