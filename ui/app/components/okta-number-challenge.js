import Component from '@glimmer/component';

/**
 * @module OktaNumberChallenge
 * OktaNumberChallenge components are used to display loading screen and correct answer for Okta Number Challenge when signing in through Okta
 *
 * @example
 * ```js
 * <OktaNumberChallenge @correctAnswer={correctAnswer}/>
 * ```
 * @param {number} correctAnswer - The correct answer to click for the okta number challenge.
 */

export default class OktaNumberChallenge extends Component {
  get oktaNumberChallengeCorrectAnswer() {
    return this.args.correctAnswer;
  }
}
