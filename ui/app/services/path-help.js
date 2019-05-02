/*
  This service is used to pull an OpenAPI document describing the
  shape of data at a specific path to hydrate a model with attrs it
  has less (or no) information about.
*/
import Service from '@ember/service';
import { encodePath } from 'vault/utils/path-encoding-helpers';
import { getOwner } from '@ember/application';
import { capitalize } from '@ember/string';
import { assign } from '@ember/polyfills';
import { expandOpenApiProps, combineAttributes } from 'vault/utils/openapi-to-attrs';
import fieldToAttrs from 'vault/utils/field-to-attrs';
import { resolve } from 'rsvp';
import DS from 'ember-data';
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

  //Makes a call to grab the OpenAPI document.
  //Returns relevant information from OpenAPI
  //as determined by the expandOpenApiProps util
  getProps(helpUrl, backend) {
    return this.ajax(helpUrl, backend).then(help => {
      //paths is an array but it will have a single entry
      // for the scope we're in
      const path = Object.keys(help.openapi.paths)[0];
      const pathInfo = help.openapi.paths[path];
      const params = pathInfo.parameters;
      let paramProp = {};

      //include url params
      if (params) {
        let label = capitalize(params[0].name);
        if (label.toLowerCase() !== 'name') {
          label += ' name';
        }
        paramProp[params[0].name] = {
          name: params[0].name,
          label: label,
          type: params[0].schema.type,
          description: params[0].description,
          isId: true,
        };
      }
      const props = pathInfo.post.requestBody.content['application/json'].schema.properties;
      //put url params (e.g. {name}, {role})
      //at the front of the props list
      const newProps = assign({}, paramProp, props);
      return expandOpenApiProps(newProps);
    });
  },

  getPaths(apiPath, backend, itemType) {
    return this.ajax(`/v1/${apiPath}?help=1`, backend).then(help => {
      const pathInfo = help.openapi.paths;
      let paths = Object.keys(pathInfo);
      paths = paths.filter(path => pathInfo[path]['x-vault-sudo'] !== true); //get rid of deprecated paths
      const configPath = paths
        .map(path => {
          if (
            pathInfo[path].post &&
            !path.includes('{') &&
            pathInfo[path].get &&
            (!pathInfo[path].get.parameters || pathInfo[path].get.parameters[0].name !== 'list')
          ) {
            return { path: path, tag: pathInfo[path].get.tags[0] };
          }
        })
        .filter(path => path != undefined);

      const listPaths = paths
        .map(path => {
          if (
            pathInfo[path].get &&
            pathInfo[path].get.parameters &&
            pathInfo[path].get.parameters[0].name == 'list'
          ) {
            return { path: path, tag: pathInfo[path].get.tags[0] };
          }
        })
        .filter(path => path != undefined);

      //we always want to keep list endpoints, but best to only use relevant post/delete endpoints
      if (itemType) {
        paths = paths.filter(path => path.includes(itemType));
      }
      const deletePaths = paths
        .map(path => {
          if (pathInfo[path].delete) {
            return { path: path, tag: pathInfo[path].delete.tags[0] };
          }
        })
        .filter(path => path != undefined);
      const createPaths = paths
        .map(path => {
          if (pathInfo[path].post && path.includes('{') && !path.includes('login')) {
            return { path: path, tag: pathInfo[path].post.tags[0] };
          }
        })
        .filter(path => path != undefined);
      return {
        apiPath: apiPath,
        configPath: configPath,
        list: listPaths,
        create: createPaths,
        delete: deletePaths,
      };
    });
  },

  getNewAdapter(backend, paths, itemType) {
    const { list, create } = paths;
    return generatedItemAdapter.extend({
      urlForItem(method, id, type) {
        let listPath = list.find(pathInfo => pathInfo.path.includes(itemType));
        let { tag, path } = listPath;
        let url = `${this.buildURL()}/${tag}/${backend}${path}/`;
        if (id) {
          url = url + encodePath(id);
        }
        return url;
      },

      urlForFindRecord(id, modelName, snapshot) {
        return this.urlForItem(null, id, null);
      },

      urlForUpdateRecord(id, modelName, snapshot) {
        let { tag, path } = create[0];
        path = path.slice(0, path.indexOf('{') - 1);
        return `${this.buildURL()}/${tag}/${backend}${path}/${id}`;
      },

      urlForCreateRecord(modelType, snapshot) {
        const { id } = snapshot;
        let { tag, path } = create[0];
        path = path.slice(0, path.indexOf('{') - 1);
        return `${this.buildURL()}/${tag}/${backend}${path}/${id}`;
      },
    });
  },

  getNewModel(modelType, owner, backend, apiPath, itemType) {
    const modelName = `model:${modelType}`;
    const modelFactory = owner.factoryFor(modelName);
    let newModel, helpUrl;
    if (modelFactory) {
      newModel = modelFactory.class;
      const modelProto = newModel.proto();
      if (newModel.merged || modelProto.useOpenAPI !== true) {
        return resolve();
      }
      helpUrl = modelProto.getHelpUrl(backend);
      return this.registerNewModel(helpUrl, backend, newModel, modelName, owner);
    } else {
      newModel = DS.Model.extend({});
      return this.getPaths(apiPath, backend, itemType).then(paths => {
        const adapterFactory = owner.factoryFor(`adapter:${modelType}`);
        if (!adapterFactory) {
          const adapter = this.getNewAdapter(backend, paths, itemType);
          owner.register(`adapter:${modelType}`, adapter);
        }

        //if we have an item we want the create info for that itemType
        if (itemType) {
          let { tag, path } = paths.create[0];
          path = path.slice(0, path.indexOf('{') - 1) + '/example';
          helpUrl = `/v1/${tag}/${backend}${path}?help=true`;
        } else {
          //we need the mount config
          let { tag, path } = paths.configPath[0];
          helpUrl = `/v1/${tag}/${backend}${path}?help=true`;
        }
        return this.registerNewModel(helpUrl, backend, newModel, modelName, owner);
      });
    }
  },

  registerNewModel(helpUrl, backend, newModel, modelName, owner) {
    return this.getProps(helpUrl, backend).then(props => {
      const { attrs, newFields } = combineAttributes(newModel.attributes, props);
      newModel = newModel.extend(attrs, { newFields });
      if (!newModel.fieldGroups) {
        const fieldGroups = fieldToAttrs(newModel, [{ default: newFields }]);
        newModel = newModel.extend({ fieldGroups });
      }
      newModel.reopenClass({ merged: true });
      owner.unregister(modelName);
      owner.register(modelName, newModel);
    });
  },
});
