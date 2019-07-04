import HoverDropdown from 'ember-basic-dropdown-hover/components/basic-dropdown-hover';
import layout from '../templates/components/tool-tip';

export default HoverDropdown.extend({
  layout,
  delay: 0,
  horizontalPosition: 'auto-right',
});
