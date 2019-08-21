/*
  This service is used to pull an OpenAPI document describing the
  shape of data at a specific path to hydrate a model with attrs it
  has less (or no) information about.
*/
import Service from '@ember/service';
import DS from 'ember-data';
import { encodePath } from 'vault/utils/path-encoding-helpers';
import { getOwner } from '@ember/application';
import { assign } from '@ember/polyfills';
import { dasherize } from '@ember/string';
import { singularize } from 'ember-inflector';
import { expandOpenApiProps, combineAttributes } from 'vault/utils/openapi-to-attrs';
import fieldToAttrs from 'vault/utils/field-to-attrs';
import { resolve, reject } from 'rsvp';
import { debug } from '@ember/debug';
import { apiPath as dynamicPath } from 'vault/utils/api-path';

const { belongsTo } = DS;

import generatedItemAdapter from 'vault/adapters/generated-item-list';
export function sanitizePath(path) {
  //remove whitespace + remove trailing and leading slashes
  return path.trim().replace(/^\/+|\/+$/g, '');
}

export default Service.extend({
  attrs: null,
  ajax(url, options = {}) {
    let appAdapter = getOwner(this).lookup(`adapter:application`);
    let { data } = options;
    return appAdapter.ajax(url, 'GET', {
      data,
    });
  },

  getNewModel(modelType, backend, apiPath, itemType, itemID) {
    let owner = getOwner(this);
    const modelName = `model:${modelType}`;
    const modelFactory = owner.factoryFor(modelName);
    let newModel, helpUrl;
    //if we have a factory, we need to take the existing model into account
    if (modelFactory) {
      debug(`Model factory found for ${modelType}`);
      newModel = modelFactory.class;
      const modelProto = newModel.proto();
      if (newModel.merged || modelProto.useOpenAPI !== true) {
        return resolve();
      }
      // helpUrl = modelProto.getHelpUrl(backend);
      // return this.registerNewModelWithProps(helpUrl, backend, newModel, modelName);
    } else {
      debug(`Creating new Model for ${modelType}`);
      newModel = DS.Model.extend({});
    }

    //use paths to dynamically create our openapi help url
    //if we have a brand new model
    return this.getPaths(apiPath, backend, itemType, itemID)
      .then(pathInfo => {
        const adapterFactory = owner.factoryFor(`adapter:${modelType}`);
        //if we have an adapter already use that, otherwise create one
        if (!adapterFactory) {
          debug(`Creating new adapter for ${modelType}`);
          const adapter = this.getNewAdapter(pathInfo, itemType);
          owner.register(`adapter:${modelType}`, adapter);
        }
        let path, paths;
        //if we have an item we want the create info for that itemType
        if (itemType) {
          paths = this.filterPathsByItemType(pathInfo, itemType);
        }
        const createPath = paths.find(path => path.operations.includes('post') && path.action !== 'Delete');
        path = createPath.path;
        path = path.includes('{') ? path.slice(0, path.indexOf('{') - 1) + '/example' : path;
        if (!path) {
          return reject();
        }

        helpUrl = `/v1/${apiPath}${path.slice(1)}?help=true` || newModel.proto().getHelpUrl(backend);
        pathInfo.paths = paths;
        newModel = newModel.extend({ paths: pathInfo });
        return this.registerNewModelWithProps(helpUrl, backend, newModel, modelName);
      })
      .catch(err => {
        debugger;
      });
  },

  reducePathsByPathName(pathInfo, currentPath) {
    //will be passed something like { apiPath, itemType }
    const pathName = currentPath[0];
    const pathDetails = currentPath[1];
    const displayAttrs = pathDetails['x-vault-displayAttrs'];

    if (!displayAttrs) {
      return pathInfo;
    }

    if (pathName.includes('{')) {
      //we need to know if there are url params
      pathName.split('{')[1].split('}')[0];
    }

    let itemType, itemName;
    if (displayAttrs.itemType) {
      itemType = displayAttrs.itemType;
      let items = itemType.split(':');
      itemName = items[items.length - 1];
      items = items.map(item => dasherize(singularize(item.toLowerCase())));
      itemType = items.join('_');
    }

    if (itemType && !pathInfo.itemTypes.includes(itemType)) {
      pathInfo.itemTypes.push(itemType);
    }

    let operations = [];
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
    return pathInfo.paths.filter(path => {
      return itemType === path.itemType || path.itemType.indexOf(`${itemType}_`) === 0;
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

  //Makes a call to grab the OpenAPI document.
  //Returns relevant information from OpenAPI
  //as determined by the expandOpenApiProps util
  getProps(helpUrl, backend) {
    debug(`Fetching schema properties for ${backend} from ${helpUrl}`);

    return this.ajax(helpUrl, backend).then(help => {
      //paths is an array but it will have a single entry
      // for the scope we're in
      const path = Object.keys(help.openapi.paths)[0];
      const pathInfo = help.openapi.paths[path];
      const params = pathInfo.parameters;
      let paramProp = {};

      //include url params
      if (params) {
        const { name, schema, description } = params[0];
        let label = name.split('_').join(' ');

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

      //TODO: handle post endpoints without requestBody
      const props = pathInfo.post.requestBody.content['application/json'].schema.properties;
      //put url params (e.g. {name}, {role})
      //at the front of the props list
      const newProps = assign({}, paramProp, props);
      return expandOpenApiProps(newProps);
    });
  },

  replaceParamInPath(path, id) {
    let pathParts = path.split('{');
    let pathBeginning = pathParts[0];
    pathParts = pathParts[1].split('}');
    let pathEnd = pathParts[1];
    return `${pathBeginning.slice(1)}${id}${pathEnd}`;
  },

  getNewAdapter(pathInfo, itemType) {
    //we need list and create paths to set the correct urls for actions
    let paths = this.filterPathsByItemType(pathInfo, itemType);
    let { apiPath } = pathInfo;
    if (pathInfo.itemID) {
      paths.forEach(path => {
        if (path.path.includes('{')) {
          path.path = this.replaceParamInPath(path.path, pathInfo.itemID);
        }
      });
    }
    const getPath = paths.find(path => path.operations.includes('get'));
    const createPath = paths.find(path => path.action === 'Create' || path.operations.includes('post'));
    const deletePath = paths.find(path => path.operations.includes('delete'));

    return generatedItemAdapter.extend({
      urlForItem(method, id) {
        debugger;
        let url = `${this.buildURL()}/${apiPath}${getPath.path.slice(1)}/`;
        if (id) {
          url = url + encodePath(id);
        }
        return url;
      },

      urlForFindRecord(id, modelName) {
        return this.urlForItem(modelName, id);
      },

      //urlForQuery if there is an id and we are listing, use the id to construct the path

      urlForUpdateRecord(id) {
        let path = createPath.path.slice(1, createPath.path.indexOf('{') - 1);
        return `${this.buildURL()}/${apiPath}${path}/${id}`;
      },

      urlForCreateRecord(modelType, snapshot) {
        const { id } = snapshot;
        let path = createPath.path.slice(1, createPath.path.indexOf('{') - 1);
        return `${this.buildURL()}/${apiPath}${path}/${id}`;
      },

      urlForDeleteRecord(id) {
        let path = deletePath.path.slice(1, deletePath.path.indexOf('{') - 1);
        return `${this.buildURL()}/${apiPath}${path}/${id}`;
      },
    });
  },

  registerNewModelWithProps(helpUrl, backend, newModel, modelName) {
    return this.getProps(helpUrl, backend).then(props => {
      const { attrs, newFields } = combineAttributes(newModel.attributes, props);
      let owner = getOwner(this);
      newModel = newModel.extend(attrs, { newFields });
      //if our newModel doesn't have fieldGroups already
      //we need to create them
      try {
        let fieldGroups = newModel.proto().fieldGroups;
        if (!fieldGroups) {
          debug(`Constructing fieldGroups for ${backend}`);
          fieldGroups = this.getFieldGroups(newModel);
          newModel = newModel.extend({ fieldGroups });
        }
      } catch (err) {
        //eat the error, fieldGroups is computed in the model definition
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
      //if the attr comes in with a fieldGroup from OpenAPI,
      //add it to that group
      if (attr.options.fieldGroup) {
        if (groups[attr.options.fieldGroup]) {
          groups[attr.options.fieldGroup].push(attr.name);
        } else {
          groups[attr.options.fieldGroup] = [attr.name];
        }
      } else {
        //otherwise just add that attr to the default group
        groups.default.push(attr.name);
      }
    });
    for (let group in groups) {
      fieldGroups.push({ [group]: groups[group] });
    }
    return fieldToAttrs(newModel, fieldGroups);
  },
});
