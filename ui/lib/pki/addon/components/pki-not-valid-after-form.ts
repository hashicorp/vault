import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { HTMLElementEvent } from 'forms';

/**
 * <PkiNotValidAfterForm /> components are used to manage two mutually exclusive role options in the form.
 */
interface Args {
  model: {
    notAfter: string;
    ttl: string;
    set: (key: string, value: string) => void;
  };
}

export default class RadioSelectTtlOrString extends Component<Args> {
  @tracked groupValue = 'ttl';
  @tracked originalNotAfter: string;
  @tracked originalTtl: string;

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    const { model } = this.args;
    this.originalNotAfter = model.notAfter;
    this.originalTtl = model.ttl;
    if (model.notAfter) {
      this.groupValue = 'specificDate';
    }
  }

  @action onRadioButtonChange(selection: string) {
    this.groupValue = selection;
    // Clear the previous selection if they have clicked the other radio button.
    if (selection === 'specificDate') {
      this.args.model.set('ttl', '');
      this.args.model.set('notAfter', this.originalNotAfter);
    }
    if (selection === 'ttl') {
      this.args.model.set('notAfter', '');
      this.args.model.set('ttl', this.originalTtl);
    }
  }

  @action setAndBroadcastTtl(ttlObject: { enabled: boolean; goSafeTimeString: string }) {
    const { enabled, goSafeTimeString } = ttlObject;
    if (this.groupValue === 'specificDate') {
      // do not save ttl on the model unless the ttl radio button is selected
      return;
    }
    this.args.model.set('ttl', enabled === true ? goSafeTimeString : '0');
  }

  @action setAndBroadcastInput(evt: HTMLElementEvent<HTMLInputElement>) {
    const setDate = evt.target.valueAsDate;
    if (!setDate) return;
    this.args.model.set('notAfter', setDate.toISOString());
  }
}
