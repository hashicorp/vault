import Component from '@ember/component';

/**
 * @module Select
 * Select components are used to render a dropdown.
 *
 * @example
 * ```js
 * <Select label='Date Range' @options={{[{ value: 'berry', label: 'Berry' }]}} @onChange={{onChange}}/>
 * ```
 *
 * @param label=null {String} - The label for the select element.
 * @param options=null {Array} - A list of items that the user will select from. This can be an array of strings or objects.
 * @param [valueAttribute=value] {String} - When `options` is an array objects, the key to check for when assigning the option elements value.
 * @param [labelAttribute=label] {String} - When `options` is an array objects, the key to check for when assigning the option elements' inner text.
 * @param [isInline=false] {Bool} - Whether or not the select should be displayed as inline-block or block.
 * @param [isFullwidth=false] {Bool} - Whether or not the select should take up the full width of the parent element.
 * @param onChange=null {Func} - The action to take once the user has selected an item.
 */

export default Component.extend({
  classNames: ['field'],
  label: null,
  options: null,
  valueAttribute: 'value',
  labelAttribute: 'label',
  isInline: false,
  isFullwidth: false,
  onChange() {},
});
