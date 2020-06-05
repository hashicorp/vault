import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import notes from './toggle.md';
import { withKnobs, text, boolean, select } from '@storybook/addon-knobs';

storiesOf('Toggle/', module)
  .addParameters({ options: { showPanel: true } })
  .addDecorator(withKnobs())
  .add(
    `Toggle`,
    () => ({
      template: hbs`
      <h5 class="title is-5">Toggle</h5>
      <Toggle
        @name={{name}}
        @checked={{checked}}
        @onChange={{onChange}}
        @disabled={{disabled}}
        @size={{size}}
        @status={{status}}
        @round={{round}}
        data-test-secret-json-toggle
      >
        {{yielded}}
      </Toggle>
    `,
      context: {
        name: text('name', 'my-checkbox'),
        checked: boolean('checked', true),
        yielded: text('yield', 'Label content here ✔️'),
        onChange() {
          this.set('checked', !this.checked);
        },
        disabled: boolean('disabled', false),
        size: select('size', ['small', 'medium'], 'small'),
        status: select('status', ['normal', 'success'], 'success'),
      },
    }),
    { notes }
  );
