import Component from '@glimmer/component';
import { action } from '@ember/object';

/**
 * @module ConfigureSshSecretComponent
 *
 * @example
 * ```js
 * <ConfigureSshSecret
 *    @model={{this.model}}
 *    @configured={{this.configured}}
 *    @saveConfig={{action "saveConfig"}}
 *    @loading={{this.loading}}
 *  />
 * ```
 *
 * @param {string} model - ssh secret engine model
 * @param {Function} saveConfig - parent action which updates the configuration
 * @param {boolean} loading - property in parent that updates depending on status of parent's action
 *
 */
export default class ConfigureSshSecretComponent extends Component {
  @action
  delete() {
    this.args.saveConfig({ delete: true });
  }

  @action
  saveConfig(event) {
    event.preventDefault();
    this.args.saveConfig({ delete: false });
  }
}
