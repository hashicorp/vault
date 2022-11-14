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
 * @param {Model} model - Ember Data model that `attr` is defined on.
 * @param {Object} attr - Usually derived from ember model `attributes` lookup, and all members of `attr.options` are optional.
 */

export default class RadioSelectTtlOrString extends Component {
  @tracked groupValue = 'ttl';
  @tracked ttlTime;
  @tracked notAfter;

  @action onRadioButtonChange(selection) {
    this.groupValue = selection;
    // Clear the previous selection if they have clicked the other radio button.
    if (selection === 'specificDate') {
      this.args.model.set('ttl', '');
      this.ttlTime = '';
    }
    if (selection === 'ttl') {
      this.args.model.set('notAfter', '');
      this.notAfter = '';
      this.args.model.set('ttl', this.ttlTime);
    }
  }

  @action setAndBroadcastTtl(value) {
    const valueToSet = value.enabled === true ? `${value.seconds}s` : 0;
    if (this.groupValue === 'specificDate') {
      // do not save ttl on the model until the ttl radio button is selected
      return;
    }
    this.args.model.set('ttl', `${valueToSet}`);
  }

  @action setAndBroadcastInput(event) {
    this.args.model.set('notAfter', event.target.value);
  }
}
