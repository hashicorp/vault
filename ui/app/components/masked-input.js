import Component from '@ember/component';
import { computed } from '@ember/object';
import autosize from 'autosize';

export default Component.extend({
  value: null,
  placeholder: 'value',
  didInsertElement() {
    this._super(...arguments);
    autosize(this.element.querySelector('textarea'));
  },
  didUpdate() {
    this._super(...arguments);
    autosize.update(this.element.querySelector('textarea'));
  },
  willDestroyElement() {
    this._super(...arguments);
    autosize.destroy(this.element.querySelector('textarea'));
  },
  shouldObscure: computed('isMasked', 'isFocused', 'value', function() {
    if (this.get('value') === '') {
      return false;
    }
    if (this.get('isFocused') === true) {
      return false;
    }
    return this.get('isMasked');
  }),
  displayValue: computed('shouldObscure', function() {
    if (this.get('shouldObscure')) {
      return '■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■ ■';
    } else {
      return this.get('value');
    }
  }),
  isMasked: true,
  isFocused: false,
  displayOnly: false,
  onKeyDown() {},
  onChange() {},
  actions: {
    toggleMask() {
      this.toggleProperty('isMasked');
    },
    updateValue(e) {
      this.set('value', e.target.value);
      this.onChange();
    },
  },
});
