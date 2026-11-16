/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { singularize } from 'ember-inflector';
import { service } from '@ember/service';
import Route from '@ember/routing/route';
import GeneratedItemForm from 'vault/forms/generated-item';

// Reads an existing item into a form; the path identifier comes from the URL since the response omits it.
export async function fetchGeneratedItemForm(route, id) {
  const { formMetadata, urls } = route.modelFor('vault.cluster.access.method.item');
  let response;
  try {
    response = await route.api.request.get(urls.read(id));
  } catch (error) {
    // The application route sends Control Group errors to the approval page.
    const { response: errorResponse } = await route.api.parseError(error);
    throw errorResponse?.isControlGroupError ? errorResponse : error;
  }
  const { data } = await response.json();
  const identifier = formMetadata.identifier ? { [formMetadata.identifier]: id } : {};
  return new GeneratedItemForm(formMetadata, { ...data, ...identifier }, { isNew: false });
}

export default Route.extend({
  api: service(),

  model(params) {
    return fetchGeneratedItemForm(this, params.item_id);
  },

  setupController(controller) {
    this._super(...arguments);
    const { itemType, urls } = this.modelFor('vault.cluster.access.method.item');
    const { path: method } = this.paramsFor('vault.cluster.access.method');
    controller.set('itemType', singularize(itemType));
    controller.set('method', method);
    controller.set('urls', urls);
  },
});
