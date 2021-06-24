import Component from '@glimmer/component';
import layout from '../../templates/components/linkable-item/content';
import { setComponentTemplate } from '@ember/component';

class ContentComponent extends Component {}

export default setComponentTemplate(layout, ContentComponent);
