/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/*
returns an instance of CodeMirror, see docs for callable functions https://codemirror.net/5/doc/manual.html#api_constructor
sample use:

  import codemirror from 'vault/tests/helpers/codemirror';

  test('it renders initial value', function (assert) {

    assert.strictEqual(codemirror.getValue(), 'some value')
  )}
*/

const invariant = (truthy, error) => {
  if (!truthy) throw new Error(error);
};

export default function () {
  const element = document.querySelector('.CodeMirror');
  invariant(element, `Selector '.CodeMirror' matched no elements`);

  const cm = element.CodeMirror;
  invariant(cm, `No registered CodeMirror instance for '.CodeMirror'`);

  return cm;
}
