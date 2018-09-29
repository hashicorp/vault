import { filterBy } from '@ember/object/computed';
import { computed } from '@ember/object';
import Controller from '@ember/controller';
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
});
