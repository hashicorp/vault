/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';
import errorMessage from 'vault/utils/error-message';

/**
 * @module RoleDetailsPage
 * RoleDetailsPage component is a child component for create and edit role pages.
 *
 * @param {object} model - role model that contains role record and backend
 * @param {array} breadcrumbs - breadcrumbs as an array of objects that contain label and route
 */

export default class RoleDetailsPageComponent extends Component {
  @service router;
  @service flashMessages;

  get extraFields() {
    const fields = [];
    if (this.args.model.extraAnnotations) {
      fields.push({ label: 'Annotations', key: 'extraAnnotations' });
    }
    if (this.args.model.extraLabels) {
      fields.push({ label: 'Labels', key: 'extraLabels' });
    }
    return fields;
  }

  @action
  async delete() {
    try {
      await this.args.model.destroyRecord();
      this.router.transitionTo('vault.cluster.secrets.backend.kubernetes.roles');
    } catch (error) {
      const message = errorMessage(error, 'Unable to delete role. Please try again or contact support');
      this.flashMessages.danger(message);
    }
  }
}
