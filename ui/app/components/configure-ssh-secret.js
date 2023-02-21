import Component from '@glimmer/component';
import { action } from '@ember/object';

/**
 * @module ConfigureSshSecretComponent
 *
 * @example
 * ```js
 *  <ConfigureSshSecret 
      @model={{model}} 
      @configured={{configured}} 
      @saveConfig={{action "saveConfig"}} />
 * ```
 *
 * @param {string} model - ssh secret engine model
 * @param {Function} saveConfig - parent action which updates the configuration
 * 
 */
export default class ConfigureSshSecretComponent extends Component {
  @action
  saveConfig(data, event) {
    event.preventDefault();
    this.args.saveConfig(data);
  }
}
