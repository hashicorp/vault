import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { parseCertificate } from 'vault/utils/parse-pki-cert';
import camelizeKeys from 'vault/utils/camelize-object-keys';
// TYPES
import Store from '@ember-data/store';
import Router from '@ember/routing/router';
import FlashMessageService from 'vault/services/flash-messages';
import PkiIssuerModel from 'vault/models/pki/issuer';
import { Breadcrumb } from 'vault/vault/app-types';
interface Args {
  oldRoot: PkiIssuerModel;
  breadcrumbs: Breadcrumb;
}

export default class PagePkiIssuerRotateRootComponent extends Component<Args> {
  @service declare readonly store: Store;
  @service declare readonly router: Router;
  @service declare readonly flashMessages: FlashMessageService;

  @tracked title = 'Generate new root';
  @tracked rotateForm = 'use-old-settings';
  @tracked showOldSettings = false;
  @tracked newRootModel;
  @tracked warningBanner = '';

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    const certData = parseCertificate(this.args.oldRoot.certificate);
    if (certData.parsing_errors && certData.parsing_errors.length > 0) {
      const errorMessage = certData.parsing_errors.map((e: Error) => e.message).join(', ');
      this.warningBanner = errorMessage;
    }
    this.newRootModel = this.store.createRecord('pki/action', {
      actionType: 'rotate-root',
      type: 'internal',
      ...camelizeKeys(certData), // copy old root settings over to new one
    });
  }

  get rotationOptions() {
    return [
      {
        key: 'use-old-settings',
        icon: 'vector',
        label: 'Use old root settings',
        description: `Provide only a new common name and issuer name, using the old root’s settings. Selecting this option generates an internal root type.`,
      },
      {
        key: 'customize',
        icon: 'vector',
        label: 'Customize new root certificate',
        description:
          'Generates a new self-signed CA certificate and private key. This generated root will sign its own CRL.',
      },
    ];
  }
}
