import Service from '@ember/service';

export default Service.extend({
  managedNamespaceRoot: null,

  setManagedNamespaceRoot(path) {
    this.set('managedNamespaceRoot', path);
  },
});
