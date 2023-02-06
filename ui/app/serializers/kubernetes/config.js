import ApplicationSerializer from '../application';

export default class KubernetesConfigSerializer extends ApplicationSerializer {
  primaryKey = 'backend';

  serialize() {
    const json = super.serialize(...arguments);
    // remove backend value from payload
    delete json.backend;
    return json;
  }
}
