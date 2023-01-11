import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

/**
 * @module ObjectListInput
 * ObjectListInput components are used to...
 *
 * @example
 * ```js
 * <ObjectListInput @requiredParam={requiredParam} @optionalParam={optionalParam} @param1={{param1}}/>
 * ```
 * @param {array} objectKeys - array of strings that correspond to object keys
 * @param {string} [optionalParam] - optionalParam is...
 * @param {string} [param1=defaultValue] - param1 is...
 */

export default class ObjectListInput extends Component {
  @tracked inputList = [];
  @tracked inputRow = {};
  constructor() {
    super(...arguments);
    this.args.objectKeys.forEach((e) => (this.inputRow[e.name] = ''));
    this.inputList.push(this.inputRow);
  }

  @action
  handleInput(idx, { target }) {
    const inputObj = this.inputList.objectAt(idx);
    inputObj[target.name] = target.value;
  }

  @action
  addInput() {
    // handle this
  }
  @action
  removeInput() {
    // handle this
  }
}
