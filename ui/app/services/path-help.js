/*
  This service is used to pull an OpenAPI document describing the
  shape of data at a specific path to hydrate a model with attrs it
  has less (or no) information about.
*/
import Service from '@ember/service';

import { getOwner } from '@ember/application';
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
      let path = Object.keys(help.openapi.paths)[0];
      path = help.openapi.paths[path];
      const params = path.parameters;
      let param = {};
      //put params at the front of the props list
      if (params) {
        let label = params[0].name;
        if (label.toLowerCase() !== 'name') {
          label += ' name';
        }
        param[params[0].name] = {
          name: params[0].name,
          label: label,
          type: params[0].schema.type,
          description: params[0].description,
          isId: true,
        };
      }
      let props = path.post.requestBody.content['application/json'].schema.properties;
      let newProps = { ...param, ...props };
      return expandOpenApiProps(newProps);
    });
  },

  //Makes a call to grab the OpenAPI document.
  //Returns relevant information from OpenAPI
  //as determined by the expandOpenApiProps util
  getPathHelp(apiPath, backend) {
    let helpUrl = `/v1/${apiPath}?help=1`;
    return this.ajax(helpUrl, backend).then(help => {
      let path = Object.keys(help.openapi.paths)[0];
      let props = help.openapi.paths[path].post.requestBody.content['application/json'].schema.properties;
      return expandOpenApiProps(props);
    });
  },

  getPaths(apiPath, backend) {
    let helpUrl = `/v1/${apiPath}?help=1`;
    return this.ajax(helpUrl, backend).then(help => {
      let paths = Object.keys(help.openapi.paths);
      let listPaths = paths
        .map(path => {
          if (
            help.openapi.paths[path].get &&
            help.openapi.paths[path].get.parameters &&
            help.openapi.paths[path].get.parameters[0].name == 'list'
          ) {
            return path;
          }
        })
        .filter(path => path != undefined);
      let createPaths = paths.filter(path => path.includes('{') && !path.includes('login'));
      return { apiPath: apiPath, list: listPaths, create: createPaths };
    });
  },

  // getPathsForModel(modelType, owner, backend){
  //   let name = `model:${modelType}`;
  //   let newModel = owner.factoryFor(name).class;
  //   let modelProto = newModel.proto();
  //   let helpUrl = modelProto.getHelpUrl(backend);
  //   return this.getPaths(helpUrl, backend).then(paths => {
  //     return paths;
  //   });
  // },

  getNewModel(modelType, owner, backend, apiPath, itemType) {
    let name = `model:${modelType}`;
    let factory = owner.factoryFor(name);
    let newModel, helpUrl;
    if (factory) {
      newModel = factory.class;
      let modelProto = newModel.proto();
      if (newModel.merged || modelProto.useOpenAPI !== true) {
        return resolve();
      }
      helpUrl = apiPath ? `/v1/${apiPath}?help=1` : modelProto.getHelpUrl(backend);
    } else {
      newModel = DS.Model.extend({});
      helpUrl = `/v1/${apiPath}?help=1`;
      let newAdapter = generatedItemAdapter.extend({ itemType });
      owner.register(`adapter:${modelType}`, newAdapter);
    }

    return this.getProps(helpUrl, backend).then(props => {
      const { attrs, newFields } = combineAttributes(newModel.attributes, props);
      newModel = newModel.extend(attrs, { newFields });
      if (!newModel.fieldGroups) {
        const fieldGroups = fieldToAttrs(newModel, [{ default: newFields }]);
        newModel = newModel.extend({ fieldGroups });
      }
      newModel.reopenClass({ merged: true });
      owner.unregister(name);
      owner.register(name, newModel);
    });
  },
});
