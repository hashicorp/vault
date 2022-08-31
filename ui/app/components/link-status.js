import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';

/**
 * @module LinkStatus
 * LinkStatus components are used to indicate link status to the hashicorp cloud platform
 *
 * @example
 * ```js
 * <LinkStatus />
 * ```
 */

export default class LinkStatus extends Component {
  @service store;
  @service version;

  @tracked status;

  constructor() {
    super(...arguments);
    // enterprise only feature at this time but will expand to OSS in future release
    if (this.version.isEnterprise) {
      this.fetchStatus();
    }
  }

  async fetchStatus() {
    try {
      const { hcp_link_status } = await this.store.adapterFor('cluster').sealStatus();
      // there are plans to handle connection failure states -- only alert if connected until further states are returned
      if (hcp_link_status === 'connected') {
        this.status = hcp_link_status;
      }
    } catch (error) {
      this.status = null;
    }
  }

  get bannerClass() {
    return this.status === 'connected' ? 'connected' : 'warning';
  }
}
