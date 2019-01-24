// Low level service that allows users to input paths to make requests to vault
// this service provides the UI synecdote to the cli commands read, write, delete, and list
import Service from '@ember/service';

import { getOwner } from '@ember/application';
import { expandOpenApiProps } from 'vault/utils/openapi-to-attrs';

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

  getProps(modelType, backend) {
    let adapter = getOwner(this).lookup(`adapter:${modelType}`);
    let path = adapter.pathForType();
    let helpUrl = `/v1/${backend}/${path}/example?help=1`;
    let wildcard;
    switch (path) {
      case 'roles':
        wildcard = 'name';
        break;
      case 'mounts':
        if (modelType === 'secret') {
          wildcard = 'path';
        } else {
          wildcard = 'config';
        }
        break;
      case 'sign':
      case 'issue':
        wildcard = 'role';
        break;
    }
    //{ roles: 'name', mounts: 'config', sign: 'role', issue: 'role' }[path];
    return this.ajax(helpUrl, backend).then(help => {
      let props =
        help.openapi.paths[`/${path}/{${wildcard}}`].post.requestBody.content['application/json'].schema
          .properties;
      return expandOpenApiProps(props);
    });
  },
});
