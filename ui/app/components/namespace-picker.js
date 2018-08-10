import Ember from 'ember';
import flat from 'flat';
import keyUtils from 'vault/lib/key-utils';

const { ancestorKeysForKey } = keyUtils;
const { unflatten } = flat;
const { Component, computed, inject } = Ember;
const DOT_REPLACEMENT = '☃';

export default Component.extend({
  namespaceService: inject.service('namespace'),
  auth: inject.service(),

  init() {
    this._super(...arguments);
    this.get('namespaceService.findNamespacesForUser').perform();
  },

  namespacePath: computed.alias('namespaceService.path'),

  accessibleNamespaces: computed.alias('namespaceService.accessibleNamespaces'),
  namespaceTree: computed('accessibleNamespaces', function() {
    let nsList = this.get('accessibleNamespaces');
    if (!nsList) {
      return [];
    }
    nsList = nsList.slice(0).reverse();
    let tree = {};
    let maxDepth;
    let nsTree = nsList.reduce((accumulator, ns) => {
      let prefixInList = accumulator.some(nsPath => nsPath.startsWith(ns));
      if (!prefixInList) {
        accumulator.push(ns);
      }
      return accumulator;
    }, []);

    for (let ns of nsTree) {
      ns = ns.replace(/\.+/g, DOT_REPLACEMENT);
      let branch = unflatten({ [ns]: null }, { delimiter: '/' });
      tree = {
        ...tree,
        ...branch,
      };
    }
    return tree;
  }),

  pathToLeaf(path) {
    // dots are allowed in namespace paths
    // so we need to preserve them, and replace slashes with dots
    // in order to use Ember's get function on the namespace tree
    // to pull out the correct level
    return (
      path
        // trim trailing slash
        .replace(/\/$/, '')
        // replace dots with snowman
        .replace(/\.+/g, DOT_REPLACEMENT)
        // replace slash with dots
        .replace(/\/+/g, '.')
    );
  },
  // an array that keeps track of what additional panels to render
  // on the menu stack
  menuLeaves: computed('namespacePath', 'namespaceTree', function() {
    let ns = this.get('namespacePath');
    let leaves = ancestorKeysForKey(ns) || [];
    leaves.push(ns);
    return leaves.map(this.pathToLeaf);
  }),

  rootLeaves: computed('namespacePath', 'namespaceTree', function() {
    //let ns = this.get('namespacePath');
    let leaves = Object.keys(this.get('namespaceTree'));
    return leaves.map(this.pathToLeaf);
  }),

  currentLeaf: computed.alias('menuLeaves.lastObject'),
  canAccessMultipleNamespaces: computed.gt('accessibleNamespaces.length', 1),
  isUserRootNamespace: computed('auth.authData.userRootNamespace', 'namespacePath', function() {
    return this.get('auth.authData.userRootNamespace') === this.get('namespacePath');
  }),

  namespaceDisplay: computed('namespacePath', 'accessibleNamespaces', 'accessibleNamespaces.[]', function() {
    let namespace = this.get('namespacePath');
    if (namespace === '') {
      return '';
    } else {
      let parts = namespace.split('/');
      return parts[parts.length - 1];
    }
  }),
});
