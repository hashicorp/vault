/**
 * @module LinkableItem
 * LinkableItem components are used to...
 *
 * @example
 * ```js
 * <LinkableItem @requiredParam={requiredParam} @optionalParam={optionalParam} @param1={{param1}}/>
 * ```
 * @param {object} requiredParam - requiredParam is...
 * @param {string} [optionalParam] - optionalParam is...
 * @param {string} [param1=defaultValue] - param1 is...
 */

import Component from '@glimmer/component';
import layout from '../templates/components/linkable-item';
import { setComponentTemplate } from '@ember/component';

class LinkableItemComponent extends Component {}

export default setComponentTemplate(layout, LinkableItemComponent);
