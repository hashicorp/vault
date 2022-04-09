#!/usr/bin/env node
/* eslint-env node */
/**
 * Codemod to transform Icon component to new API to accomodate FlightIcon
 * example execution from ui directory -> npx ember-template-recast ./templates -t ./scripts/codemods/icon-transform.js
 * above will run transform on all files in templates directory
 */

module.exports = (env) => {
  const {
    syntax: { builders },
  } = env;
  const hsSizes = ['s', 'm', 'l', 'xlm', 'xl', 'xxl'];

  // find attribute by name
  const findAttribute = (attrs, name) => {
    for (let i = 0; i < attrs.length; i++) {
      if (attrs[i].name === name) {
        return [attrs[i], i, attrs[i].value.chars];
      }
    }
    return [];
  };

  // possibly a bug with ember-template-recast for multi line components with attributes on their own lines
  // when removing an attribute the one on the line below will jump to the same line as the previous one
  // this does not happen when removing the first attribute -- doing so may add unnecessary quotes to the first shifted attribute
  // example: class="{{foo}}" -> class=""{{foo}}""
  const preserveFormatting = (attributes, removeIndex) => {
    if (removeIndex > 0) {
      // shift the location of the attributes that appear after the one being removed to preserve formatting
      for (let i = attributes.length - 1; i > removeIndex; i--) {
        attributes[i].loc = attributes[i - 1].loc;
      }
    }
  };

  // transform structure icon size letter to flight icon supported size
  const transformSize = (attributes, attrName) => {
    const [attr, attrIndex, value] = findAttribute(attributes, attrName);

    if (hsSizes.includes(value)) {
      if (['s', 'm', 'l'].includes(value)) {
        // before removing attribute set the location of the remaining attributes
        preserveFormatting(attributes, attrIndex);
        // since 16 is the default in the component we can remove the attribute
        attributes.splice(attrIndex, 1);
      } else {
        attr.value = builders.text('24');
        // rename attribute
        if (attrName === '@sizeClass') {
          attr.name = '@size';
        }
      }
    }
  };

  return {
    ElementNode(node) {
      if (node.tag === 'Icon') {
        const { attributes } = node;
        // the inital refactor of the component introduced a sizeClass attribute
        // this can now be mapped to size and removed
        transformSize(attributes, '@sizeClass');
        // check for old component instances that may still have a letter for size value
        transformSize(attributes, '@size');
      }
    },
  };
};
