/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
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
