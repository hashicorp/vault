import Component from '@glimmer/component';
import { HdsBadge, HdsTextBody, HdsTextDisplay, HdsCardContainer } from '@hashicorp/design-system-components/components';
import { precompileTemplate } from '@ember/template-compilation';
import { setComponentTemplate } from '@ember/component';

/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
class PluginCard extends Component {
  get isEnterprisePlugin() {
    return this.args.plugin.tags === 'enterprise';
  }
  get pluginPublishDate() {
    return new Date(this.args.plugin.publishDate).toDateString();
  }
  static {
    setComponentTemplate(precompileTemplate("\n    <HdsCardContainer class=\"card-container\" @hasBorder={{true}} data-test-card-plugin-container>\n      <div class=\"is-flex\">\n        <div class=\"is-flex-column is-flex-1\">\n          <div class=\"has-bottom-padding-s\">\n            <HdsTextDisplay>{{@plugin.pluginName}}</HdsTextDisplay>\n            <HdsTextBody @tag=\"p\" @color=\"foreground-faint\" @size=\"200\">by\n              {{@plugin.author}}</HdsTextBody>\n            <HdsTextBody @tag=\"p\" @color=\"foreground-faint\" @size=\"200\">{{@plugin.pluginType}}</HdsTextBody>\n          </div>\n          <HdsTextBody @tag=\"p\">{{@plugin.description}}t</HdsTextBody>\n        </div>\n        <div class=\"is-flex-column is-flex-1\">\n          <HdsTextBody class=\"has-text-right\" @tag=\"p\">\n            <span>\n              32k downloads\n              <HdsBadge @text={{@plugin.official.tags}} @icon={{@plugin.official.author}} />\n              <HdsBadge @text={{@plugin.tags}} @color={{if this.isEnterprisePlugin \"highlight\"}} />\n            </span>\n          </HdsTextBody>\n          <HdsTextBody class=\"has-text-right\" @tag=\"p\">{{@plugin.pluginVersion}}\n            published\n            {{this.pluginPublishDate}}\n          </HdsTextBody>\n          <HdsTextBody class=\"has-text-right\" @tag=\"p\">\n            Compatible with Vault\n            {{@plugin.pluginVersion}}\n          </HdsTextBody>\n        </div>\n      </div>\n    </HdsCardContainer>\n  ", {
      strictMode: true,
      scope: () => ({
        HdsCardContainer,
        HdsTextDisplay,
        HdsTextBody,
        HdsBadge
      })
    }), this);
  }
}

export { PluginCard as default };
//# sourceMappingURL=plugin-card.js.map
