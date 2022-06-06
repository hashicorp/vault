import Component from '@glimmer/component';
import layout from '../../templates/components/linkable-item/menu';
import { setComponentTemplate } from '@ember/component';

/**
 * @module Menu
 * Menu components are contextual components of LinkableItem, used to display a menu on the right side of a LinkableItem component.
 *
 * @example
 * ```js
 * <LinkableItem as |Li|>
 *  <Li.menu>
 *     Some menu here
 *  </Li.menu>
 * </LinkableItem>
 * ```
 */

/* eslint ember/no-empty-glimmer-component-classes: 'warn' */
class MenuComponent extends Component {}

export default setComponentTemplate(layout, MenuComponent);
