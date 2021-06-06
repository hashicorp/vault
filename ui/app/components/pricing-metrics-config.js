/**
 * @module PricingMetricsConfig
 * PricingMetricsConfig components are used to show and edit the pricing metrics config information.
 *
 * @example
 * ```js
 * <PricingMetricsConfig @model={{model}} @mode="edit" />
 * ```
 * @param {object} model - model is the DS metrics/config model which should be passed in
 * @param {string} [mode=show] - mode is either show or edit. Show results in a table with the config, show has a form.
 */

import { computed } from '@ember/object';
import Component from '@ember/component';
import { inject as service } from '@ember/service';
import { task } from 'ember-concurrency';

export default Component.extend({
  router: service(),
  mode: 'show',
  model: null,

  error: null,
  modalOpen: false,
  infoRows: computed(function() {
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
      {
        label: 'Default display',
        helperText: 'The number of months we’ll display in the Vault usage dashboard by default.',
        valueKey: 'defaultReportMonths',
      },
    ];
  }),
  modalTitle: computed('model.enabled', function() {
    let content = 'Turn usage tracking off?';
    if (this.model.enabled === 'On') {
      content = 'Turn usage tracking on?';
    }
    return content;
  }),

  save: task(function*() {
    let model = this.model;
    try {
      yield model.save();
    } catch (err) {
      this.set('error', err.message);
      return;
    }
    this.router.transitionTo('vault.cluster.metrics.config');
  }).drop(),

  actions: {
    onSaveChanges: function(evt) {
      evt.preventDefault();
      const changed = this.model.changedAttributes();
      if (!changed.enabled) {
        this.save.perform();
        return;
      }
      this.set('modalOpen', true);
    },
  },
});
