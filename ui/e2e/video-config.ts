/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

/**
 * Opt-in demo video recording for the Playwright suite.
 *
 * Off by default so normal runs stay fast and produce nothing.
 * Set PW_VIDEO to record a human-paced walkthrough, typically to attach to a pull request:
 *   pnpm test:e2e:video -- <playwright args>
 *
 *   PW_VIDEO           any value except "0"/"false" enables recording
 *   PW_VIDEO_SLOW_MO   ms paused between actions (default 250)
 *   PW_VIDEO_GLIDE_MS  ms of pointer travel before each click (default 320)
 *   PW_VIDEO_PLAYBACK  playback slowdown factor applied by ffmpeg (default 1.5)
 */

const flag = process.env.PW_VIDEO ?? '';

export const isVideoEnabled = flag !== '' && flag !== '0' && flag !== 'false';

// Fail loudly on a bad value; NaN would otherwise reach waitForTimeout and the CSS
// transition, where it degrades silently into no animation at all.
const readMs = (name: string, fallback: number) => {
  const raw = process.env[name];
  if (raw === undefined || raw === '') return fallback;

  const value = Number(raw);
  if (!Number.isFinite(value) || value < 0 || value > 10_000) {
    throw new Error(`${name} must be a number between 0 and 10000, got "${raw}".`);
  }
  return value;
};

// Pauses before each Playwright call, so it spaces out actions without slowing page loads.
const slowMo = readMs('PW_VIDEO_SLOW_MO', 250);

const viewport = { width: 1280, height: 800 };

/**
 * `click()` dispatches mousemove and mousedown 1ms apart, so the pointer can never reach a
 * target on its own. The demo fixture hovers first and waits out `glideMs`, which is what
 * lets the pointer travel visibly while the ripple still fires the moment the click lands.
 */
export const cursorTimings = {
  glideMs: readMs('PW_VIDEO_GLIDE_MS', 320),
  rippleMs: 650,
};

export const videoTimeout = 180_000;

/**
 * Recording options for a project, or `{}` when recording is off.
 *
 * Only `chrome:` projects are recorded. `setup:` projects enter credentials (see
 * e2e/init.setup.ts), which must never end up in a file bound for a pull request.
 */
export const videoOptionsFor = (projectName: string) => {
  if (!isVideoEnabled || !projectName.startsWith('chrome:')) {
    return {};
  }

  return {
    viewport,
    video: { mode: 'on' as const, size: viewport },
    launchOptions: { slowMo },
  };
};

/**
 * Draws a synthetic pointer into the page, since browsers never paint the real cursor into a
 * recording. Tracks Playwright's mousemove/mousedown events: the dot follows the pointer and
 * a ripple marks every click.
 *
 * Styling goes through CSSOM and the Web Animations API rather than a <style> tag, because
 * Vault serves `style-src 'self'` with no `'unsafe-inline'` (see vault/ui.go) — an injected
 * stylesheet is silently discarded, leaving the element in the DOM with a null `.sheet`.
 *
 * Injected via addInitScript so it survives navigations.
 */
export const cursorInitScript = (timings: { glideMs: number; rippleMs: number }) => {
  const POSITION_KEY = '__demo_cursor_pos';
  const { glideMs, rippleMs } = timings;

  // HDS dropdowns use the popover API and modals use <dialog>, both of which render in the
  // top layer. No z-index can paint above that, so the overlay has to join it too.
  const supportsTopLayer = typeof HTMLElement !== 'undefined' && 'popover' in HTMLElement.prototype;

  const joinTopLayer = (el: HTMLElement) => {
    if (!supportsTopLayer) return;
    try {
      el.setAttribute('popover', 'manual');
      (el as unknown as { showPopover: () => void }).showPopover();
    } catch {
      // the z-index below still covers ordinary content
    }
  };

  // The top layer stacks in promotion order, so anything opened after the cursor covers it.
  const raiseToTop = (el: HTMLElement) => {
    if (!supportsTopLayer || !el.isConnected) return;
    try {
      const popover = el as unknown as { hidePopover: () => void; showPopover: () => void };
      popover.hidePopover();
      popover.showPopover();
      // Re-showing resets style resolution; read layout so the next move still animates.
      void el.offsetWidth;
    } catch {
      // keeps its current stacking position
    }
  };

  const install = () => {
    if (document.getElementById('__demo_cursor')) return;

    const glide = `left ${glideMs}ms ease-in-out, top ${glideMs}ms ease-in-out, opacity 150ms ease`;

    const dot = document.createElement('div');
    dot.id = '__demo_cursor';
    Object.assign(dot.style, {
      // `inset` is a shorthand for top/right/bottom/left and must precede them, or it
      // clears the position. It is here at all to undo the UA styles that come with
      // [popover].
      inset: 'auto',
      padding: '0',
      overflow: 'visible',
      position: 'fixed',
      left: '0px',
      top: '0px',
      width: '20px',
      height: '20px',
      marginLeft: '-10px',
      marginTop: '-10px',
      borderRadius: '50%',
      background: 'rgba(28,125,255,0.35)',
      border: '2px solid #1c7dff',
      boxShadow: '0 0 0 2px rgba(255,255,255,0.9)',
      pointerEvents: 'none',
      zIndex: '2147483647',
      opacity: '0',
      transition: glide,
    });
    document.body.appendChild(dot);
    joinTopLayer(dot);

    let x = 0;
    let y = 0;

    // Restore across navigations so the pointer does not blink out on a full page load.
    try {
      const saved = JSON.parse(sessionStorage.getItem(POSITION_KEY) || 'null');
      if (saved) {
        x = saved.x;
        y = saved.y;
        dot.style.left = `${x}px`;
        dot.style.top = `${y}px`;
        dot.style.opacity = '1';
      }
    } catch {
      // sessionStorage can be unavailable; the cursor simply starts hidden
    }

    window.addEventListener(
      'mousemove',
      (event) => {
        x = event.clientX;
        y = event.clientY;
        raiseToTop(dot);
        dot.style.opacity = '1';
        dot.style.left = `${x}px`;
        dot.style.top = `${y}px`;
        try {
          sessionStorage.setItem(POSITION_KEY, JSON.stringify({ x, y }));
        } catch {
          // ignore
        }
      },
      true
    );

    window.addEventListener(
      'mousedown',
      () => {
        raiseToTop(dot);

        const ripple = document.createElement('div');
        Object.assign(ripple.style, {
          // see the pointer above — `inset` must precede left/top
          inset: 'auto',
          padding: '0',
          overflow: 'visible',
          position: 'fixed',
          left: `${x}px`,
          top: `${y}px`,
          // Wider than the 20px pointer, which is parked on the target when this fires; a
          // smaller ring would spend its brightest frames hidden underneath it.
          width: '28px',
          height: '28px',
          marginLeft: '-14px',
          marginTop: '-14px',
          borderRadius: '50%',
          border: '3px solid #1c7dff',
          pointerEvents: 'none',
          zIndex: '2147483647',
          opacity: '0',
        });
        document.body.appendChild(ripple);
        // Promoted after the pointer so it draws on top; the ring is hollow, so the pointer
        // still reads through the middle.
        joinTopLayer(ripple);

        const animation = ripple.animate(
          [
            { transform: 'scale(1)', opacity: 0.95 },
            { transform: 'scale(3)', opacity: 0 },
          ],
          { duration: rippleMs, easing: 'ease-out' }
        );
        animation.onfinish = () => ripple.remove();
      },
      true
    );
  };

  if (document.readyState === 'loading') {
    window.addEventListener('DOMContentLoaded', install);
  } else {
    install();
  }
};
