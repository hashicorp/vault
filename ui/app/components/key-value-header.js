import { computed } from '@ember/object';
import Component from '@ember/component';
import utils from 'vault/lib/key-utils';
import { encodePath } from 'vault/utils/path-encoding-helpers';

export default Component.extend({
  tagName: 'nav',
  classNames: 'key-value-header breadcrumb',
  ariaLabel: 'breadcrumbs',
  attributeBindings: ['ariaLabel:aria-label', 'aria-hidden'],

  baseKey: null,
  path: null,
  showCurrent: true,
  linkToPaths: true,

  stripTrailingSlash(str) {
    return str[str.length - 1] === '/' ? str.slice(0, -1) : str;
  },

  currentPath: computed('mode', 'path', 'showCurrent', function() {
    const mode = this.get('mode');
    const path = this.get('path');
    const showCurrent = this.get('showCurrent');
    if (!mode || showCurrent === false) {
      return path;
    }
    return `vault.cluster.secrets.backend.${mode}`;
  }),

  secretPath: computed('baseKey', 'baseKey.{display,id}', 'root', 'showCurrent', function() {
    let crumbs = [];
    const root = this.get('root');
    const baseKey = this.get('baseKey.display') || this.get('baseKey.id');
    const baseKeyModel = encodePath(this.get('baseKey.id'));

    if (root) {
      crumbs.push(root);
    }

    if (!baseKey) {
      return crumbs;
    }

    const path = this.get('path');
    const currentPath = this.get('currentPath');
    const showCurrent = this.get('showCurrent');
    const ancestors = utils.ancestorKeysForKey(baseKey);
    const parts = utils.keyPartsForKey(baseKey);
    if (!ancestors) {
      crumbs.push({
        label: baseKey,
        text: this.stripTrailingSlash(baseKey),
        path: currentPath,
        model: baseKeyModel,
      });

      if (!showCurrent) {
        crumbs.pop();
      }

      return crumbs;
    }

    ancestors.forEach((ancestor, index) => {
      crumbs.push({
        label: parts[index],
        text: this.stripTrailingSlash(parts[index]),
        path: path,
        model: ancestor,
      });
    });

    crumbs.push({
      label: utils.keyWithoutParentKey(baseKey),
      text: this.stripTrailingSlash(utils.keyWithoutParentKey(baseKey)),
      path: currentPath,
      model: baseKeyModel,
    });

    if (!showCurrent) {
      crumbs.pop();
    }

    return crumbs;
  }),
});
