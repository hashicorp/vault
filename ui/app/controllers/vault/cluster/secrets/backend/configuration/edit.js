/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { isPresent } from '@ember/utils';
import { service } from '@ember/service';
import Controller from '@ember/controller';

const CONFIG_ATTRS = {
  // ssh
  configured: false,

  // aws root config
  iamEndpoint: null,
  stsEndpoint: null,
  accessKey: null,
  secretKey: null,
  region: '',
};

export default Controller.extend(CONFIG_ATTRS, {
  queryParams: ['tab'],
  tab: '',
  flashMessages: service(),
  loading: false,
  reset() {
    this.model.rollbackAttributes();
    this.setProperties(CONFIG_ATTRS);
  },
  actions: {
    save(method, data) {
      this.set('loading', true);
      const hasData = Object.keys(data).some((key) => {
        return isPresent(data[key]);
      });
      if (!hasData) {
        return;
      }
      this.model
        .save({
          adapterOptions: {
            adapterMethod: method,
            data,
          },
        })
        .then(() => {
          this.reset();
          this.flashMessages.success('The backend configuration saved successfully!');
        })
        .finally(() => {
          this.set('loading', false);
        });
    },
  },
});
