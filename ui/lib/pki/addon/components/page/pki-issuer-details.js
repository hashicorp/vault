import { action } from '@ember/object';
import { service } from '@ember/service';
import Component from '@glimmer/component';
import errorMessage from 'vault/utils/error-message';

export default class PagePkiIssuerDetailsComponent extends Component {
  @service download;
  @service flashMessages;
  @service store;

  get win() {
    return this.window || window;
  }

  fetchCertByFormat(backend, issuerId, format) {
    const endpoint = `/v1/${backend}/issuer/${issuerId}/${format}`;
    const adapter = this.store.adapterFor('application');
    return adapter.rawRequest(endpoint, 'GET', { unauthenticated: true }).then(function (response) {
      if (format === 'der') {
        return response.blob();
      }
      return response.text();
    });
  }

  @action
  async downloadCert(dropdown, format) {
    const { issuer } = this.args;
    try {
      const contents = await this.fetchCertByFormat(issuer.backend, issuer.id, format);
      this.download.download(`${issuer.backend}-issuer`, contents, format);
    } catch (e) {
      this.flashMessages.danger(errorMessage(e, 'Your certificate download failed.'));
    }
  }
}
