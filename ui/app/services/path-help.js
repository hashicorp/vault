/*
  This service is used to pull an OpenAPI document describing the
  shape of data at a specific path to hydrate a model with attrs it
  has less (or no) information about.
*/
import Service from '@ember/service';
import DS from 'ember-data';
import { encodePath } from 'vault/utils/path-encoding-helpers';
import { getOwner } from '@ember/application';
import { capitalize } from '@ember/string';
import { assign } from '@ember/polyfills';
import { expandOpenApiProps, combineAttributes } from 'vault/utils/openapi-to-attrs';
import { supportedAuthBackends } from 'vault/helpers/supported-auth-backends';
import fieldToAttrs from 'vault/utils/field-to-attrs';
import { resolve } from 'rsvp';
import { debug } from '@ember/debug';

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

  getNewModel(modelType, backend, apiPath, itemType) {
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
      helpUrl = modelProto.getHelpUrl(backend);
      return this.registerNewModelWithProps(helpUrl, backend, newModel, modelName);
    } else {
      debug(`Creating new Model for ${modelType}`);
      newModel = DS.Model.extend({});
      //use paths to dynamically create our openapi help url
      //if we have a brand new model
      return this.getPaths(apiPath, backend, itemType).then(paths => {
        const adapterFactory = owner.factoryFor(`adapter:${modelType}`);
        //if we have an adapter already use that, otherwise create one
        if (!adapterFactory) {
          debug(`Creating new adapter for ${modelType}`);
          const adapter = this.getNewAdapter(paths, itemType);
          owner.register(`adapter:${modelType}`, adapter);
        }
        //if we have an item we want the create info for that itemType
        let path;
        if (itemType) {
          const createPath = paths.create.find(path => path.path.includes(itemType));
          path = createPath.path;
          path = path.slice(0, path.indexOf('{') - 1) + '/example';
        } else {
          //we need the mount config
          path = paths.configPath[0].path;
        }
        helpUrl = `/v1/${apiPath}${path}?help=true`;
        return this.registerNewModelWithProps(helpUrl, backend, newModel, modelName);
      });
    }
  },

  reducePaths(paths, currentPath) {
    const pathName = currentPath[0];
    const pathInfo = currentPath[1];

    //config is a get/post endpoint that doesn't take route params
    //and isn't also a list endpoint and has an Action of Configure
    if (
      pathInfo.post &&
      pathInfo.get &&
      (pathInfo['x-vault-displayAttrs'] && pathInfo['x-vault-displayAttrs'].action === 'Configure')
    ) {
      paths.configPath.push({ path: pathName });
      return paths; //config path should only be config path
    }

    //list endpoints all have { name: "list" } in their get parameters
    if (pathInfo.get && pathInfo.get.parameters && pathInfo.get.parameters[0].name === 'list') {
      paths.list.push({ path: pathName });
    }

    if (pathInfo.delete) {
      paths.delete.push({ path: pathName });
    }

    //create endpoints have path an action (e.g. "Create" or "Generate")
    if (pathInfo.post && pathInfo['x-vault-displayAttrs'] && pathInfo['x-vault-displayAttrs'].action) {
      paths.create.push({
        path: pathName,
        action: pathInfo['x-vault-displayAttrs'].action,
      });
    }

    if (pathInfo['x-vault-displayAttrs'] && pathInfo['x-vault-displayAttrs'].navigation) {
      paths.navPaths.push({ path: pathName });
    }

    return paths;
  },

  getPaths(apiPath, backend) {
    debug(`Fetching relevant paths for ${backend} from ${apiPath}`);
    return this.ajax(`/v1/${apiPath}?help=1`, backend).then(help => {
      const pathInfo = help.openapi.paths;
      let paths = Object.entries(pathInfo);

      return paths.reduce(this.reducePaths, {
        apiPath: apiPath,
        configPath: [],
        list: [],
        create: [],
        delete: [],
        navPaths: [],
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

  getNewAdapter(paths, itemType) {
    //we need list and create paths to set the correct urls for actions
    const { list, create, apiPath } = paths;
    const createPath = create.find(path => path.path.includes(itemType));
    const listPath = list.find(pathInfo => pathInfo.path.includes(itemType));
    const deletePath = paths.delete.find(path => path.path.includes(itemType));
    return generatedItemAdapter.extend({
      urlForItem(method, id) {
        let { path } = listPath;
        let url = `${this.buildURL()}/${apiPath}${path}/`;
        if (id) {
          url = url + encodePath(id);
        }
        return url;
      },

      urlForFindRecord(id, modelName, snapshot) {
        return this.urlForItem(modelName, id, snapshot);
      },

      urlForUpdateRecord(id) {
        let { path } = createPath;
        path = path.slice(0, path.indexOf('{') - 1);
        return `${this.buildURL()}/${apiPath}${path}/${id}`;
      },

      urlForCreateRecord(modelType, snapshot) {
        const { id } = snapshot;
        let { path } = createPath;
        path = path.slice(0, path.indexOf('{') - 1);
        return `${this.buildURL()}/${apiPath}${path}/${id}`;
      },

      urlForDeleteRecord(id) {
        let { path } = deletePath;
        path = path.slice(0, path.indexOf('{') - 1);
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
