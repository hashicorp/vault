import Ember from 'ember';
import { supportedSecretBackends } from 'vault/helpers/supported-secret-backends';
const { computed, Controller } = Ember;
const LINKED_BACKENDS = supportedSecretBackends();

export default Controller.extend({
  displayableBackends: computed.filterBy('model', 'shouldIncludeInList'),

  supportedBackends: computed('displayableBackends', 'displayableBackends.[]', function() {
    return (this.get('displayableBackends') || [])
      .filter(backend => LINKED_BACKENDS.includes(backend.get('type')))
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
