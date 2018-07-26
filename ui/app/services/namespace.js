import Ember from 'ember';

const { Service, computed } = Ember;
const DEFAULT_NAMESPACE = '';
export default Service.extend({
  //populated by the query param on the cluster route
  path: null,
  isDefault: computed.equal('namespace', DEFAULT_NAMESPACE),
  setNamespace(path) {
    this.set('path', path);
  },
});
