/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';

export default class PreviewImage extends Component {
  @service theme;

  get imageSrc() {
    const base = this.args.message.authenticated
      ? '~/custom-messages-dashboard.png'
      : '~/custom-messages-login.png';
    if (this.theme.isDarkMode) {
      return base.replace(/(\.[^.]+)$/, '-dark$1');
    }
    return base;
  }
}
