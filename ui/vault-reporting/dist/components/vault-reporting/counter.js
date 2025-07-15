import Component from '@glimmer/component';
import { HdsLinkInline, HdsIcon, HdsTooltipButton, HdsTextBody } from '@hashicorp/design-system-components/components';
import { precompileTemplate } from '@ember/template-compilation';
import { setComponentTemplate } from '@ember/component';

/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
class SSUReportingCounter extends Component {
  get shouldShowEmptyState() {
    return this.args.count === 0 && this.args.emptyText;
  }
  get count() {
    if (this.shouldShowEmptyState) {
      return this.args.emptyText;
    }
    if (this.args.suffix) {
      return `${this.args.count} ${this.args.suffix}`;
    }
    return this.args.count;
  }
  get icon() {
    return this.args.icon || 'info';
  }
  get link() {
    if (this.shouldShowEmptyState && this.args.emptyLink) {
      return this.args.emptyLink;
    }
    return this.args.link;
  }
  static {
    setComponentTemplate(precompileTemplate("\n    <div ...attributes data-test-vault-reporting-counter={{@title}} class=\"ssu-counter\" aria-label=\"{{@title}} {{this.count}}\">\n      <div class=\"ssu-counter__title-row\">\n        <HdsTextBody @weight=\"semibold\" @size=\"200\" @color=\"primary\">{{@title}}\n          {{#if @tooltipMessage}}\n            <HdsTooltipButton data-test-vault-reporting-counter-tooltip-button class=\"ssu-counter__title-row__tooltip\" @text={{@tooltipMessage}} aria-label=\"Tooltip for {{@title}}\" @isInline={{true}}>\n              <HdsIcon @name=\"help\" @isInline={{true}} />\n            </HdsTooltipButton>\n          {{/if}}\n        </HdsTextBody>\n      </div>\n\n      <HdsTextBody>\n        {{#if this.link}}\n          <HdsLinkInline @href={{this.link}} @color=\"secondary\" class=\"ssu-counter__link\" target=\"_self\">{{this.count}}\n          </HdsLinkInline>\n        {{else}}\n          {{this.count}}\n        {{/if}}\n      </HdsTextBody>\n    </div>\n  ", {
      strictMode: true,
      scope: () => ({
        HdsTextBody,
        HdsTooltipButton,
        HdsIcon,
        HdsLinkInline
      })
    }), this);
  }
}

export { SSUReportingCounter as default };
//# sourceMappingURL=counter.js.map
