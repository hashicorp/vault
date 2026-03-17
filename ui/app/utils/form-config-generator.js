/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

/**
 * Utility functions for generating form configurations from OpenAPI specs.
 * These functions have no side effects and can be safely imported and used in unit tests.
 */

import { dasherize, classify } from '@ember/string';

const TYPE_DEFAULTS = {
  string: '',
  number: 0,
  integer: 0,
  boolean: false,
  object: {},
  array: [],
};

const API_CLASS_FROM_TAG = {
  auth: 'auth',
  identity: 'identity',
  secrets: 'secrets',
  system: 'sys',
};

const formatLabel = (str) => {
  // HDS guidelines suggest using sentence case for labels
  // "some_example_text" → "Some example text"
  // "client-id" → "Client id"
  const sentence = str.replace(/[_-]/g, ' ').toLowerCase();
  return sentence.charAt(0).toUpperCase() + sentence.slice(1);
};

/**
 * Construct request type name to match the vault-client-typescript SDK conventions.
 * Example: tag='system', methodName='mountsEnableSecretsEngine'
 *        → 'SystemApiMountsEnableSecretsEngineOperationRequest'
 */
const getRequestType = (tag, methodName) => {
  const apiClassName = `${classify(tag)}Api`;
  const methodPascal = classify(methodName);
  return `${apiClassName}${methodPascal}OperationRequest`;
};

const getOperationDetails = (spec, operationId) => {
  const { components } = spec;

  for (const pathUrl in spec.paths) {
    const pathItem = spec.paths[pathUrl];
    const { post } = pathItem;

    if (post?.operationId !== operationId) continue;

    const ref = post.requestBody?.content?.['application/json']?.schema?.$ref;
    const schemaName = ref?.split('/').pop();
    const requestBody = components.schemas[schemaName];

    // There are typically two parts to the API operations -
    // the parameters and the request body. These are in different places
    // in the spec so we need to address them separately.
    const params = [];
    for (const param of pathItem.parameters) {
      if (param.deprecated) continue;
      params.push({
        name: param.name,
        description: param.description,
        required: param.required,
        type: param.schema.type,
      });
    }

    const properties = {};
    for (const [propName, prop] of Object.entries(requestBody.properties)) {
      if (prop.deprecated) continue;
      properties[propName] = prop;
    }

    return {
      operationId: post.operationId,
      tag: post.tags?.[0],
      description: pathItem.description || post.summary || '',
      parameters: params,
      requestBody: [schemaName, properties],
    };
  }

  return null;
};

const buildPayloadFromOperation = (operation) => {
  const [requestSchemaName, requestProperties] = operation.requestBody;
  const payload = {};

  for (const param of operation.parameters) {
    payload[param.name] = TYPE_DEFAULTS[param.type];
  }

  const requestPayload = {};
  for (const [propName, prop] of Object.entries(requestProperties)) {
    requestPayload[propName] = TYPE_DEFAULTS[prop.type];
  }
  payload[requestSchemaName] = requestPayload;

  return payload;
};

const buildSectionsFromOperation = (operation) => {
  const [requestSchemaName, requestProperties] = operation.requestBody;
  // { default: [...], Advanced: [...] }
  const groups = {};

  // Group all of the parameter fields together
  for (const param of operation.parameters) {
    // `??=` is used to initialize the group if it doesn't exist
    (groups.params ??= []).push({
      name: param.name,
      type: 'TextInput',
      label: formatLabel(param.name),
      helperText: param.description,
    });
  }

  // Add request body fields to their respective groups
  for (const [propName, prop] of Object.entries(requestProperties)) {
    const group = prop['x-vault-displayAttrs']?.group || 'default';
    (groups[group] ??= []).push({
      name: `${requestSchemaName}.${propName}`,
      type: 'TextInput',
      label: prop['x-vault-displayAttrs']?.name || formatLabel(propName),
      helperText: prop.description,
    });
  }

  // Returns a converted `group` object to the correct format for sections:
  // [{ name: 'default', fields: [...] }, { name: 'Advanced', fields: [...] }, ...]
  return Object.entries(groups).map(([name, fields]) => ({ name, fields }));
};

export const prepFormConfig = (spec, methodName) => {
  const operation = getOperationDetails(spec, dasherize(methodName));

  if (!operation) return null;

  return {
    name: methodName,
    description: operation.description,
    payload: buildPayloadFromOperation(operation),
    sections: buildSectionsFromOperation(operation),
    apiClass: API_CLASS_FROM_TAG[operation.tag],
    requestType: getRequestType(operation.tag, methodName),
  };
};

export const generateConfigContent = (config) => {
  return `
    /**
      * Copyright IBM Corp. 2016, 2025
      * SPDX-License-Identifier: BUSL-1.1
    */

    // ⚠️ AUTO-GENERATED FILE - DO NOT EDIT
    // This file is generated from openapi.json
    // To customize this form, create an override in
    // forms/v2/overrides/

    import type ApiService from 'vault/services/api';
    import type { FormConfig } from '../form-config';
    import type { ${config.requestType} } from '@hashicorp/vault-client-typescript';

    /**
     * Form configuration for ${config.name}
     * Auto-generated from OpenAPI specification
     */
    const ${config.name}Config: FormConfig<${config.requestType},unknown> = {
      name: '${config.name}',
      description: '${config.description}',
      submit: async (api: ApiService, payload: ${config.requestType}) => {
        return await api.${config.apiClass}.${config.name}Raw(payload);
      },
      payload: ${JSON.stringify(config.payload, null, 2)},
      sections: ${JSON.stringify(config.sections, null, 2)},
    };

    export default ${config.name}Config;
  `;
};
