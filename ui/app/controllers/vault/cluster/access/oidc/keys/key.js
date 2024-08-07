/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';

export default class OidcKeyController extends Controller {
  @service router;
  @tracked isEditRoute;

  constructor() {
    super(...arguments);
    this.router.on('routeDidChange', ({ targetName }) => {
      return (this.isEditRoute = targetName.includes('edit') ? true : false);
    });
  }

  get showHeader() {
    // hide header when rendering the edit form
    return !this.isEditRoute;
  }
}
