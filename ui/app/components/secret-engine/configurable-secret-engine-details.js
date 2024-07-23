/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Component from '@glimmer/component';
import errorMessage from 'vault/utils/error-message';
import { tracked } from '@glimmer/tracking';
// ARG TODO add documentation
export default class ConfigurableSecretEngineDetails extends Component {
  @service store;
  @tracked configModel = null;
  @tracked configError = null;

  constructor() {
    super(...arguments);
    const { model } = this.args;
    // Currently two secret engines that return configuration data and that can be configured by the user on the ui: aws, and ssh.
    if (model.type === 'aws') {
      this.fetchAwsRootConfig(model.id);
    }
    if (model.type === 'ssh') {
      this.fetchSshCaConfig(model.id);
    }
  }

  async fetchAwsRootConfig(backend) {
    try {
      this.configModel = await this.store.queryRecord('aws/root-config', { backend });
    } catch (e) {
      this.configError = errorMessage(e);
    }
  }

  async fetchSshCaConfig(backend) {
    try {
      this.configModel = await this.store.queryRecord('ssh/ca-config', { backend });
    } catch (e) {
      this.configError = errorMessage(e);
    }
  }

  get typeDisplay() {
    // TODO will eventually handle GCP and Azure.
    // Did not use capitalize helper because some are all caps and some only title case.
    const { type } = this.args.model;
    switch (type) {
      case 'aws':
        return 'AWS';
      case 'ssh':
        return 'SSH';
      default:
        return type;
    }
  }
}
