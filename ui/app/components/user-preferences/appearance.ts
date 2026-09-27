/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';

import { ThemeChoice } from 'vault/services/theme';
import type ThemeService from 'vault/services/theme';

export default class Appearance extends Component {
  @service declare readonly theme: ThemeService;

  get currentTheme() {
    return this.theme.theme;
  }

  @action
  selectTheme(choice: ThemeChoice) {
    this.theme.setTheme(choice);
  }
}
