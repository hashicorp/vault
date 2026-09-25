/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

/* eslint-env node */
/* eslint-disable no-console */
/* eslint-disable no-process-exit */

/**
 * Records a human-paced Playwright demo video for a pull request.
 *
 *   node scripts/record-e2e-video.js [--name <file-name>] -- <playwright args>
 *
 * Runs the specs with PW_VIDEO=1, concatenates the per-test WebM clips into a single MP4
 * (GitHub rejects WebM uploads), and leaves exactly one file in ui/e2e/videos/. That
 * directory is intentionally NOT gitignored, so never `git add -A` while a video exists.
 */

const { spawnSync } = require('child_process');
const fs = require('fs');
const os = require('os');
const path = require('path');

const UI_ROOT = path.join(__dirname, '..');
const VIDEO_DIR = path.join(UI_ROOT, 'e2e', 'videos');
const TEST_RESULTS = path.join(UI_ROOT, 'test-results');

const EXIT_NO_FFMPEG = 2;

const fail = (message, code = 1) => {
  console.error(`\n[record-e2e-video] ${message}\n`);
  process.exit(code);
};

// --name is consumed here; everything else is forwarded to playwright untouched.
const parseArgs = (argv) => {
  const passthrough = [];
  let name = '';

  for (let i = 0; i < argv.length; i++) {
    if (argv[i] === '--name') {
      name = argv[++i] || '';
    } else if (argv[i] !== '--') {
      passthrough.push(argv[i]);
    }
  }
  return { name, passthrough };
};

const deriveName = (explicit, passthrough) => {
  if (explicit) {
    const name = explicit.replace(/\.mp4$/, '');
    // ffmpeg runs with -y, so `../../foo` would silently overwrite a file elsewhere.
    if (!name || name !== path.basename(name) || name.startsWith('.')) {
      fail(`--name must be a plain file name without path separators, got "${explicit}".`);
    }
    return name;
  }
  const spec = passthrough.find((arg) => arg.endsWith('.spec.ts'));
  return spec ? `${path.basename(spec, '.spec.ts')}-demo` : 'e2e-demo';
};

// Walks the JSON report so clips keep declaration order; directory names are derived from
// test titles, so the filesystem would return them alphabetically.
const collectVideos = (suite, found = []) => {
  for (const spec of suite.specs || []) {
    for (const test of spec.tests || []) {
      for (const result of test.results || []) {
        // CI retries keep a video per attempt; only the passing one belongs in the demo.
        if (result.status !== 'passed') continue;
        for (const attachment of result.attachments || []) {
          if (attachment.name === 'video' && attachment.path) {
            found.push({ project: test.projectName || '', path: attachment.path, title: spec.title });
          }
        }
      }
    }
  }
  for (const child of suite.suites || []) collectVideos(child, found);
  return found;
};

const { name: nameArg, passthrough } = parseArgs(process.argv.slice(2));

if (!passthrough.length) {
  fail(
    'No Playwright arguments given.\nUsage: node scripts/record-e2e-video.js [--name <file>] -- <playwright args>'
  );
}

if (spawnSync('ffmpeg', ['-version'], { stdio: 'ignore' }).status !== 0) {
  fail(
    'ffmpeg is required to combine the clips into an MP4, and it is not installed.\n' +
      'Ask the user for permission before installing it, then run:  brew install ffmpeg',
    EXIT_NO_FFMPEG
  );
}

// Resolved before the run so a bad --name fails immediately, not after a full suite.
const outputFile = path.join(VIDEO_DIR, `${deriveName(nameArg, passthrough)}.mp4`);

const playback = Number(process.env.PW_VIDEO_PLAYBACK ?? 1.5);
if (!Number.isFinite(playback) || playback < 1 || playback > 4) {
  fail(`PW_VIDEO_PLAYBACK must be a number between 1 and 4, got "${process.env.PW_VIDEO_PLAYBACK}".`);
}

const reportFile = path.join(os.tmpdir(), `pw-video-report-${Date.now()}.json`);

console.log(
  '[record-e2e-video] Recording — actions are slowed deliberately, so this takes longer than a normal run.\n'
);

const run = spawnSync('npx', ['playwright', 'test', ...passthrough, '--reporter=list,json'], {
  cwd: UI_ROOT,
  stdio: 'inherit',
  env: { ...process.env, PW_VIDEO: '1', PLAYWRIGHT_JSON_OUTPUT_NAME: reportFile },
});

if (run.status !== 0) {
  fail(
    'Tests failed, so no video was produced. Fix the failures and re-run — a demo video must show a passing run.'
  );
}

if (!fs.existsSync(reportFile)) {
  fail('Playwright did not write a JSON report, so clip order cannot be determined.');
}

const report = JSON.parse(fs.readFileSync(reportFile, 'utf-8'));
fs.rmSync(reportFile, { force: true });

const videos = (report.suites || []).flatMap((suite) => collectVideos(suite));

if (!videos.length) {
  fail('No video was recorded. Check that the run targeted a `chrome:` project.');
}

// setup projects enter credentials. The config already refuses to record them; abort
// loudly if that ever regresses.
const leaked = videos.filter((video) => !video.project.startsWith('chrome:'));
if (leaked.length) {
  // Destroy the clips first, or the guard leaves the very file it rejected on disk.
  fs.rmSync(TEST_RESULTS, { recursive: true, force: true });
  fail(
    `Refusing to continue: a non-browser project was recorded (${leaked
      .map((video) => video.project)
      .join(', ')}). Setup projects enter credentials and must never be filmed. ` +
      `${TEST_RESULTS} has been deleted.`
  );
}

fs.rmSync(VIDEO_DIR, { recursive: true, force: true });
fs.mkdirSync(VIDEO_DIR, { recursive: true });

const listFile = path.join(os.tmpdir(), `pw-video-list-${Date.now()}.txt`);
fs.writeFileSync(listFile, videos.map((video) => `file '${video.path}'`).join('\n'));

// Slows the finished recording. Elapsed time stays proportional, unlike slowMo padding.
const ffmpeg = spawnSync(
  'ffmpeg',
  [
    ...['-y', '-f', 'concat', '-safe', '0', '-i', listFile],
    ...['-c:v', 'libx264', '-preset', 'slow', '-crf', '26', '-pix_fmt', 'yuv420p'],
    // setpts before fps: retime first, then resample to a smooth 30fps
    ...['-vf', `setpts=${playback}*PTS,fps=30`],
    // -an: there is no audio, and an empty track upsets some players
    ...['-movflags', '+faststart', '-an', outputFile],
  ],
  { stdio: ['ignore', 'ignore', 'pipe'] }
);

fs.rmSync(listFile, { force: true });

if (ffmpeg.status !== 0) {
  fail(`ffmpeg failed to combine the clips:\n${ffmpeg.stderr?.toString() || 'no stderr'}`);
}

// Leave only the MP4 behind.
fs.rmSync(TEST_RESULTS, { recursive: true, force: true });

console.log(`\n[record-e2e-video] Combined ${videos.length} clip(s):`);
videos.forEach((video) => console.log(`  - ${video.title}`));
console.log(`\n[record-e2e-video] Video saved to:\n  ${outputFile}\n`);
console.log('[record-e2e-video] This file is NOT gitignored so you can see it in `git status`.');
console.log('[record-e2e-video] Attach it to the PR by dragging it into the description on github.com');
console.log('[record-e2e-video] (`gh` cannot upload media), then delete it yourself.\n');
