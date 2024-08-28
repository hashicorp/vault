/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ApplicationAdapter from './application';
import { v4 as uuidv4 } from 'uuid';

export default class AwsRootConfig extends ApplicationAdapter {
  namespace = 'v1';

  queryRecord() {
    return this.ajax(`${this.buildURL()}/identity/oidc/config`, 'GET').then((resp) => {
      return {
        ...resp,
        id: uuidv4(), // generate a random id for ember data
      };
    });
  }

  // ARG TODO return to this
  // createOrUpdate(store, type, snapshot) {
  //   const serializer = store.serializerFor(type.modelName);
  //   const data = serializer.serialize(snapshot);
  //   const backend = snapshot.record.backend;
  //   return this.ajax(`${this.buildURL()}/${backend}/config/root`, 'POST', { data }).then((resp) => {
  //     // ember data requires an id on the response
  //     return {
  //       ...resp,
  //       id: backend,
  //     };
  //   });
  // }

  // createRecord() {
  //   return this.createOrUpdate(...arguments);
  // }

  // updateRecord() {
  //   return this.createOrUpdate(...arguments);
  // }
}
