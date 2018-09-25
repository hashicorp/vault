import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';

module('Unit | Serializer | policy', function(hooks) {
  setupTest(hooks);

  const POLICY_LIST_RESPONSE = {
    keys: ['default', 'root'],
    policies: ['default', 'root'],
    request_id: '3a6a3d67-dc3b-a086-2fc7-902bdc4dec3a',
    lease_id: '',
    renewable: false,
    lease_duration: 0,
    data: {
      keys: ['default', 'root'],
      policies: ['default', 'root'],
    },
    wrap_info: null,
    warnings: null,
    auth: null,
  };

  const EMBER_DATA_EXPECTS_FOR_POLICY_LIST = [{ name: 'default' }, { name: 'root' }];

  const POLICY_SHOW_RESPONSE = {
    name: 'default',
    rules:
      '\n# Allow tokens to look up their own properties\npath "auth/token/lookup-self" {\n    capabilities = ["read"]\n}\n\n# Allow tokens to renew themselves\npath "auth/token/renew-self" {\n    capabilities = ["update"]\n}\n\n# Allow tokens to revoke themselves\npath "auth/token/revoke-self" {\n    capabilities = ["update"]\n}\n\n# Allow a token to look up its own capabilities on a path\npath "sys/capabilities-self" {\n    capabilities = ["update"]\n}\n\n# Allow a token to renew a lease via lease_id in the request body\npath "sys/renew" {\n    capabilities = ["update"]\n}\n\n# Allow a token to manage its own cubbyhole\npath "cubbyhole/*" {\n    capabilities = ["create", "read", "update", "delete", "list"]\n}\n\n# Allow a token to list its cubbyhole (not covered by the splat above)\npath "cubbyhole" {\n    capabilities = ["list"]\n}\n\n# Allow a token to wrap arbitrary values in a response-wrapping token\npath "sys/wrapping/wrap" {\n    capabilities = ["update"]\n}\n\n# Allow a token to look up the creation time and TTL of a given\n# response-wrapping token\npath "sys/wrapping/lookup" {\n    capabilities = ["update"]\n}\n\n# Allow a token to unwrap a response-wrapping token. This is a convenience to\n# avoid client token swapping since this is also part of the response wrapping\n# policy.\npath "sys/wrapping/unwrap" {\n    capabilities = ["update"]\n}\n',
    request_id: '890eabf8-d418-07af-f978-928d328a7e64',
    lease_id: '',
    renewable: false,
    lease_duration: 0,
    data: {
      name: 'default',
      rules:
        '\n# Allow tokens to look up their own properties\npath "auth/token/lookup-self" {\n    capabilities = ["read"]\n}\n\n# Allow tokens to renew themselves\npath "auth/token/renew-self" {\n    capabilities = ["update"]\n}\n\n# Allow tokens to revoke themselves\npath "auth/token/revoke-self" {\n    capabilities = ["update"]\n}\n\n# Allow a token to look up its own capabilities on a path\npath "sys/capabilities-self" {\n    capabilities = ["update"]\n}\n\n# Allow a token to renew a lease via lease_id in the request body\npath "sys/renew" {\n    capabilities = ["update"]\n}\n\n# Allow a token to manage its own cubbyhole\npath "cubbyhole/*" {\n    capabilities = ["create", "read", "update", "delete", "list"]\n}\n\n# Allow a token to list its cubbyhole (not covered by the splat above)\npath "cubbyhole" {\n    capabilities = ["list"]\n}\n\n# Allow a token to wrap arbitrary values in a response-wrapping token\npath "sys/wrapping/wrap" {\n    capabilities = ["update"]\n}\n\n# Allow a token to look up the creation time and TTL of a given\n# response-wrapping token\npath "sys/wrapping/lookup" {\n    capabilities = ["update"]\n}\n\n# Allow a token to unwrap a response-wrapping token. This is a convenience to\n# avoid client token swapping since this is also part of the response wrapping\n# policy.\npath "sys/wrapping/unwrap" {\n    capabilities = ["update"]\n}\n',
    },
    wrap_info: null,
    warnings: null,
    auth: null,
  };

  const EMBER_DATA_EXPECTS_FOR_POLICY_SHOW = {
    name: 'default',
    rules:
      '\n# Allow tokens to look up their own properties\npath "auth/token/lookup-self" {\n    capabilities = ["read"]\n}\n\n# Allow tokens to renew themselves\npath "auth/token/renew-self" {\n    capabilities = ["update"]\n}\n\n# Allow tokens to revoke themselves\npath "auth/token/revoke-self" {\n    capabilities = ["update"]\n}\n\n# Allow a token to look up its own capabilities on a path\npath "sys/capabilities-self" {\n    capabilities = ["update"]\n}\n\n# Allow a token to renew a lease via lease_id in the request body\npath "sys/renew" {\n    capabilities = ["update"]\n}\n\n# Allow a token to manage its own cubbyhole\npath "cubbyhole/*" {\n    capabilities = ["create", "read", "update", "delete", "list"]\n}\n\n# Allow a token to list its cubbyhole (not covered by the splat above)\npath "cubbyhole" {\n    capabilities = ["list"]\n}\n\n# Allow a token to wrap arbitrary values in a response-wrapping token\npath "sys/wrapping/wrap" {\n    capabilities = ["update"]\n}\n\n# Allow a token to look up the creation time and TTL of a given\n# response-wrapping token\npath "sys/wrapping/lookup" {\n    capabilities = ["update"]\n}\n\n# Allow a token to unwrap a response-wrapping token. This is a convenience to\n# avoid client token swapping since this is also part of the response wrapping\n# policy.\npath "sys/wrapping/unwrap" {\n    capabilities = ["update"]\n}\n',
  };

  test('it transforms a list request payload', function(assert) {
    let serializer = this.owner.lookup('serializer:policy');

    let transformedPayload = serializer.normalizePolicies(POLICY_LIST_RESPONSE);

    assert.deepEqual(
      transformedPayload,
      EMBER_DATA_EXPECTS_FOR_POLICY_LIST,
      'transformed payload matches the expected payload'
    );
  });

  test('it transforms a list request payload', function(assert) {
    let serializer = this.owner.lookup('serializer:policy');

    let transformedPayload = serializer.normalizePolicies(POLICY_SHOW_RESPONSE);

    assert.deepEqual(
      transformedPayload,
      EMBER_DATA_EXPECTS_FOR_POLICY_SHOW,
      'transformed payload matches the expected payload'
    );
  });
});
