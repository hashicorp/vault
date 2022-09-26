import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

/**
 * @module RadioSelectTtlOrString
 * `RadioSelectTtlOrString` components are yielded out within the formField component when the editType on the model is yield.
 * The component is two radio buttons, where the first option is a TTL, and the second option is an input field without a title.
 * This component is used in the PKI engine inside various forms.
 *
 * @example
 * ```js
 * {{#each @model.fields as |attr|}}
 *  <RadioSelectTtlOrString @attr={{attr}} @model={{this.model}} />
 * {{/each}}
 * ```
 * @callback onChange
 * @param {Model} model - Ember Data model that `attr` is defined on.
 * @param {Object} attr - Usually derived from ember model `attributes` lookup, and all members of `attr.options` are optional.
 * @param {onChange} [onChange] - Callback triggered on save success.
 */

export default class RadioSelectTtlOrString extends Component {
  @tracked groupValue = 'ttl';
  @tracked ttlTime = '';

  @action selectionChange(selection) {
    this.groupValue = selection;
    // Clear the TTL time selection if they have clicked the specificDate radio button.
    if (selection === 'specificDate') {
      this.ttlTime = '';
    }
  }

  @action setAndBroadcastTtl(value) {
    let valueToSet = value.enabled === true ? `${value.seconds}s` : 0;
    this.setAndBroadcast('ttl', `${valueToSet}`);
  }

  @action setAndBroadcastInput(event) {
    const prop = event.target.type === 'checkbox' ? 'checked' : 'value';
    this.setAndBroadcast('not_after', event.target[prop]);
  }

  // Send off the new value and the param name to the parent.
  @action setAndBroadcast(modelParam, value) {
    this.args.onChange(modelParam, value);
  }
}
