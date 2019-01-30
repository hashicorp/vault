import { inject as service } from '@ember/service';
import { alias, gt } from '@ember/object/computed';
import Component from '@ember/component';
import { computed } from '@ember/object';
import keyUtils from 'vault/lib/key-utils';
import pathToTree from 'vault/lib/path-to-tree';
import { task, timeout } from 'ember-concurrency';

const { ancestorKeysForKey } = keyUtils;
const DOT_REPLACEMENT = '☃';
const ANIMATION_DURATION = 250;

export default Component.extend({
  tagName: '',
  namespaceService: service('namespace'),
  auth: service(),
  store: service(),
  namespace: null,
  listCapability: null,
  canList: alias('listCapability.canList'),

  init() {
    this._super(...arguments);
    this.get('namespaceService.findNamespacesForUser').perform();
  },

  didReceiveAttrs() {
    this._super(...arguments);

    let ns = this.get('namespace');
    let oldNS = this.get('oldNamespace');
    if (!oldNS || ns !== oldNS) {
      this.get('setForAnimation').perform();
      this.get('fetchListCapability').perform();
    }
    this.set('oldNamespace', ns);
  },

  fetchListCapability: task(function*() {
    try {
      if (this.listCapability) {
        yield this.listCapability.reload();
        return;
      }
      let capability = yield this.store.findRecord('capabilities', 'sys/namespaces/');
      this.set('listCapability', capability);
    } catch (e) {
      //do nothing here
    }
  }),
  setForAnimation: task(function*() {
    let leaves = this.get('menuLeaves');
    let lastLeaves = this.get('lastMenuLeaves');
    if (!lastLeaves) {
      this.set('lastMenuLeaves', leaves);
      yield timeout(0);
      return;
    }
    let isAdding = leaves.length > lastLeaves.length;
    let changedLeaf = (isAdding ? leaves : lastLeaves).get('lastObject');
    this.set('isAdding', isAdding);
    this.set('changedLeaf', changedLeaf);

    // if we're adding we want to render immediately an animate it in
    // if we're not adding, we need time to move the item out before
    // a rerender removes it
    if (isAdding) {
      this.set('lastMenuLeaves', leaves);
      yield timeout(0);
      return;
    }
    yield timeout(ANIMATION_DURATION);
    this.set('lastMenuLeaves', leaves);
  }).drop(),

  isAnimating: alias('setForAnimation.isRunning'),

  namespacePath: alias('namespaceService.path'),

  // this is an array of namespace paths that the current user
  // has access to
  accessibleNamespaces: alias('namespaceService.accessibleNamespaces'),
  inRootNamespace: alias('namespaceService.inRootNamespace'),

  namespaceTree: computed('accessibleNamespaces', function() {
    let nsList = this.get('accessibleNamespaces');

    if (!nsList) {
      return [];
    }
    return pathToTree(nsList);
  }),

  maybeAddRoot(leaves) {
    let userRoot = this.get('auth.authData.userRootNamespace');
    if (userRoot === '') {
      leaves.unshift('');
    }

    return leaves.uniq();
  },

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
  // if you're in  'foo/bar/baz',
  // this array will be: ['foo', 'foo.bar', 'foo.bar.baz']
  // the template then iterates over this, and does  Ember.get(namespaceTree, leaf)
  // to render the nodes of each leaf

  // gets set as  'lastMenuLeaves' in the ember concurrency task above
  menuLeaves: computed('namespacePath', 'namespaceTree', function() {
    let ns = this.get('namespacePath');
    let leaves = ancestorKeysForKey(ns) || [];
    leaves.push(ns);
    leaves = this.maybeAddRoot(leaves);

    leaves = leaves.map(this.pathToLeaf);
    return leaves;
  }),

  // the nodes at the root of the namespace tree
  // these will get rendered as the bottom layer
  rootLeaves: computed('namespaceTree', function() {
    let tree = this.get('namespaceTree');
    let leaves = Object.keys(tree);
    return leaves;
  }),

  currentLeaf: alias('lastMenuLeaves.lastObject'),
  canAccessMultipleNamespaces: gt('accessibleNamespaces.length', 1),
  isUserRootNamespace: computed('auth.authData.userRootNamespace', 'namespacePath', function() {
    return this.get('auth.authData.userRootNamespace') === this.get('namespacePath');
  }),

  namespaceDisplay: computed('namespacePath', 'accessibleNamespaces', 'accessibleNamespaces.[]', function() {
    let namespace = this.get('namespacePath');
    if (namespace === '') {
      return '';
    }
    let parts = namespace.split('/');
    return parts[parts.length - 1];
  }),
});
