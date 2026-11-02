/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

/** @module path-help Runtime OpenAPI discovery and form metadata for mounted plugins. */
import Model, { attr } from '@ember-data/model';
import Service from '@ember/service';
import { getOwner } from '@ember/owner';
import { resolve, reject } from 'rsvp';
import { debug } from '@ember/debug';
import { capitalize } from '@ember/string';
import { computed } from '@ember/object'; // eslint-disable-line

import {
  filterPathsByItemType,
  pathToHelpUrlSegment,
  reducePathsByPathName,
  combineOpenApiAttrs,
  expandOpenApiProps,
} from 'vault/utils/openapi-helpers';
import GeneratedItemModel from 'vault/models/generated-item';
import GeneratedItemListAdapter from 'vault/adapters/generated-item-list';
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
   * Registers new ModelClass at specified model type, and busts cache
   */
  _registerModel(owner, NewKlass, modelType, isNew = false) {
    const store = owner.lookup('service:store');
    // bust cache in ember's registry
    if (!isNew) {
      owner.unregister('model:' + modelType);
    }
    owner.register('model:' + modelType, NewKlass);

    // bust cache in EmberData's model lookup
    delete store._modelFactoryCache[modelType];
  }

  /**
   * upgradeModelSchema takes an existing ModelClass and hydrates it with the passed attributes
   * @param {ModelClass} Klass model class retrieved with store.modelFor(modelType)
   * @param {Attribute[]} attrs array of attributes {name, type, options}
   * @returns new ModelClass extended from passed one, with the passed attributes added
   */
  _upgradeModelSchema(Klass, attrs, newFields) {
    // extending the class will ensure that static schema lookups regenerate
    const NewKlass = class extends Klass {};

    for (const { name, type, options } of attrs) {
      const decorator = attr(type, options);
      const descriptor = decorator(NewKlass.prototype, name, {});
      Object.defineProperty(NewKlass.prototype, name, descriptor);
    }

    // newFields is used in combineFieldGroups within various models
    if (newFields) {
      NewKlass.prototype.newFields = newFields;
    }

    // Ensure this class doesn't get re-hydrated
    NewKlass.merged = true;

    return NewKlass;
  }

  /**
   * Temporary compatibility bridge for generated-item consumers until VAULT-49929.
   * getNewModel instantiates models which use OpenAPI to generate the model fully
   * @param {string} modelType
   * @param {string} backend
   * @param {string} apiPath this method will call getPaths and build submodels for item types
   * @param {*} itemType (optional) used in getPaths for additional models
   * @returns void - as side effect, registers model via registerNewModelWithAttrs
   */
  getNewModel(modelType, backend, apiPath, itemType) {
    const owner = getOwner(this);
    const modelName = `model:${modelType}`;

    const modelFactory = owner.factoryFor(modelName);

    if (modelFactory) {
      // if the modelFactory already exists, it means either this model was already
      // generated or the model exists in the code already. In either case resolve

      if (!modelFactory.class.merged) {
        // no merged flag means this model was not previously generated
        debug(`Model exists for ${modelType} -- using its existing schema`);
      }
      return resolve();
    }
    debug(`Creating new Model for ${modelType}`);
    let newModel = Model.extend({});

    // use paths to dynamically create our openapi help url
    // if we have a brand new model
    return this.getPaths(apiPath, backend, itemType)
      .then((pathInfo) => {
        const adapterFactory = owner.factoryFor(`adapter:${modelType}`);
        // if we have an adapter already use that, otherwise create one
        if (!adapterFactory) {
          debug(`Creating new adapter for ${modelType}`);
          const adapter = this.getNewAdapter(pathInfo, itemType);
          owner.register(`adapter:${modelType}`, adapter);
        }
        // if we have an item we want the create info for that itemType
        const paths = itemType ? filterPathsByItemType(pathInfo, itemType) : pathInfo.paths;
        const createPath = paths.find((path) => path.operations.includes('post') && path.action !== 'Delete');
        const path = pathToHelpUrlSegment(createPath.path);
        if (!path) {
          // TODO: we don't know if path will ever be falsey
          // if it is never falsey we can remove this.
          return reject();
        }

        const helpUrl = `/v1/${apiPath}${path.slice(1)}?help=true`;
        pathInfo.paths = paths;
        newModel = newModel.extend({ paths: pathInfo });
        return this.registerNewModelWithAttrs(helpUrl, modelType);
      })
      .catch((err) => {
        // TODO: we should handle the error better here
        console.error(err); // eslint-disable-line
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
   * @param {string} helpUrl Runtime help URL for a plugin's create endpoint.
   * @returns {Promise<import('vault/services/path-help').FormMetadata>}
   */
  async getFormMetadata(helpUrl) {
    const props = await this.getProps(helpUrl, 'form');
    const formFields = [];
    const groups = new Map([['default', []]]);
    for (const [name, prop] of Object.entries(props)) {
      // Path identifiers bind to their own form data, not Ember Data's mutableId.
      // eslint-disable-next-line no-unused-vars
      const { type, fieldGroup, fieldValue, readOnly, ...options } = prop;
      const field = new FormField(name, type, options);
      formFields.push(field);
      if (!groups.has(fieldGroup)) {
        groups.set(fieldGroup, []);
      }
      groups.get(fieldGroup).push(field);
    }
    const formFieldGroups = [...groups].map(([name, fields]) => new FormFieldGroup(name, fields));
    return { formFields, formFieldGroups };
  }

  getNewAdapter(pathInfo, itemType) {
    // we need list and create paths to set the correct urls for actions
    const paths = filterPathsByItemType(pathInfo, itemType);
    const { apiPath } = pathInfo;
    const getPath = paths.find((path) => path.operations.includes('get'));

    // the action might be "Generate" or something like that so we'll grab the first post endpoint if there
    // isn't one with "Create"
    // TODO: look into a more sophisticated way to determine the create endpoint
    const createPath = paths.find((path) => path.action === 'Create' || path.operations.includes('post'));
    const deletePath = paths.find((path) => path.operations.includes('delete'));

    return class NewAdapter extends GeneratedItemListAdapter {
      apiPath = apiPath;

      paths = {
        createPath: createPath?.path,
        deletePath: deletePath?.path,
        getPath: getPath?.path,
      };
    };
  }

  /**
   * registerNewModelWithAttrs takes the helpUrl of the given model type,
   * fetches props, and registers the model hydrated with the provided attrs
   * @param {string} helpUrl like /v1/auth/userpass2/users/example?help=true
   * @param {string} modelType like generated-user-userpass
   */
  async registerNewModelWithAttrs(helpUrl, modelType) {
    const owner = getOwner(this);
    const props = await this.getProps(helpUrl);
    const { attrs, newFields } = combineOpenApiAttrs(new Map(), props);
    const NewKlass = this._upgradeModelSchema(GeneratedItemModel, attrs, newFields);
    this._registerModel(owner, NewKlass, modelType, true);
  }
}
