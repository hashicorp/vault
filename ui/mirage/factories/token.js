import { Factory, trait } from 'miragejs';
import timestamp from 'core/utils/timestamp';
import { addHours, addDays } from 'date-fns';
import uuid from 'core/utils/uuid';

export default Factory.extend({
  // Core token information
  accessor: () => uuid(),
  creation_time: () => Math.floor(timestamp.now() / 1000),
  creation_ttl: 0,
  display_name: 'root',
  entity_id: '',
  expire_time: null,
  explicit_max_ttl: 0,
  issue_time: () => new Date(timestamp.now()).toISOString(),
  meta: () => ({}),
  num_uses: 0,
  orphan: true,
  path: 'auth/token/root',
  policies: () => ['root'],
  renewable: false,
  ttl: 0,
  type: 'service',

  // Additional fields
  bound_cidrs: () => [],
  external_namespace_policies: () => ({}),
  identity_policies: () => [],
  namespace_path: '',

  // Traits for different token types
  isService: trait({
    entity_id: uuid(),
    display_name: 'token',
    policies: () => ['default'],
    orphan: false,
    renewable: true,
    ttl: 2764800, // 32 days
    creation_ttl: 2764800,
    explicit_max_ttl: 0,
    path: 'auth/token/create',
  }),

  isBatch: trait({
    entity_id: uuid(),
    type: 'batch',
    display_name: 'batch-token',
    policies: () => ['default'],
    orphan: false,
    renewable: false,
    ttl: 600, // 10 minutes
    creation_ttl: 600,
    explicit_max_ttl: 0,
    path: 'auth/token/create',
  }),

  isUserpass: trait({
    entity_id: uuid(),
    display_name: 'userpass-test-user',
    path: 'auth/userpass/login/test-user',
    policies: () => ['default', 'user-policy'],
    orphan: false,
    renewable: true,
    ttl: 2764800,
    creation_ttl: 2764800,
    meta: () => ({ username: 'test-user' }),
  }),

  isLdap: trait({
    entity_id: uuid(),
    display_name: 'ldap-john.doe',
    path: 'auth/ldap/login/john.doe',
    policies: () => ['default', 'ldap-users'],
    orphan: false,
    renewable: true,
    ttl: 2764800,
    creation_ttl: 2764800,
    meta: () => ({ username: 'john.doe' }),
  }),

  withExpiry: trait({
    expire_time: () => addHours(timestamp.now(), 24).toISOString(),
    ttl: 86400, // 24 hours
    creation_ttl: 86400,
    renewable: true,
  }),

  isShortLived: trait({
    expire_time: () => addHours(timestamp.now(), 1).toISOString(),
    ttl: 3600, // 1 hour
    creation_ttl: 3600,
    renewable: true,
  }),

  withUses: trait({
    num_uses: 10,
    ttl: 0,
    creation_ttl: 0,
    renewable: false,
  }),

  withEntity: trait({
    entity_id: () => `entity_${Math.random().toString(36).substr(2, 16)}`,
    identity_policies: () => ['entity-policy'],
  }),

  withBoundCidrs: trait({
    bound_cidrs: () => ['192.168.1.0/24', '10.0.0.0/8'],
  }),

  isExpired: trait({
    expire_time: () => addDays(timestamp.now(), -1).toISOString(),
    ttl: 0,
    renewable: false,
  }),

  isNearExpiry: trait({
    expire_time: () => addHours(timestamp.now(), 1).toISOString(),
    ttl: 3600,
    renewable: true,
  }),

  withNamespace: trait({
    namespace_path: 'admin/',
    policies: () => ['admin-policy'],
    external_namespace_policies: () => ({
      'team-a/': ['team-a-policy'],
    }),
  }),
});
