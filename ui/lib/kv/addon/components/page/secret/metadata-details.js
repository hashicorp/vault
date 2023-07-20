/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { action } from '@ember/object';
import { getOwner } from '@ember/application';

/**
 * @module MetadataDetails
 * MetadataDetails component xx
 *
 * @param {array} model - An xx
 * @param {array} breadcrumbs - Breadcrumbs as an array of objects that contain label, route, and modelId. They are updated via the util kv-breadcrumbs to handle dynamic *pathToSecret on the list-directory route.
 */

export default class KvListPageComponent extends Component {
  @service flashMessages;

  get mountPoint() {
    // mountPoint tells the LinkedBlock component where to start the transition. In this case, mountPoint will always be vault.cluster.secrets.backend.kv.
    return getOwner(this).mountPoint;
  }

  @action
  onDelete() {
    // todo
  }
}
