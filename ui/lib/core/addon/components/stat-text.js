/**
 * @module StatText
 * StatText components are used to display a label and associated value beneath, with the option to include a description.
 *
 * @example
 * ```js
 * <StatText @label="Active Clients" @stat="4,198" @size="l" @subText="These are the active client counts"/>
 * ```
 * @param {string} label=null - The label for the statistic
 * @param {string} value=null - Value passed in, usually a number or statistic
 * @param {string} size=null - Sizing changes whether or not there is subtext. If there is subtext 's' and 'l' are valid sizes. If no subtext, then 'm' is also acceptable.
 * @param {string} [subText] - SubText is optional and will display below the label
 */

import Component from '@glimmer/component';
import layout from '../templates/components/stat-text';
import { setComponentTemplate } from '@ember/component';

class StatTextComponent extends Component {}

export default setComponentTemplate(layout, StatTextComponent);
