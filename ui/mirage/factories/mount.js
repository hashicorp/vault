import { Factory, trait } from 'miragejs';
import uuid from 'core/utils/uuid';

export default Factory.extend({
  // Core mount information
  type: 'cubbyhole',
  path: 'cubbyhole/',
  description: 'per-token private secret storage',
  accessor: '', // set in afterCreate based on type

  // Configuration
  config: () => ({
    default_lease_ttl: 2764800,
    max_lease_ttl: 2764800,
    force_no_cache: false,
    listing_visibility: 'hidden',
  }),

  // Mount options
  options: null,

  // Seal wrapping
  seal_wrap: false,
  external_entropy_access: false,

  // Plugin information
  plugin_version: '',
  running_plugin_version: '',
  running_sha256: '',

  // Additional metadata
  local: false,
  uuid: () => uuid(),

  afterCreate(mount) {
    mount.accessor = `${mount.type}_${uuid().split('-')[0]}`;
  },

  // Traits for different mount types
  isKv: trait({
    type: 'kv',
    path: 'kv/',
    description: 'key/value secret storage',
    options: () => ({ version: '2' }),
  }),

  isKvV1: trait({
    type: 'kv',
    path: 'kv/',
    description: 'key/value secret storage v1',
    options: () => ({ version: '1' }),
  }),

  isPki: trait({
    type: 'pki',
    path: 'pki/',
    description: 'PKI certificate authority',
    config: () => ({
      default_lease_ttl: 0,
      max_lease_ttl: 315360000, // 10 years
      force_no_cache: false,
      listing_visibility: 'unauth',
    }),
  }),

  isDatabase: trait({
    type: 'database',
    path: 'database/',
    description: 'Database secrets engine',
    config: () => ({
      default_lease_ttl: 0,
      max_lease_ttl: 0,
      force_no_cache: false,
      listing_visibility: 'unauth',
    }),
  }),

  isAws: trait({
    type: 'aws',
    path: 'aws/',
    description: 'AWS secrets engine',
    config: () => ({
      default_lease_ttl: 0,
      max_lease_ttl: 0,
      force_no_cache: false,
      listing_visibility: 'unauth',
      identity_token_key: null,
      identity_token_ttl: 3600,
    }),
  }),

  isTransit: trait({
    type: 'transit',
    path: 'transit/',
    description: 'Transit encryption as a service',
    config: () => ({
      default_lease_ttl: 0,
      max_lease_ttl: 0,
      force_no_cache: false,
      listing_visibility: 'unauth',
    }),
  }),

  isLdap: trait({
    type: 'ldap',
    path: 'ldap/',
    description: 'LDAP secrets engine',
    config: () => ({
      default_lease_ttl: 0,
      max_lease_ttl: 0,
      force_no_cache: false,
      listing_visibility: 'unauth',
    }),
  }),

  isSsh: trait({
    type: 'ssh',
    path: 'ssh/',
    description: 'SSH secrets engine',
    config: () => ({
      default_lease_ttl: 0,
      max_lease_ttl: 0,
      force_no_cache: false,
      listing_visibility: 'unauth',
    }),
  }),

  withTtl: trait({
    config: () => ({
      default_lease_ttl: 3600, // 1 hour
      max_lease_ttl: 86400, // 24 hours
      force_no_cache: false,
      listing_visibility: 'unauth',
    }),
  }),

  isSealWrapped: trait({
    seal_wrap: true,
  }),

  isLocal: trait({
    local: true,
  }),

  withIdentityToken: trait({
    config: () => ({
      default_lease_ttl: 0,
      max_lease_ttl: 0,
      force_no_cache: false,
      listing_visibility: 'unauth',
      identity_token_key: 'test-key',
      identity_token_ttl: 3600,
    }),
  }),

  isHidden: trait({
    config: () => ({
      default_lease_ttl: 0,
      max_lease_ttl: 0,
      force_no_cache: false,
      listing_visibility: 'hidden',
    }),
  }),
});
