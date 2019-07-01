/**
 * @module ToolbarFilters
 * `ToolbarFilters` components are containers for Toolbar filters and toggles.
 * It should only be used inside of `Toolbar`.
 *
 * @example
 * ```js
 * <Toolbar>
 *   <ToolbarFilters>
 *     <div class="control has-icons-left">
 *       <input class="filter input" placeholder="Filter keys" type="text">
 *       <Icon @glyph="search" @size="l" class="search-icon has-text-grey-light" />
 *     </div>
 *   </ToolbarFilters>
 * </Toolbar>
 ```
 *
 */

import Component from '@ember/component';
import layout from '../templates/components/toolbar-filters';

export default Component.extend({
  layout,
  tagName: '',
});
