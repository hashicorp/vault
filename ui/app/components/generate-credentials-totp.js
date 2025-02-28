/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { task, timeout } from 'ember-concurrency';

export default class GenerateCredentialsTotp extends Component {
  @tracked elapsedTime = 0;

  title = 'Generate TOTP code';

  constructor() {
    super(...arguments);
    this.startTimer.perform();
  }

  get remainingTime() {
    const { model } = this.args;

    if (!model.period) {
      return 0; // TODO improve this
    }

    return model.period - this.elapsedTime;
  }

  @task({ restartable: true })
  *startTimer() {
    const { model } = this.args;
    if (model.period) {
      while (this.elapsedTime <= model.period) {
        yield timeout(1000);
        this.elapsedTime += 1;
      }

      if (this.elapsedTime > model.period) {
        this.elapsedTime = 0;
        this.startTimer.perform();
      }
    }
  }

  willDestroy() {
    super.willDestroy();
    this.startTimer.cancelAll();
  }
}

// TODO this isn't perfect and currently doesn't reset the code at zero nor when refreshing
