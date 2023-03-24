import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
// TYPES
import Store from '@ember-data/store';
import Router from '@ember/routing/router';
import FlashMessageService from 'vault/services/flash-messages';
import PkiActionModel from 'vault/models/pki/action';
import { Breadcrumb } from 'vault/vault/app-types';

interface Args {
  config: PkiActionModel;
  onCancel: CallableFunction;
  breadcrumbs: Breadcrumb;
}

export default class PagePkiIssuerRotateRootComponent extends Component<Args> {
  @service declare readonly store: Store;
  @service declare readonly router: Router;
  @service declare readonly flashMessages: FlashMessageService;

  @tracked title = 'Generate new root';

  get generateOptions() {
    return [
      {
        key: 'rotate-root',
        icon: 'vector',
        label: 'Use old root settings',
        description: 'Provide only a new common name and issuer name, using the old root’s settings. ',
      },
      {
        key: 'generate-root',
        icon: 'vector',
        label: 'Generate root',
        description:
          'Generates a new self-signed CA certificate and private key. This generated root will sign its own CRL.',
      },
    ];
  }
}
