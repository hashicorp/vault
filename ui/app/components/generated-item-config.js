import { inject as service } from '@ember/service';
import Component from '@ember/component';
import { task } from 'ember-concurrency';
import DS from 'ember-data';

/**
 * @module GeneratedItemConfig
 * The `GeneratedItemConfig` is the form to configure generated items related to mounts (e.g. groups, roles, users)
 *
 * @example
 * ```js
 * <GeneratedItemConfig @model={{model}} @mode={{mode}} />
 * ```
 *
 * @property model=null {DS.Model} - The corresponding item model that is being configured.
 * @property mode {String} - which config mode to use. either `edit` or `create`
 *
 */

const ItemConfig = Component.extend({
  model: null,

  flashMessages: service(),
  router: service(),
  wizard: service(),
  saveModel: task(function*() {
    try {
      yield this.model.save();
    } catch (err) {
      // AdapterErrors are handled by the error-message component
      // in the form
      if (err instanceof DS.AdapterError === false) {
        throw err;
      }
      return;
    }
    this.router.transitionTo('vault.cluster.access.method.item').followRedirects();
    this.flashMessages.success('The configuration was saved successfully.');
  }).withTestWaiter(),
});

ItemConfig.reopenClass({
  positionalParams: ['model'],
});

export default ItemConfig;
