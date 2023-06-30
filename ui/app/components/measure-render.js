import { action } from '@ember/object';
import Component from '@glimmer/component';

export default class MeasureRenderComponent extends Component {
  constructor() {
    super(...arguments);
    performance.mark(`start-${this.args.name}`);
  }
  @action didRender() {
    console.log('did render', this.args.capture);
    performance.mark(`end-${this.args.name}`);
    if (this.args.capture) {
      const captureMeasure = performance.measure(...this.args.capture);
      // console.log('Capture measure', ...this.args.capture, captureMeasure.duration);
      const dataFetchMeasure = performance.measure('data-load-duration', 'list-load', 'list-end');
      // const renderMeasure = performance.measure('render-duration', 'start-top', 'end-bottom');
      const fullMeasure = performance.measure('full-measure', 'list-load', 'end-bottom');
      console.log('Data Fetch | Render Page | Total time');
      console.log(`| ${dataFetchMeasure.duration} | ${captureMeasure.duration} | ${fullMeasure.duration} |`);
    }
  }
}
