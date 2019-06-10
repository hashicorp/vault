/* eslint-disable import/extensions */
import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import notes from './layout-loading.md';

storiesOf('Loading/LayoutLoading/', module)
  .addParameters({ options: { showPanel: true } })
  .add(`LayoutLoading`, () => ({
    template: hbs`
        <h5 class="title is-5">Layout Loading</h5>
        <LayoutLoading/>
    `,
    context: {},
  }),
  {notes}
);
