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
  get linkColor() {
    if (this.shouldShowEmptyState) {
      return 'secondary';
    }
    return 'primary';
  }
  static {
    setComponentTemplate(precompileTemplate("\n    <div ...attributes data-test-counter={{@title}} class=\"ssu-counter\">\n      <div class=\"ssu-counter__title-row\">\n        <HdsTextBody @weight=\"semibold\">{{@title}}\n          {{#if @tooltipMessage}}\n            <HdsTooltipButton data-test-counter-tooltip-button class=\"ssu-counter__title-row__tooltip\" @text={{@tooltipMessage}} aria-label=\"tooltip\" @isInline={{true}}>\n              <HdsIcon @name=\"help\" @isInline={{true}} />\n            </HdsTooltipButton>\n          {{/if}}\n        </HdsTextBody>\n      </div>\n\n      <HdsTextBody>\n        {{#if @link}}\n          <HdsLinkInline @href={{@link}} @color={{this.linkColor}} class=\"ssu-counter__link\">{{this.count}}\n          </HdsLinkInline>\n        {{else}}\n          {{this.count}}\n        {{/if}}\n      </HdsTextBody>\n    </div>\n  ", {
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
