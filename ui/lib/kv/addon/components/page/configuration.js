/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { toLabel } from 'core/helpers/to-label';
import { duration } from 'core/helpers/format-duration';

/**
 * @module KvConfigPageComponent
 * KvConfigPageComponent is a component to show secrets mount and engine configuration data
 *
 * @param {object} config - config data for mount and engine
 * @param {string} backend - The name of the kv secret engine.
 * @param {array} breadcrumbs - Breadcrumbs as an array of objects that contain label, route, and modelId. They are updated via the util kv-breadcrumbs to handle dynamic *pathToSecret on the list-directory route.
 */

export default class KvConfigPageComponent extends Component {
  label = (key) => {
    const label = toLabel([key]);
    // map specific fields to custom labels
    return (
      {
        cas_required: 'Require check and set',
        delete_version_after: 'Automate secret deletion',
        max_versions: 'Maximum number of versions',
        default_lease_ttl: 'Default Lease TTL',
        max_lease_ttl: 'Max Lease TTL',
      }[key] || label
    );
  };

  value = (key, value) => {
    if (key === 'delete_version_after') {
      return value === '0s' ? 'Never delete' : duration([value]);
    }
    return value;
  };
}
