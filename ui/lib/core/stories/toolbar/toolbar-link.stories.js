import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import { withKnobs, select, text, boolean } from '@storybook/addon-knobs';
import notes from './toolbar-link.md';

storiesOf('Toolbar', module)
  .addParameters({ options: { showPanel: true } })
  .addDecorator(withKnobs())
  .add(
    `ToolbarLink`,
    () => ({
      template: hbs`
      <h5 class="title is-5">ToolbarLink</h5>
      <div style="width: 400px;">
        <Toolbar>
          <ToolbarActions>
            <ToolbarLink
              @params={{array '#'}}
              @type={{type}}
              @disabled={{disabled}}
              @disabledTooltip={{disabledTooltip}}
            >
              {{label}}
            </ToolbarLink>
          </ToolbarActions>
        </Toolbar>
      </div>
    `,
      context: {
        type: select('Type', ['', 'add']),
        label: text('Button text', 'Edit secret'),
        disabled: boolean('disabled', false),
        disabledTooltip: text('Tooltip to display when disabled', ''),
      },
    }),
    { notes }
  );
