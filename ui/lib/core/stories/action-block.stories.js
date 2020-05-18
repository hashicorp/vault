import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import notes from './action-block.md';

storiesOf('ActionBlock', module)
  .addParameters({ options: { showPanel: true } })
  .add(
    `ActionBlock`,
    () => ({
      template: hbs`
      <h5 class="title is-5">Action Block</h5>
      <ActionBlock>
        <h1>Any of your own content goes in here!</h1>
      </ActionBlock>
    `,
      context: {},
    }),
    { notes }
  );
