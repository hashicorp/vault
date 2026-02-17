/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { action } from '@ember/object';

import type ApiService from 'vault/services/api';
import type FlashMessageService from 'vault/services/flash-messages';
import type RouterService from '@ember/routing/router-service';
import type WizardService from 'vault/services/wizard';

interface Args {
  isIntroModal: boolean;
  onRefresh: CallableFunction;
}

export const WIZARD_ID = 'secret-engines';

export default class WizardSecretEnginesWizardComponent extends Component<Args> {
  @service declare readonly api: ApiService;
  @service declare readonly router: RouterService;
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly wizard: WizardService;

  wizardId = WIZARD_ID;

  @action
  onDismiss() {
    this.wizard.dismiss(this.wizardId);
    this.args.onRefresh();
  }

  @action
  onIntroChange(visible: boolean) {
    this.wizard.setIntroVisible(this.wizardId, visible);
  }
}
