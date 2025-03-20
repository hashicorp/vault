import Component from '@glimmer/component';
import PluginCard from './plugin-card.js';
import { precompileTemplate } from '@ember/template-compilation';
import { setComponentTemplate } from '@ember/component';

/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
class Plugins extends Component {
  get plugins() {
    return this.args.plugins;
  }
  static {
    setComponentTemplate(precompileTemplate("\n    {{#each @plugins as |plugin|}}\n      <PluginCard @plugin={{plugin}} />\n    {{/each}}\n  ", {
      strictMode: true,
      scope: () => ({
        PluginCard
      })
    }), this);
  }
}

export { Plugins as default };
//# sourceMappingURL=plugins.js.map
