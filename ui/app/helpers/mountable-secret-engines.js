import { helper as buildHelper } from '@ember/component/helper';

const MOUNTABLE_SECRET_ENGINES = [
  {
    displayName: 'Active Directory',
    value: 'ad',
    type: 'ad',
    category: 'cloud',
  },
  {
    displayName: 'AliCloud',
    value: 'alicloud',
    type: 'alicloud',
    category: 'cloud',
  },
  {
    displayName: 'AWS',
    value: 'aws',
    type: 'aws',
    category: 'cloud',
  },
  {
    displayName: 'Azure',
    value: 'azure',
    type: 'azure',
    category: 'cloud',
  },
  {
    displayName: 'Consul',
    value: 'consul',
    type: 'consul',
    category: 'infra',
  },
  {
    displayName: 'Databases',
    value: 'database',
    type: 'database',
    category: 'infra',
  },
  {
    displayName: 'Google Cloud',
    value: 'gcp',
    type: 'gcp',
    category: 'cloud',
  },
  {
    displayName: 'Google Cloud KMS',
    value: 'gcpkms',
    type: 'gcpkms',
    category: 'cloud',
  },
  {
    displayName: 'KV',
    value: 'kv',
    type: 'kv',
    category: 'generic',
  },
  {
    displayName: 'Nomad',
    value: 'nomad',
    type: 'nomad',
    category: 'infra',
  },
  {
    displayName: 'PKI Certificates',
    value: 'pki',
    type: 'pki',
    category: 'generic',
  },
  {
    displayName: 'RabbitMQ',
    value: 'rabbitmq',
    type: 'rabbitmq',
    category: 'infra',
  },
  {
    displayName: 'SSH',
    value: 'ssh',
    type: 'ssh',
    category: 'generic',
  },
  {
    displayName: 'Transit',
    value: 'transit',
    type: 'transit',
    category: 'generic',
  },
  {
    displayName: 'TOTP',
    value: 'totp',
    type: 'totp',
    category: 'generic',
  },
];

export function engines() {
  return MOUNTABLE_SECRET_ENGINES.slice();
}

export default buildHelper(engines);
