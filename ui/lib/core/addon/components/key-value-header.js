/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { ancestorKeysForKey, keyPartsForKey, keyWithoutParentKey } from 'core/utils/key-utils';
import { encodePath } from 'vault/utils/path-encoding-helpers';

/**
 * @module KeyValueHeader
 * KeyValueHeader components show breadcrumbs for secret engines.
 *
 * @example
 * ```js
 <KeyValueHeader @path="vault.cluster.secrets.backend.show" @mode={{this.mode}} @root={{@root}}/>
 * ```
 * @param {string} [mode=null] - Used to set the currentPath.
 * @param {string} [baseKey=null] - Used to generate the path backward.
 * @param {string} [path=null] - The fallback path.
 * @param {string} [root=null] - Used to set the secretPath.
 * @param {boolean} [showCurrent=true] - Boolean to show the second part of the breadcrumb, ex: the secret's name.
 * @param {boolean} [linkToPaths=true] - If true link to the path.
 */

export default class KeyValueHeader extends Component {
  get showCurrent() {
    return this.args.showCurrent || true;
  }

  get linkToPaths() {
    return this.args.linkToPaths || true;
  }

  stripTrailingSlash(str) {
    return str[str.length - 1] === '/' ? str.slice(0, -1) : str;
  }

  get currentPath() {
    if (!this.args.mode || this.showCurrent === false) {
      return this.args.path;
    }
    return `vault.cluster.secrets.backend.${this.args.mode}`;
  }

  get secretPath() {
    const crumbs = [];
    const root = this.args.root;
    const baseKey = this.args.baseKey?.display || this.args.baseKey?.id;
    const baseKeyModel = encodePath(this.args.baseKey?.id);

    if (root) {
      crumbs.push(root);
    }

    if (!baseKey) {
      return crumbs;
    }

    const path = this.args.path;
    const currentPath = this.currentPath;
    const showCurrent = this.showCurrent;
    const ancestors = ancestorKeysForKey(baseKey);
    const parts = keyPartsForKey(baseKey);
    if (ancestors.length === 0) {
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
        model: encodePath(ancestor),
      });
    });

    crumbs.push({
      label: keyWithoutParentKey(baseKey),
      text: this.stripTrailingSlash(keyWithoutParentKey(baseKey)),
      path: currentPath,
      model: baseKeyModel,
    });

    if (!showCurrent) {
      crumbs.pop();
    }

    return crumbs;
  }
}
