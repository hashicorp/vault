import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

/**
 * @module RadioSelectTtlOrString
 * `RadioSelectTtlOrString` components are form field type that is yielded out by the model editType: yield.
 * The component is two radio buttons, the first a ttl and the second something similar to editType optionalText: an input field without a title.
 * This component is used in the PKI engine in various forms.
 *
 * @example
 * ```js
 * {{#each @model.fields as |attr|}}
 *  <RadioSelectTtlOrString @attr={{attr}} @model={{this.model}} />
 * {{/each}}
 * ```
 * @param {Model} model - Ember Data model that `attr` is defined on.
 * @param {Object} attr - usually derived from ember model `attributes` lookup, and all members of `attr.options` are optional.
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
