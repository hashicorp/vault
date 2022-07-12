import Component from '@glimmer/component';

/**
 * @module OktaNumberChallenge
 * OktaNumberChallenge components are used to display loading screen and correct answer for Okta Number Challenge when signing in through Okta
 *
 * @example
 * ```js
 * <OktaNumberChallenge @correctAnswer={this.oktaNumberChallengeAnswer} @hasError={this.error} @onReturnToLogin={this.returnToLoginFromOktaNumberChallenge}/>
 * ```
 * @param {number} correctAnswer - The correct answer to click for the okta number challenge.
 * @param {boolean} hasError - Determines if there is an error being thrown.
 * @param {function} onReturnToLogin - Sets waitingForOktaNumberChallenge to false if want to return to main login.
 */

export default class OktaNumberChallenge extends Component {
  get oktaNumberChallengeCorrectAnswer() {
    return this.args.correctAnswer;
  }

  get errorThrown() {
    return this.args.hasError;
  }
}
