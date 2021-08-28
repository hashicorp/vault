/**
 * @module DummyParentComponent
 * DummyParentComponent components are used to...
 *
 * @example
 * ```js
 * <DummyParentComponent @requiredParam={requiredParam} @optionalParam={optionalParam} @param1={{param1}}/>
 * ```
 * @param {object} requiredParam - requiredParam is...
 * @param {string} [optionalParam] - optionalParam is...
 * @param {string} [param1=defaultValue] - param1 is...
 */

import Component from '@glimmer/component';
import layout from '../templates/components/bar-chart';
import { setComponentTemplate } from '@ember/component';

class DummyParentComponent extends Component {}

export default setComponentTemplate(layout, DummyParentComponent);
