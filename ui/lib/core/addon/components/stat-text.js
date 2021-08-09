/**
 * @module StatText
 * StatText components are used to display a label and associated statistic below, with the option to add a description.
 *
 * @example
 * ```js
 *
 * <StatText @label="label" @stat="number" @subText="I am optional subtext" />
 * <StatText
 *  @label={label}
 *  @stat={stat}
 *  @size={size}
 *  @subText={subText}
 * ```
 * @param {string} label - Label is the name of the statistic
 * @param {string} stat - stat is the integer passed in to be displayed
 * @param {string} size - sizes component as large or small
 * @param {string} [subText] - subText is optional and will display below the label
 */

import Component from '@glimmer/component';
import layout from '../templates/components/stat-text';
import { setComponentTemplate } from '@ember/component';

class StatTextComponent extends Component {}

export default setComponentTemplate(layout, StatTextComponent);
