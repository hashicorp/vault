import Ember from 'ember';
import { supportedSecretBackends } from 'vault/helpers/supported-secret-backends';

const SUPPORTED_BACKENDS = supportedSecretBackends();

const { computed } = Ember;

export default Ember.Controller.extend({
  mountTypes: [
    { label: 'Active Directory', value: 'ad' },
    { label: 'AWS', value: 'aws' },
    { label: 'Consul', value: 'consul' },
    { label: 'Databases', value: 'database' },
    { label: 'Google Cloud', value: 'gcp' },
    { label: 'KV', value: 'kv' },
    { label: 'Nomad', value: 'nomad' },
    { label: 'PKI', value: 'pki' },
    { label: 'RabbitMQ', value: 'rabbitmq' },
    { label: 'SSH', value: 'ssh' },
    { label: 'Transit', value: 'transit' },
    { label: 'TOTP', value: 'totp' },
    { label: 'Cassandra', value: 'cassandra', deprecated: true },
    { label: 'MongoDB', value: 'mongodb', deprecated: true },
    { label: 'MSSQL', value: 'mssql', deprecated: true },
    { label: 'MySQL', value: 'mysql', deprecated: true },
    { label: 'PostgreSQL', value: 'postgresql', deprecated: true },
  ],

  selectedType: null,
  selectedPath: null,
  description: null,
  default_lease_ttl: null,
  max_lease_ttl: null,
  showConfig: false,
  local: false,
  sealWrap: false,
  version: 2,

  selection: computed('selectedType', function() {
    return this.get('mountTypes').findBy('value', this.get('selectedType'));
  }),

  flashMessages: Ember.inject.service(),

  reset() {
    const defaultBackend = this.get('mountTypes.firstObject.value');
    this.setProperties({
      selectedPath: defaultBackend,
      selectedType: defaultBackend,
      description: null,
      default_lease_ttl: null,
      max_lease_ttl: null,
      local: false,
      showConfig: false,
      sealWrap: false,
      version: 2,
    });
  },

  init() {
    this._super(...arguments);
    this.reset();
  },

  actions: {
    onTypeChange(val) {
      const { selectedPath, selectedType } = this.getProperties('selectedPath', 'selectedType');
      this.set('selectedType', val);
      if (selectedPath === selectedType) {
        this.set('selectedPath', val);
      }
    },

    toggleShowConfig() {
      this.toggleProperty('showConfig');
    },

    mountBackend() {
      const {
        selectedPath: path,
        selectedType: type,
        description,
        default_lease_ttl,
        local,
        max_lease_ttl,
        sealWrap,
        version,
      } = this.getProperties(
        'selectedPath',
        'selectedType',
        'description',
        'default_lease_ttl',
        'local',
        'max_lease_ttl',
        'sealWrap',
        'version'
      );
      const currentModel = this.get('model');
      if (currentModel && currentModel.rollbackAttributes) {
        currentModel.rollbackAttributes();
      }
      let attrs = {
        path,
        type,
        description,
        local,
        sealWrap,
      };

      if (this.get('showConfig')) {
        attrs.config = {
          defaultLeaseTtl: default_lease_ttl,
          maxLeaseTtl: max_lease_ttl,
        };
      }

      if (type === 'kv') {
        attrs.options = {
          version,
        };
      }

      const model = this.store.createRecord('secret-engine', attrs);

      this.set('model', model);
      model.save().then(() => {
        this.reset();
        let transition;
        if (SUPPORTED_BACKENDS.includes(type)) {
          transition = this.transitionToRoute('vault.cluster.secrets.backend.index', path);
        } else {
          transition = this.transitionToRoute('vault.cluster.secrets.backends');
        }
        transition.followRedirects().then(() => {
          this.get('flashMessages').success(`Successfully mounted '${type}' at '${path}'!`);
        });
      });
    },
  },
});
