/*
  This service is used to pull an OpenAPI document describing the
  shape of data at a specific path to hydrate a model with attrs it
  has less (or no) information about.
*/
import Service from '@ember/service';

import { getOwner } from '@ember/application';
import { expandOpenApiProps, combineAttributes } from 'vault/utils/openapi-to-attrs';
import { resolve } from 'rsvp';

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
      let props = help.openapi.paths[path].post.requestBody.content['application/json'].schema.properties;
      return expandOpenApiProps(props);
    });
  },

  getNewModel(modelType, owner, backend) {
    let name = `model:${modelType}`;
    let newModel = owner.factoryFor(name).class;
    let modelProto = newModel.proto();
    if (newModel.merged || modelProto.useOpenAPI !== true) {
      return resolve();
    }
    let helpUrl = modelProto.getHelpUrl(backend);

    return this.getProps(helpUrl, backend).then(props => {
      if (owner.hasRegistration(name) && !newModel.merged) {
        let { attrs, newFields } = combineAttributes(newModel.attributes, props);
        newModel = newModel.extend(attrs, { newFields });
      } else {
        //generate a whole new model
      }

      newModel.reopenClass({ merged: true });
      owner.unregister(name);
      owner.register(name, newModel);
    });
  },
});
