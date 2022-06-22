import Component from '@glimmer/component';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';

/**
 * @module OidcClientForm
 * OidcClientForm components are used to...
 *
 * @example
 * ```js
 * <OidcClientForm @requiredParam={requiredParam} @optionalParam={optionalParam} @param1={{param1}}/>
 * ```
 * @param {object} requiredParam - requiredParam is...
 * @param {string} [optionalParam] - optionalParam is...
 * @param {string} [param1=defaultValue] - param1 is...
 */

export default class OidcClientForm extends Component {
  @service flashMessages;
  @tracked showMoreOptions = false;
  @tracked radioCardGroupValue = 'allow_all';

  @action
  radioCardSelect(name, value) {
    console.log(name);
    console.log(value);
  }
}
