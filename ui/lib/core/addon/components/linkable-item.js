/**
 * @module LinkableItem
 * LinkableItem components are used to show information on the left with a menu on the right, aligned vertically centered. If passed a link, the block will be clickable
 *
 * @example
 * ```js
 * <LinkableItem @link={{hash route='vault.backends' model='my-backend-path'}} data-test-row="my-backend-path" as |Li| />
 * ```
 * @param {object} [link] - link should have route and model
 */

import Component from '@glimmer/component';
import layout from '../templates/components/linkable-item';
import { setComponentTemplate } from '@ember/component';

class LinkableItemComponent extends Component {}

export default setComponentTemplate(layout, LinkableItemComponent);
