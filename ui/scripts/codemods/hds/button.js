/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/* eslint-env node */

/**
 * codemod to transform button html element to Hds::Button component
 * transformation is skipped if is-ghost or is-transparent is found in class list
 * if loading or is-loading is found to be a conditionally applied class the loading icon will be conditionally applied instead
 * if the text arg cannot be built from the child nodes (chained if block or multiple nodes that cannot be easily combined) the transformation will be skipped
 * classes relevant to the legacy button will be removed (see classesToRemove array)
 * html onclick event handler will be replaced with the {{on "click"}} modifier
 *
 * example execution from ui directory:
 ** -> npx ember-template-recast ./app/templates -t ./scripts/codemods/hds/button.js
 * for best results run prettier after:
 ** -> npx ember-template-recast ./app/templates -t ./scripts/codemods/hds/button.js && npx prettier --config .prettierrc.js --write ./app/templates
 */

class Transforms {
  // button classes that will be removed from attribute
  classesToRemove = [
    'button',
    'is-compact',
    'is-danger',
    'is-danger-outlined',
    'is-flat',
    'is-icon',
    'is-loading',
    'is-link',
    'is-primary',
    'tool-tip-trigger',
    'is-secondary',
  ];
  classesToTransform = [{ current: 'toolbar-link', updated: 'toolbar-button' }];

  constructor(node, builders) {
    this.node = node;
    this.attrs = [];
    this.modifiers = [...node.modifiers];
    this.builders = builders;
    this.hasIcon = false;
    this.hasText = false;
  }

  shouldTransform() {
    // buttons that have the is-ghost and/or is-transparent class will not be transformed
    // these usages have unclear mappings to tertiary buttons and in some cases will be replaced with Hds::Interactive
    const classAttr = this.node.attributes.find((attr) => attr.name === 'class');
    if (classAttr) {
      const shouldTransform = (chars) => {
        return chars.includes('is-ghost') || chars.includes('is-transparent') ? false : true;
      };
      if (classAttr.value.type === 'ConcatStatement') {
        for (const part of classAttr.value.parts) {
          if (part.type === 'TextNode' && !shouldTransform(part.chars)) {
            return false;
          }
        }
      } else {
        return shouldTransform(classAttr.value.chars);
      }
    }
    return true;
  }

  addAttr(name, value) {
    this.attrs.push(this.builders.attr(name, value));
  }

  filterClassTextNode(value) {
    // map color related classes to @color args
    let color = 'secondary'; // currently the default for .button class
    for (const colorClass of ['is-primary', 'is-danger', 'is-danger-outlined']) {
      if (value.chars.includes(colorClass)) {
        color = colorClass === 'is-primary' ? null : 'critical';
        break;
      }
    }
    if (color) {
      this.addAttr('@color', this.builders.text(color));
    }
    // remove button related classes no longer needed
    // map unused classes to new ones
    const classArray = value.chars.split(' ');
    const chars = classArray
      .filter((className) => !this.classesToRemove.includes(className))
      .map((className) => {
        const transform = this.classesToTransform.find((classHash) => classHash.current === className);
        return transform?.updated || className;
      })
      .join(' ');
    return chars ? { ...value, chars } : null;
  }

  convertIsLoadingMustache(part, filteredParts) {
    let isLoading = false;
    const filteredParams = part.params.map((param) => {
      if (param.type === 'StringLiteral' && param.value.includes('loading')) {
        // rebuild param since icon name is loading and class name could be is-loading
        isLoading = true;
        return this.builders.string('loading');
      }
      return param;
    });
    if (isLoading) {
      this.addAttr('@icon', this.builders.mustache('if', filteredParams));
    } else {
      filteredParts.push(part);
    }
  }

  filterClassConcatStatement(attr) {
    const filteredParts = [];
    attr.value.parts.forEach((part) => {
      if (part.type === 'TextNode') {
        const value = this.filterClassTextNode(part);
        if (value) {
          filteredParts.push(value);
        }
      } else if (part.type === 'MustacheStatement') {
        this.convertIsLoadingMustache(part, filteredParts);
      } else {
        filteredParts.push(part);
      }
    });
    if (filteredParts.length) {
      return filteredParts.length === 1 ? filteredParts[0] : { ...attr.value, parts: filteredParts };
    }
  }

  filterClasses(attr) {
    if (attr.name === 'class') {
      let attrValue = attr.value;
      const { type } = attrValue;
      if (type === 'ConcatStatement') {
        attrValue = this.filterClassConcatStatement(attr);
      } else if (type === 'TextNode') {
        attrValue = this.filterClassTextNode(attr.value);
      }
      if (attrValue) {
        this.addAttr('class', attrValue);
      }
    }
  }

  convertOnClick(attr) {
    const params = [this.builders.string('click')];
    if (!attr.value.params.length) {
      params.push(attr.value.path);
    } else {
      params.push(this.builders.sexpr(attr.value.path, attr.value.params));
    }
    const onClickModifier = this.builders.elementModifier('on', params);
    this.modifiers.push(onClickModifier);
  }

  filterAttributes() {
    this.node.attributes.forEach((attr) => {
      if (attr.name === 'class') {
        return this.filterClasses(attr);
      } else if (attr.name === 'onclick') {
        return this.convertOnClick(attr);
      } else if (attr.name === 'type' && attr.value.chars === 'button') {
        // remove type="button" attribute since it is default
        return;
      }
      this.attrs.push(attr);
    });
  }

  textToString(node) {
    // filter out escape charaters like \n and whitespace from TextNode and rebuild as StringLiteral
    const text = decodeURI(node.chars).trim();
    if (text) {
      return this.builders.string(text);
    }
  }

  filterTextNode(node, parts) {
    if (node.type === 'TextNode') {
      const text = this.textToString(node);
      if (text) {
        parts.push(text);
      }
    }
  }

  convertBlockStatementNode(node, parts) {
    // convert if/else block statement to inline if mustache
    if (node.type === 'BlockStatement' && node.path.original === 'if' && !node.inverse.chained) {
      // only deal with text nodes -- more complex expressions should be converted to getter on component
      const program = node.program.body;
      const ifValueNode = program.length === 1 && program[0].type === 'TextNode' ? program[0] : null;
      const inverse = node.inverse.body;
      const elseValueNode = inverse.length === 1 && inverse[0].type === 'TextNode' ? inverse[0] : null;

      if (ifValueNode && elseValueNode) {
        const params = [...node.params, this.textToString(ifValueNode), this.textToString(elseValueNode)];
        parts.push(this.builders.mustache(node.path, params));
      }
    }
  }

  convertIconNode(node) {
    if (node.tag === 'Icon') {
      const nameAttr = node.attributes.find((attr) => attr.name === '@name');
      this.addAttr('@icon', this.builders.string(nameAttr.value.chars));
      // Hds::Button has @iconPosition arg when used with text
      // it seems most usages with button are leading which is default and recommended
      this.hasIcon = true;
    }
  }

  pushAcceptedNodes(node, parts) {
    // some nodes may not need conversion and can be added to the @text assembly as is
    const acceptedNodes = ['MustacheStatement'];
    if (acceptedNodes.includes(node.type)) {
      parts.push(node);
    }
  }

  childNodesToArgs() {
    // convert child nodes to a format supported by an attr value for @text arg
    const parts = [];
    this.node.children.forEach((node) => {
      // following methods are used to build the @text arg
      this.filterTextNode(node, parts);
      this.convertBlockStatementNode(node, parts);
      this.pushAcceptedNodes(node, parts);
      // we also need to set the icon related args
      this.convertIconNode(node);
    });

    // filter out ignored text nodes (\n) and compare with out compiled parts
    // if the lengths do not match then we were unable to transform a part and we must abort text build
    const relevantParts = this.node.children.filter((node) => {
      if (node.type === 'TextNode' && !this.textToString(node)) {
        return false;
      }
      return true;
    });
    if (parts.length && relevantParts.length === parts.length) {
      const value = parts.length === 1 ? parts[0] : this.builders.concat(parts);
      this.addAttr('@text', value);
      this.hasText = true;
    } else if (this.hasIcon) {
      // if there was an icon node but no text we need to add the @isIconOnly arg
      this.addAttr('@isIconOnly', this.builders.mustache(this.builders.boolean(true)));
    }
  }

  buildElement() {
    if (this.hasText || this.hasIcon) {
      return this.builders.element(
        { name: 'Hds::Button', selfClosing: true },
        { attrs: this.attrs, modifiers: this.modifiers }
      );
    }
  }
}

module.exports = (env) => {
  const { builders } = env.syntax;

  return {
    ElementNode(node) {
      if (node.tag === 'button') {
        try {
          const transforms = new Transforms(node, builders);
          if (transforms.shouldTransform()) {
            transforms.childNodesToArgs();
            transforms.filterAttributes();
            return transforms.buildElement();
          }
        } catch (error) {
          console.log(`\nError caught transforming button in ${env.filePath}\n`, error); // eslint-disable-line
        }
      }
    },
  };
};
