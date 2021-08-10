/**
 * @module StatText
 * StatText components are used to display a label and associated statistic below, with the option to add a description.
 *
 * @example
 * ```js
 * <StatText @label="Active Clients" @stat="4,198" @size="l" @subText="These are the active client counts"/>
 * ```
 * @param {string} label=null - the label for the statistic
 * @param {string} value=null - value passed in, usually a number or statistic
 * @param {string} size=null - size the component as small or large, 's', 'm' or 'l'
 * @param {string} [subText] - subText is optional and will display below the label
 */

import Component from '@glimmer/component';
import layout from '../templates/components/stat-text';
import { setComponentTemplate } from '@ember/component';

class StatTextComponent extends Component {}

export default setComponentTemplate(layout, StatTextComponent);
