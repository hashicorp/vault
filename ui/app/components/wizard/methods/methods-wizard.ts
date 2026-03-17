/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import { action } from '@ember/object';
import Component from '@glimmer/component';

import type WizardService from 'vault/services/wizard';

interface Args {
  isIntroModal: boolean;
  onRefresh: CallableFunction;
}

export const WIZARD_ID = 'auth-methods';
export default class WizardMethodsWizardComponent extends Component<Args> {
  @service declare readonly wizard: WizardService;

  wizardId = WIZARD_ID;

  @action
  async onDismiss() {
    this.wizard.dismiss(this.wizardId);
    await this.args.onRefresh();
  }
}
