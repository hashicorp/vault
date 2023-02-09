import { helper as buildHelper } from '@ember/component/helper';
import { assert } from '@ember/debug';

export const MESSAGE_TYPES = {
  info: {
    class: 'is-info',
    glyphClass: 'has-text-info',
    glyph: 'info',
    text: 'Info',
  },
  success: {
    class: 'is-success',
    glyphClass: 'has-text-success',
    glyph: 'check-circle-fill',
    text: 'Success',
  },
  danger: {
    class: 'is-danger',
    glyphClass: 'has-text-danger',
    glyph: 'x-square-fill',
    text: 'Error',
  },
  warning: {
    class: 'is-highlight',
    glyphClass: 'has-text-highlight',
    glyph: 'alert-triangle-fill',
    text: 'Warning',
  },
  loading: {
    class: 'is-success',
    glyphClass: 'has-text-success',
    glyph: 'loading',
    text: 'Loading',
  },
};

export function messageTypes([type]) {
  if (!(type in MESSAGE_TYPES)) {
    assert('type is not a valid message type.');
  }
  return MESSAGE_TYPES[type];
}

export default buildHelper(messageTypes);
