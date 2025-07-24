import { Factory, trait } from 'miragejs';

export default Factory.extend({
  // Core seal status
  type: 'shamir',
  initialized: true,
  sealed: false,
  t: 3,
  n: 5,
  progress: 0,
  nonce: '',
  version: '1.15.0+ent',
  build_date: '2023-09-01T10:00:00Z',
  migration: false,
  cluster_name: 'vault-cluster-e779cd7c',
  cluster_id: '5f20f5ab-acea-0481-787e-71ec2ff5a60b',
  recovery_seal: false,
  storage_type: 'consul',

  // Traits for different seal states
  isSealed: trait({
    sealed: true,
    progress: 0,
    nonce: '',
  }),

  isUnsealing: trait({
    sealed: true,
    progress: 1,
    nonce: 'af9b5c9c-1c2a-4b3d-8e9f-0a1b2c3d4e5f',
  }),

  isAutoUnseal: trait({
    type: 'awskms',
    sealed: false,
    recovery_seal: true,
  }),

  isTransitUnseal: trait({
    type: 'transit',
    sealed: false,
    recovery_seal: true,
  }),

  isHsmUnseal: trait({
    type: 'pkcs11',
    sealed: false,
    recovery_seal: true,
  }),

  isMigrating: trait({
    migration: true,
  }),
});
