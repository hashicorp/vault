/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationAdapter from './application';

export default class NamespaceAdapter extends ApplicationAdapter {
  pathForType() {
    return 'namespaces';
  }
  urlForFindAll(modelName, snapshot) {
    if (snapshot.adapterOptions && snapshot.adapterOptions.forUser) {
      return `/${this.urlPrefix()}/internal/ui/namespaces`;
    }
    return `/${this.urlPrefix()}/namespaces?list=true`;
  }

  urlForCreateRecord(modelName, snapshot) {
    const id = snapshot.attr('path');
    return this.buildURL(modelName, id);
  }

  createRecord(store, type, snapshot) {
    const id = snapshot.attr('path');
    return super.createRecord(...arguments).then(() => {
      return { id };
    });
  }

  findAll(store, type, sinceToken, snapshot) {
    if (snapshot.adapterOptions && typeof snapshot.adapterOptions.namespace !== 'undefined') {
      return this.ajax(this.urlForFindAll('namespace', snapshot), 'GET', {
        namespace: snapshot.adapterOptions.namespace,
      });
    }
    return super.findAll(...arguments);
  }
  query() {
    return this.ajax(`/${this.urlPrefix()}/namespaces?list=true`);
  }
}
