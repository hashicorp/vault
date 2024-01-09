import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { task } from 'ember-concurrency';
import errorMessage from 'vault/utils/error-message';

interface Args {
  onSubmit: () => Promise<void>;
  onCancel: CallableFunction;
  onValidate?: CallableFunction;
}
interface AllByKey {
  [key: string]: unknown;
}
interface Group {
  fields: string[];
  title?: string;
  toggles?: boolean;
  hasDivider?: boolean;
}

export default class NewFormComponent extends Component<Args> {
  @tracked errorMessage = '';
  @tracked validations = null;

  isValid() {
    if (this.args.onValidate) {
      const { isValid, state } = this.args.onValidate();
      this.validations = state;
      return isValid;
    }
    // No validation method provided, so assume valid
    return true;
  }
  @task *handleSubmit(evt: FormDataEvent) {
    evt.preventDefault();
    this.errorMessage = '';
    if (!this.isValid()) {
      return;
    }
    try {
      yield this.args.onSubmit();
    } catch (e) {
      // TODO: handle control group error
      this.errorMessage = errorMessage(e);
    }
  }

  fieldValidity = (validations: any) => {
    return (fieldName: string) => {
      return validations && validations[fieldName];
    };
  };

  modelFields = () => {
    return (fieldNames: string[], allByKey: AllByKey) => {
      return fieldNames.map((fieldName) => allByKey[fieldName]).filter((f) => !!f);
    };
  };

  modelGroups = () => {
    return (groups: Group[], allByKey: AllByKey) => {
      return groups.map((group) => {
        const expanded = group.fields.map((fieldName) => allByKey[fieldName]).filter((f) => !!f);
        return { ...group, fields: expanded };
      });
    };
  };
}
