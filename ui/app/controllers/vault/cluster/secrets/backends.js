import { filterBy } from '@ember/object/computed';
import { computed } from '@ember/object';
import Controller from '@ember/controller';
import { task } from 'ember-concurrency';
import { supportedSecretBackends } from 'vault/helpers/supported-secret-backends';
const LINKED_BACKENDS = supportedSecretBackends();

export default Controller.extend({
  displayableBackends: filterBy('model', 'shouldIncludeInList'),

  supportedBackends: computed('displayableBackends', 'displayableBackends.[]', function() {
    return (this.get('displayableBackends') || [])
      .filter(backend => LINKED_BACKENDS.includes(backend.get('engineType')))
      .sortBy('id');
  }),

  unsupportedBackends: computed(
    'displayableBackends',
    'displayableBackends.[]',
    'supportedBackends',
    'supportedBackends.[]',
    function() {
      return (this.get('displayableBackends') || [])
        .slice()
        .removeObjects(this.get('supportedBackends'))
        .sortBy('id');
    }
  ),

  disableEngine: task(function*(engine) {
    const { engineType, path } = engine;
    try {
      yield engine.destroyRecord();
      this.get('flashMessages').success(`The ${engineType} Secrets Engine at ${path} has been disabled.`);
    } catch (err) {
      this.get('flashMessages').danger(
        `There was an error disabling the ${engineType} Secrets Engine at ${path}: ${err.errors.join(' ')}.`
      );
    }
  }).drop(),
});
