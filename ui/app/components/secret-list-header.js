import Component from '@glimmer/component';

export default class SecretListHeader extends Component {
  tagName = '';

  // api
  isCertTab = false;
  isConfigure = false;
  baseKey = null;
  backendCrumb = null;
  model = null;
  options = null;
}
