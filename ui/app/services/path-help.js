/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

/** @module path-help Runtime OpenAPI discovery and form metadata for mounted plugins. */
import Service from '@ember/service';
import { getOwner } from '@ember/owner';
import { debug } from '@ember/debug';
import { capitalize } from '@ember/string';

import { reducePathsByPathName, expandOpenApiProps } from 'vault/utils/openapi-helpers';
import FormField from 'vault/utils/forms/field';
import FormFieldGroup from 'vault/utils/forms/field-group';

export default class PathHelpService extends Service {
  ajax(url, options = {}) {
    const appAdapter = getOwner(this).lookup(`adapter:application`);
    const { data } = options;
    return appAdapter.ajax(url, 'GET', {
      data,
    });
  }

  /**
   * getPaths is used to fetch all the openAPI paths available for an auth method,
   * to populate the tab navigation in each specific method page
   * @param {string} apiPath path of openApi
   * @param {string} backend backend name, mostly for debug purposes
   * @param {string} itemType optional
   * @param {string} itemID optional - ID of specific item being fetched
   * @returns PathsInfo
   */
  getPaths(apiPath, backend, itemType, itemID) {
    const debugString =
      itemID && itemType
        ? `Fetching relevant paths for ${backend} ${itemType} ${itemID} from ${apiPath}`
        : `Fetching relevant paths for ${backend} ${itemType} from ${apiPath}`;
    debug(debugString);
    return this.ajax(`/v1/${apiPath}?help=1`, backend).then((help) => {
      const pathInfo = help.openapi.paths;
      const paths = Object.entries(pathInfo);

      return paths.reduce(reducePathsByPathName, {
        apiPath,
        itemType,
        itemTypes: [],
        paths: [],
        itemID,
      });
    });
  }

  /**
   * Fetches runtime schema properties; legacy consumers retain camelCase model metadata.
   * @param {string} helpUrl
   * @param {'model' | 'form'} outputFormat
   */
  getProps(helpUrl, outputFormat = 'model') {
    // add name of thing you want
    debug(`Fetching schema properties from ${helpUrl}`);

    return this.ajax(helpUrl).then((help) => {
      // paths is an array but it will have a single entry
      // for the scope we're in
      const path = Object.keys(help.openapi.paths)[0]; // do this or look at name
      const pathInfo = help.openapi.paths[path];
      const schema = pathInfo?.post?.requestBody?.content?.['application/json']?.schema;
      if (!schema) {
        throw new Error(`OpenAPI request schema not found for ${helpUrl}`);
      }
      const params = pathInfo.parameters;
      const paramProp = {};

      // include url params
      if (params) {
        const { name, schema, description } = params[0];
        const label = capitalize(name.split('_').join(' '));

        paramProp[name] = {
          'x-vault-displayAttrs': {
            name: label,
            group: 'default',
          },
          type: schema.type,
          description: description,
          isId: true,
        };
      }

      let props;
      if (schema.$ref) {
        // $ref will be shaped like `#/components/schemas/MyResponseType
        // which maps to the location of the item within the openApi response
        const loc = schema.$ref.replace('#/', '').split('/');
        props = loc.reduce((prev, curr) => prev?.[curr], help.openapi)?.properties;
      } else if (schema.properties) {
        props = schema.properties;
      }
      if (!props) {
        throw new Error(`OpenAPI request properties not found for ${helpUrl}`);
      }
      // put url params (e.g. {name}, {role}) at the front of the props list
      const newProps = { ...paramProp, ...props };
      return expandOpenApiProps(newProps, outputFormat);
    });
  }

  /**
   * Returns snake_case fields and ordered groups without registering a model or adapter.
   * The path identifier (e.g. `username`) is marked read-only and returned as `identifier`,
   * since it is sent in the request URL rather than the body.
   * @param {string} helpUrl Runtime help URL for a plugin's create endpoint.
   * @returns {Promise<import('vault/services/path-help').FormMetadata>}
   */
  async getFormMetadata(helpUrl) {
    const props = await this.getProps(helpUrl, 'form');
    const formFields = [];
    const groups = new Map([['default', []]]);
    let identifier;
    for (const [name, prop] of Object.entries(props)) {
      // Path identifiers bind to their own form data, not Ember Data's mutableId.
      // eslint-disable-next-line no-unused-vars
      const { type, fieldGroup, fieldValue, readOnly, ...options } = prop;
      const field = new FormField(name, type, readOnly ? { ...options, readOnly } : options);
      if (readOnly) {
        identifier = name;
      }
      formFields.push(field);
      if (!groups.has(fieldGroup)) {
        groups.set(fieldGroup, []);
      }
      groups.get(fieldGroup).push(field);
    }
    const formFieldGroups = [...groups].map(([name, fields]) => new FormFieldGroup(name, fields));
    return { formFields, formFieldGroups, identifier };
  }
}
