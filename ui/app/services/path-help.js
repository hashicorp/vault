/*
  This service is used to pull an OpenAPI document describing the
  shape of data at a specific path to hydrate a model with attrs it
  has less (or no) information about.
*/
import Model from '@ember-data/model';
import Service from '@ember/service';
import { encodePath } from 'vault/utils/path-encoding-helpers';
import { getOwner } from '@ember/application';
import { assign } from '@ember/polyfills';
import { expandOpenApiProps, combineAttributes } from 'vault/utils/openapi-to-attrs';
import fieldToAttrs from 'vault/utils/field-to-attrs';
import { resolve, reject } from 'rsvp';
import { debug } from '@ember/debug';
import { dasherize, capitalize } from '@ember/string';
import { singularize } from 'ember-inflector';
import buildValidations from 'vault/utils/build-api-validators';

import generatedItemAdapter from 'vault/adapters/generated-item-list';
export function sanitizePath(path) {
  // remove whitespace + remove trailing and leading slashes
  return path.trim().replace(/^\/+|\/+$/g, '');
}

export default Service.extend({
  attrs: null,
  dynamicApiPath: '',
  ajax(url, options = {}) {
    let appAdapter = getOwner(this).lookup(`adapter:application`);
    let { data } = options;
    return appAdapter.ajax(url, 'GET', {
      data,
    });
  },

  getNewModel(modelType, backend, apiPath, itemType) {
    let owner = getOwner(this);
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
      .then(pathInfo => {
        const adapterFactory = owner.factoryFor(`adapter:${modelType}`);
        // if we have an adapter already use that, otherwise create one
        if (!adapterFactory) {
          debug(`Creating new adapter for ${modelType}`);
          const adapter = this.getNewAdapter(pathInfo, itemType);
          owner.register(`adapter:${modelType}`, adapter);
        }
        let path, paths;
        // if we have an item we want the create info for that itemType
        paths = itemType ? this.filterPathsByItemType(pathInfo, itemType) : pathInfo.paths;
        const createPath = paths.find(path => path.operations.includes('post') && path.action !== 'Delete');
        path = createPath.path;
        path = path.includes('{') ? path.slice(0, path.indexOf('{') - 1) + '/example' : path;
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
      .catch(err => {
        // TODO: we should handle the error better here
        console.error(err);
      });
  },

  reducePathsByPathName(pathInfo, currentPath) {
    const pathName = currentPath[0];
    const pathDetails = currentPath[1];
    const displayAttrs = pathDetails['x-vault-displayAttrs'];

    if (!displayAttrs) {
      return pathInfo;
    }

    let itemType, itemName;
    if (displayAttrs.itemType) {
      itemType = displayAttrs.itemType;
      let items = itemType.split(':');
      itemName = items[items.length - 1];
      items = items.map(item => dasherize(singularize(item.toLowerCase())));
      itemType = items.join('~*');
    }

    if (itemType && !pathInfo.itemTypes.includes(itemType)) {
      pathInfo.itemTypes.push(itemType);
    }

    const operations = [];
    if (pathDetails.get) {
      operations.push('get');
    }
    if (pathDetails.post) {
      operations.push('post');
    }
    if (pathDetails.delete) {
      operations.push('delete');
    }
    if (pathDetails.get && pathDetails.get.parameters && pathDetails.get.parameters[0].name === 'list') {
      operations.push('list');
    }

    pathInfo.paths.push({
      path: pathName,
      itemType: itemType || displayAttrs.itemType,
      itemName: itemName || pathInfo.itemType || displayAttrs.itemType,
      operations,
      action: displayAttrs.action,
      navigation: displayAttrs.navigation === true,
      param: pathName.includes('{') ? pathName.split('{')[1].split('}')[0] : false,
    });

    return pathInfo;
  },

  filterPathsByItemType(pathInfo, itemType) {
    if (!itemType) {
      return pathInfo.paths;
    }
    return pathInfo.paths.filter(path => {
      return itemType === path.itemType;
    });
  },

  getPaths(apiPath, backend, itemType, itemID) {
    let debugString =
      itemID && itemType
        ? `Fetching relevant paths for ${backend} ${itemType} ${itemID} from ${apiPath}`
        : `Fetching relevant paths for ${backend} ${itemType} from ${apiPath}`;
    debug(debugString);
    return this.ajax(`/v1/${apiPath}?help=1`, backend).then(help => {
      const pathInfo = help.openapi.paths;
      let paths = Object.entries(pathInfo);

      return paths.reduce(this.reducePathsByPathName, {
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

    return this.ajax(helpUrl, backend).then(help => {
      // paths is an array but it will have a single entry
      // for the scope we're in
      const path = Object.keys(help.openapi.paths)[0]; // do this or look at name
      const pathInfo = help.openapi.paths[path];
      const params = pathInfo.parameters;
      let paramProp = {};

      // include url params
      if (params) {
        const { name, schema, description } = params[0];
        let label = capitalize(name.split('_').join(' '));

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

      // TODO: handle post endpoints without requestBody
      const props = pathInfo.post
        ? pathInfo.post.requestBody.content['application/json'].schema.properties
        : {};
      // put url params (e.g. {name}, {role})
      // at the front of the props list
      const newProps = assign({}, paramProp, props);
      return expandOpenApiProps(newProps);
    });
  },

  getNewAdapter(pathInfo, itemType) {
    // we need list and create paths to set the correct urls for actions
    let paths = this.filterPathsByItemType(pathInfo, itemType);
    let { apiPath } = pathInfo;
    const getPath = paths.find(path => path.operations.includes('get'));

    // the action might be "Generate" or something like that so we'll grab the first post endpoint if there
    // isn't one with "Create"
    // TODO: look into a more sophisticated way to determine the create endpoint
    const createPath = paths.find(path => path.action === 'Create' || path.operations.includes('post'));
    const deletePath = paths.find(path => path.operations.includes('delete'));

    return generatedItemAdapter.extend({
      urlForItem(id, isList, dynamicApiPath) {
        const itemType = getPath.path.slice(1);
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
        const { id } = snapshot;
        const path = createPath.path.slice(1, createPath.path.indexOf('{') - 1);
        return `${this.buildURL()}/${apiPath}${path}/${id}`;
      },

      urlForDeleteRecord(id) {
        const path = deletePath.path.slice(1, deletePath.path.indexOf('{') - 1);
        return `${this.buildURL()}/${apiPath}${path}/${id}`;
      },
    });
  },

  registerNewModelWithProps(helpUrl, backend, newModel, modelName) {
    return this.getProps(helpUrl, backend).then(props => {
      const { attrs, newFields } = combineAttributes(newModel.attributes, props);
      let owner = getOwner(this);
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
          let validations = buildValidations(fieldGroups);
          newModel = newModel.extend(validations);
        }
      } catch (err) {
        // eat the error, fieldGroups is computed in the model definition
      }
      newModel.reopenClass({ merged: true });
      owner.unregister(modelName);
      owner.register(modelName, newModel);
    });
  },
  getFieldGroups(newModel) {
    let groups = {
      default: [],
    };
    let fieldGroups = [];
    newModel.attributes.forEach(attr => {
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
    for (let group in groups) {
      fieldGroups.push({ [group]: groups[group] });
    }
    return fieldToAttrs(newModel, fieldGroups);
  },
});
