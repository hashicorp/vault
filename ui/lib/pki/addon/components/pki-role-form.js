import Component from '@glimmer/component';
import { task } from 'ember-concurrency';
import { tracked } from '@glimmer/tracking';

/**
 * @module PkiRoleForm
//  * ARG TODO
 * PkiRoleForm components are used to...
 *
 * @example
 * ```js
 * <PkiRoleForm @requiredParam={requiredParam} @optionalParam={optionalParam} @param1={{param1}}/>
 * ```
 * @param {object} requiredParam - requiredParam is...
 * @param {string} [optionalParam] - optionalParam is...
 * @param {string} [param1=defaultValue] - param1 is...
 */

export default class PkiRoleForm extends Component {
  @tracked errorBanner;

  @task
  *save(event) {
    event.preventDefault();
    // ARG TODO see client assignment-form
    try {
      const { isValid, state } = this.args.model.validate();
      this.modelValidations = isValid ? null : state;
      if (isValid) {
        const { isNew, name } = this.args.model;
        yield this.args.model.save();
        this.flashMessages.success(`Successfully ${isNew ? 'created' : 'updated'} the role ${name}.`);
      }
    } catch (error) {
      const message = error.errors ? error.errors.join('. ') : error.message;
      this.errorBanner = message;
    }
  }
}
