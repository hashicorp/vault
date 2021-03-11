/**
 * @module BoxRadioSet
 * BoxRadioSet components are used to wrap radio options, specifically sets of BoxRadio
 *
 * @example
 * ```js
 * <BoxRadioSet @title="My Category">
 *   <BoxRadio ... />
 *   <BoxRadio ... />
 *   <BoxRadio ... />
 * </BoxRadioSet>
 * ```
 * @param {string} [title] - Optional text that is the heading for the section
 */

import Component from '@glimmer/component';
import layout from '../templates/components/box-radio-set';
import { setComponentTemplate } from '@ember/component';

class BoxRadioSet extends Component {}

export default setComponentTemplate(layout, BoxRadioSet);
