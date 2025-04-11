import Component from '@glimmer/component';
import { HdsTextBody, HdsLinkStandalone, HdsTextDisplay } from '@hashicorp/design-system-components/components';
import { precompileTemplate } from '@ember/template-compilation';
import { setComponentTemplate } from '@ember/component';

/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
class TitleRow extends Component {
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
  static {
    setComponentTemplate(precompileTemplate("\n    <div class=\"ssu-title-row\" data-test-dashboard-card-title-row>\n      <div class=\"ssu-title-row__container\">\n        <HdsTextDisplay data-test-dashboard-card-title @size=\"300\">\n          {{@title}}\n        </HdsTextDisplay>\n\n        {{#if this.hasLink}}\n          <HdsLinkStandalone data-test-dashboard-card-title-link class=\"ssu-title-row__container__link\" @text={{this.linkText}} @href={{this.linkUrl}} @icon={{this.linkIcon}} target={{this.linkTarget}} @iconPosition=\"trailing\" />\n        {{/if}}\n      </div>\n\n      {{#if @description}}\n        <HdsTextBody class=\"ssu-title-row__description\" data-test-dashboard-card-description>\n          {{@description}}\n        </HdsTextBody>\n      {{/if}}\n    </div>\n  ", {
      strictMode: true,
      scope: () => ({
        HdsTextDisplay,
        HdsLinkStandalone,
        HdsTextBody
      })
    }), this);
  }
}

export { TitleRow as default };
//# sourceMappingURL=title-row.js.map
