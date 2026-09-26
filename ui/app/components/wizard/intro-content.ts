/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';

import type ThemeService from 'vault/services/theme';

interface IntroContentSignature {
  Args: {
    imageSrc: string;
  };
}

export default class IntroContent extends Component<IntroContentSignature> {
  @service declare readonly theme: ThemeService;

  get introPageImgSrc(): string {
    if (!this.theme.isDarkMode) {
      return this.args.imageSrc;
    }

    return this.args.imageSrc.replace(/(\.[^.]+)$/, '-dark$1');
  }
}
