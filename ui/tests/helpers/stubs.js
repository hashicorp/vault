/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { Response } from 'miragejs';

export function capabilitiesStub(requestPath, capabilitiesArray) {
  // sample of capabilitiesArray: ['read', 'update']
  return {
    request_id: '40f7e44d-af5c-9b60-bd20-df72eb17e294',
    lease_id: '',
    renewable: false,
    lease_duration: 0,
    data: {
      [requestPath]: capabilitiesArray,
      capabilities: capabilitiesArray,
    },
    wrap_info: null,
    warnings: null,
    auth: null,
  };
}

export const noopStub = (response) => {
  return function () {
    return [response, { 'Content-Type': 'application/json' }, JSON.stringify({})];
  };
};

/**
 * allowAllCapabilitiesStub mocks the response from capabilities-self
 * that allows the user to do any action (root user)
 * Example usage assuming setupMirage(hooks) was called:
 * this.server.post('/sys/capabilities-self', allowAllCapabilitiesStub(['read']));
 */
export function allowAllCapabilitiesStub(capabilitiesList = ['root']) {
  return function (_, { requestBody }) {
    const { paths } = JSON.parse(requestBody);
    const specificCapabilities = paths.reduce((obj, path) => {
      return {
        ...obj,
        [path]: capabilitiesList,
      };
    }, {});
    return {
      ...specificCapabilities,
      capabilities: capabilitiesList,
      request_id: 'mirage-795dc9e1-0321-9ac6-71fc',
      lease_id: '',
      renewable: false,
      lease_duration: 0,
      data: { ...specificCapabilities, capabilities: capabilitiesList },
      wrap_info: null,
      warnings: null,
      auth: null,
    };
  };
}

/**
 * returns a response with the given httpStatus and data based on status
 * @param {number} httpStatus 403, 404, 204, or 200 (default)
 * @param {object} payload what to return in the response if status is 200
 * @returns {Response}
 */
export function overrideResponse(httpStatus = 200, payload = {}) {
  if (httpStatus === 403) {
    return new Response(403, { 'Content-Type': 'application/json' }, formatError('permission denied'));
  }
  if (httpStatus === 404) {
    return new Response(404, { 'Content-Type': 'application/json' });
  }
  if (httpStatus === 204) {
    return new Response(204, { 'Content-Type': 'application/json' });
  }
  return new Response(httpStatus, { 'Content-Type': 'application/json' }, payload);
}

export const formatError = (msg) => JSON.stringify({ errors: [msg] });

/**
 * Minimal OpenAPI spec fixture for testing.
 * Contains only the mounts-enable-secrets-engine operation.
 */
export const OAS_STUB = {
  openapi: '3.0.0',
  info: {
    title: 'Vault API',
    version: '1.0.0',
  },
  paths: {
    '/sys/mounts/{path}': {
      description: 'Mount a new backend at a new path.',
      parameters: [
        {
          name: 'path',
          description: 'The path to mount to. Example: "aws/east"',
          in: 'path',
          schema: { type: 'string' },
          required: true,
        },
      ],
      post: {
        summary: 'Enable a new secrets engine at the given path.',
        operationId: 'mounts-enable-secrets-engine',
        tags: ['system'],
        requestBody: {
          required: true,
          content: {
            'application/json': {
              schema: {
                $ref: '#/components/schemas/MountsEnableSecretsEngineRequest',
              },
            },
          },
        },
        responses: {
          204: { description: 'OK' },
        },
      },
    },
  },
  components: {
    schemas: {
      MountsEnableSecretsEngineRequest: {
        type: 'object',
        properties: {
          config: {
            type: 'object',
            description: 'Configuration for this mount, such as default_lease_ttl and max_lease_ttl.',
          },
          description: {
            type: 'string',
            description: 'User-friendly description for this mount.',
          },
          external_entropy_access: {
            type: 'boolean',
            description: "Whether to give the mount access to Vault's external entropy.",
            default: false,
            deprecated: true,
          },
          local: {
            type: 'boolean',
            description:
              'Mark the mount as a local mount, which is not replicated and is unaffected by replication.',
            default: false,
          },
          options: {
            type: 'object',
            description:
              'The options to pass into the backend. Should be a json object with string keys and values.',
          },
          plugin_name: {
            type: 'string',
            description: 'Name of the plugin to mount based from the name registered in the plugin catalog.',
          },
          plugin_version: {
            type: 'string',
            description: 'The semantic version of the plugin to use, or image tag if oci_image is provided.',
          },
          seal_wrap: {
            type: 'boolean',
            description: 'Whether to turn on seal wrapping for the mount.',
            default: false,
            'x-vault-displayAttrs': {
              name: 'Seal Wrap',
              group: 'Advanced',
            },
          },
          type: {
            type: 'string',
            description: 'The type of the backend. Example: "passthrough"',
          },
          allowed_managed_keys: {
            type: 'array',
            description: 'List of managed key names allowed for this mount.',
            'x-vault-displayAttrs': {
              name: 'Allowed Managed Keys',
              group: 'Advanced',
            },
          },
        },
      },
    },
  },
};
