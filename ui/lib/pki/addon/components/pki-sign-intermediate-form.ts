import { action } from '@ember/object';
import { service } from '@ember/service';
import { waitFor } from '@ember/test-waiters';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { task } from 'ember-concurrency';
import PkiIssuerModel from 'vault/models/pki/issuer';
import FlashMessageService from 'vault/services/flash-messages';
import errorMessage from 'vault/utils/error-message';

interface Args {
  onCancel: CallableFunction;
  model: PkiIssuerModel;
}

export default class PkiSignIntermediateFormComponent extends Component<Args> {
  @service declare readonly flashMessages: FlashMessageService;
  @tracked errorBanner = '';
  @tracked inlineFormAlert = '';

  @action cancel() {
    this.args.model.unloadRecord();
    this.args.onCancel();
  }
  @task
  @waitFor
  *save(event: Event) {
    event.preventDefault();
    try {
      yield this.args.model.save();
      this.flashMessages.success('Successfully signed CSR.');
      // this.args.onSave();
    } catch (e) {
      this.errorBanner = errorMessage(e);
      this.inlineFormAlert = 'There was a problem signing the CSR.';
    }
  }

  get groups() {
    return {
      'Signing options': ['usePss', 'skid', 'signatureBits'],
      'Subject Alternative Name (SAN) Options': ['altNames', 'ipSans', 'uriSans', 'otherSans'],
      'Additional subject fields': [
        'ou',
        'organization',
        'country',
        'locality',
        'province',
        'streetAddress',
        'postalCode',
      ],
    };
  }

  get showFields() {
    return ['serialNumber', 'certificate', 'issuingCa', 'caChain'];
  }
}
