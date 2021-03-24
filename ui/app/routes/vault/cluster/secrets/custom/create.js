import Route from '@ember/routing/route';
import { expandOpenApiProps, nonModelAttrs } from '../../../../../utils/openapi-to-attrs';

export default class VaultClusterSecretsCustomCreateRoute extends Route {
  async model(params) {
    const parentModel = this.modelFor('vault.cluster.secrets.custom');
    const { itempath } = params;
    let queryUrl = `/v1/${parentModel.backend}`;
    const info = parentModel.pathInfo[itempath];
    const expanded = expandOpenApiProps(info.post.requestBody.content['application/json'].schema.properties);
    return {
      fields: nonModelAttrs(expanded),
      path: `${queryUrl}${info.urlPath}`,
    };
  }
}
