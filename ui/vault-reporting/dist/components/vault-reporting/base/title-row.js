import Component from '@glimmer/component';
import { HdsTextBody, HdsLinkStandalone, HdsTextDisplay } from '@hashicorp/design-system-components/components';
import { on } from '@ember/modifier';
import { service } from '@ember/service';
import { precompileTemplate } from '@ember/template-compilation';
import { setComponentTemplate } from '@ember/component';
import { g, i } from 'decorator-transforms/runtime';

/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
class TitleRow extends Component {
  static {
    g(this.prototype, "reportingAnalytics", [service]);
  }
  #reportingAnalytics = (i(this, "reportingAnalytics"), void 0);
  get hasLink() {
    return this.args.linkUrl;
  }
  get linkText() {
    return this.args.linkText || 'View all';
  }
  get linkUrl() {
    return this.args.linkUrl || '#';
  }
  get linkIcon() {
    return this.args.linkIcon || 'arrow-right';
  }
  get linkTarget() {
    return this.args.linkTarget || '_self';
  }
  handleLinkClick = () => {
    this.reportingAnalytics.trackEvent(`card_link`, {
      card: this.args.title,
      link: this.linkText,
      target: this.linkTarget
    });
  };
  static {
    setComponentTemplate(precompileTemplate("\n    <div class=\"ssu-title-row\" data-test-vault-reporting-dashboard-card-title-row>\n      <div class=\"ssu-title-row__container\">\n        <HdsTextDisplay data-test-vault-reporting-dashboard-card-title @size=\"300\">\n          {{@title}}\n        </HdsTextDisplay>\n\n        {{#if this.hasLink}}\n          <HdsLinkStandalone data-test-vault-reporting-dashboard-card-title-link class=\"ssu-title-row__container__link\" @text={{this.linkText}} @href={{this.linkUrl}} @icon={{this.linkIcon}} target={{this.linkTarget}} @iconPosition=\"trailing\" {{on \"click\" this.handleLinkClick}} />\n        {{/if}}\n      </div>\n\n      {{#if @description}}\n        <HdsTextBody class=\"ssu-title-row__description\" data-test-vault-reporting-dashboard-card-description>\n          {{@description}}\n        </HdsTextBody>\n      {{/if}}\n    </div>\n  ", {
      strictMode: true,
      scope: () => ({
        HdsTextDisplay,
        HdsLinkStandalone,
        on,
        HdsTextBody
      })
    }), this);
  }
}

export { TitleRow as default };
//# sourceMappingURL=title-row.js.map
