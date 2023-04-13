/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { attr } from '@ember-data/model';
import { expandOpenApiProps, combineAttributes, combineFieldGroups } from 'vault/utils/openapi-to-attrs';
import { module, test } from 'qunit';
import { camelize } from '@ember/string';

module('Unit | Util | OpenAPI Data Utilities', function () {
  const OPENAPI_RESPONSE_PROPS = {
    ttl: {
      type: 'string',
      format: 'seconds',
      description: 'this is a TTL!',
      'x-vault-displayAttrs': {
        name: 'TTL',
      },
    },
    'awesome-people': {
      type: 'array',
      items: {
        type: 'string',
      },
      'x-vault-displayAttrs': {
        value: 'Grace Hopper,Lady Ada',
      },
    },
    'favorite-ice-cream': {
      type: 'string',
      enum: ['vanilla', 'chocolate', 'strawberry'],
    },
    'default-value': {
      default: 30,
      'x-vault-displayAttrs': {
        value: 300,
      },
      type: 'integer',
    },
    default: {
      'x-vault-displayAttrs': {
        value: 30,
      },
      type: 'integer',
    },
    'super-secret': {
      type: 'string',
      'x-vault-displayAttrs': {
        sensitive: true,
      },
      description: 'A really secret thing',
    },
  };
  const EXPANDED_PROPS = {
    ttl: {
      helpText: 'this is a TTL!',
      editType: 'ttl',
      label: 'TTL',
      fieldGroup: 'default',
    },
    awesomePeople: {
      editType: 'stringArray',
      defaultValue: 'Grace Hopper,Lady Ada',
      fieldGroup: 'default',
    },
    favoriteIceCream: {
      editType: 'string',
      type: 'string',
      possibleValues: ['vanilla', 'chocolate', 'strawberry'],
      fieldGroup: 'default',
    },
    defaultValue: {
      editType: 'number',
      type: 'number',
      defaultValue: 300,
      fieldGroup: 'default',
    },
    default: {
      editType: 'number',
      type: 'number',
      defaultValue: 30,
      fieldGroup: 'default',
    },
    superSecret: {
      type: 'string',
      editType: 'string',
      sensitive: true,
      helpText: 'A really secret thing',
      fieldGroup: 'default',
    },
  };

  const EXISTING_MODEL_ATTRS = [
    {
      key: 'name',
      value: {
        isAttribute: true,
        name: 'name',
        options: {
          editType: 'string',
          label: 'Role name',
        },
      },
    },
    {
      key: 'awesomePeople',
      value: {
        isAttribute: true,
        name: 'awesomePeople',
        options: {
          label: 'People Who Are Awesome',
        },
      },
    },
  ];

  const COMBINED_ATTRS = {
    name: attr('string', {
      editType: 'string',
      type: 'string',
      label: 'Role name',
    }),
    ttl: attr('string', {
      editType: 'ttl',
      label: 'TTL',
      helpText: 'this is a TTL!',
    }),
    awesomePeople: attr({
      label: 'People Who Are Awesome',
      editType: 'stringArray',
      defaultValue: 'Grace Hopper,Lady Ada',
    }),
    favoriteIceCream: attr('string', {
      type: 'string',
      editType: 'string',
      possibleValues: ['vanilla', 'chocolate', 'strawberry'],
    }),
    superSecret: attr('string', {
      type: 'string',
      editType: 'string',
      sensitive: true,
      description: 'A really secret thing',
    }),
  };

  const NEW_FIELDS = ['one', 'two', 'three'];

  const OPENAPI_DESCRIPTIONS = {
    token_bound_cidrs: {
      type: 'array',
      description: `Comma separated string or JSON list of CIDR blocks. If set, specifies the blocks of IP addresses which are allowed to use the generated token.`,
      items: {
        type: 'string',
      },
      'x-vault-displayAttrs': {
        name: "Generated Token's Bound CIDRs",
        group: 'Tokens',
      },
    },
    token_policies: {
      type: 'array',
      description: `Comma-separated list of policies`,
      items: {
        type: 'string',
      },
      'x-vault-displayAttrs': {
        name: "Generated Token's Policies",
        group: 'Tokens',
      },
    },
    secret_id_bound_cidrs: {
      type: 'array',
      description: `Comma separated string or list of CIDR blocks. If set, specifies the blocks of IP addresses which can perform the login operation.`,
      items: {
        type: 'string',
      },
    },
    bound_cidr_list: {
      type: 'array',
      description: `Deprecated: Please use "secret_id_bound_cidrs" instead. Comma separated string or list of CIDR blocks. If set, specifies the blocks of IP addresses which can perform the login operation.`,
      items: {
        type: 'string',
      },
    },
    allowed_roles: {
      type: 'array',
      description: `Comma separated string or array of the role names allowed to get creds from this database connection. If empty no roles are allowed. If "*" all roles are allowed.`,
      items: {
        type: 'string',
      },
    },
    cidr_list: {
      type: 'array',
      description: `[Optional for OTP type] [Not applicable for CA type] Comma separated list of CIDR blocks for which the role is applicable for. CIDR blocks can belong to more than one role.`,
      items: {
        type: 'string',
      },
    },
    ocsp_servers_override: {
      type: 'array',
      description: `A comma-separated list of OCSP server addresses.  If unset, the OCSP server is determined from the AuthorityInformationAccess extension on the certificate being inspected.`,
      items: {
        type: 'string',
      },
    },
    key_usage: {
      type: 'array',
      description: `A comma-separated string or list of key usages (not extended key usages). Valid values can be found at https://golang.org/pkg/crypto/x509/#KeyUsage -- simply drop the "KeyUsage" part of the name. To remove all key usages from being set, set this value to an empty list.`,
      items: {
        type: 'string',
      },
    },
    required_extensions: {
      type: 'array',
      description: `A comma-separated string or array of extensions formatted as "oid:value". Expects the extension value to be some type of ASN1 encoded string. All values much match. Supports globbing on "value".`,
      items: {
        type: 'string',
      },
    },
    crl_distribution_points: {
      type: 'array',
      description: `Comma-separated list of URLs to be used for the CRL distribution points attribute. See also RFC 5280 Section 4.2.1.13.`,
      items: {
        type: 'string',
      },
    },
  };

  const STRING_ARRAY_DESCRIPTIONS = {
    token_bound_cidrs: {
      helpText: `List of CIDR blocks. If set, specifies the blocks of IP addresses which are allowed to use the generated token.`,
    },
    token_policies: {
      helpText: `List of policies`,
    },
    secret_id_bound_cidrs: {
      helpText: `List of CIDR blocks. If set, specifies the blocks of IP addresses which can perform the login operation.`,
    },
    bound_cidr_list: {
      helpText: `Deprecated: Please use "secret_id_bound_cidrs" instead. List of CIDR blocks. If set, specifies the blocks of IP addresses which can perform the login operation.`,
    },
    allowed_roles: {
      helpText: `List of the role names allowed to get creds from this database connection. If empty no roles are allowed. If "*" all roles are allowed.`,
    },
    cidr_list: {
      helpText: `[Optional for OTP type] [Not applicable for CA type] List of CIDR blocks for which the role is applicable for. CIDR blocks can belong to more than one role.`,
    },
    ocsp_servers_override: {
      helpText: `A list of OCSP server addresses.  If unset, the OCSP server is determined from the AuthorityInformationAccess extension on the certificate being inspected.`,
    },
    key_usage: {
      helpText: `A list of key usages (not extended key usages). Valid values can be found at https://golang.org/pkg/crypto/x509/#KeyUsage -- simply drop the "KeyUsage" part of the name. To remove all key usages from being set, set this value to an empty list.`,
    },
    required_extensions: {
      helpText: `A list of extensions formatted as "oid:value". Expects the extension value to be some type of ASN1 encoded string. All values much match. Supports globbing on "value".`,
    },
    crl_distribution_points: {
      helpText: `List of URLs to be used for the CRL distribution points attribute. See also RFC 5280 Section 4.2.1.13.`,
    },
  };

  test('it creates objects from OpenAPI schema props', function (assert) {
    assert.expect(6);
    const generatedProps = expandOpenApiProps(OPENAPI_RESPONSE_PROPS);
    for (const propName in EXPANDED_PROPS) {
      assert.deepEqual(EXPANDED_PROPS[propName], generatedProps[propName], `correctly expands ${propName}`);
    }
  });

  test('it combines OpenAPI props with existing model attrs', function (assert) {
    assert.expect(3);
    const combined = combineAttributes(EXISTING_MODEL_ATTRS, EXPANDED_PROPS);
    for (const propName in EXISTING_MODEL_ATTRS) {
      assert.deepEqual(COMBINED_ATTRS[propName], combined[propName]);
    }
  });

  test('it adds new fields from OpenAPI to fieldGroups except for exclusions', function (assert) {
    assert.expect(3);
    const modelFieldGroups = [
      { default: ['name', 'awesomePeople'] },
      {
        Options: ['ttl'],
      },
    ];
    const excludedFields = ['two'];
    const expectedGroups = [
      { default: ['name', 'awesomePeople', 'one', 'three'] },
      {
        Options: ['ttl'],
      },
    ];
    const newFieldGroups = combineFieldGroups(modelFieldGroups, NEW_FIELDS, excludedFields);
    for (const groupName in modelFieldGroups) {
      assert.deepEqual(
        newFieldGroups[groupName],
        expectedGroups[groupName],
        'it incorporates all new fields except for those excluded'
      );
    }
  });
  test('it adds all new fields from OpenAPI to fieldGroups when excludedFields is empty', function (assert) {
    assert.expect(3);
    const modelFieldGroups = [
      { default: ['name', 'awesomePeople'] },
      {
        Options: ['ttl'],
      },
    ];
    const excludedFields = [];
    const expectedGroups = [
      { default: ['name', 'awesomePeople', 'one', 'two', 'three'] },
      {
        Options: ['ttl'],
      },
    ];
    const nonExcludedFieldGroups = combineFieldGroups(modelFieldGroups, NEW_FIELDS, excludedFields);
    for (const groupName in modelFieldGroups) {
      assert.deepEqual(
        nonExcludedFieldGroups[groupName],
        expectedGroups[groupName],
        'it incorporates all new fields'
      );
    }
  });
  test('it keeps fields the same when there are no brand new fields from OpenAPI', function (assert) {
    assert.expect(3);
    const modelFieldGroups = [
      { default: ['name', 'awesomePeople', 'two', 'one', 'three'] },
      {
        Options: ['ttl'],
      },
    ];
    const excludedFields = [];
    const expectedGroups = [
      { default: ['name', 'awesomePeople', 'two', 'one', 'three'] },
      {
        Options: ['ttl'],
      },
    ];
    const fieldGroups = combineFieldGroups(modelFieldGroups, NEW_FIELDS, excludedFields);
    for (const groupName in modelFieldGroups) {
      assert.deepEqual(fieldGroups[groupName], expectedGroups[groupName], 'it incorporates all new fields');
    }
  });

  test('it removes references to comma separation in help text for string array attrs', async function (assert) {
    assert.expect(10);
    const generatedProps = expandOpenApiProps(OPENAPI_DESCRIPTIONS);
    for (const propName in STRING_ARRAY_DESCRIPTIONS) {
      assert.strictEqual(
        generatedProps[camelize(propName)].helpText,
        STRING_ARRAY_DESCRIPTIONS[propName].helpText,
        `correctly updates helpText for ${propName}`
      );
    }
  });
});
