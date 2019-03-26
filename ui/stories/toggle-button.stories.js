/* eslint-disable import/extensions */
import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import notes from './toggle-button.md';

storiesOf('ToggleButton', module)
  .addParameters({ options: { showPanel: false } })
  .add(
    'ToggleButton',
    () => ({
      template: hbs`
        <ToggleButton
          @toggleAttr="showOptions"
          @toggleTarget={{this}}
          />
        `,
    }),
    { notes }
  )
  .add(
    'ToggleButton with content',
    () => ({
      template: hbs`
        <ToggleButton
          @openLabel="Hide me!"
          @closedLabel="Show me!"
          @toggleTarget={{this}}
          @toggleAttr="showOptions"
        />

        {{#if showOptions}}
          <div>
            <p>
              I will be toggled!
            </p>
          </div>
        {{/if}}
        `,
    }),
    { notes }
  );
