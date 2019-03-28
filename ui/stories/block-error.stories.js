/* eslint-disable import/extensions */
import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
//import notes from './block-error.md';

storiesOf('BlockError/', module)
  .addParameters({ options: { showPanel: true } })
  .add(
    `BlockError`,
    () => ({
      template: hbs`
        <h5 class="title is-5">Block Error</h5>
        <BlockError/>
    `,
      context: {},
    })
    //{ notes }
  );
