import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { runTask } from 'ember-lifeline';
import { HdsTextDisplay, HdsCardContainer } from '@hashicorp/design-system-components/components';
import TitleRow from './base/title-row.js';
import { precompileTemplate } from '@ember/template-compilation';
import { setComponentTemplate } from '@ember/component';
import { g, i, n } from 'decorator-transforms/runtime';

/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
class GlobalLease extends Component {
  static {
    g(this.prototype, "animatedPercentage", [tracked], function () {
      return 0;
    });
  }
  #animatedPercentage = (i(this, "animatedPercentage"), undefined);
  static {
    g(this.prototype, "displayPercentage", [tracked], function () {
      return 0;
    });
  }
  #displayPercentage = (i(this, "displayPercentage"), undefined);
  static {
    g(this.prototype, "initialState", [tracked], function () {
      return true;
    });
  }
  #initialState = (i(this, "initialState"), undefined);
  constructor(owner, args) {
    super(owner, args);
    runTask(this, () => this.startAnimation(), 50);
  }
  get actualPercentage() {
    return Math.round(Math.min(this.args.count / this.args.quota * 100, 100));
  }
  get progressFillClass() {
    if (this.initialState) {
      return 'ssu-global-lease__progress-fill--initial';
    }
    if (this.actualPercentage < 50) {
      return 'ssu-global-lease__progress-fill--low';
    }
    if (this.actualPercentage < 100) {
      return 'ssu-global-lease__progress-fill--medium';
    }
    return 'ssu-global-lease__progress-fill--high';
  }
  get formattedCount() {
    const formatter = new Intl.NumberFormat('en-US', {
      notation: 'compact',
      compactDisplay: 'short'
    });
    const count = formatter.format(this.args.count);
    const total = formatter.format(this.args.quota);
    return `${count} / ${total}`;
  }
  startAnimation() {
    this.animatedPercentage = 0;
    this.displayPercentage = 0;
    this.initialState = true;
    runTask(this, () => {
      this.animatePercentageText();
      this.animatedPercentage = this.actualPercentage;
      this.initialState = false;
    }, 100);
  }
  static {
    n(this.prototype, "startAnimation", [action]);
  }
  animatePercentageText() {
    const targetPercentage = this.actualPercentage;
    const duration = 1000;
    const steps = 20;
    const stepDuration = duration / steps;
    let currentStep = 0;
    const updatePercentage = () => {
      currentStep++;
      this.displayPercentage = Math.round(currentStep / steps * targetPercentage);
      if (currentStep < steps) {
        runTask(this, updatePercentage, stepDuration);
      } else {
        this.displayPercentage = targetPercentage;
      }
    };
    updatePercentage();
  }
  static {
    n(this.prototype, "animatePercentageText", [action]);
  }
  static {
    setComponentTemplate(precompileTemplate("\n    <HdsCardContainer data-test-global-lease @hasBorder={{true}} class=\"ssu-global-lease\">\n      <TitleRow @title=\"Global lease count quota\" @description=\"Snapshot of global lease count quota consumption\" @linkText=\"Documentation\" @linkIcon=\"docs-link\" @linkUrl=\"https://developer.hashicorp.com/vault/docs/enterprise/lease-count-quotas\" />\n\n      <HdsTextDisplay class=\"ssu-global-lease__percentage-text\">{{this.displayPercentage}}%</HdsTextDisplay>\n\n      <div class=\"ssu-global-lease__progress-wrapper\">\n        <div class=\"ssu-global-lease__progress-bar\">\n          <div class=\"ssu-global-lease__progress-fill {{this.progressFillClass}}\" {{!-- template-lint-disable no-inline-styles style-concatenation --}} style=\"width: {{this.animatedPercentage}}%;\"></div>\n        </div>\n        <span class=\"ssu-global-lease__count-text\">\n          <HdsTextDisplay @size=\"400\" @weight=\"medium\">\n            {{this.formattedCount}}\n          </HdsTextDisplay>\n        </span>\n      </div>\n\n    </HdsCardContainer>\n  ", {
      strictMode: true,
      scope: () => ({
        HdsCardContainer,
        TitleRow,
        HdsTextDisplay
      })
    }), this);
  }
}

export { GlobalLease as default };
//# sourceMappingURL=global-lease.js.map
