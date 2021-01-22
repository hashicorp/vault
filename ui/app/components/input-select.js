// ARG TODO only keep what you need, then glimmerize

import Component from '@ember/component';

export default Component.extend({
  /*
   * @public
   * @param Function
   *
   * Function called when any of the inputs change
   *
   */
  onChange: () => {},

  actions: {
    inputChanged(ind, val) {
      const onChange = this.onChange;
      onChange(val);
    },
  },
});

// import Component from '@glimmer/component';
// import { action } from '@ember/object';
// import { tracked } from '@glimmer/tracking';

// export default class inputSelect extends Component {
//   /*
//    * @public
//    * @param Function
//    *
//    * Function called when any of the inputs change
//    *
//    */
//   onChange = () => {};
//   @tracked inputValue = '';
//   @action
//   inputChanged(val) {
//     const onChange = this.onChange;
//     onChange(val);
//   }
// }
