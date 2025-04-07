import Component from '@glimmer/component';
import { HdsIcon, HdsTooltipButton, HdsTextBody } from '@hashicorp/design-system-components/components';
import { precompileTemplate } from '@ember/template-compilation';
import { setComponentTemplate } from '@ember/component';

/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
class SSUReportingCounter extends Component {
  get count() {
    if (this.args.suffix) {
      return `${this.args.count} ${this.args.suffix}`;
    }
    return this.args.count;
  }
  get icon() {
    return this.args.icon || 'info';
  }
  static {
    setComponentTemplate(precompileTemplate("\n    <div ...attributes data-test-counter={{@title}} class=\"ssu-counter\">\n      <div class=\"ssu-counter__title-row\">\n        <HdsTextBody @weight=\"semibold\">{{@title}}\n          {{#if @tooltipMessage}}\n            <HdsTooltipButton data-test-counter-tooltip-button class=\"ssu-counter__title-row__tooltip\" @text={{@tooltipMessage}} aria-label=\"tooltip\" @isInline={{true}}>\n              <HdsIcon @name=\"help\" @isInline={{true}} />\n            </HdsTooltipButton>\n          {{/if}}\n        </HdsTextBody>\n      </div>\n\n      {{!-- Render count as a link if a link is provided --}}\n      {{#if @link}}\n        <a href={{@link}} class=\"ssu-counter__link\">\n          <HdsTextBody>{{this.count}} </HdsTextBody>\n        </a>\n      {{else}}\n        <HdsTextBody>{{this.count}}</HdsTextBody>\n      {{/if}}\n    </div>\n  ", {
      strictMode: true,
      scope: () => ({
        HdsTextBody,
        HdsTooltipButton,
        HdsIcon
      })
    }), this);
  }
}

export { SSUReportingCounter as default };
//# sourceMappingURL=counter.js.map
