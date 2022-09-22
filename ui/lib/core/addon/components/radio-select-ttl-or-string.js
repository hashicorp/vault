import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

/**
 * @module RadioSelectTtlOrString
 * `RadioSelectTtlOrString` components are ARG TODO
 *
 * @example
 * ```js
 *  <RadioSelectTtlOrString/>
 * ```
//  ARG TODO
 */

export default class RadioSelectTtlOrString extends Component {
  @tracked groupValue = 'ttl'; // ARG TODO some conditional here to change?

  @action onChange(selection) {
    this.groupValue = selection; // ARG TODO work on.
    // this.args.model.set(this.valuePath, selection);
    // this.onChange(this.valuePath, selection);
  }
  @action ttlPickerChange() {
    // ARG TODO something.
  }
  @action
  onChangeWithEvent(event) {
    // ARG TODO finish
    const prop = event.target.type === 'checkbox' ? 'checked' : 'value';
    this.setAndBroadcast(event.target[prop]);
  }
}
