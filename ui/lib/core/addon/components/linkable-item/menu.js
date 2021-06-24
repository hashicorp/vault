import Component from '@glimmer/component';
import layout from '../../templates/components/linkable-item/menu';
import { setComponentTemplate } from '@ember/component';

class MenuComponent extends Component {}

export default setComponentTemplate(layout, MenuComponent);
