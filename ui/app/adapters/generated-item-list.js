/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationAdapter from './application';
import { service } from '@ember/service';
import { sanitizePath } from 'core/utils/sanitize-path';
import { encodePath } from 'vault/utils/path-encoding-helpers';
import { tracked } from '@glimmer/tracking';

export default class GeneratedItemListAdapter extends ApplicationAdapter {
  @service store;
  namespace = 'v1';

  // these items are set by calling getNewAdapter in the path-help service.
  @tracked apiPath = '';
  paths = {};

  // These are the paths used for the adapter actions
  get getPath() {
    return this.paths.getPath || '';
  }
  get createPath() {
    return this.paths.createPath || '';
  }
  get deletePath() {
    return this.paths.deletePath || '';
  }

  getDynamicApiPath(id) {
    const result = this.store.peekRecord('auth-method', id);
    this.apiPath = result.apiPath;
    return result.apiPath;
  }

  async fetchByQuery(store, query, isList) {
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
  }

  query(store, type, query) {
    return this.fetchByQuery(store, query, true);
  }

  queryRecord(store, type, query) {
    return this.fetchByQuery(store, query);
  }

  urlForItem(id, isList, dynamicApiPath) {
    const itemType = sanitizePath(this.getPath);
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
    const itemType = this.createPath.slice(1, this.createPath.indexOf('{') - 1);
    return `${this.buildURL()}/${this.apiPath}${itemType}/${id}`;
  }

  urlForCreateRecord(modelType, snapshot) {
    const id = snapshot.record.mutableId; // computed property that returns either id or private settable _id value
    const path = this.createPath.slice(1, this.createPath.indexOf('{') - 1);
    return `${this.buildURL()}/${this.apiPath}${path}/${id}`;
  }

  urlForDeleteRecord(id) {
    const path = this.deletePath.slice(1, this.deletePath.indexOf('{') - 1);
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
