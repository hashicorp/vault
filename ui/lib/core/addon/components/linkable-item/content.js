import Component from '@glimmer/component';
import layout from '../../templates/components/linkable-item/content';
import { setComponentTemplate } from '@ember/component';

/**
 * @module Content
 * Content components are contextual components of LinkableItem, used to display content on the left side of a LinkableItem component.
 *
 * @example
 * ```js
 * <LinkableItem as |Li|>
 *  <Li.content
 *    @accessor="cubbyhole_e21f8ee6"
 *    @description="per-token private secret storage"
 *    @glyphText="tooltip text"
 *    @glyph=glyph
 *    @title="title"
 *  />
 * </LinkableItem>
 * ```
 * @param {string} accessor=null - formatted as HTML <code> tag
 * @param {string} description=null - will truncate if wider than parent div
 * @param {string} glyphText=null - tooltip for glyph
 * @param {string} glyph=null - will display as icon beside the title
 * @param {string} title=null - if @link object is passed in then title will link to @link.route
 */

/* eslint ember/no-empty-glimmer-component-classes: 'warn' */
class ContentComponent extends Component {}

export default setComponentTemplate(layout, ContentComponent);
