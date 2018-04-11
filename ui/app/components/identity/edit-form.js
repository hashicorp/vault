import Ember from 'ember';
import { task } from 'ember-concurrency';
import { humanize } from 'vault/helpers/humanize';

const { computed } = Ember;
export default Ember.Component.extend({
  model: null,
  mode: 'create',
  /*
   * @param Function
   * @public
   *
   * Optional param to call a function upon successfully mounting a backend
   *
   */
  onSave: () => {},

  cancelLink: computed('mode', 'model', function() {
    let { model, mode } = this.getProperties('model', 'mode');
    let key = `${mode}-${model.get('identityType')}`;
    let routes = {
      'create-entity': 'vault.cluster.access.identity',
      'edit-entity': 'vault.cluster.access.identity.show',
      'merge-entity-merge': 'vault.cluster.access.identity',
      'create-entity-alias': 'vault.cluster.access.identity.aliases',
      'edit-entity-alias': 'vault.cluster.access.identity.aliases.show',
      'create-group': 'vault.cluster.access.identity',
      'edit-group': 'vault.cluster.access.identity.show',
      'create-group-alias': 'vault.cluster.access.identity.aliases',
      'edit-group-alias': 'vault.cluster.access.identity.aliases.show',
    };

    return routes[key];
  }),

  getMessage(model) {
    let mode = this.get('mode');
    let typeDisplay = humanize([model.get('identityType')]);
    if (mode === 'merge') {
      return 'Successfully merged entities';
    }
    if (model.get('id')) {
      return `Successfully saved ${typeDisplay} ${model.id}.`;
    }
    return `Successfully saved ${typeDisplay}.`;
  },

  save: task(function*() {
    let model = this.get('model');
    let message = this.getMessage(model);

    try {
      yield model.save();
    } catch (err) {
      // err will display via model state
      return;
    }
    this.get('flashMessages').success(message);
    yield this.get('onSave')(model);
  }).drop(),

  willDestroy() {
    let model = this.get('model');
    if ((model.get('isDirty') && !model.isDestroyed) || !model.isDestroying) {
      model.rollbackAttributes();
    }
  },
});
