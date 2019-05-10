import { helper as buildHelper } from '@ember/component/helper';

export const MESSAGE_TYPES = {
  info: {
    class: 'is-info',
    glyphClass: 'has-text-info',
    glyph: 'info-circle-outline',
    text: 'Info',
  },
  success: {
    class: 'is-success',
    glyphClass: 'has-text-success',
    glyph: 'check-circle-outline',
    text: 'Success',
  },
  danger: {
    class: 'is-danger',
    glyphClass: 'has-text-danger',
    glyph: 'cancel-square-fill',
    text: 'Error',
  },
  warning: {
    class: 'is-highlight',
    glyphClass: 'has-text-highlight',
    glyph: 'alert-triangle',
    text: 'Warning',
  },
};

export function messageTypes([type]) {
  return MESSAGE_TYPES[type];
}

export default buildHelper(messageTypes);
