/**
 * @module ClientsConfig
 * ClientsConfig components are used to show and edit the client count config information.
 *
 * @example
 * ```js
 * <Clients::Config @model={{model}} @mode="edit" />
 * ```
 * @param {object} model - model is the DS clients/config model which should be passed in
 * @param {string} [mode=show] - mode is either show or edit. Show results in a table with the config, show has a form.
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { task } from 'ember-concurrency';

export default class ConfigComponent extends Component {
  @service router;
  @tracked mode = 'show';
  @tracked modalOpen = false;
  error = null;

  get infoRows() {
    return [
      {
        label: 'Usage data collection',
        helperText: 'Enable or disable collecting data to track clients.',
        valueKey: 'enabled',
      },
      {
        label: 'Retention period',
        helperText: 'The number of months of activity logs to maintain for client tracking.',
        valueKey: 'retentionMonths',
      },
    ];
  }

  get modalTitle() {
    let content = 'Turn usage tracking off?';
    if (this.args.model && this.args.model.enabled === 'On') {
      content = 'Turn usage tracking on?';
    }
    return content;
  }

  @(task(function* () {
    try {
      yield this.args.model.save();
    } catch (err) {
      this.error = err.message;
      return;
    }
    this.router.transitionTo('vault.cluster.clients.config');
  }).drop())
  save;

  @action
  updateBooleanValue(attr, value) {
    const valueToSet = value === true ? attr.options.trueValue : attr.options.falseValue;
    this.args.model[attr.name] = valueToSet;
  }

  @action
  onSaveChanges(evt) {
    evt.preventDefault();
    const changed = this.args.model.changedAttributes();
    if (!changed.enabled) {
      this.save.perform();
      return;
    }
    this.modalOpen = true;
  }
}
