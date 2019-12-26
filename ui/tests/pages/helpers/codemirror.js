import getCodeMirrorInstance from 'vault/tests/helpers/codemirror';
// Like fillable, but for the CodeMirror editor
//
// Usage: fillIn: codeFillable('[data-test-editor]')
//        Page.fillIn(code);
export function codeFillable(selector) {
  return {
    isDescriptor: true,

    get() {
      return function(context, code) {
        const cm = getCodeMirrorInstance(context, selector);
        cm.setValue(code);
        return this;
      };
    },
  };
}

// Like text, but for the CodeMirror editor
//
// Usage: content: code('[data-test-editor]')
//        Page.code(); // some = [ 'string', 'of', 'code' ]
export function code(selector) {
  return {
    isDescriptor: true,
    get() {
      return function(context) {
        const cm = getCodeMirrorInstance(context, selector);
        return cm.getValue();
      };
    },
  };
}
