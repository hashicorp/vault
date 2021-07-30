import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import { withKnobs } from '@storybook/addon-knobs';

storiesOf('ReadMore', module)
  .addParameters({ options: { showPanel: true } })
  .addDecorator(withKnobs())
  .add(`ReadMore`, () => ({
    template: hbs`
      <h5 class="title is-5">Read More</h5>
      <h6 class="has-text-grey">NOTE: normally this component has a "See More" button when it truncates text (which is calculated in the .js file). The button is clickable and will expand the height of the outer box (shown here with a black outline).</h6>
      <div style="width: 500px; border: 1px solid black;">
        <ReadMore>
          <strong>Anything can go in here</strong> Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nam volutpat vulputate lacus sit amet lobortis. Nulla fermentum porta consequat. Mauris porttitor libero nibh, ac facilisis ex molestie non. Nulla dolor est, pharetra et maximus vel, varius eu augue. Maecenas eget nisl convallis, vehicula massa quis, pharetra justo.
        </ReadMore>
      </div>
    `,
    context: {},
  }));
