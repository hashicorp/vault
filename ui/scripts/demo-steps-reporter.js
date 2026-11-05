/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

/* eslint-env node */

/**
 * Playwright reporter used only by scripts/record-e2e-video.js. For each test it writes when the page fixture
 * was set up (a recording starts when its page is created) and when each top-level `test.step` began, so the
 * recorder can place step titles on the demo video's timeline. The JSON reporter has no step start times.
 *
 * Output: the JSON file named by PW_VIDEO_STEPS_FILE.
 */

const fs = require('fs');
const path = require('path');

class DemoStepsReporter {
  constructor() {
    this.tests = new Map();
  }

  // Tests are matched with the JSON report by project, file name, line and column.
  key(test, result) {
    const { file, line, column } = test.location;
    return `${test.parent.project()?.name || ''}|${path.basename(file)}|${line}|${column}|${result.retry}`;
  }

  entry(test, result) {
    const key = this.key(test, result);
    if (!this.tests.has(key)) this.tests.set(key, { key, pageSetup: [], steps: [] });
    return this.tests.get(key);
  }

  onStepBegin(test, result, step) {
    const entry = this.entry(test, result);
    if (step.category === 'fixture' && step.title === 'Fixture "page"') {
      entry.pageSetup.push({ start: step.startTime.getTime() });
    } else if (step.category === 'test.step' && !step.parent) {
      entry.steps.push({ title: step.title, start: step.startTime.getTime() });
    }
  }

  onStepEnd(test, result, step) {
    if (step.category === 'fixture' && step.title === 'Fixture "page"') {
      const setup = this.entry(test, result).pageSetup.find((s) => s.start === step.startTime.getTime());
      if (setup) setup.end = setup.start + step.duration;
    }
  }

  onEnd() {
    const file = process.env.PW_VIDEO_STEPS_FILE;
    if (file) fs.writeFileSync(file, `${JSON.stringify([...this.tests.values()], null, 2)}\n`);
  }

  printsToStdio() {
    return false;
  }
}

module.exports = DemoStepsReporter;
