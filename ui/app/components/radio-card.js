import Component from '@glimmer/component';

/**
 * @module RadioCard
 * RadioCard components are used to...
 *
 * @example
 * ```js
 * <RadioCard @requiredParam={requiredParam} @optionalParam={optionalParam} @param1={{param1}}/>
 * ```
 * @param {string} title - Title inside the card.
 * @param {string} [description] - Description text under the title.
 * @param {string} icon - Icon to the left of the title.
 * @param {string} value - The value associated with the selected Card.
 * @param {string} groupValue - The selected value of the radio card group in inline mode -- new, existing or skip are the accepted values. ARG TODO XX.
 * @param {function} onChange - The callback function triggered on selecting a radio card.
 * @param {boolean} isHalfWidth - If there are two cards only, change width from 19rem to 38rem.
 */
/* eslint ember/no-empty-glimmer-component-classes: 'warn' */
export default class RadioCard extends Component {}
