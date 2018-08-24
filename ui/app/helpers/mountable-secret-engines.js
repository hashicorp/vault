import Ember from 'ember';

const MOUNTABLE_SECRET_ENGINES = [
  {
    displayName: 'Active Directory',
    value: 'ad',
    type: 'ad',
  },
  {
    displayName: 'AWS',
    value: 'aws',
    type: 'aws',
  },
  {
    displayName: 'Consul',
    value: 'consul',
    type: 'consul',
  },
  {
    displayName: 'Databases',
    value: 'database',
    type: 'database',
  },
  {
    displayName: 'Google Cloud',
    value: 'gcp',
    type: 'gcp',
  },
  {
    displayName: 'KV',
    value: 'kv',
    type: 'kv',
  },
  {
    displayName: 'Nomad',
    value: 'nomad',
    type: 'nomad',
  },
  {
    displayName: 'PKI',
    value: 'pki',
    type: 'pki',
  },
  {
    displayName: 'RabbitMQ',
    value: 'rabbitmq',
    type: 'rabbitmq',
  },
  {
    displayName: 'SSH',
    value: 'ssh',
    type: 'ssh',
  },
  {
    displayName: 'Transit',
    value: 'transit',
    type: 'transit',
  },
  {
    displayName: 'TOTP',
    value: 'totp',
    type: 'totp',
  },
];

export function engines() {
  return MOUNTABLE_SECRET_ENGINES;
}

export default Ember.Helper.helper(engines);
