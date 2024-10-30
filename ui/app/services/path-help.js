/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/*
  This service is used to pull an OpenAPI document describing the
  shape of data at a specific path to hydrate a model with attrs it
  has less (or no) information about.
*/
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
  getHelpUrlForModel,
  combineOpenApiAttrs,
  expandOpenApiProps,
} from 'vault/utils/openapi-helpers';
import GeneratedItemModel from 'vault/models/generated-item';
import GeneratedItemListAdapter from 'vault/adapters/generated-item-list';

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
   * hydrateModel instantiates models which use OpenAPI partially
   * @param {string} modelType path for model, eg pki/role
   * @param {string} backend path, which will be used for the generated helpUrl
   * @returns void - as side effect, re-registers model via upgradeModelSchema
   */
  async hydrateModel(modelType, backend) {
    const owner = getOwner(this);
    const helpUrl = getHelpUrlForModel(modelType, backend);
    const store = owner.lookup('service:store');
    const Klass = store.modelFor(modelType);

    if (Klass?.merged || !helpUrl) {
      // if the model is already merged, we don't need to do anything
      return resolve();
    }
    debug(`Hydrating model ${modelType} at backend ${backend}`);

    // fetch props from openAPI
    const props = await this.getProps(helpUrl);
    // combine existing attributes with openAPI data
    const { attrs, newFields } = combineOpenApiAttrs(Klass.attributes, props);
    debug(`${modelType} has ${newFields.length} new fields: ${newFields.join(', ')}`);

    // hydrate model
    const HydratedKlass = this._upgradeModelSchema(Klass, attrs, newFields);

    this._registerModel(owner, HydratedKlass, modelType);
  }

  /**
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
        debug(`Model exists for ${modelType} -- use hydrateModel instead`);
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

  // Makes a call to grab the OpenAPI document.
  // Returns relevant information from OpenAPI
  // as determined by the expandOpenApiProps util
  getProps(helpUrl) {
    // add name of thing you want
    debug(`Fetching schema properties from ${helpUrl}`);

    return this.ajax(helpUrl).then((help) => {
      // paths is an array but it will have a single entry
      // for the scope we're in
      const path = Object.keys(help.openapi.paths)[0]; // do this or look at name
      const pathInfo = help.openapi.paths[path];
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

      let props = {};
      const schema = pathInfo?.post?.requestBody?.content['application/json'].schema;
      if (schema.$ref) {
        // $ref will be shaped like `#/components/schemas/MyResponseType
        // which maps to the location of the item within the openApi response
        const loc = schema.$ref.replace('#/', '').split('/');
        props = loc.reduce((prev, curr) => {
          return prev[curr] || {};
        }, help.openapi).properties;
      } else if (schema.properties) {
        props = schema.properties;
      }
      // put url params (e.g. {name}, {role}) at the front of the props list
      const newProps = { ...paramProp, ...props };
      return expandOpenApiProps(newProps);
    });
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
