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
 *   node scripts/record-e2e-video.js [--name <file-name>] [--title "<intro title>"] [--no-cards] -- <playwright args>
 *
 * Runs the specs with PW_VIDEO=1, joins the per-test WebM clips into a single MP4 (GitHub rejects
 * WebM uploads), and writes it to ui/e2e/videos/<name>.mp4, replacing an earlier video with that name.
 * Other videos in the folder are left alone. The folder is intentionally NOT gitignored, so never
 * `git add -A` while a video exists.
 *
 * When the video has more than one test, or a --title, it gets silent cards: an intro card from
 * --title, a card named after each test before its clip, and a short fade at the end. The blank
 * white page each clip opens on is cut. --no-cards keeps the plain joined clips.
 *
 * The card-free recording, each clip's start and end, and the start of each top-level test.step are kept
 * outside the repo, in $TMPDIR/vault-ui-recordings/<name>-<path hash>/, with the video's SHA-256. A voice-over
 * can then narrate the plain recording, place its own cards at the same boundaries, and find it again after
 * the video is moved. Folders older than a week are removed at the start of each run.
 */

const { spawnSync } = require('child_process');
const { createHash } = require('crypto');
const fs = require('fs');
const os = require('os');
const path = require('path');
const { CARD_COLOR, cardTitle, renderCards } = require('./video-cards');

const UI_ROOT = path.join(__dirname, '..');
const VIDEO_DIR = path.join(UI_ROOT, 'e2e', 'videos');
const TEST_RESULTS = path.join(UI_ROOT, 'test-results');
const RECORDINGS = path.join(os.tmpdir(), 'vault-ui-recordings');
const STEPS_REPORTER = path.join(__dirname, 'demo-steps-reporter.js');
const KEEP_RECORDINGS_DAYS = 7;
// Names other tools can use as file and folder names without quoting: narrate.mjs accepts the same set.
const NAME_PATTERN = /^[\w.-]+$/;

// Silent cards are short: long enough to read a title, never long enough to feel like a pause.
const CARD_SECONDS = 1.8;
const INTRO_SECONDS = { base: 3, perItem: 0.4, max: 5 };
const FADE_SECONDS = 0.3;
const END_FADE_SECONDS = 1;
// A blank page frame is near-white everywhere; YLOW is the darkest luma in the frame.
const WHITE_YLOW = 230;
const FRAME = 1 / 30;
// Frames of a mid-test page load that are kept, so a navigation still reads as one without a white flash.
const KEEP_WHITE_FRAMES = 2;

const EXIT_NO_FFMPEG = 2;

const fail = (message, code = 1) => {
  console.error(`\n[record-e2e-video] ${message}\n`);
  process.exit(code);
};

// --name, --title and --no-cards are consumed here, before the `--`; everything else is forwarded to playwright.
const parseArgs = (argv) => {
  const passthrough = [];
  let name = '';
  let title = '';
  let cards = true;
  let forwarding = false;

  for (let i = 0; i < argv.length; i++) {
    if (argv[i] === '--') {
      forwarding = true;
    } else if (!forwarding && argv[i] === '--name') {
      name = argv[++i] || '';
    } else if (!forwarding && argv[i] === '--title') {
      title = (argv[++i] || '').trim();
    } else if (!forwarding && argv[i] === '--no-cards') {
      cards = false;
    } else {
      passthrough.push(argv[i]);
    }
  }
  return { name, title, cards, passthrough };
};

const deriveName = (explicit, passthrough) => {
  const spec = passthrough.find((arg) => arg.endsWith('.spec.ts'));
  const name = explicit
    ? explicit.replace(/\.mp4$/, '')
    : spec
    ? `${path.basename(spec, '.spec.ts')}-demo`
    : 'e2e-demo';
  // ffmpeg runs with -y, so `../../foo` would silently overwrite a file elsewhere, and a quote would break
  // the clip list. Letters, digits, `.`, `_` and `-` only, and no leading dot.
  if (!NAME_PATTERN.test(name) || name.startsWith('.')) {
    fail(`--name may only use letters, digits, ".", "_" and "-", and can't start with ".". Got "${name}".`);
  }
  return name;
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
            found.push({
              project: test.projectName || '',
              path: attachment.path,
              title: spec.title,
              file: spec.file || '',
              // Matches the key written by demo-steps-reporter.js.
              key: `${test.projectName || ''}|${path.basename(spec.file || '')}|${spec.line}|${spec.column}|${
                result.retry
              }`,
            });
          }
        }
      }
    }
  }
  for (const child of suite.suites || []) collectVideos(child, found);
  return found;
};

const {
  name: nameArg,
  title: introTitle,
  cards: cardsWanted,
  passthrough,
} = parseArgs(process.argv.slice(2));

if (!passthrough.length) {
  fail(
    'No Playwright arguments given.\n' +
      'Usage: node scripts/record-e2e-video.js [--name <file>] [--title "<intro title>"] [--no-cards] -- <playwright args>'
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

// Recordings kept for narration are only useful for a few days; remove the old ones so $TMPDIR doesn't fill up.
const pruneRecordings = () => {
  if (!fs.existsSync(RECORDINGS)) return;
  const cutoff = Date.now() - KEEP_RECORDINGS_DAYS * 24 * 60 * 60 * 1000;
  for (const entry of fs.readdirSync(RECORDINGS, { withFileTypes: true })) {
    const dir = path.join(RECORDINGS, entry.name);
    if (entry.isDirectory() && fs.statSync(dir).mtimeMs < cutoff)
      fs.rmSync(dir, { recursive: true, force: true });
  }
};
pruneRecordings();

const reportFile = path.join(os.tmpdir(), `pw-video-report-${Date.now()}.json`);
const stepsFile = path.join(os.tmpdir(), `pw-video-steps-${Date.now()}.json`);

console.log(
  '[record-e2e-video] Recording — actions are slowed deliberately, so this takes longer than a normal run.\n'
);

const run = spawnSync(
  'npx',
  ['playwright', 'test', ...passthrough, `--reporter=list,json,${STEPS_REPORTER}`],
  {
    cwd: UI_ROOT,
    stdio: 'inherit',
    env: {
      ...process.env,
      PW_VIDEO: '1',
      PLAYWRIGHT_JSON_OUTPUT_NAME: reportFile,
      PW_VIDEO_STEPS_FILE: stepsFile,
    },
  }
);

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
// Step timing is optional: without it, a voice-over splits scenes at page loads only.
const stepTimes = new Map(
  (fs.existsSync(stepsFile) ? JSON.parse(fs.readFileSync(stepsFile, 'utf-8')) : []).map((t) => [t.key, t])
);
fs.rmSync(stepsFile, { force: true });

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

fs.mkdirSync(VIDEO_DIR, { recursive: true });
fs.rmSync(outputFile, { force: true });

const ffmpeg = (args, what) => {
  const result = spawnSync('ffmpeg', ['-v', 'error', '-y', ...args], {
    stdio: ['ignore', 'pipe', 'pipe'],
    maxBuffer: 1 << 28,
  });
  if (result.status !== 0) fail(`ffmpeg failed to ${what}:\n${result.stderr?.toString() || 'no stderr'}`);
  return result.stdout.toString();
};
const probe = (file, entries) =>
  spawnSync('ffprobe', ['-v', 'error', '-show_entries', entries, '-of', 'csv=p=0', file])
    .stdout.toString()
    .trim();

// Runs of blank white frames: the page each test opens on, and any page load during the test.
const whiteRuns = (file) => {
  const stats = ffmpeg(
    [
      '-i',
      file,
      '-an',
      '-vf',
      'scale=64:40,signalstats,metadata=print:key=lavfi.signalstats.YLOW:file=-',
    ].concat(['-f', 'null', '-']),
    'find blank page frames'
  );
  const runs = [];
  for (const [, at, ylow] of stats.matchAll(/pts_time:([\d.]+)\s+lavfi\.signalstats\.YLOW=(\d+)/g)) {
    if (+ylow < WHITE_YLOW) continue;
    const last = runs[runs.length - 1];
    if (last && +at - last.end < FRAME * 1.5) last.end = +at;
    else runs.push({ start: +at, end: +at });
  }
  return runs;
};

// The parts of a clip to keep: everything except the blank page it opens on (a card already marks the change)
// and all but the first frames of each later page load. Cuts sit half a frame between frames.
const keptParts = (clip) => {
  const half = FRAME / 2;
  const parts = [];
  let cursor = 0;
  for (const run of whiteRuns(clip.file)) {
    if (run.start < half) {
      cursor = run.end + half;
      continue;
    }
    const keepUntil = run.start + (KEEP_WHITE_FRAMES - 0.5) * FRAME;
    if (keepUntil >= run.end + half) continue;
    parts.push([cursor, keepUntil]);
    cursor = run.end + half;
  }
  if (clip.seconds - cursor >= FRAME) parts.push([cursor, clip.seconds]);
  const kept = parts.filter(([from, to]) => to - from > half);
  // A test that never paints anything but a blank page is shown as it is, rather than dropped.
  return kept.length ? kept : [[0, clip.seconds]];
};

// Seconds from the start of a test's retimed clip to the start of each of its top-level test.steps. Recording
// starts when the page fixture is set up, so step times are measured from there.
const stepsIn = (video, seconds) => {
  const timing = stepTimes.get(video.key);
  const anchor = timing?.pageSetup?.[0]?.start;
  if (!anchor) return [];
  return timing.steps
    .map((step) => ({ title: step.title, offset: ((step.start - anchor) / 1000) * playback }))
    .filter((step) => step.offset >= 0 && step.offset < seconds);
};

const combine = async () => {
  // One folder per output path: <name>-<first 8 hex of the path's SHA-256>, a key other tools can derive.
  const work = path.join(
    RECORDINGS,
    `${path.basename(outputFile, '.mp4')}-${createHash('sha256')
      .update(outputFile)
      .digest('hex')
      .slice(0, 8)}`
  );
  fs.rmSync(work, { recursive: true, force: true });
  fs.mkdirSync(work, { recursive: true });

  // Retime each clip on its own, so every test's start and end in the joined video is known exactly.
  // Slowing the finished recording keeps elapsed time proportional, unlike slowMo padding.
  const clips = videos.map((video, index) => {
    const file = path.join(work, `clip-${index}.mp4`);
    ffmpeg(
      // setpts before fps: retime first, then resample to a smooth 30fps
      ['-i', video.path, '-an', '-vf', `setpts=${playback}*PTS,fps=30`].concat([
        '-c:v',
        'libx264',
        '-preset',
        'veryfast',
        '-crf',
        '16',
        '-pix_fmt',
        'yuv420p',
        file,
      ]),
      'retime the clips'
    );
    const seconds = Number(probe(file, 'format=duration'));
    return { title: video.title, file, seconds, steps: stepsIn(video, seconds) };
  });
  const [width, height] = probe(clips[0].file, 'stream=width,height').split(',').map(Number);

  // The card-free recording, kept for narration.
  const list = path.join(work, 'clips.txt');
  // The concat list quotes each path; a quote inside a path is closed, escaped and reopened.
  fs.writeFileSync(list, clips.map((clip) => `file '${clip.file.replace(/'/g, "'\\''")}'`).join('\n'));
  const raw = path.join(work, 'raw.mp4');
  ffmpeg(['-f', 'concat', '-safe', '0', '-i', list, '-c', 'copy', raw], 'join the clips');
  let at = 0;
  const timeline = clips.map((clip) => {
    const entry = {
      title: clip.title,
      start: +at.toFixed(3),
      end: +(at + clip.seconds).toFixed(3),
      steps: clip.steps.map((step) => ({ title: step.title, start: +(at + step.offset).toFixed(3) })),
    };
    at += clip.seconds;
    return entry;
  });

  const encode = ['-c:v', 'libx264', '-preset', 'slow', '-crf', '26', '-pix_fmt', 'yuv420p'];
  // -an: there is no audio, and an empty track upsets some players
  const finish = ['-movflags', '+faststart', '-an', outputFile];
  const withCards = cardsWanted && (clips.length > 1 || Boolean(introTitle));
  if (!withCards) {
    ffmpeg(['-i', raw, ...encode, ...finish], 'encode the video');
  } else {
    const cards = [];
    const several = clips.length > 1;
    if (introTitle) {
      const items = several ? clips.map((clip) => cardTitle(clip.title)) : [];
      cards.push({
        before: 0,
        eyebrow: 'What this demo covers',
        title: introTitle,
        subtitle: several ? '' : cardTitle(clips[0].title),
        items,
        seconds: Math.min(
          INTRO_SECONDS.max,
          INTRO_SECONDS.base + INTRO_SECONDS.perItem * Math.max(0, items.length - 3)
        ),
      });
    }
    if (several) {
      clips.forEach((clip, index) =>
        cards.push({
          before: index,
          eyebrow: `Scenario ${index + 1} of ${clips.length}`,
          title: cardTitle(clip.title),
          seconds: CARD_SECONDS,
        })
      );
    }
    cards.forEach((card, index) => (card.file = path.join(work, `card-${index}.png`)));
    await renderCards(cards, { width, height });

    const inputs = [];
    const graph = [];
    const segments = [];
    let total = 0;
    const nextInput = () => inputs.filter((arg) => arg === '-i').length - 1;
    clips.forEach((clip, index) => {
      for (const card of cards.filter((c) => c.before === index)) {
        inputs.push('-loop', '1', '-framerate', '30', '-t', card.seconds.toFixed(3), '-i', card.file);
        const fadeOut = (card.seconds - FADE_SECONDS).toFixed(3);
        graph.push(
          `[${nextInput()}:v]format=yuv420p,fade=t=in:st=0:d=${FADE_SECONDS}:color=${CARD_COLOR},` +
            `fade=t=out:st=${fadeOut}:d=${FADE_SECONDS}:color=${CARD_COLOR},setsar=1[p${segments.length}]`
        );
        segments.push(`[p${segments.length}]`);
        total += card.seconds;
      }
      const parts = keptParts(clip);
      inputs.push('-i', clip.file);
      const k = nextInput();
      const labels = parts.map((_, n) => `[c${k}_${n}]`);
      graph.push(`[${k}:v]split=${parts.length}${labels.join('')}`);
      parts.forEach(([from, to], n) => {
        graph.push(
          `${labels[n]}trim=start=${from.toFixed(4)}:end=${to.toFixed(4)},setpts=PTS-STARTPTS,fps=30,` +
            `format=yuv420p,setsar=1[p${segments.length}]`
        );
        segments.push(`[p${segments.length}]`);
        total += to - from;
      });
    });
    const fadeAt = Math.max(0, total - END_FADE_SECONDS).toFixed(3);
    graph.push(
      `${segments.join('')}concat=n=${segments.length}:v=1:a=0,` +
        `fade=t=out:st=${fadeAt}:d=${END_FADE_SECONDS}:color=${CARD_COLOR}[out]`
    );
    ffmpeg(
      [...inputs, '-filter_complex', graph.join(';'), '-map', '[out]', ...encode, ...finish],
      'add the cards'
    );
  }

  const { size, mtimeMs } = fs.statSync(outputFile);
  // The content hash lets a voice-over find this recording after the video is moved or renamed.
  const sha256 = createHash('sha256').update(fs.readFileSync(outputFile)).digest('hex');
  fs.writeFileSync(
    path.join(work, 'clips.json'),
    `${JSON.stringify(
      {
        output: outputFile,
        size,
        mtimeMs,
        sha256,
        playback,
        title: introTitle,
        cards: withCards,
        clips: timeline,
      },
      null,
      2
    )}\n`
  );
  for (const file of fs.readdirSync(work)) {
    if (!['raw.mp4', 'clips.json'].includes(file)) fs.rmSync(path.join(work, file), { force: true });
  }

  // Leave only the MP4 behind.
  fs.rmSync(TEST_RESULTS, { recursive: true, force: true });

  console.log(`\n[record-e2e-video] Combined ${videos.length} clip(s):`);
  videos.forEach((video) => console.log(`  - ${video.title}`));
  if (withCards) {
    console.log(
      `[record-e2e-video] Added ${introTitle ? 'an intro card, ' : ''}${
        clips.length > 1 ? `${clips.length} scenario cards, ` : ''
      }` + 'and a fade at the end. Use --no-cards for the plain clips.'
    );
  }
  console.log(`\n[record-e2e-video] Video saved to:\n  ${outputFile}\n`);
  const others = fs
    .readdirSync(VIDEO_DIR)
    .filter((file) => file.endsWith('.mp4') && file !== path.basename(outputFile));
  if (others.length)
    console.log(`[record-e2e-video] Other videos left in ${VIDEO_DIR}: ${others.join(', ')}`);
  console.log('[record-e2e-video] This file is NOT gitignored so you can see it in `git status`.');
  console.log('[record-e2e-video] Attach it to the PR by dragging it into the description on github.com');
  console.log('[record-e2e-video] (`gh` cannot upload media), then delete it yourself.\n');
};

// A spec without the demo fixture records with no cursor and no pause before clicks, so viewers can't
// follow it. Warn rather than fail: the video still shows a passing run.
const warnNotDemoReady = () => {
  const rootDir = report.config?.rootDir || path.join(UI_ROOT, 'e2e');
  const DEMO_IMPORT = /from\s+['"](?:\.\.?\/)+(?:fixtures\/demo|mocks\/fixture)['"]/;
  const notDemoReady = [...new Set(videos.map((video) => video.file).filter(Boolean))].filter((file) => {
    const spec = path.resolve(rootDir, file);
    return fs.existsSync(spec) && !DEMO_IMPORT.test(fs.readFileSync(spec, 'utf-8'));
  });
  if (notDemoReady.length) {
    console.warn(
      '[record-e2e-video] Warning: these specs are not demo-ready, so the video may be too fast to follow:'
    );
    notDemoReady.forEach((file) =>
      console.warn(`  - ${path.relative(UI_ROOT, path.resolve(rootDir, file))}`)
    );
    console.warn('[record-e2e-video] They import `test` from @playwright/test instead of e2e/fixtures/demo,');
    console.warn('[record-e2e-video] so there is no cursor or pause before clicks. Record a demo-ready spec');
    console.warn('[record-e2e-video] instead (playwright skill, references/video.md § Demo-ready specs).\n');
  }
};

combine()
  .then(warnNotDemoReady)
  .catch((error) => fail(error.stack || String(error)));
