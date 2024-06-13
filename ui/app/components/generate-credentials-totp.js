/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { later, cancel } from '@ember/runloop';
import timestamp from 'core/utils/timestamp';
import { getUnixTime } from 'date-fns';

export default class GenerateCredentialsTotp extends Component {
  title = "Generate TOTP code";

  @tracked elapsedTime = 0;
  nextTick = null;

  get remainingTime() {
    const { model } = this.args;
    return model.period - this.elapsedTime;
  }

  @action
  cancelTimer() {
    cancel(this.nextTick);
  }

  @action
  startTimer() {
    this.nextTick = later(
      this,
      function () {
        const { model } = this.args;
        this.elapsedTime = getUnixTime(timestamp.now()) % model.period;
        this.startTimer();
      },
      1000
    );
  }
}
