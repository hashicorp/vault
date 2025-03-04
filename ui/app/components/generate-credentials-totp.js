/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { task, timeout } from 'ember-concurrency';
import { service } from '@ember/service';

export default class GenerateCredentialsTotp extends Component {
  @tracked elapsedTime = 0;
  @tracked totpCode = null;
  @service store;

  title = 'Generate TOTP code';

  constructor() {
    super(...arguments);
    const { model } = this.args;
    this.generateTotpCode(model.backendPath, model.roleName);
    this.startTimer.perform();
  }

  get remainingTime() {
    const { model } = this.args;

    if (!model.totpCodePeriod) {
      return 0; // TODO improve this
    }

    return model.totpCodePeriod - this.elapsedTime;
  }

  @task({ restartable: true })
  *startTimer() {
    const { model } = this.args;
    if (model.totpCodePeriod) {
      while (this.elapsedTime <= model.totpCodePeriod) {
        yield timeout(1000);
        this.elapsedTime += 1;
      }

      if (this.elapsedTime > model.totpCodePeriod) {
        this.elapsedTime = 0;
        this.generateTotpCode(model.backendPath, model.roleName);
        this.startTimer.perform();
      }
    }
  }

  async generateTotpCode(backend, keyName) {
    // TODO improvement: refreshing does not currently result in a new code
    try {
      const totpCode = await this.store.adapterFor('totp-key').generateCode(backend, keyName);
      this.totpCode = totpCode.code;
    } catch (e) {
      // swallow error, non-essential data
      return;
    }
  }

  willDestroy() {
    super.willDestroy();
    this.startTimer.cancelAll();
  }
}

// TODO this isn't perfect and currently doesn't reset the code at zero nor when refreshing
