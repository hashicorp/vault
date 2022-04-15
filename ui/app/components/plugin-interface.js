import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';

/**
 * @module PluginInterface
 * PluginInterface components are used to display dynamic plugin UI's served by the plugin itself
 *
 * @example
 * ```js
 * <PluginInterface @tabs={arrayOfTabs} @wrappedToken={token} />
 * ```
 * @param {array} tabs - tabs is the tab objects to display
 * @param {string} wrappedToken - wrappedToken is passed in the URL to the iframeUrl
 */

export default class PluginInterface extends Component {
  @tracked activeTab;

  constructor() {
    super(...arguments);
    console.log('constructing', this.args.tabs);
    this.activeTab = this.args.tabs[0] || null;
  }

  @action
  selectTab(tab) {
    console.log({ tab });
    this.activeTab = tab;
  }
}
