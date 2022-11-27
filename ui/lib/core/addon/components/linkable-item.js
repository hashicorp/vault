import Component from '@glimmer/component';
import layout from '../templates/components/linkable-item';
import { setComponentTemplate } from '@ember/component';

/**
 * @module LinkableItem
 * LinkableItem components have two contextual components, a Content component used to show information on the left with a Menu component on the right, all aligned vertically centered. If passed a link, the block will be clickable.
 *
 * @example
 * ```js
 * <LinkableItem @link={{hash route='vault.backends' model='my-backend-path'}} data-test-row="my-backend-path">
 * // Use <LinkableItem.content> and <LinkableItem.menu> here
 * </LinkableItem>
 * ```
 *
 * @param {object} [link=null] - Link should have route and model
 * @param {boolean} [disabled=false] - If no link then should be given a disabled attribute equal to true
 */

/* eslint ember/no-empty-glimmer-component-classes: 'warn' */
class LinkableItemComponent extends Component {}

export default setComponentTemplate(layout, LinkableItemComponent);
