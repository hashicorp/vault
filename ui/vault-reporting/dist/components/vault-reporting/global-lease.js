import Component from '@glimmer/component';
import { HdsApplicationState, HdsAlert, HdsTextDisplay, HdsCardContainer } from '@hashicorp/design-system-components/components';
import TitleRow from './base/title-row.js';
import cssCustomProperty from '../../modifiers/css-custom-property.js';
import { htmlSafe } from '@ember/template';
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
    if (this.percentage >= 100) {
      return 'ssu-global-lease__progress-fill--exceeded';
    }
    return '';
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
      return htmlSafe('Total number of active <a class="hds-link-inline--color-secondary" href="https://developer.hashicorp.com/vault/docs/concepts/lease" target="_blank" data-test-vault-reporting-global-lease-description-link>leases</a> for this quota.');
    }
  }
  get linkUrl() {
    if (this.hasData) {
      return 'https://developer.hashicorp.com/vault/tutorials/operations/resource-quotas#global-default-lease-count-quota';
    }
  }
  get alert() {
    if (this.percentage >= 100) {
      return {
        color: 'warning',
        description: 'Global lease quota limit reached. If lease creation is blocked, reduce usage or increase the limit.'
      };
    }
    if (this.percentage >= 95) {
      return {
        color: 'neutral',
        description: 'Approaching quota limit. Reduce usage or increase the lease limit to avoid blocking new leases.'
      };
    }
  }
  static {
    setComponentTemplate(precompileTemplate("\n    <HdsCardContainer data-test-vault-reporting-global-lease @hasBorder={{true}} class=\"ssu-global-lease\" {{cssCustomProperty \"--vault-reporting-global-lease-percentage\" this.percentageString}} ...attributes>\n      <TitleRow @title=\"Global lease count quota\" @description={{this.description}} @linkText=\"Documentation\" @linkIcon=\"docs-link\" @linkUrl={{this.linkUrl}} @linkTarget=\"_blank\" />\n      {{#if this.hasData}}\n        <HdsTextDisplay @size=\"300\" @weight=\"medium\" data-test-vault-reporting-global-lease-percentage-text>{{this.percentage}}%</HdsTextDisplay>\n\n        {{#if this.alert}}\n          <HdsAlert data-test-vault-reporting-global-lease-alert class=\"ssu-global-lease__alert\" @type=\"compact\" @color={{this.alert.color}} as |A|>\n            <A.Description>{{this.alert.description}}</A.Description>\n          </HdsAlert>\n        {{/if}}\n\n        <div class=\"ssu-global-lease__progress-wrapper\">\n          <div class=\"ssu-global-lease__progress-bar\">\n            <div class=\"ssu-global-lease__progress-fill {{this.progressFillClass}}\" data-test-vault-reporting-global-lease-fill></div>\n          </div>\n          <span>\n            <HdsTextDisplay @size=\"200\" @weight=\"semibold\" data-test-vault-reporting-global-lease-count-text>\n              {{this.formattedCount}}\n            </HdsTextDisplay>\n          </span>\n        </div>\n      {{else}}\n\n        <HdsApplicationState data-test-vault-reporting-global-lease-empty-state class=\"ssu-global-lease__empty-state\" as |A|>\n          {{#if (has-block \"empty\")}}\n            {{yield A to=\"empty\"}}\n          {{else}}\n            <A.Body data-test-vault-reporting-global-lease-empty-state-description @text=\"Lease quotas enforce limits on active secrets and tokens. It's recommended to enable this to protect stability for this Vault cluster.\" />\n\n            <A.Footer as |F|>\n              <F.LinkStandalone data-test-vault-reporting-global-lease-empty-state-link @icon=\"docs-link\" @iconPosition=\"trailing\" @text=\"Global lease count quota\" @href=\"https://developer.hashicorp.com/vault/tutorials/operations/resource-quotas#global-default-lease-count-quota\" target=\"_blank\" />\n            </A.Footer>\n          {{/if}}\n        </HdsApplicationState>\n      {{/if}}\n\n    </HdsCardContainer>\n  ", {
      strictMode: true,
      scope: () => ({
        HdsCardContainer,
        cssCustomProperty,
        TitleRow,
        HdsTextDisplay,
        HdsAlert,
        HdsApplicationState
      })
    }), this);
  }
}

export { GlobalLease as default };
//# sourceMappingURL=global-lease.js.map
