/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';

/**
 * @module Overview
 * OverviewPage component is a child component to overview kubernetes secrets engine.
 *
 * @param {boolean} promptConfig - whether or not to display config cta
 * @param {object} backend - backend model that contains kubernetes configuration
 * @param {array} roles - array of roles
 * @param {array} breadcrumbs - breadcrumbs as an array of objects that contain label and route
 */

export default class OverviewPageComponent extends Component {
  @service router;

  @tracked selectedRole = null;
  @tracked roleOptions = [];

  constructor() {
    super(...arguments);
    this.roleOptions = this.args.roles.map((role) => {
      return { name: role.name, id: role.name };
    });
  }

  @action
  selectRole([roleName]) {
    this.selectedRole = roleName;
  }

  @action
  generateCredential() {
    this.router.transitionTo(
      'vault.cluster.secrets.backend.kubernetes.roles.role.credentials',
      this.selectedRole
    );
  }
}
