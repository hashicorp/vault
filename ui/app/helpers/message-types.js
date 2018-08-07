import Ember from 'ember';

const MESSAGE_TYPES = {
  info: {
    class: 'is-info',
    glyphClass: 'has-text-info',
    glyph: 'information-circled',
    text: 'Info',
  },
  success: {
    class: 'is-success',
    glyphClass: 'has-text-success',
    glyph: 'checkmark-circled',
    text: 'Success',
  },
  danger: {
    class: 'is-danger',
    glyphClass: 'has-text-danger',
    glyph: 'close-circled',
    text: 'Error',
  },
  warning: {
    class: 'is-highlight',
    glyphClass: 'has-text-highlight',
    glyph: 'alert-circled',
    text: 'Warning',
  },
};

export function messageTypes([type]) {
  return MESSAGE_TYPES[type];
}

export default Ember.Helper.helper(messageTypes);
