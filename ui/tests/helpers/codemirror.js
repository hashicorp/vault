/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/*
returns an instance of CodeMirror, see docs for callable functions https://codemirror.net/5/doc/manual.html#api_constructor
If you are targeting a specific CodeMirror instance, pass the selector of the parent element as an argument.
sample use:

  import codemirror from 'vault/tests/helpers/codemirror';

  test('it renders initial value', function (assert) {
    // General use
    assert.strictEqual(codemirror().getValue(), 'some other value')
    // Specific selector
    codemirror('#my-control').setValue('some value');
    assert.strictEqual(codemirror('#my-control').getValue(), 'some value')
  )}
*/
export default function (parent) {
  const selector = parent ? `${parent} .hds-code-editor__editor` : '.hds-code-editor__editor';
  const element = document.querySelector(selector);
  invariant(element, `Selector '${selector}' matched no elements`);

  const cm = element.editor;
  invariant(cm, `No registered CodeMirror instance for ''${selector}'`);

  return cm;
}

export function setCodeEditorValue(editorView, value, { from, to } = {}) {
  invariant(editorView, 'No editor view provided');
  invariant(value, 'No value provided');

  editorView.dispatch({
    changes: [
      {
        from: from ?? 0,
        to: to ?? editorView.state.doc.length,
        insert: value,
      },
    ],
  });
}

export function getCodeEditorValue(editorView) {
  invariant(editorView, 'No editor view provided');

  return editorView.state.doc.toString();
}

const invariant = (truthy, error) => {
  if (!truthy) throw new Error(error);
};
