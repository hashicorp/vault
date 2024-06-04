/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/*
  This service is used to pull an OpenAPI document describing the
  shape of data at a specific path to hydrate a model with attrs it
  has less (or no) information about.
*/
import Model from '@ember-data/model';
import Service from '@ember/service';
import { encodePath } from 'vault/utils/path-encoding-helpers';
import { getOwner } from '@ember/application';
import { expandOpenApiProps, combineAttributes } from 'vault/utils/openapi-to-attrs';
import fieldToAttrs from 'vault/utils/field-to-attrs';
import { resolve, reject } from 'rsvp';
import { debug } from '@ember/debug';
import { capitalize } from '@ember/string';
import { computed } from '@ember/object'; // eslint-disable-line
import { withModelValidations } from 'vault/decorators/model-validations';

import generatedItemAdapter from 'vault/adapters/generated-item-list';
import { sanitizePath } from 'core/utils/sanitize-path';
import {
  filterPathsByItemType,
  pathToHelpUrlSegment,
  reducePathsByPathName,
} from 'vault/utils/openapi-helpers';
import { isPresent } from '@ember/utils';

export default Service.extend({
  attrs: null,
  dynamicApiPath: '',
  ajax(url, options = {}) {
    const appAdapter = getOwner(this).lookup(`adapter:application`);
    const { data } = options;
    return appAdapter.ajax(url, 'GET', {
      data,
    });
  },

  /**
   * getNewModel instantiates models which use OpenAPI fully or partially
   * @param {string} modelType
   * @param {string} backend
   * @param {string} apiPath (optional) if passed, this method will call getPaths and build submodels for item types
   * @param {*} itemType (optional) used in getPaths for additional models
   * @returns void - as side effect, registers model via registerNewModelWithProps
   */
  getNewModel(modelType, backend, apiPath, itemType) {
    const owner = getOwner(this);
    const modelName = `model:${modelType}`;

    const modelFactory = owner.factoryFor(modelName);
    let newModel, helpUrl;
    // if we have a factory, we need to take the existing model into account
    if (modelFactory) {
      debug(`Model factory found for ${modelType}`);
      newModel = modelFactory.class;
      const modelProto = newModel.proto();
      if (newModel.merged || modelProto.useOpenAPI !== true) {
        return resolve();
      }

      helpUrl = modelProto.getHelpUrl(backend);
      return this.registerNewModelWithProps(helpUrl, backend, newModel, modelName);
    } else {
      debug(`Creating new Model for ${modelType}`);
      newModel = Model.extend({});
    }

    // we don't have an apiPath for dynamic secrets
    // and we don't need paths for them yet
    if (!apiPath) {
      helpUrl = newModel.proto().getHelpUrl(backend);
      return this.registerNewModelWithProps(helpUrl, backend, newModel, modelName);
    }

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

        helpUrl = `/v1/${apiPath}${path.slice(1)}?help=true` || newModel.proto().getHelpUrl(backend);
        pathInfo.paths = paths;
        newModel = newModel.extend({ paths: pathInfo });
        return this.registerNewModelWithProps(helpUrl, backend, newModel, modelName);
      })
      .catch((err) => {
        // TODO: we should handle the error better here
        console.error(err); // eslint-disable-line
      });
  },

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
  },

  // Makes a call to grab the OpenAPI document.
  // Returns relevant information from OpenAPI
  // as determined by the expandOpenApiProps util
  getProps(helpUrl, backend) {
    // add name of thing you want
    debug(`Fetching schema properties for ${backend} from ${helpUrl}`);

    return this.ajax(helpUrl, backend).then((help) => {
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
      // put url params (e.g. {name}, {role})
      // at the front of the props list
      const newProps = { ...paramProp, ...props };
      return expandOpenApiProps(newProps);
    });
  },

  getNewAdapter(pathInfo, itemType) {
    // we need list and create paths to set the correct urls for actions
    const paths = filterPathsByItemType(pathInfo, itemType);
    let { apiPath } = pathInfo;
    const getPath = paths.find((path) => path.operations.includes('get'));

    // the action might be "Generate" or something like that so we'll grab the first post endpoint if there
    // isn't one with "Create"
    // TODO: look into a more sophisticated way to determine the create endpoint
    const createPath = paths.find((path) => path.action === 'Create' || path.operations.includes('post'));
    const deletePath = paths.find((path) => path.operations.includes('delete'));

    return generatedItemAdapter.extend({
      urlForItem(id, isList, dynamicApiPath) {
        const itemType = sanitizePath(getPath.path);
        let url;
        id = encodePath(id);
        // the apiPath changes when you switch between routes but the apiPath variable does not unless the model is reloaded
        // overwrite apiPath if dynamicApiPath exist.
        // dynamicApiPath comes from the model->adapter
        if (dynamicApiPath) {
          apiPath = dynamicApiPath;
        }
        // isList indicates whether we are viewing the list page
        // of a top-level item such as userpass
        if (isList) {
          url = `${this.buildURL()}/${apiPath}${itemType}/`;
        } else {
          // build the URL for the show page of a nested item
          // such as a userpass group
          url = `${this.buildURL()}/${apiPath}${itemType}/${id}`;
        }

        return url;
      },

      urlForQueryRecord(id, modelName) {
        return this.urlForItem(id, modelName);
      },

      urlForUpdateRecord(id) {
        const itemType = createPath.path.slice(1, createPath.path.indexOf('{') - 1);
        return `${this.buildURL()}/${apiPath}${itemType}/${id}`;
      },

      urlForCreateRecord(modelType, snapshot) {
        const id = snapshot.record.mutableId; // computed property that returns either id or private settable _id value
        const path = createPath.path.slice(1, createPath.path.indexOf('{') - 1);
        return `${this.buildURL()}/${apiPath}${path}/${id}`;
      },

      urlForDeleteRecord(id) {
        const path = deletePath.path.slice(1, deletePath.path.indexOf('{') - 1);
        return `${this.buildURL()}/${apiPath}${path}/${id}`;
      },

      createRecord(store, type, snapshot) {
        return this._super(...arguments).then((response) => {
          // if the server does not return an id and one has not been set on the model we need to set it manually from the mutableId value
          if (!response?.id && !snapshot.record.id) {
            snapshot.record.id = snapshot.record.mutableId;
            snapshot.id = snapshot.record.id;
          }
          return response;
        });
      },
    });
  },

  registerNewModelWithProps(helpUrl, backend, newModel, modelName) {
    return this.getProps(helpUrl, backend).then((props) => {
      const { attrs, newFields } = combineAttributes(newModel.attributes, props);
      const owner = getOwner(this);
      newModel = newModel.extend(attrs, { newFields });
      // if our newModel doesn't have fieldGroups already
      // we need to create them
      try {
        // Initialize prototype to access field groups
        let fieldGroups = newModel.proto().fieldGroups;
        if (!fieldGroups) {
          debug(`Constructing fieldGroups for ${backend}`);
          fieldGroups = this.getFieldGroups(newModel);
          newModel = newModel.extend({ fieldGroups });
          // Build and add validations on model
          // NOTE: For initial phase, initialize validations only for user pass auth
          if (backend === 'userpass') {
            const validations = {
              password: [
                {
                  validator(model) {
                    return (
                      !(isPresent(model.password) && isPresent(model.passwordHash)) &&
                      (isPresent(model.password) || isPresent(model.passwordHash))
                    );
                  },
                  message: 'You must provide either password or password hash, but not both.',
                },
              ],
            };
            @withModelValidations(validations)
            class GeneratedItemModel extends newModel {}
            newModel = GeneratedItemModel;
          }
        }
      } catch (err) {
        // eat the error, fieldGroups is computed in the model definition
      }
      // attempting to set the id prop on a model will trigger an error
      // this computed will be used in place of the the id fieldValue -- see openapi-to-attrs
      newModel.reopen({
        mutableId: computed('id', '_id', {
          get() {
            return this._id || this.id;
          },
          set(key, value) {
            return (this._id = value);
          },
        }),
      });
      newModel.reopenClass({ merged: true });
      owner.unregister(modelName);
      owner.register(modelName, newModel);
    });
  },
  getFieldGroups(newModel) {
    const groups = {
      default: [],
    };
    const fieldGroups = [];
    newModel.attributes.forEach((attr) => {
      // if the attr comes in with a fieldGroup from OpenAPI,
      // add it to that group
      if (attr.options.fieldGroup) {
        if (groups[attr.options.fieldGroup]) {
          groups[attr.options.fieldGroup].push(attr.name);
        } else {
          groups[attr.options.fieldGroup] = [attr.name];
        }
      } else {
        // otherwise just add that attr to the default group
        groups.default.push(attr.name);
      }
    });
    for (const group in groups) {
      fieldGroups.push({ [group]: groups[group] });
    }
    return fieldToAttrs(newModel, fieldGroups);
  },
});
