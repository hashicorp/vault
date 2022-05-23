import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { dropTask } from 'ember-concurrency';
import { action } from '@ember/object';

/**
 * @module MountAccessorSelect
 * The MountAccessorSelect component is used to selectDrop down mount options.
 *
 * @example
 * ```js
 * <MountAccessorSelect @value={this.aliasMountAccessor} @onChange={this.onChange} />
 * ```
 * @param {string} value - the selected value.
 * @param {function} onChange - the parent function that handles when a new value is selected.
 * @param {boolean} [showAccessor] - whether or not you should show the value or the more detailed accessor off the class
 * @param {boolean} [noDefault] - whether or not there is a default value.
 * @param {boolean} [filterToken] -whether or not you should filter out type "token".
 * @param {string} [name] - name on the label
 * @param {string} [label] - label above the select input
 * @param {string} [helpText] - text shown in tooltip.
 */

export default class MountAccessorSelect extends Component {
  @service store;

  filterToken = false;
  noDefault = false;
  value = '';

  constructor() {
    super(...arguments);
    this.authMethods.perform();
  }

  @dropTask *authMethods() {
    let methods = yield this.store.findAll('auth-method');
    if (!this.args.value && !this.args.noDefault) {
      // ARG can't set value as args.
      this.value = methods.get('firstObject.accessor');
      this.args.onChange(this.value);
    }
    return methods;
  }

  @action change(event) {
    this.args.onChange(event.target.value);
  }
}
