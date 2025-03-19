/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import AdapterError from '@ember-data/adapter/error';
import { set } from '@ember/object';
import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { MANAGED_AUTH_BACKENDS } from 'vault/utils/mountable-auth-methods';

export default Route.extend({
  store: service(),
  pathHelp: service('path-help'),

  model(params) {
    const { path } = params;
    return this.store.query('auth-method', {}).then((modelArray) => {
      const model = modelArray.find((m) => m.id === path);
      if (!model) {
        const error = new AdapterError();
        set(error, 'httpStatus', 404);
        throw error;
      }
      if (!MANAGED_AUTH_BACKENDS.includes(model.methodType)) {
        // do not fetch path-help for unmanaged auth types
        model.set('paths', {
          apiPath: model.apiPath,
          paths: [],
        });
        return model;
      }
      return this.pathHelp.getPaths(model.apiPath, path).then((paths) => {
        model.set('paths', paths);
        return model;
      });
    });
  },
});
