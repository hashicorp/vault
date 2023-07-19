/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

/**
 * @module KvPageSecretDetails
 * KvPageSecretDetails shows detail of kv secrets.
 *
 * @param {array} model - An array of models generated form kv/metadata query.
 */

export default class KvPageSecretDetails extends Component {
  @tracked showJsonView = false;

  @action
  toggleJsonView() {
    this.showJsonView = !this.showJsonView;
  }
}
