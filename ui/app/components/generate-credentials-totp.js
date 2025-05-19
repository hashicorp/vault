/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { task, timeout } from 'ember-concurrency';
import { service } from '@ember/service';
import { action } from '@ember/object';

const SHOW_ROUTE = 'vault.cluster.secrets.backend.show';

export default class GenerateCredentialsTotp extends Component {
  @tracked elapsedTime = 0;
  @tracked totpCode = null;
  @service store;
  @service router;

  title = 'Generate TOTP code';

  constructor() {
    super(...arguments);
    this.startTimer.perform();
  }

  get remainingTime() {
    const { totpCodePeriod } = this.args;

    return totpCodePeriod - this.elapsedTime;
  }

  @task({ restartable: true })
  *startTimer() {
    const { backendPath, keyName, totpCodePeriod } = this.args;
    this.generateTotpCode(backendPath, keyName);
    while (this.elapsedTime <= totpCodePeriod) {
      yield timeout(1000);
      this.elapsedTime += 1;
    }

    if (this.elapsedTime > totpCodePeriod) {
      this.elapsedTime = 0;
      this.generateTotpCode(backendPath, keyName);
      this.startTimer.perform();
    }
  }

  async generateTotpCode(backend, keyName) {
    // refreshing will generate a new code if the period has expired.
    try {
      const totpCode = await this.store.adapterFor('totp-key').generateCode(backend, keyName);
      this.totpCode = totpCode.code;
    } catch (e) {
      // swallow error, non-essential data
      return;
    }
  }

  @action redirectPreviousPage() {
    const { backRoute, keyName } = this.args;
    if (backRoute === SHOW_ROUTE) {
      this.router.transitionTo(this.args.backRoute, keyName);
    } else {
      this.router.transitionTo(this.args.backRoute);
    }
  }
}
