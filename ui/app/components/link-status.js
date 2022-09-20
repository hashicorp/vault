import Component from '@glimmer/component';
import { inject as service } from '@ember/service';

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
  @service store;
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
      return this.state !== 'connected' ? this.args.status.split('since')[1].split('m=')[0].trim() : null;
    } catch {
      return null;
    }
  }

  get message() {
    if (this.args.status) {
      const error = this.args.status.split('error:')[1];
      const time = `[${this.timestamp}]`;
      if (this.state === 'disconnected') {
        // if generally disconnected hide the banner
        return !error || error.includes('UNKNOWN')
          ? null
          : `Vault has been disconnected from the Hashicorp Cloud Platform since ${time}. Error: ${error}`;
      } else if (this.state === 'connecting') {
        if (error.includes('connection refused')) {
          return `Vault has been trying to connect to the Hashicorp Cloud Platform since ${time}, but the Scada provider is down. Vault will try again soon.`;
        } else if (error.includes('principal does not have permission to register as provider')) {
          return `Vault tried connecting to the Hashicorp Cloud Platform, but the Resource ID is invalid. Check your resource ID. ${time}`;
        } else if (error.includes('cannot fetch token: 401 Unauthorized')) {
          return `Vault tried connecting to the Hashicorp Cloud Platform, but the authorization information is wrong. Update it and try again. ${time}`;
        } else {
          // catch all for any unknown errors
          return `Vault has been trying to connect to the Hashicorp Cloud Platform since ${time}. Vault will try again soon. Error: ${error}`;
        }
      }
    }
    return null;
  }

  get showStatus() {
    // enterprise only feature at this time but will expand to OSS in future release
    if (!this.version.isEnterprise || !this.args.status) {
      return false;
    }
    if (this.state === 'disconnected' && !this.message) {
      return false;
    }
    return true;
  }
}
