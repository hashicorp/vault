/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { encodePath } from 'vault/utils/path-encoding-helpers';
import ApplicationAdapter from '../application';

export default class PkiUrlsAdapter extends ApplicationAdapter {
  namespace = 'v1';

  _url(backend) {
    return `${this.buildURL()}/${encodePath(backend)}/config/urls`;
  }

  urlForCreateRecord(modelName, snapshot) {
    return this._url(snapshot.record.id);
  }
  urlForFindRecord(id) {
    return this._url(id);
  }
  urlForUpdateRecord(store, type, snapshot) {
    return this._url(snapshot.record.id);
  }
}
