import Component from '@glimmer/component';
import { HdsTextDisplay, HdsCardContainer } from '@hashicorp/design-system-components/components';
import TitleRow from './base/title-row.js';
import cssCustomProperty from '../../modifiers/css-custom-property.js';
import { precompileTemplate } from '@ember/template-compilation';
import { setComponentTemplate } from '@ember/component';

/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
class GlobalLease extends Component {
  get percentage() {
    return Math.round(Math.min(this.args.count / this.args.quota * 100, 100));
  }
  get progressFillClass() {
    if (this.percentage < 50) {
      return 'ssu-global-lease__progress-fill--low';
    }
    if (this.percentage < 100) {
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
  get percentageString() {
    return `${this.percentage}%`;
  }
  static {
    setComponentTemplate(precompileTemplate("\n    <HdsCardContainer data-test-global-lease @hasBorder={{true}} class=\"ssu-global-lease\" {{cssCustomProperty \"--vault-reporting-global-lease-percentage\" this.percentageString}} ...attributes>\n      <TitleRow @title=\"Global lease count quota\" @description=\"Snapshot of global lease count quota consumption\" @linkText=\"Documentation\" @linkIcon=\"docs-link\" @linkUrl=\"https://developer.hashicorp.com/vault/docs/enterprise/lease-count-quotas\" @linkTarget=\"_blank\" />\n\n      <HdsTextDisplay class=\"ssu-global-lease__percentage-text\">{{this.percentage}}%</HdsTextDisplay>\n\n      <div class=\"ssu-global-lease__progress-wrapper\">\n        <div class=\"ssu-global-lease__progress-bar\">\n          <div class=\"ssu-global-lease__progress-fill {{this.progressFillClass}}\" data-test-global-lease-fill></div>\n        </div>\n        <span class=\"ssu-global-lease__count-text\">\n          <HdsTextDisplay @size=\"400\" @weight=\"medium\">\n            {{this.formattedCount}}\n          </HdsTextDisplay>\n        </span>\n      </div>\n\n    </HdsCardContainer>\n  ", {
      strictMode: true,
      scope: () => ({
        HdsCardContainer,
        cssCustomProperty,
        TitleRow,
        HdsTextDisplay
      })
    }), this);
  }
}

export { GlobalLease as default };
//# sourceMappingURL=global-lease.js.map
