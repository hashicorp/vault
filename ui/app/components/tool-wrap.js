import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

export default class HashTool extends Component {
  @tracked data = '{\n}';
  @tracked buttonDisabled = false;

  @action
  onClear() {
    this.args.onClear();
  }
  @action
  updateTtl(evt) {
    if (!evt) return;
    const ttl = evt.enabled ? `${evt.seconds}s` : '30m';
    this.args.updateTtl(ttl);
  }
  @action
  codemirrorUpdated(val, codemirror) {
    codemirror.performLint();
    const hasErrors = codemirror.state.lint.marked.length > 0;
    this.data = val;
    this.buttonDisabled = hasErrors;
    this.args.codemirrorUpdated(val, hasErrors);
  }
}
