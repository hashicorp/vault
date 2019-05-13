const invariant = (truthy, error) => {
  if (!truthy) throw new Error(error);
};

export default function(context, selector) {
  let cmService = context.owner.lookup('service:code-mirror');

  let element = document.querySelector(selector);
  invariant(element, `Selector ${selector} matched no elements`);

  let cm = cmService.instanceFor(element.id);
  invariant(cm, `No registered CodeMirror instance for ${selector}`);

  return cm;
}
