import { action } from '@ember/object';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import errorMessage from 'vault/utils/error-message';

/**
 * Pki Config/Generate Root shows only the fields valid for the generate root endpoint.
 * This form handles the model save and rollback actions, and will call the passed
 * onSave and onCancel args for transition (passed from parent).
 * NOTE: not TS because decorator-added items aren't recognized on the model.
 */
export default class PkiGenerateRootComponent extends Component {
  @tracked showGroup = null;
  @tracked errorBanner;
  @tracked invalidFormAlert;

  @action
  toggleGroup(group, isOpen) {
    this.showGroup = isOpen ? group : null;
  }

  get defaultFields() {
    return [
      'type',
      'commonName',
      'issuerName',
      'customTtl',
      'notBeforeDuration',
      'format',
      'permittedDnsDomains',
      'maxPathLength',
    ];
  }
  get keyParamFields() {
    const { type } = this.args.model;
    let fields = ['keyName', 'keyType', 'keyBits'];
    if (type === 'existing') {
      fields = ['keyReference'];
    } else if (type === 'kms') {
      fields = ['keyName', 'managedKeyName', 'managedKeyId'];
    }
    return fields.map((fieldName) => {
      return this.args.model.allFields.find((attr) => attr.name === fieldName);
    });
  }

  @action cancel() {
    this.args.model.unloadRecord();
    this.args.onCancel();
  }
  get groups() {
    return {
      'Key parameters': this.keyParamFields,
      'Subject Alternative Name (SAN) Options': [],
      'Other subject data': [],
    };
  }

  @action
  async generateRoot(event) {
    event.preventDefault();
    const useIssuer = this.args.model.canGenerateIssuerRoot;
    try {
      await this.args.model.save({ adapterOptions: { formType: 'generate-root', useIssuer } });
      this.args.onSave();
    } catch (e) {
      this.errorBanner = errorMessage(e);
      this.invalidFormAlert = 'There was a problem importing key.';
    }
  }
}
