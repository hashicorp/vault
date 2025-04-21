import Component from '@glimmer/component';
import { HdsApplicationState, HdsTextDisplay, HdsCardContainer } from '@hashicorp/design-system-components/components';
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
    const {
      count = 0,
      quota = 0
    } = this.args;
    return Math.round(Math.min(count / quota * 100, 100));
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
    const {
      count = 0,
      quota = 0
    } = this.args;
    const formattedCount = formatter.format(count);
    const formattedTotal = formatter.format(quota);
    return `${formattedCount} / ${formattedTotal}`;
  }
  get percentageString() {
    return `${this.percentage}%`;
  }
  get hasData() {
    return this.args.quota && typeof this.args.quota === 'number';
  }
  get description() {
    if (this.hasData) {
      return 'Snapshot of global lease count quota consumption';
    }
  }
  get linkUrl() {
    if (this.hasData) {
      return 'https://developer.hashicorp.com/vault/docs/enterprise/lease-count-quotas';
    }
  }
  static {
    setComponentTemplate(precompileTemplate("\n    <HdsCardContainer data-test-global-lease @hasBorder={{true}} class=\"ssu-global-lease\" {{cssCustomProperty \"--vault-reporting-global-lease-percentage\" this.percentageString}} ...attributes>\n      <TitleRow @title=\"Global lease count quota\" @description={{this.description}} @linkText=\"Documentation\" @linkIcon=\"docs-link\" @linkUrl={{this.linkUrl}} @linkTarget=\"_blank\" />\n      {{#if this.hasData}}\n        <HdsTextDisplay class=\"ssu-global-lease__percentage-text\">{{this.percentage}}%</HdsTextDisplay>\n\n        <div class=\"ssu-global-lease__progress-wrapper\">\n          <div class=\"ssu-global-lease__progress-bar\">\n            <div class=\"ssu-global-lease__progress-fill {{this.progressFillClass}}\" data-test-global-lease-fill></div>\n          </div>\n          <span class=\"ssu-global-lease__count-text\">\n            <HdsTextDisplay @size=\"400\" @weight=\"medium\">\n              {{this.formattedCount}}\n            </HdsTextDisplay>\n          </span>\n        </div>\n      {{else}}\n\n        <HdsApplicationState data-test-global-lease-empty-state class=\"ssu-global-lease__empty-state\" as |A|>\n          {{#if (has-block \"empty\")}}\n            {{yield A to=\"empty\"}}\n          {{else}}\n            <A.Header data-test-global-lease-empty-state-title @title=\"None enforced\" />\n            <A.Body data-test-global-lease-empty-state-description @text=\"Global lease count quota is disabled. Enable it to manage active leases.\" />\n\n            <A.Footer as |F|>\n              <F.LinkStandalone data-test-global-lease-empty-state-link @icon=\"docs-link\" @iconPosition=\"trailing\" @text=\"Global lease count quota\" @href=\"https://developer.hashicorp.com/vault/tutorials/operations/resource-quotas#global-default-lease-count-quota\" target=\"_blank\" />\n            </A.Footer>\n          {{/if}}\n        </HdsApplicationState>\n      {{/if}}\n\n    </HdsCardContainer>\n  ", {
      strictMode: true,
      scope: () => ({
        HdsCardContainer,
        cssCustomProperty,
        TitleRow,
        HdsTextDisplay,
        HdsApplicationState
      })
    }), this);
  }
}

export { GlobalLease as default };
//# sourceMappingURL=global-lease.js.map
