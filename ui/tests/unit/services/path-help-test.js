/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';

const openapiStub = {
  openapi: {
    openapi: '3.0.2',
    info: {
      title: 'HashiCorp Vault API',
      description: 'HTTP API that gives you full access to Vault. All API routes are prefixed with `/v1/`.',
      version: '1.16.0',
      license: {
        name: 'Mozilla Public License 2.0',
        url: 'https://www.mozilla.org/en-US/MPL/2.0',
      },
    },
    paths: {
      '/login/{username}': {
        description: 'Log in with a username and password.',
        parameters: [
          {
            name: 'username',
            description: 'Username of the user.',
            in: 'path',
            schema: {
              type: 'string',
            },
            required: true,
          },
        ],
        'x-vault-unauthenticated': true,
        post: {
          summary: 'Log in with a username and password.',
          operationId: 'userpass-login',
          tags: ['auth'],
          requestBody: {
            required: true,
            content: {
              'application/json': {
                schema: {
                  $ref: '#/components/schemas/UserpassLoginRequest',
                },
              },
            },
          },
          responses: {
            200: {
              description: 'OK',
            },
          },
        },
      },
      '/users/': {
        description: 'Manage users allowed to authenticate.',
        'x-vault-displayAttrs': {
          navigation: true,
          itemType: 'User',
        },
        get: {
          summary: 'Manage users allowed to authenticate.',
          operationId: 'userpass-list-users',
          tags: ['auth'],
          parameters: [
            {
              name: 'list',
              description: 'Must be set to `true`',
              in: 'query',
              schema: {
                type: 'string',
                enum: ['true'],
              },
              required: true,
            },
          ],
          responses: {
            200: {
              description: 'OK',
              content: {
                'application/json': {
                  schema: {
                    $ref: '#/components/schemas/StandardListResponse',
                  },
                },
              },
            },
          },
        },
      },
      '/users/{username}': {
        description: 'Manage users allowed to authenticate.',
        parameters: [
          {
            name: 'username',
            description: 'Username for this user.',
            in: 'path',
            schema: {
              type: 'string',
            },
            required: true,
          },
        ],
        'x-vault-createSupported': true,
        'x-vault-displayAttrs': {
          itemType: 'User',
          action: 'Create',
        },
        get: {
          summary: 'Manage users allowed to authenticate.',
          operationId: 'userpass-read-user',
          tags: ['auth'],
          responses: {
            200: {
              description: 'OK',
            },
          },
        },
        post: {
          summary: 'Manage users allowed to authenticate.',
          operationId: 'userpass-write-user',
          tags: ['auth'],
          requestBody: {
            required: true,
            content: {
              'application/json': {
                schema: {
                  $ref: '#/components/schemas/UserpassWriteUserRequest',
                },
              },
            },
          },
          responses: {
            200: {
              description: 'OK',
            },
          },
        },
        delete: {
          summary: 'Manage users allowed to authenticate.',
          operationId: 'userpass-delete-user',
          tags: ['auth'],
          responses: {
            204: {
              description: 'empty body',
            },
          },
        },
      },
      '/users/{username}/password': {
        description: "Reset user's password.",
        parameters: [
          {
            name: 'username',
            description: 'Username for this user.',
            in: 'path',
            schema: {
              type: 'string',
            },
            required: true,
          },
        ],
        post: {
          summary: "Reset user's password.",
          operationId: 'userpass-reset-password',
          tags: ['auth'],
          requestBody: {
            required: true,
            content: {
              'application/json': {
                schema: {
                  $ref: '#/components/schemas/UserpassResetPasswordRequest',
                },
              },
            },
          },
          responses: {
            200: {
              description: 'OK',
            },
          },
        },
      },
      '/users/{username}/policies': {
        description: 'Update the policies associated with the username.',
        parameters: [
          {
            name: 'username',
            description: 'Username for this user.',
            in: 'path',
            schema: {
              type: 'string',
            },
            required: true,
          },
        ],
        post: {
          summary: 'Update the policies associated with the username.',
          operationId: 'userpass-update-policies',
          tags: ['auth'],
          requestBody: {
            required: true,
            content: {
              'application/json': {
                schema: {
                  $ref: '#/components/schemas/UserpassUpdatePoliciesRequest',
                },
              },
            },
          },
          responses: {
            200: {
              description: 'OK',
            },
          },
        },
      },
    },
    components: {
      schemas: {
        StandardListResponse: {
          type: 'object',
          properties: {
            keys: {
              type: 'array',
              items: {
                type: 'string',
              },
            },
          },
        },
        UserpassLoginRequest: {
          type: 'object',
          properties: {
            password: {
              type: 'string',
              description: 'Password for this user.',
            },
          },
        },
        UserpassResetPasswordRequest: {
          type: 'object',
          properties: {
            password: {
              type: 'string',
              description: 'Password for this user.',
            },
          },
        },
        UserpassUpdatePoliciesRequest: {
          type: 'object',
          properties: {
            policies: {
              type: 'array',
              description:
                'Use "token_policies" instead. If this and "token_policies" are both specified, only "token_policies" will be used.',
              items: {
                type: 'string',
              },
              deprecated: true,
            },
            token_policies: {
              type: 'array',
              description: 'Comma-separated list of policies',
              items: {
                type: 'string',
              },
              'x-vault-displayAttrs': {
                description: 'A list of policies that will apply to the generated token for this user.',
              },
            },
          },
        },
        UserpassWriteUserRequest: {
          type: 'object',
          properties: {
            bound_cidrs: {
              type: 'array',
              description:
                'Use "token_bound_cidrs" instead. If this and "token_bound_cidrs" are both specified, only "token_bound_cidrs" will be used.',
              items: {
                type: 'string',
              },
              deprecated: true,
            },
            max_ttl: {
              type: 'string',
              description:
                'Use "token_max_ttl" instead. If this and "token_max_ttl" are both specified, only "token_max_ttl" will be used.',
              format: 'duration',
              deprecated: true,
            },
            password: {
              type: 'string',
              description: 'Password for this user.',
              'x-vault-displayAttrs': {
                sensitive: true,
              },
            },
            policies: {
              type: 'array',
              description:
                'Use "token_policies" instead. If this and "token_policies" are both specified, only "token_policies" will be used.',
              items: {
                type: 'string',
              },
              deprecated: true,
            },
            token_bound_cidrs: {
              type: 'array',
              description:
                'Comma separated string or JSON list of CIDR blocks. If set, specifies the blocks of IP addresses which are allowed to use the generated token.',
              items: {
                type: 'string',
              },
              'x-vault-displayAttrs': {
                name: "Generated Token's Bound CIDRs",
                description:
                  'A list of CIDR blocks. If set, specifies the blocks of IP addresses which are allowed to use the generated token.',
                group: 'Tokens',
              },
            },
            token_explicit_max_ttl: {
              type: 'string',
              description:
                'If set, tokens created via this role carry an explicit maximum TTL. During renewal, the current maximum TTL values of the role and the mount are not checked for changes, and any updates to these values will have no effect on the token being renewed.',
              format: 'duration',
              'x-vault-displayAttrs': {
                name: "Generated Token's Explicit Maximum TTL",
                group: 'Tokens',
              },
            },
            token_max_ttl: {
              type: 'string',
              description: 'The maximum lifetime of the generated token',
              format: 'duration',
              'x-vault-displayAttrs': {
                name: "Generated Token's Maximum TTL",
                group: 'Tokens',
              },
            },
            token_no_default_policy: {
              type: 'boolean',
              description:
                "If true, the 'default' policy will not automatically be added to generated tokens",
              'x-vault-displayAttrs': {
                name: "Do Not Attach 'default' Policy To Generated Tokens",
                group: 'Tokens',
              },
            },
            token_num_uses: {
              type: 'integer',
              description: 'The maximum number of times a token may be used, a value of zero means unlimited',
              'x-vault-displayAttrs': {
                name: 'Maximum Uses of Generated Tokens',
                group: 'Tokens',
              },
            },
            token_period: {
              type: 'string',
              description:
                'If set, tokens created via this role will have no max lifetime; instead, their renewal period will be fixed to this value. This takes an integer number of seconds, or a string duration (e.g. "24h").',
              format: 'duration',
              'x-vault-displayAttrs': {
                name: "Generated Token's Period",
                group: 'Tokens',
              },
            },
            token_policies: {
              type: 'array',
              description: 'Comma-separated list of policies',
              items: {
                type: 'string',
              },
              'x-vault-displayAttrs': {
                name: "Generated Token's Policies",
                description: 'A list of policies that will apply to the generated token for this user.',
                group: 'Tokens',
              },
            },
            token_ttl: {
              type: 'string',
              description: 'The initial ttl of the token to generate',
              format: 'duration',
              'x-vault-displayAttrs': {
                name: "Generated Token's Initial TTL",
                group: 'Tokens',
              },
            },
            token_type: {
              type: 'string',
              description: 'The type of token to generate, service or batch',
              default: 'default-service',
              'x-vault-displayAttrs': {
                name: "Generated Token's Type",
                group: 'Tokens',
              },
            },
            ttl: {
              type: 'string',
              description:
                'Use "token_ttl" instead. If this and "token_ttl" are both specified, only "token_ttl" will be used.',
              format: 'duration',
              deprecated: true,
            },
          },
        },
      },
    },
  },
};

module('Unit | Service | path-help', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.pathHelp = this.owner.lookup('service:path-help');
    this.store = this.owner.lookup('service:store');
  });

  test('it should generate model with mutableId', async function (assert) {
    assert.expect(2);

    this.server.get('/auth/userpass/', () => openapiStub);
    this.server.get('/auth/userpass/users/example', () => openapiStub);
    this.server.post('/auth/userpass/users/test', () => {
      assert.ok(true, 'POST request made to correct endpoint');
      return;
    });

    const modelType = 'generated-user-userpass';
    await this.pathHelp.newModelFromOpenApi(modelType, 'userpass', 'auth/userpass/', 'user');
    const model = this.store.createRecord(modelType);
    model.set('mutableId', 'test');
    await model.save();
    assert.strictEqual(model.get('id'), 'test', 'model id is set to mutableId value on save success');
  });

  test('it should return correct data for given path and item', async function (assert) {
    assert.expect(2);
    this.server.get('auth/cert-auth/', (_, request) => {
      assert.strictEqual(request.queryParams.help, '1', 'calls given path with help query param');
      return {
        openapi: {
          paths: {
            '/certs/': {
              description: 'Manage trusted certificates used for authentication.',
              'x-vault-displayAttrs': {
                navigation: true,
                itemType: 'Certificate',
              },
              get: {
                summary: 'Manage trusted certificates used for authentication.',
                operationId: 'cert-list-certificates',
                tags: ['auth'],
                parameters: [
                  {
                    name: 'list',
                    description: 'Must be set to `true`',
                    in: 'query',
                    schema: {
                      type: 'string',
                      enum: ['true'],
                    },
                    required: true,
                  },
                ],
                responses: {
                  200: {
                    description: 'OK',
                    content: {
                      'application/json': {
                        schema: {
                          $ref: '#/components/schemas/StandardListResponse',
                        },
                      },
                    },
                  },
                },
              },
            },
            '/certs/{name}': {
              description: 'Manage trusted certificates used for authentication.',
              parameters: [
                {
                  name: 'name',
                  description: 'The name of the certificate',
                  in: 'path',
                  schema: {
                    type: 'string',
                  },
                  required: true,
                },
              ],
              'x-vault-displayAttrs': {
                itemType: 'Certificate',
                action: 'Create',
              },
              get: {
                summary: 'Manage trusted certificates used for authentication.',
                operationId: 'cert-read-certificate',
                tags: ['auth'],
                responses: {
                  200: {
                    description: 'OK',
                  },
                },
              },
              post: {
                summary: 'Manage trusted certificates used for authentication.',
                operationId: 'cert-write-certificate',
                tags: ['auth'],
                requestBody: {
                  required: true,
                  content: {
                    'application/json': {
                      schema: {
                        $ref: '#/components/schemas/CertWriteCertificateRequest',
                      },
                    },
                  },
                },
                responses: {
                  200: {
                    description: 'OK',
                  },
                },
              },
              delete: {
                summary: 'Manage trusted certificates used for authentication.',
                operationId: 'cert-delete-certificate',
                tags: ['auth'],
                responses: {
                  204: {
                    description: 'empty body',
                  },
                },
              },
            },
            '/config': {
              get: {
                operationId: 'cert-read-configuration',
                tags: ['auth'],
                responses: {
                  200: {
                    description: 'OK',
                  },
                },
              },
              post: {
                operationId: 'cert-configure',
                tags: ['auth'],
                requestBody: {
                  required: true,
                  content: {
                    'application/json': {
                      schema: {
                        $ref: '#/components/schemas/CertConfigureRequest',
                      },
                    },
                  },
                },
                responses: {
                  200: {
                    description: 'OK',
                  },
                },
              },
            },
            '/crls/': {
              description: 'Manage Certificate Revocation Lists checked during authentication.',
              get: {
                operationId: 'cert-list-crls',
                tags: ['auth'],
                parameters: [
                  {
                    name: 'list',
                    description: 'Must be set to `true`',
                    in: 'query',
                    schema: {
                      type: 'string',
                      enum: ['true'],
                    },
                    required: true,
                  },
                ],
                responses: {
                  200: {
                    description: 'OK',
                    content: {
                      'application/json': {
                        schema: {
                          $ref: '#/components/schemas/StandardListResponse',
                        },
                      },
                    },
                  },
                },
              },
            },
            '/crls/{name}': {
              description: 'Manage Certificate Revocation Lists checked during authentication.',
              parameters: [
                {
                  name: 'name',
                  description: 'The name of the certificate',
                  in: 'path',
                  schema: {
                    type: 'string',
                  },
                  required: true,
                },
              ],
              get: {
                summary: 'Manage Certificate Revocation Lists checked during authentication.',
                operationId: 'cert-read-crl',
                tags: ['auth'],
                responses: {
                  200: {
                    description: 'OK',
                  },
                },
              },
              post: {
                summary: 'Manage Certificate Revocation Lists checked during authentication.',
                operationId: 'cert-write-crl',
                tags: ['auth'],
                requestBody: {
                  required: true,
                  content: {
                    'application/json': {
                      schema: {
                        $ref: '#/components/schemas/CertWriteCrlRequest',
                      },
                    },
                  },
                },
                responses: {
                  200: {
                    description: 'OK',
                  },
                },
              },
              delete: {
                summary: 'Manage Certificate Revocation Lists checked during authentication.',
                operationId: 'cert-delete-crl',
                tags: ['auth'],
                responses: {
                  204: {
                    description: 'empty body',
                  },
                },
              },
            },
            '/login': {
              'x-vault-unauthenticated': true,
              post: {
                operationId: 'cert-login',
                tags: ['auth'],
                requestBody: {
                  required: true,
                  content: {
                    'application/json': {
                      schema: {
                        $ref: '#/components/schemas/CertLoginRequest',
                      },
                    },
                  },
                },
                responses: {
                  200: {
                    description: 'OK',
                  },
                },
              },
            },
          },
        },
      };
    });
    const result = await this.pathHelp.getPaths('auth/cert-auth/', 'cert-auth');
    assert.deepEqual(
      result,
      {
        apiPath: 'auth/cert-auth/',
        itemID: undefined,
        itemType: undefined,
        itemTypes: ['certificate'],
        paths: [
          {
            action: undefined,
            path: '/certs/',
            itemType: 'certificate',
            itemName: 'Certificate',
            operations: ['get', 'list'],
            navigation: true,
            param: false,
          },
          {
            path: '/certs/{name}',
            itemType: 'certificate',
            itemName: 'Certificate',
            operations: ['get', 'post', 'delete'],
            action: 'Create',
            navigation: false,
            param: 'name',
          },
        ],
      },
      'returns correct data for given path and item'
    );
  });
});
