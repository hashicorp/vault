import Component from '@glimmer/component';
import { inject as service } from '@ember/service';

/**
 * @module LinkStatus
 * LinkStatus components are used to indicate link status to the hashicorp cloud platform
 *
 * @example
 * ```js
 * <LinkStatus @status={{this.currentCluser.cluster.hcpLinkStatus}} />
 * ```
 *
 * @param {string} status - cluster.hcpLinkStatus value from currentCluster service
 */

export default class LinkStatus extends Component {
  @service store;
  @service version;

  get showBanner() {
    // enterprise only feature at this time but will expand to OSS in future release
    // there are plans to handle connection failure states -- only alert if connected until further states are returned
    return this.version.isEnterprise && this.args.status === 'connected';
  }

  get bannerClass() {
    return this.args.status === 'connected' ? 'connected' : 'warning';
  }
}
