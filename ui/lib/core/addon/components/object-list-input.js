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
  @tracked inputKeys;
  @tracked disableAdd = true;

  constructor() {
    super(...arguments);
    this.inputKeys = this.args.objectKeys.map((e) => e.key);
    this.inputList = [this.createInputRow(this.inputKeys)];
  }

  createInputRow(keys) {
    // create a new object with from input keys with empty values
    return Object.fromEntries(keys.map((key) => [key, '']));
  }

  @action
  handleInput(idx, { target }) {
    const inputObj = this.inputList.objectAt(idx);
    inputObj[target.name] = target.value;
    this.handleChange();
  }

  @action
  addRow() {
    const newRow = this.createInputRow(this.inputKeys);
    this.inputList.pushObject(newRow);
    this.disableAdd = true;
  }

  @action
  removeRow(idx) {
    const row = this.inputList.objectAt(idx);
    this.inputList.removeObject(row);
    this.handleChange();
  }

  @action
  handleChange() {
    // disable/enable "add" button
    const lastObject = this.inputList[this.inputList.length - 1];
    this.disableAdd = Object.values(lastObject).any((input) => input === '') ? true : false;

    // don't send an empty last object to parent
    if (Object.values(lastObject).every((input) => input === '')) {
      this.args.onChange(this.inputList.slice(0, -1));
    } else {
      this.args.onChange(this.inputList);
    }
  }
}
