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
    saveConfig(options = { delete: false }) {
      const isDelete = options.delete;
      if (this.model.type === 'ssh') {
        this.set('loading', true);
        this.model
          .saveCA({ isDelete })
          .then(() => {
            this.send('refreshRoute');
            this.set('configured', !isDelete);
            if (isDelete) {
              this.flashMessages.success('SSH Certificate Authority Configuration deleted!');
            } else {
              this.flashMessages.success('SSH Certificate Authority Configuration saved!');
            }
          })
          .catch((error) => {
            const errorMessage = error.errors ? error.errors.join('. ') : error;
            this.flashMessages.danger(errorMessage);
          })
          .finally(() => {
            this.set('loading', false);
          });
      }
    },

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
