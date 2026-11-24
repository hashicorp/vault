/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';

/**
 * @module KvVersionDropdown
 * Toolbar menu for switching between versions of a KV secret.
 *
 * @param {number} displayVersion - the version currently being viewed
 * @param {object} metadata - secret metadata, supplies `versions` and `current_version`
 * @param {function} [onSelect] - when passed, items are buttons that call this instead of linking
 * @param {function} [onClose] - closes the menu after a link is followed
 * @param {string} [route] - route for the link variant, e.g. "secret.details"
 * @param {array} [models] - models for that route
 */
export default class KvVersionDropdown extends Component {
  // Returns the icon, colour class and disabled state together so the template resolves
  // one condition set per version instead of nesting {{if}} helpers per argument.
  // The colour comes from the same `has-text-*` helper classes the old BasicDropdown used:
  // they carry `!important`, so unlike HDS's `@color` they survive the disabled styling.
  status = (versionData) => {
    const { destroyed, isSecretDeleted, version } = versionData;
    // version keys are strings, current_version is a number
    const isCurrent = Number(version) === Number(this.args.metadata?.current_version);

    if (destroyed) return { icon: 'x-square-fill', class: 'has-text-danger', disabled: true };
    if (isSecretDeleted) return { icon: 'x-square-fill', class: 'has-text-grey', disabled: true };
    if (isCurrent) return { icon: 'check-circle', class: 'has-text-success', disabled: false };
    return { icon: null, class: '', disabled: false };
  };
}
