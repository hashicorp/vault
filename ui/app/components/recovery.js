/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import config from 'vault/config/environment';

export default class AlphabetEditComponent extends Component {
  get isNotProduction() {
    return config.environment !== 'production';
  }
}
