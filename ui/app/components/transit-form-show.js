/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { action } from '@ember/object';

export default class TransitFormShow extends Component {
  @service store;
  @service router;
  @service flashMessages;

  @action async rotateKey() {
    const { backend, id } = this.args.key;
    try {
      await this.store.adapterFor('transit-key').keyAction('rotate', { backend, id });
      this.flashMessages.success('Key rotated.');
      // must refresh to see the updated versions, a model refresh does not trigger the change.
      await this.router.refresh();
    } catch (e) {
      this.flashMessages.danger(e.errors);
    }
  }
}
