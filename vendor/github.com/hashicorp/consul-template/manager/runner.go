package manager

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/hashicorp/consul-template/child"
	"github.com/hashicorp/consul-template/config"
	dep "github.com/hashicorp/consul-template/dependency"
	"github.com/hashicorp/consul-template/renderer"
	"github.com/hashicorp/consul-template/template"
	"github.com/hashicorp/consul-template/watch"
	multierror "github.com/hashicorp/go-multierror"
	shellwords "github.com/mattn/go-shellwords"
	"github.com/pkg/errors"
)

const (
	// saneViewLimit is the number of views that we consider "sane" before we
	// warn the user that they might be DDoSing their Consul cluster.
	saneViewLimit = 128
)

// Runner responsible rendering Templates and invoking Commands.
type Runner struct {
	// ErrCh and DoneCh are channels where errors and finish notifications occur.
	ErrCh  chan error
	DoneCh chan struct{}

	// config is the Config that created this Runner. It is used internally to
	// construct other objects and pass data.
	config *config.Config

	// signals sending output to STDOUT instead of to a file.
	dry bool

	// outStream and errStream are the io.Writer streams where the runner will
	// write information. These can be modified by calling SetOutStream and
	// SetErrStream accordingly.

	// inStream is the ioReader where the runner will read information.
	outStream, errStream io.Writer
	inStream             io.Reader

	// ctemplatesMap is a map of each template ID to the TemplateConfigs
	// that made it.
	ctemplatesMap map[string]config.TemplateConfigs

	// templates is the list of calculated templates.
	templates []*template.Template

	// renderEvents is a mapping of a template ID to the render event.
	renderEvents map[string]*RenderEvent

	// renderEventLock protects access into the renderEvents map
	renderEventsLock sync.RWMutex

	// renderedCh is used to signal that a template has been rendered
	renderedCh chan struct{}

	// renderEventCh is used to signal that there is a new render event. A
	// render event doesn't necessarily mean that a template has been rendered,
	// only that templates attempted to render and may have updated their
	// dependency sets.
	renderEventCh chan struct{}

	// dependencies is the list of dependencies this runner is watching.
	dependencies map[string]dep.Dependency

	// dependenciesLock is a lock around touching the dependencies map.
	dependenciesLock sync.Mutex

	// watcher is the watcher this runner is using.
	watcher *watch.Watcher

	// brain is the internal storage database of returned dependency data.
	brain *template.Brain

	// child is the child process under management. This may be nil if not running
	// in exec mode.
	child *child.Child

	// childLock is the internal lock around the child process.
	childLock sync.RWMutex

	// quiescenceMap is the map of templates to their quiescence timers.
	// quiescenceCh is the channel where templates report returns from quiescence
	// fires.
	quiescenceMap map[string]*quiescence
	quiescenceCh  chan *template.Template

	// dedup is the deduplication manager if enabled
	dedup *DedupManager

	// Env represents a custom set of environment variables to populate the
	// template and command runtime with. These environment variables will be
	// available in both the command's environment as well as the template's
	// environment.
	// NOTE this is only used when CT is being used as a library.
	Env map[string]string

	// stopLock is the lock around checking if the runner can be stopped
	stopLock sync.Mutex

	// stopped is a boolean of whether the runner is stopped
	stopped bool
}

// RenderEvent captures the time and events that occurred for a template
// rendering.
type RenderEvent struct {
	// Missing is the list of dependencies that we do not yet have data for, but
	// are contained in the watcher. This is different from unwatched dependencies,
	// which includes dependencies the watcher has not yet started querying for
	// data.
	MissingDeps *dep.Set

	// Template is the template attempting to be rendered.
	Template *template.Template

	// Contents is the raw, rendered contents from the template.
	Contents []byte

	// TemplateConfigs is the list of template configs that correspond to this
	// template.
	TemplateConfigs []*config.TemplateConfig

	// Unwatched is the list of dependencies that are not present in the watcher.
	// This value may change over time due to the n-pass evaluation.
	UnwatchedDeps *dep.Set

	// UpdatedAt is the last time this render event was updated.
	UpdatedAt time.Time

	// Used is the full list of dependencies seen in the template. Because of
	// the n-pass evaluation, this number can change over time. The dependencies
	// in this list may or may not have data. This just contains the list of all
	// dependencies parsed out of the template with the current data.
	UsedDeps *dep.Set

	// WouldRender determines if the template would have been rendered. A template
	// would have been rendered if all the dependencies are satisfied, but may
	// not have actually rendered if the file was already present or if an error
	// occurred when trying to write the file.
	WouldRender bool

	// LastWouldRender marks the last time the template would have rendered.
	LastWouldRender time.Time

	// DidRender determines if the Template was actually written to disk. In dry
	// mode, this will always be false, since templates are not written to disk
	// in dry mode. A template is only rendered to disk if all dependencies are
	// satisfied and the template is not already in place with the same contents.
	DidRender bool

	// LastDidRender marks the last time the template was written to disk.
	LastDidRender time.Time

	// ForQuiescence determines if this event is returned early in the
	// render loop due to quiescence. When evaluating if all templates have
	// been rendered we need to know if the event is triggered by quiesence
	// and if we can skip evaluating it as a render event for those purposes
	ForQuiescence bool
}

// NewRunner accepts a slice of TemplateConfigs and returns a pointer to the new
// Runner and any error that occurred during creation.
func NewRunner(config *config.Config, dry bool) (*Runner, error) {
	log.Printf("[INFO] (runner) creating new runner (dry: %v, once: %v)",
		dry, config.Once)

	runner := &Runner{
		config: config,
		dry:    dry,
	}

	if err := runner.init(); err != nil {
		return nil, err
	}

	return runner, nil
}

// Start begins the polling for this runner. Any errors that occur will cause
// this function to push an item onto the runner's error channel and the halt
// execution. This function is blocking and should be called as a goroutine.
func (r *Runner) Start() {
	log.Printf("[INFO] (runner) starting")

	// Create the pid before doing anything.
	if err := r.storePid(); err != nil {
		r.ErrCh <- err
		return
	}

	// Start the de-duplication manager
	var dedupCh <-chan struct{}
	if r.dedup != nil {
		if err := r.dedup.Start(); err != nil {
			r.ErrCh <- err
			return
		}
		dedupCh = r.dedup.UpdateCh()
	}

	// Setup the child process exit channel
	var childExitCh <-chan int

	// Fire an initial run to parse all the templates and setup the first-pass
	// dependencies. This also forces any templates that have no dependencies to
	// be rendered immediately (since they are already renderable).
	log.Printf("[DEBUG] (runner) running initial templates")
	if err := r.Run(); err != nil {
		r.ErrCh <- err
		return
	}

	for {
		// Warn the user if they are watching too many dependencies.
		if r.watcher.Size() > saneViewLimit {
			log.Printf("[WARN] (runner) watching %d dependencies - watching this "+
				"many dependencies could DDoS your servers", r.watcher.Size())
		} else {
			log.Printf("[DEBUG] (runner) watching %d dependencies", r.watcher.Size())
		}

		if r.allTemplatesRendered() {
			log.Printf("[DEBUG] (runner) all templates rendered")
			// Enable quiescence for all templates if we have specified wait
			// intervals.
		NEXT_Q:
			for _, t := range r.templates {
				if _, ok := r.quiescenceMap[t.ID()]; ok {
					continue NEXT_Q
				}

				for _, c := range r.templateConfigsFor(t) {
					if *c.Wait.Enabled {
						log.Printf("[DEBUG] (runner) enabling template-specific "+
							"quiescence for %q", t.ID())
						r.quiescenceMap[t.ID()] = newQuiescence(
							r.quiescenceCh, *c.Wait.Min, *c.Wait.Max, t)
						continue NEXT_Q
					}
				}

				if *r.config.Wait.Enabled {
					log.Printf("[DEBUG] (runner) enabling global quiescence for %q",
						t.ID())
					r.quiescenceMap[t.ID()] = newQuiescence(
						r.quiescenceCh, *r.config.Wait.Min, *r.config.Wait.Max, t)
					continue NEXT_Q
				}
			}

			// If an exec command was given and a command is not currently running,
			// spawn the child process for supervision.
			if config.StringPresent(r.config.Exec.Command) {
				// Lock the child because we are about to check if it exists.
				r.childLock.Lock()

				log.Printf("[TRACE] (runner) acquired child lock for command, spawning")

				if r.child == nil {
					env := r.config.Exec.Env.Copy()
					env.Custom = append(r.childEnv(), env.Custom...)
					child, err := spawnChild(&spawnChildInput{
						Stdin:        r.inStream,
						Stdout:       r.outStream,
						Stderr:       r.errStream,
						Command:      config.StringVal(r.config.Exec.Command),
						Env:          env.Env(),
						ReloadSignal: config.SignalVal(r.config.Exec.ReloadSignal),
						KillSignal:   config.SignalVal(r.config.Exec.KillSignal),
						KillTimeout:  config.TimeDurationVal(r.config.Exec.KillTimeout),
						Splay:        config.TimeDurationVal(r.config.Exec.Splay),
					})
					if err != nil {
						r.ErrCh <- err
						r.childLock.Unlock()
						return
					}
					r.child = child
				}

				// Unlock the child, we are done now.
				r.childLock.Unlock()

				// It's possible that we didn't start a process, in which case no
				// channel is returned. If we did get a new exitCh, that means a child
				// was spawned, so we need to watch a new exitCh. It is also possible
				// that during a run, the child process was restarted, which means a
				// new exit channel should be used.
				nexitCh := r.child.ExitCh()
				if nexitCh != nil {
					childExitCh = nexitCh
				}
			}

			// If we are running in once mode and all our templates are rendered,
			// then we should exit here.
			if r.config.Once {
				log.Printf("[INFO] (runner) once mode and all templates rendered")

				if r.child != nil {
					r.stopDedup()
					r.stopWatcher()

					log.Printf("[INFO] (runner) waiting for child process to exit")
					select {
					case c := <-childExitCh:
						log.Printf("[INFO] (runner) child process died")
						r.ErrCh <- NewErrChildDied(c)
						return
					case <-r.DoneCh:
					}
				}

				r.Stop()
				return
			}
		}

	OUTER:
		select {
		case view := <-r.watcher.DataCh():
			// Receive this update
			r.Receive(view.Dependency(), view.Data())

			// Drain all dependency data. Given a large number of dependencies, it is
			// feasible that we have data for more than one of them. Instead of
			// wasting CPU cycles rendering templates when we have more dependencies
			// waiting to be added to the brain, we drain the entire buffered channel
			// on the watcher and then reports when it is done receiving new data
			// which the parent select listens for.
			//
			// Please see https://github.com/hashicorp/consul-template/issues/168 for
			// more information about this optimization and the entire backstory.
			for {
				select {
				case view := <-r.watcher.DataCh():
					r.Receive(view.Dependency(), view.Data())
				default:
					break OUTER
				}
			}

		case <-dedupCh:
			// We may get triggered by the de-duplication manager for either a change
			// in leadership (acquired or lost lock), or an update of data for a template
			// that we are watching.
			log.Printf("[INFO] (runner) watcher triggered by de-duplication manager")
			break OUTER

		case err := <-r.watcher.ErrCh():
			// Push the error back up the stack
			log.Printf("[ERR] (runner) watcher reported error: %s", err)
			r.ErrCh <- err
			return

		case tmpl := <-r.quiescenceCh:
			// Remove the quiescence for this template from the map. This will force
			// the upcoming Run call to actually evaluate and render the template.
			log.Printf("[DEBUG] (runner) received template %q from quiescence", tmpl.ID())
			delete(r.quiescenceMap, tmpl.ID())

		case c := <-childExitCh:
			log.Printf("[INFO] (runner) child process died")
			r.ErrCh <- NewErrChildDied(c)
			return

		case <-r.DoneCh:
			log.Printf("[INFO] (runner) received finish")
			return
		}

		// If we got this far, that means we got new data or one of the timers
		// fired, so attempt to re-render.
		if err := r.Run(); err != nil {
			r.ErrCh <- err
			return
		}
	}
}

// Stop halts the execution of this runner and its subprocesses.
func (r *Runner) Stop() {
	r.internalStop(false)
}

// StopImmediately behaves like Stop but won't wait for any splay on any child
// process it may be running.
func (r *Runner) StopImmediately() {
	r.internalStop(true)
}

// TemplateRenderedCh returns a channel that will be triggered when one or more
// templates are rendered.
func (r *Runner) TemplateRenderedCh() <-chan struct{} {
	return r.renderedCh
}

// RenderEventCh returns a channel that will be triggered when there is a new
// render event.
func (r *Runner) RenderEventCh() <-chan struct{} {
	return r.renderEventCh
}

// RenderEvents returns the render events for each template was rendered. The
// map is keyed by template ID.
func (r *Runner) RenderEvents() map[string]*RenderEvent {
	r.renderEventsLock.RLock()
	defer r.renderEventsLock.RUnlock()

	times := make(map[string]*RenderEvent, len(r.renderEvents))
	for k, v := range r.renderEvents {
		times[k] = v
	}
	return times
}

func (r *Runner) internalStop(immediately bool) {
	r.stopLock.Lock()
	defer r.stopLock.Unlock()

	if r.stopped {
		return
	}

	log.Printf("[INFO] (runner) stopping")
	r.stopDedup()
	r.stopWatcher()
	r.stopChild(immediately)

	if err := r.deletePid(); err != nil {
		log.Printf("[WARN] (runner) could not remove pid at %q: %s",
			*r.config.PidFile, err)
	}

	r.stopped = true

	close(r.DoneCh)
}

func (r *Runner) stopDedup() {
	if r.dedup != nil {
		log.Printf("[DEBUG] (runner) stopping de-duplication manager")
		r.dedup.Stop()
	}
}

func (r *Runner) stopWatcher() {
	if r.watcher != nil {
		log.Printf("[DEBUG] (runner) stopping watcher")
		r.watcher.Stop()
	}
}

func (r *Runner) stopChild(immediately bool) {
	r.childLock.RLock()
	defer r.childLock.RUnlock()

	if r.child != nil {
		if immediately {
			log.Printf("[DEBUG] (runner) stopping child process immediately")
			r.child.StopImmediately()
		} else {
			log.Printf("[DEBUG] (runner) stopping child process")
			r.child.Stop()
		}
	}
}

// Receive accepts a Dependency and data for that dep. This data is
// cached on the Runner. This data is then used to determine if a Template
// is "renderable" (i.e. all its Dependencies have been downloaded at least
// once).
func (r *Runner) Receive(d dep.Dependency, data interface{}) {
	r.dependenciesLock.Lock()
	defer r.dependenciesLock.Unlock()

	// Just because we received data, it does not mean that we are actually
	// watching for that data. How is that possible you may ask? Well, this
	// Runner's data channel is pooled, meaning it accepts multiple data views
	// before actually blocking. Whilest this runner is performing a Run() and
	// executing diffs, it may be possible that more data was pushed onto the
	// data channel pool for a dependency that we no longer care about.
	//
	// Accepting this dependency would introduce stale data into the brain, and
	// that is simply unacceptable. In fact, it is a fun little bug:
	//
	//     https://github.com/hashicorp/consul-template/issues/198
	//
	// and by "little" bug, I mean really big bug.
	if _, ok := r.dependencies[d.String()]; ok {
		log.Printf("[DEBUG] (runner) receiving dependency %s", d)
		r.brain.Remember(d, data)
	}
}

// Signal sends a signal to the child process, if it exists. Any errors that
// occur are returned.
func (r *Runner) Signal(s os.Signal) error {
	r.childLock.RLock()
	defer r.childLock.RUnlock()
	if r.child == nil {
		return nil
	}
	return r.child.Signal(s)
}

// Run iterates over each template in this Runner and conditionally executes
// the template rendering and command execution.
//
// The template is rendered atomically. If and only if the template render
// completes successfully, the optional commands will be executed, if given.
// Please note that all templates are rendered **and then** any commands are
// executed.
func (r *Runner) Run() error {
	log.Printf("[DEBUG] (runner) initiating run")

	var newRenderEvent, wouldRenderAny, renderedAny bool
	runCtx := &templateRunCtx{
		depsMap: make(map[string]dep.Dependency),
	}

	for _, tmpl := range r.templates {
		event, err := r.runTemplate(tmpl, runCtx)
		if err != nil {
			return err
		}

		// If there was a render event store it
		if event != nil {
			r.renderEventsLock.Lock()
			r.renderEvents[tmpl.ID()] = event
			r.renderEventsLock.Unlock()

			// Record that there is at least one new render event
			newRenderEvent = true

			// Record that at least one template would have been rendered.
			if event.WouldRender {
				wouldRenderAny = true
			}

			// Record that at least one template was rendered.
			if event.DidRender {
				renderedAny = true
			}
		}
	}

	// Perform the diff and update the known dependencies.
	r.diffAndUpdateDeps(runCtx.depsMap)

	// Execute each command in sequence, collecting any errors that occur - this
	// ensures all commands execute at least once.
	var errs []error
	for _, t := range runCtx.commands {
		command := config.StringVal(t.Exec.Command)
		log.Printf("[INFO] (runner) executing command %q from %s", command, t.Display())
		env := t.Exec.Env.Copy()
		env.Custom = append(r.childEnv(), env.Custom...)
		if _, err := spawnChild(&spawnChildInput{
			Stdin:        r.inStream,
			Stdout:       r.outStream,
			Stderr:       r.errStream,
			Command:      command,
			Env:          env.Env(),
			Timeout:      config.TimeDurationVal(t.Exec.Timeout),
			ReloadSignal: config.SignalVal(t.Exec.ReloadSignal),
			KillSignal:   config.SignalVal(t.Exec.KillSignal),
			KillTimeout:  config.TimeDurationVal(t.Exec.KillTimeout),
			Splay:        config.TimeDurationVal(t.Exec.Splay),
		}); err != nil {
			s := fmt.Sprintf("failed to execute command %q from %s", command, t.Display())
			errs = append(errs, errors.Wrap(err, s))
		}
	}

	// Check if we need to deliver any rendered signals
	if wouldRenderAny || renderedAny {
		// Send the signal that a template got rendered
		select {
		case r.renderedCh <- struct{}{}:
		default:
		}
	}

	// Check if we need to deliver any event signals
	if newRenderEvent {
		select {
		case r.renderEventCh <- struct{}{}:
		default:
		}
	}

	// If we got this far and have a child process, we need to send the reload
	// signal to the child process.
	if renderedAny && r.child != nil {
		r.childLock.RLock()
		if err := r.child.Reload(); err != nil {
			errs = append(errs, err)
		}
		r.childLock.RUnlock()
	}

	// If any errors were returned, convert them to an ErrorList for human
	// readability.
	if len(errs) != 0 {
		var result *multierror.Error
		for _, err := range errs {
			result = multierror.Append(result, err)
		}
		return result.ErrorOrNil()
	}

	return nil
}

type templateRunCtx struct {
	// commands is the set of commands that will be executed after all templates
	// have run. When adding to the commands, care should be taken not to
	// duplicate any existing command from a previous template.
	commands []*config.TemplateConfig

	// depsMap is the set of dependencies shared across all templates.
	depsMap map[string]dep.Dependency
}

// runTemplate is used to run a particular template. It takes as input the
// template to run and a shared run context that allows sharing of information
// between templates. The run returns a potentially nil render event and any
// error that occured. The render event is nil in the case that the template has
// been already rendered and is a once template or if there is an error.
func (r *Runner) runTemplate(tmpl *template.Template, runCtx *templateRunCtx) (*RenderEvent, error) {
	log.Printf("[DEBUG] (runner) checking template %s", tmpl.ID())

	// Grab the last event
	r.renderEventsLock.RLock()
	lastEvent := r.renderEvents[tmpl.ID()]
	r.renderEventsLock.RUnlock()

	// Create the event
	event := &RenderEvent{
		Template:        tmpl,
		TemplateConfigs: r.templateConfigsFor(tmpl),
	}

	if lastEvent != nil {
		event.LastWouldRender = lastEvent.LastWouldRender
		event.LastDidRender = lastEvent.LastDidRender
	}

	// Check if we are currently the leader instance
	isLeader := true
	if r.dedup != nil {
		isLeader = r.dedup.IsLeader(tmpl)
	}

	// If we are in once mode and this template was already rendered, move
	// onto the next one. We do not want to re-render the template if we are
	// in once mode, and we certainly do not want to re-run any commands.
	if r.config.Once {
		r.renderEventsLock.RLock()
		event, ok := r.renderEvents[tmpl.ID()]
		r.renderEventsLock.RUnlock()
		if ok && (event.WouldRender || event.DidRender) {
			log.Printf("[DEBUG] (runner) once mode and already rendered")
			return nil, nil
		}
	}

	// Attempt to render the template, returning any missing dependencies and
	// the rendered contents. If there are any missing dependencies, the
	// contents cannot be rendered or trusted!
	result, err := tmpl.Execute(&template.ExecuteInput{
		Brain: r.brain,
		Env:   r.childEnv(),
	})
	if err != nil {
		return nil, errors.Wrap(err, tmpl.Source())
	}

	// Grab the list of used and missing dependencies.
	missing, used := result.Missing, result.Used

	if l := missing.Len(); l > 0 {
		log.Printf("[DEBUG] (runner) missing data for %d dependencies", l)
		for _, missingDependency := range missing.List() {
			log.Printf("[DEBUG] (runner) missing dependency: %s", missingDependency)
		}
	}

	// Add the dependency to the list of dependencies for this runner.
	for _, d := range used.List() {
		// If we've taken over leadership for a template, we may have data
		// that is cached, but not have the watcher. We must treat this as
		// missing so that we create the watcher and re-run the template.
		if isLeader && !r.watcher.Watching(d) {
			log.Printf("[DEBUG] (runner) add used dependency %s to missing since isLeader but do not have a watcher", d)
			missing.Add(d)
		}
		if _, ok := runCtx.depsMap[d.String()]; !ok {
			runCtx.depsMap[d.String()] = d
		}
	}

	// Diff any missing dependencies the template reported with dependencies
	// the watcher is watching.
	unwatched := new(dep.Set)
	for _, d := range missing.List() {
		if !r.watcher.Watching(d) {
			unwatched.Add(d)
		}
	}

	// Update the event with the new dependency information
	event.MissingDeps = missing
	event.UnwatchedDeps = unwatched
	event.UsedDeps = used
	event.UpdatedAt = time.Now().UTC()

	// If there are unwatched dependencies, start the watcher and exit since we
	// won't have data.
	if l := unwatched.Len(); l > 0 {
		log.Printf("[DEBUG] (runner) was not watching %d dependencies", l)
		for _, d := range unwatched.List() {
			// If we are deduplicating, we must still handle non-sharable
			// dependencies, since those will be ignored.
			if isLeader || !d.CanShare() {
				r.watcher.Add(d)
			}
		}
		return event, nil
	}

	// If the template is missing data for some dependencies then we are not
	// ready to render and need to move on to the next one.
	if l := missing.Len(); l > 0 {
		log.Printf("[DEBUG] (runner) missing data for %d dependencies", l)
		return event, nil
	}

	// Trigger an update of the de-duplication manager
	if r.dedup != nil && isLeader {
		if err := r.dedup.UpdateDeps(tmpl, used.List()); err != nil {
			log.Printf("[ERR] (runner) failed to update dependency data for de-duplication: %v", err)
		}
	}

	// If quiescence is activated, start/update the timers and loop back around.
	// We do not want to render the templates yet.
	if q, ok := r.quiescenceMap[tmpl.ID()]; ok {
		q.tick()
		// This event is being returned early for quiescence
		event.ForQuiescence = true
		return event, nil
	}

	// For each template configuration that is tied to this template, attempt to
	// render it to disk and accumulate commands for later use.
	for _, templateConfig := range r.templateConfigsFor(tmpl) {
		log.Printf("[DEBUG] (runner) rendering %s", templateConfig.Display())

		// Render the template, taking dry mode into account
		result, err := renderer.Render(&renderer.RenderInput{
			Backup:         config.BoolVal(templateConfig.Backup),
			Contents:       result.Output,
			CreateDestDirs: config.BoolVal(templateConfig.CreateDestDirs),
			Dry:            r.dry,
			DryStream:      r.outStream,
			Path:           config.StringVal(templateConfig.Destination),
			Perms:          config.FileModeVal(templateConfig.Perms),
		})
		if err != nil {
			return nil, errors.Wrap(err, "error rendering "+templateConfig.Display())
		}

		renderTime := time.Now().UTC()

		// If we would have rendered this template (but we did not because the
		// contents were the same or something), we should consider this template
		// rendered even though the contents on disk have not been updated. We
		// will not fire commands unless the template was _actually_ rendered to
		// disk though.
		if result.WouldRender {
			// This event would have rendered
			event.WouldRender = true
			event.LastWouldRender = renderTime
		}

		// If we _actually_ rendered the template to disk, we want to run the
		// appropriate commands.
		if result.DidRender {
			log.Printf("[INFO] (runner) rendered %s", templateConfig.Display())

			// This event did render
			event.DidRender = true
			event.LastDidRender = renderTime

			// Update the contents
			event.Contents = result.Contents

			if !r.dry {
				// If the template was rendered (changed) and we are not in dry-run mode,
				// aggregate commands, ignoring previously known commands
				//
				// Future-self Q&A: Why not use a map for the commands instead of an
				// array with an expensive lookup option? Well I'm glad you asked that
				// future-self! One of the API promises is that commands are executed
				// in the order in which they are provided in the TemplateConfig
				// definitions. If we inserted commands into a map, we would lose that
				// relative ordering and people would be unhappy.
				// if config.StringPresent(ctemplate.Command)
				if c := config.StringVal(templateConfig.Exec.Command); c != "" {
					existing := findCommand(templateConfig, runCtx.commands)
					if existing != nil {
						log.Printf("[DEBUG] (runner) skipping command %q from %s (already appended from %s)",
							c, templateConfig.Display(), existing.Display())
					} else {
						log.Printf("[DEBUG] (runner) appending command %q from %s",
							c, templateConfig.Display())
						runCtx.commands = append(runCtx.commands, templateConfig)
					}
				}
			}
		}
	}

	return event, nil
}

// init() creates the Runner's underlying data structures and returns an error
// if any problems occur.
func (r *Runner) init() error {
	// Ensure default configuration values
	r.config = config.DefaultConfig().Merge(r.config)
	r.config.Finalize()

	// Print the final config for debugging
	result, err := json.Marshal(r.config)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] (runner) final config: %s", result)

	// Create the clientset
	clients, err := newClientSet(r.config)
	if err != nil {
		return fmt.Errorf("runner: %s", err)
	}

	// Create the watcher
	watcher, err := newWatcher(r.config, clients, r.config.Once)
	if err != nil {
		return fmt.Errorf("runner: %s", err)
	}
	r.watcher = watcher

	numTemplates := len(*r.config.Templates)
	templates := make([]*template.Template, 0, numTemplates)
	ctemplatesMap := make(map[string]config.TemplateConfigs)

	// Iterate over each TemplateConfig, creating a new Template resource for each
	// entry. Templates are parsed and saved, and a map of templates to their
	// config templates is kept so templates can lookup their commands and output
	// destinations.
	for _, ctmpl := range *r.config.Templates {
		leftDelim := config.StringVal(ctmpl.LeftDelim)
		if leftDelim == "" {
			leftDelim = config.StringVal(r.config.DefaultDelims.Left)
		}
		rightDelim := config.StringVal(ctmpl.RightDelim)
		if rightDelim == "" {
			rightDelim = config.StringVal(r.config.DefaultDelims.Right)
		}

		tmpl, err := template.NewTemplate(&template.NewTemplateInput{
			Source:           config.StringVal(ctmpl.Source),
			Contents:         config.StringVal(ctmpl.Contents),
			ErrMissingKey:    config.BoolVal(ctmpl.ErrMissingKey),
			LeftDelim:        leftDelim,
			RightDelim:       rightDelim,
			FunctionDenylist: ctmpl.FunctionDenylist,
			SandboxPath:      config.StringVal(ctmpl.SandboxPath),
		})
		if err != nil {
			return err
		}

		if _, ok := ctemplatesMap[tmpl.ID()]; !ok {
			templates = append(templates, tmpl)
		}

		if _, ok := ctemplatesMap[tmpl.ID()]; !ok {
			ctemplatesMap[tmpl.ID()] = make([]*config.TemplateConfig, 0, 1)
		}
		ctemplatesMap[tmpl.ID()] = append(ctemplatesMap[tmpl.ID()], ctmpl)
	}

	// Convert the map of templates (which was only used to ensure uniqueness)
	// back into an array of templates.
	r.templates = templates

	r.renderEvents = make(map[string]*RenderEvent, numTemplates)
	r.dependencies = make(map[string]dep.Dependency)

	r.renderedCh = make(chan struct{}, 1)
	r.renderEventCh = make(chan struct{}, 1)

	r.ctemplatesMap = ctemplatesMap
	r.inStream = os.Stdin
	r.outStream = os.Stdout
	r.errStream = os.Stderr
	r.brain = template.NewBrain()

	r.ErrCh = make(chan error)
	r.DoneCh = make(chan struct{})

	r.quiescenceMap = make(map[string]*quiescence)
	r.quiescenceCh = make(chan *template.Template)

	if *r.config.Dedup.Enabled {
		if r.config.Once {
			log.Printf("[INFO] (runner) disabling de-duplication in once mode")
		} else {
			r.dedup, err = NewDedupManager(r.config.Dedup, clients, r.brain, r.templates)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// diffAndUpdateDeps iterates through the current map of dependencies on this
// runner and stops the watcher for any deps that are no longer required.
//
// At the end of this function, the given depsMap is converted to a slice and
// stored on the runner.
func (r *Runner) diffAndUpdateDeps(depsMap map[string]dep.Dependency) {
	r.dependenciesLock.Lock()
	defer r.dependenciesLock.Unlock()

	// Diff and up the list of dependencies, stopping any unneeded watchers.
	log.Printf("[DEBUG] (runner) diffing and updating dependencies")

	for key, d := range r.dependencies {
		if _, ok := depsMap[key]; !ok {
			log.Printf("[DEBUG] (runner) %s is no longer needed", d)
			r.watcher.Remove(d)
			r.brain.Forget(d)
		} else {
			log.Printf("[DEBUG] (runner) %s is still needed", d)
		}
	}

	r.dependencies = depsMap
}

// TemplateConfigFor returns the TemplateConfig for the given Template
func (r *Runner) templateConfigsFor(tmpl *template.Template) []*config.TemplateConfig {
	return r.ctemplatesMap[tmpl.ID()]
}

// TemplateConfigMapping returns a mapping between the template ID and the set
// of TemplateConfig represented by the template ID
func (r *Runner) TemplateConfigMapping() map[string][]*config.TemplateConfig {
	// this method is primarily used to support embedding consul-template
	// in other applications (ex. Nomad)
	m := make(map[string][]*config.TemplateConfig, len(r.ctemplatesMap))

	for id, set := range r.ctemplatesMap {
		ctmpls := make([]*config.TemplateConfig, len(set))
		m[id] = ctmpls
		for i, ctmpl := range set {
			ctmpls[i] = ctmpl
		}
	}

	return m
}

// allTemplatesRendered returns true if all the templates in this Runner have
// been rendered at least one time.
func (r *Runner) allTemplatesRendered() bool {
	r.renderEventsLock.RLock()
	defer r.renderEventsLock.RUnlock()

	for _, tmpl := range r.templates {
		event, rendered := r.renderEvents[tmpl.ID()]
		if !rendered {
			return false
		}

		// Skip evaluation of events from quiescence as they will
		// be default unrendered as we are still waiting for the
		// specified period
		if event.ForQuiescence {
			continue
		}

		// The template might already exist on disk with the exact contents, but
		// we still want to count that as "rendered" [GH-1000].
		if !event.DidRender && !event.WouldRender {
			return false
		}
	}

	return true
}

// childEnv creates a map of environment variables for child processes to have
// access to configurations in Consul Template's configuration.
func (r *Runner) childEnv() []string {
	var m = make(map[string]string)

	if config.StringPresent(r.config.Consul.Address) {
		m["CONSUL_HTTP_ADDR"] = config.StringVal(r.config.Consul.Address)
	}

	if config.BoolVal(r.config.Consul.Auth.Enabled) {
		m["CONSUL_HTTP_AUTH"] = r.config.Consul.Auth.String()
	}

	m["CONSUL_HTTP_SSL"] = strconv.FormatBool(config.BoolVal(r.config.Consul.SSL.Enabled))
	m["CONSUL_HTTP_SSL_VERIFY"] = strconv.FormatBool(config.BoolVal(r.config.Consul.SSL.Verify))

	if config.StringPresent(r.config.Vault.Address) {
		m["VAULT_ADDR"] = config.StringVal(r.config.Vault.Address)
	}

	if !config.BoolVal(r.config.Vault.SSL.Verify) {
		m["VAULT_SKIP_VERIFY"] = "true"
	}

	if config.StringPresent(r.config.Vault.SSL.Cert) {
		m["VAULT_CLIENT_CERT"] = config.StringVal(r.config.Vault.SSL.Cert)
	}

	if config.StringPresent(r.config.Vault.SSL.Key) {
		m["VAULT_CLIENT_KEY"] = config.StringVal(r.config.Vault.SSL.Key)
	}

	if config.StringPresent(r.config.Vault.SSL.CaPath) {
		m["VAULT_CAPATH"] = config.StringVal(r.config.Vault.SSL.CaPath)
	}

	if config.StringPresent(r.config.Vault.SSL.CaCert) {
		m["VAULT_CACERT"] = config.StringVal(r.config.Vault.SSL.CaCert)
	}

	if config.StringPresent(r.config.Vault.SSL.ServerName) {
		m["VAULT_TLS_SERVER_NAME"] = config.StringVal(r.config.Vault.SSL.ServerName)
	}

	// Append runner-supplied env (this is supplied programmatically).
	for k, v := range r.Env {
		m[k] = v
	}

	e := make([]string, 0, len(m))
	for k, v := range m {
		e = append(e, k+"="+v)
	}
	return e
}

// storePid is used to write out a PID file to disk.
func (r *Runner) storePid() error {
	path := config.StringVal(r.config.PidFile)
	if path == "" {
		return nil
	}

	log.Printf("[INFO] creating pid file at %q", path)

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("runner: could not open pid file: %s", err)
	}
	defer f.Close()

	pid := os.Getpid()
	_, err = f.WriteString(fmt.Sprintf("%d", pid))
	if err != nil {
		return fmt.Errorf("runner: could not write to pid file: %s", err)
	}
	return nil
}

// deletePid is used to remove the PID on exit.
func (r *Runner) deletePid() error {
	path := config.StringVal(r.config.PidFile)
	if path == "" {
		return nil
	}

	log.Printf("[DEBUG] removing pid file at %q", path)

	stat, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("runner: could not remove pid file: %s", err)
	}
	if stat.IsDir() {
		return fmt.Errorf("runner: specified pid file path is directory")
	}

	err = os.Remove(path)
	if err != nil {
		return fmt.Errorf("runner: could not remove pid file: %s", err)
	}
	return nil
}

// SetOutStream modifies runner output stream. Defaults to stdout.
func (r *Runner) SetOutStream(out io.Writer) {
	r.outStream = out
}

// SetErrStream modifies runner error stream. Defaults to stderr.
func (r *Runner) SetErrStream(err io.Writer) {
	r.errStream = err
}

// spawnChildInput is used as input to spawn a child process.
type spawnChildInput struct {
	Stdin        io.Reader
	Stdout       io.Writer
	Stderr       io.Writer
	Command      string
	Timeout      time.Duration
	Env          []string
	ReloadSignal os.Signal
	KillSignal   os.Signal
	KillTimeout  time.Duration
	Splay        time.Duration
}

// spawnChild spawns a child process with the given inputs and returns the
// resulting child.
func spawnChild(i *spawnChildInput) (*child.Child, error) {
	p := shellwords.NewParser()
	p.ParseEnv = true
	p.ParseBacktick = true
	args, err := p.Parse(i.Command)
	if err != nil {
		return nil, errors.Wrap(err, "failed parsing command")
	}

	child, err := child.New(&child.NewInput{
		Stdin:        i.Stdin,
		Stdout:       i.Stdout,
		Stderr:       i.Stderr,
		Command:      args[0],
		Args:         args[1:],
		Env:          i.Env,
		Timeout:      i.Timeout,
		ReloadSignal: i.ReloadSignal,
		KillSignal:   i.KillSignal,
		KillTimeout:  i.KillTimeout,
		Splay:        i.Splay,
	})
	if err != nil {
		return nil, errors.Wrap(err, "error creating child")
	}

	if err := child.Start(); err != nil {
		return nil, errors.Wrap(err, "child")
	}
	return child, nil
}

// quiescence is an internal representation of a single template's quiescence
// state.
type quiescence struct {
	template *template.Template
	min      time.Duration
	max      time.Duration
	ch       chan *template.Template
	timer    *time.Timer
	deadline time.Time
}

// newQuiescence creates a new quiescence timer for the given template.
func newQuiescence(ch chan *template.Template, min, max time.Duration, t *template.Template) *quiescence {
	return &quiescence{
		template: t,
		min:      min,
		max:      max,
		ch:       ch,
	}
}

// tick updates the minimum quiescence timer.
func (q *quiescence) tick() {
	now := time.Now()

	// If this is the first tick, set up the timer and calculate the max
	// deadline.
	if q.timer == nil {
		q.timer = time.NewTimer(q.min)
		go func() {
			select {
			case <-q.timer.C:
				q.ch <- q.template
			}
		}()

		q.deadline = now.Add(q.max)
		return
	}

	// Snooze the timer for the min time, or snooze less if we are coming
	// up against the max time. If the timer has already fired and the reset
	// doesn't work that's ok because we guarantee that the channel gets our
	// template which means that we are obsolete and a fresh quiescence will
	// be set up.
	if now.Add(q.min).Before(q.deadline) {
		q.timer.Reset(q.min)
	} else if dur := q.deadline.Sub(now); dur > 0 {
		q.timer.Reset(dur)
	}
}

// findCommand searches the list of template configs for the given command and
// returns it if it exists.
func findCommand(c *config.TemplateConfig, templates []*config.TemplateConfig) *config.TemplateConfig {
	needle := config.StringVal(c.Exec.Command)
	for _, t := range templates {
		if needle == config.StringVal(t.Exec.Command) {
			return t
		}
	}
	return nil
}

// newClientSet creates a new client set from the given config.
func newClientSet(c *config.Config) (*dep.ClientSet, error) {
	clients := dep.NewClientSet()

	if err := clients.CreateConsulClient(&dep.CreateConsulClientInput{
		Address:                      config.StringVal(c.Consul.Address),
		Namespace:                    config.StringVal(c.Consul.Namespace),
		Token:                        config.StringVal(c.Consul.Token),
		AuthEnabled:                  config.BoolVal(c.Consul.Auth.Enabled),
		AuthUsername:                 config.StringVal(c.Consul.Auth.Username),
		AuthPassword:                 config.StringVal(c.Consul.Auth.Password),
		SSLEnabled:                   config.BoolVal(c.Consul.SSL.Enabled),
		SSLVerify:                    config.BoolVal(c.Consul.SSL.Verify),
		SSLCert:                      config.StringVal(c.Consul.SSL.Cert),
		SSLKey:                       config.StringVal(c.Consul.SSL.Key),
		SSLCACert:                    config.StringVal(c.Consul.SSL.CaCert),
		SSLCAPath:                    config.StringVal(c.Consul.SSL.CaPath),
		ServerName:                   config.StringVal(c.Consul.SSL.ServerName),
		TransportDialKeepAlive:       config.TimeDurationVal(c.Consul.Transport.DialKeepAlive),
		TransportDialTimeout:         config.TimeDurationVal(c.Consul.Transport.DialTimeout),
		TransportDisableKeepAlives:   config.BoolVal(c.Consul.Transport.DisableKeepAlives),
		TransportIdleConnTimeout:     config.TimeDurationVal(c.Consul.Transport.IdleConnTimeout),
		TransportMaxIdleConns:        config.IntVal(c.Consul.Transport.MaxIdleConns),
		TransportMaxIdleConnsPerHost: config.IntVal(c.Consul.Transport.MaxIdleConnsPerHost),
		TransportTLSHandshakeTimeout: config.TimeDurationVal(c.Consul.Transport.TLSHandshakeTimeout),
	}); err != nil {
		return nil, fmt.Errorf("runner: %s", err)
	}

	if err := clients.CreateVaultClient(&dep.CreateVaultClientInput{
		Address:                      config.StringVal(c.Vault.Address),
		Namespace:                    config.StringVal(c.Vault.Namespace),
		Token:                        config.StringVal(c.Vault.Token),
		UnwrapToken:                  config.BoolVal(c.Vault.UnwrapToken),
		SSLEnabled:                   config.BoolVal(c.Vault.SSL.Enabled),
		SSLVerify:                    config.BoolVal(c.Vault.SSL.Verify),
		SSLCert:                      config.StringVal(c.Vault.SSL.Cert),
		SSLKey:                       config.StringVal(c.Vault.SSL.Key),
		SSLCACert:                    config.StringVal(c.Vault.SSL.CaCert),
		SSLCAPath:                    config.StringVal(c.Vault.SSL.CaPath),
		ServerName:                   config.StringVal(c.Vault.SSL.ServerName),
		TransportDialKeepAlive:       config.TimeDurationVal(c.Vault.Transport.DialKeepAlive),
		TransportDialTimeout:         config.TimeDurationVal(c.Vault.Transport.DialTimeout),
		TransportDisableKeepAlives:   config.BoolVal(c.Vault.Transport.DisableKeepAlives),
		TransportIdleConnTimeout:     config.TimeDurationVal(c.Vault.Transport.IdleConnTimeout),
		TransportMaxIdleConns:        config.IntVal(c.Vault.Transport.MaxIdleConns),
		TransportMaxIdleConnsPerHost: config.IntVal(c.Vault.Transport.MaxIdleConnsPerHost),
		TransportTLSHandshakeTimeout: config.TimeDurationVal(c.Vault.Transport.TLSHandshakeTimeout),
	}); err != nil {
		return nil, fmt.Errorf("runner: %s", err)
	}

	return clients, nil
}

// newWatcher creates a new watcher.
func newWatcher(c *config.Config, clients *dep.ClientSet, once bool) (*watch.Watcher, error) {
	log.Printf("[INFO] (runner) creating watcher")

	w, err := watch.NewWatcher(&watch.NewWatcherInput{
		Clients:             clients,
		MaxStale:            config.TimeDurationVal(c.MaxStale),
		Once:                c.Once,
		BlockQueryWaitTime:  config.TimeDurationVal(c.BlockQueryWaitTime),
		RenewVault:          clients.Vault().Token() != "" && config.BoolVal(c.Vault.RenewToken),
		VaultAgentTokenFile: config.StringVal(c.Vault.VaultAgentTokenFile),
		RetryFuncConsul:     watch.RetryFunc(c.Consul.Retry.RetryFunc()),
		// TODO: Add a sane default retry - right now this only affects "local"
		// dependencies like reading a file from disk.
		RetryFuncDefault: nil,
		RetryFuncVault:   watch.RetryFunc(c.Vault.Retry.RetryFunc()),
		VaultToken:       clients.Vault().Token(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "runner")
	}
	return w, nil
}
