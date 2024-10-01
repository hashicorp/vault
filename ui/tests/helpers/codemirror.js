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
  const selector = parent ? `${parent} .CodeMirror` : '.CodeMirror';
  const element = document.querySelector(selector);
  invariant(element, `Selector '.CodeMirror' matched no elements`);

  const cm = element.CodeMirror;
  invariant(cm, `No registered CodeMirror instance for '.CodeMirror'`);

  return cm;
}

const invariant = (truthy, error) => {
  if (!truthy) throw new Error(error);
};
