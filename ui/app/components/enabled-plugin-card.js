/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { action } from '@ember/object';
import keys from 'core/utils/keys';

/**
 * @module EnabledPluginCard
 * EnabledPluginCard components are used to display available secret and auth engines in the mount backend type form
 *
 * @example
 * ```js
 * <EnabledPluginCard @type={{type}} @setMountType={{@setMountType}} />
 * ```
 * @param {Object} type - The engine metadata object with properties like displayName, requiresEnterprise, requiredFeature, etc.
 * @param {Function} setMountType - Function to call when the card is clicked to select this mount type
 */

export default class EnabledPluginCard extends Component {
  @service version;

  get isDisabled() {
    const { type } = this.args;
    return (
      (type.requiresEnterprise && !this.version.isEnterprise) ||
      (type.requiredFeature && !this.hasFeature(type.requiredFeature))
    );
  }

  get showEnterpriseBadge() {
    const { type } = this.args;
    return (
      (type.requiresEnterprise && !this.version.isEnterprise) ||
      (type.requiredFeature && !this.hasFeature(type.requiredFeature))
    );
  }

  get showDeprecationBadge() {
    const { type } = this.args;
    return type.deprecationStatus && type.deprecationStatus !== 'supported';
  }

  hasFeature(featureName) {
    return this.version.features?.includes(featureName) || false;
  }

  @action
  handleSelection() {
    if (!this.isDisabled) {
      this.args.setMountType(this.args.type.type);
    }
  }

  @action
  handleKeyDown(event) {
    // Only handle Enter and Space keys for accessibility
    if (event.key === keys.ENTER || event.key === keys.SPACE) {
      event.preventDefault();
      this.handleSelection();
    }
  }
}
