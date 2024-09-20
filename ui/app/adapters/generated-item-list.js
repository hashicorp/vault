/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationAdapter from './application';
import { task } from 'ember-concurrency';
import { service } from '@ember/service';
import { sanitizePath } from 'core/utils/sanitize-path';
import { encodePath } from 'vault/utils/path-encoding-helpers';

export default class GeneratedItemListAdapter extends ApplicationAdapter {
  @service store;
  namespace = 'v1';

  // these items are set within getNewAdapter in path-help service
  apiPath = '';
  paths = {};

  getDynamicApiPath(id) {
    const result = this.store.peekRecord('auth-method', id);
    return result.apiPath;
  }

  fetchByQuery = task(async (store, query, isList) => {
    const { id } = query;
    const payload = {};
    if (isList) {
      payload.list = true;
    }
    const path = isList ? this.getDynamicApiPath(id) : '';

    const resp = await this.ajax(this.urlForItem(id, isList, path), 'GET', { data: payload });
    const data = {
      id,
      method: id,
    };
    return { ...resp, ...data };
  });

  query(store, type, query) {
    return this.fetchByQuery.perform(store, query, true);
  }

  queryRecord(store, type, query) {
    return this.fetchByQuery.perform(store, query);
  }

  urlForItem(id, isList, dynamicApiPath) {
    const itemType = sanitizePath(this.paths.getPath);
    let url;
    id = encodePath(id);
    // the apiPath changes when you switch between routes but the apiPath variable does not unless the model is reloaded
    // overwrite apiPath if dynamicApiPath exist.
    // dynamicApiPath comes from the model->adapter
    let apiPath = this.apiPath;
    if (dynamicApiPath) {
      apiPath = dynamicApiPath;
    }
    // isList indicates whether we are viewing the list page
    // of a top-level item such as userpass
    if (isList) {
      url = `${this.buildURL()}/${apiPath}${itemType}/`;
    } else {
      // build the URL for the show page of a nested item
      // such as a userpass group
      url = `${this.buildURL()}/${apiPath}${itemType}/${id}`;
    }

    return url;
  }

  urlForQueryRecord(id, modelName) {
    return this.urlForItem(id, modelName);
  }

  urlForUpdateRecord(id) {
    const itemType = this.paths.createPath.slice(1, this.paths.createPath.indexOf('{') - 1);
    return `${this.buildURL()}/${this.apiPath}${itemType}/${id}`;
  }

  urlForCreateRecord(modelType, snapshot) {
    const id = snapshot.record.mutableId; // computed property that returns either id or private settable _id value
    const path = this.paths.createPath.slice(1, this.paths.createPath.indexOf('{') - 1);
    return `${this.buildURL()}/${this.apiPath}${path}/${id}`;
  }

  urlForDeleteRecord(id) {
    const path = this.paths.deletePath.slice(1, this.paths.deletePath.indexOf('{') - 1);
    return `${this.buildURL()}/${this.apiPath}${path}/${id}`;
  }

  createRecord(store, type, snapshot) {
    return super.createRecord(...arguments).then((response) => {
      // if the server does not return an id and one has not been set on the model we need to set it manually from the mutableId value
      if (!response?.id && !snapshot.record.id) {
        snapshot.record.id = snapshot.record.mutableId;
        snapshot.id = snapshot.record.id;
      }
      return response;
    });
  }
}
