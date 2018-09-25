import { inject as service } from '@ember/service';
import Component from '@ember/component';
import { computed } from '@ember/object';
import { task } from 'ember-concurrency';
import { humanize } from 'vault/helpers/humanize';

export default Component.extend({
  flashMessages: service(),
  'data-test-component': 'identity-edit-form',
  model: null,

  // 'create', 'edit', 'merge'
  mode: 'create',
  /*
   * @param Function
   * @public
   *
   * Optional param to call a function upon successfully saving an entity
   */
  onSave: () => {},

  cancelLink: computed('mode', 'model.identityType', function() {
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

  getMessage(model, isDelete = false) {
    let mode = this.get('mode');
    let typeDisplay = humanize([model.get('identityType')]);
    let action = isDelete ? 'deleted' : 'saved';
    if (mode === 'merge') {
      return 'Successfully merged entities';
    }
    if (model.get('id')) {
      return `Successfully ${action} ${typeDisplay} ${model.id}.`;
    }
    return `Successfully ${action} ${typeDisplay}.`;
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
    yield this.get('onSave')({ saveType: 'save', model });
  }).drop(),

  willDestroy() {
    let model = this.get('model');
    if ((model.get('isDirty') && !model.isDestroyed) || !model.isDestroying) {
      model.rollbackAttributes();
    }
  },

  actions: {
    deleteItem(model) {
      let message = this.getMessage(model, true);
      let flash = this.get('flashMessages');
      model.destroyRecord().then(() => {
        flash.success(message);
        return this.get('onSave')({ saveType: 'delete', model });
      });
    },
  },
});
