import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { inject as service } from '@ember/service';
import { action } from '@ember/object';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';
import FlashMessageService from 'vault/services/flash-messages';
import PkiActionModel from 'vault/models/pki/action';
import errorMessage from 'vault/utils/error-message';

interface Args {
  model: PkiActionModel;
  useIssuer: boolean;
  onSave: CallableFunction;
  onCancel: CallableFunction;
}

export default class PkiGenerateCsrComponent extends Component<Args> {
  @service declare readonly flashMessages: FlashMessageService;

  @tracked modelValidations = null;
  @tracked error: string | null = null;
  @tracked alert: string | null = null;

  formFields;
  // fields rendered after CSR generation
  showFields = ['csr', 'keyId', 'privateKey', 'privateKeyType'];

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    this.formFields = expandAttributeMeta(this.args.model, [
      'type',
      'commonName',
      'excludeCnFromSans',
      'format',
      'serialNumber',
      'addBasicConstraints',
    ]);
  }

  @action
  cancel() {
    this.args.model.unloadRecord();
    this.args.onCancel();
  }

  async getCapability(): Promise<boolean> {
    try {
      const issuerCapabilities = await this.args.model.generateIssuerCsrPath;
      return issuerCapabilities.get('canCreate') === true;
    } catch (error) {
      return false;
    }
  }

  @task
  @waitFor
  *save(event: Event): Generator<Promise<boolean | PkiActionModel>> {
    event.preventDefault();
    try {
      const { model } = this.args;
      const { isValid, state, invalidFormMessage } = model.validate();
      if (isValid) {
        const useIssuer = yield this.getCapability();
        yield model.save({ adapterOptions: { actionType: 'generate-csr', useIssuer } });
        this.flashMessages.success('Successfully generated CSR.');
      } else {
        this.modelValidations = state;
        this.alert = invalidFormMessage;
      }
    } catch (e) {
      this.error = errorMessage(e);
      this.alert = 'There was a problem generating the CSR.';
    }
  }
}
