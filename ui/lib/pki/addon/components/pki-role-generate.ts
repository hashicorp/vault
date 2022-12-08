import Component from '@glimmer/component';
import { action } from '@ember/object';
import { task } from 'ember-concurrency';
import { service } from '@ember/service';
import Router from '@ember/routing/router';
import Store from '@ember-data/store';
import { tracked } from '@glimmer/tracking';

interface Args {
  onSuccess: CallableFunction;
  model: CertModel;
}

// pki/certificate/generate model
interface CertModel {
  name: string;
  backend: string;
  formFields: FormField;
  formFieldsGroup: {
    [k: string]: FormField[];
  }[];
  save: () => void;
  rollbackAttributes: () => void;
  destroyRecord: () => void;
}
interface FormField {
  name: string;
  type: string;
  options: unknown;
}
export default class PkiRoleGenerate extends Component<Args> {
  @service declare readonly router: Router;
  @service declare readonly store: Store;

  @tracked errorBanner = '';

  @task
  *save(evt: Event) {
    evt.preventDefault();
    this.errorBanner = '';
    const { model, onSuccess } = this.args;
    try {
      yield model.save();
      onSuccess();
    } catch (e) {
      const err = e as {
        errors: string[];
        message: string;
      };
      this.errorBanner = err.errors.join('. ') || err.message;
    }
  }

  @action revoke() {
    // TODO
  }

  @action download() {
    // TODO
  }

  @action cancel() {
    this.args.model.destroyRecord();
    this.router.transitionTo('vault.cluster.secrets.backend.pki.roles.role.details');
  }
}
