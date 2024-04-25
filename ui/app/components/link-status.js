/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';

/**
 * @module LinkStatus
 * LinkStatus components are used to indicate link status to the hashicorp cloud platform
 *
 * @example
 * ```js
 * <LinkStatus @status={{this.currentCluster.cluster.hcpLinkStatus}} />
 * ```
 *
 * @param {string} status - cluster.hcpLinkStatus value from currentCluster service -- returned from seal-status endpoint
 */

export default class LinkStatus extends Component {
  @service version;

  get state() {
    if (!this.args.status) return null;
    // connected state is returned with no further information
    if (this.args.status === 'connected') return this.args.status;
    // disconnected and connecting states are returned with a timestamp and error
    // state is always the first word of the string
    return this.args.status.split(' ', 1).toString();
  }

  get timestamp() {
    try {
      return this.state !== 'connected' ? this.args.status.split('since')[1].split(';')[0].trim() : null;
    } catch {
      return null;
    }
  }

  get error() {
    const status = this.args.status;
    return status && status !== 'connected' ? status.split('error:')[1] : null;
  }
}
