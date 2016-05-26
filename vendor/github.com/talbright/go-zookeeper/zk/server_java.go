package zk

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type retryError struct {
	msg      string
	duration time.Duration
	attempts int
}

func newRetryError(msg string, duration time.Duration, attempts int) *retryError {
	return &retryError{msg: msg, duration: duration, attempts: attempts}
}

func (e *retryError) Error() string {
	return fmt.Sprintf("retry failed: %s (retried %d times over %v)", e.msg, e.attempts, e.duration*time.Duration(e.attempts))
}

var maxStartStopPolls = 30
var startStopPollInterval = time.Second

type ErrMissingServerConfigField string

func (e ErrMissingServerConfigField) Error() string {
	return fmt.Sprintf("zk: missing server config field '%s'", string(e))
}

const (
	DefaultServerTickTime                 = 2000
	DefaultServerInitLimit                = 10
	DefaultServerSyncLimit                = 5
	DefaultServerAutoPurgeSnapRetainCount = 3
	DefaultPeerPort                       = 2888
	DefaultLeaderElectionPort             = 3888
)

type ServerConfigServer struct {
	ID                 int
	Host               string
	PeerPort           int
	LeaderElectionPort int
}

type ServerConfig struct {
	TickTime                 int    // Number of milliseconds of each tick
	InitLimit                int    // Number of ticks that the initial synchronization phase can take
	SyncLimit                int    // Number of ticks that can pass between sending a request and getting an acknowledgement
	DataDir                  string // Direcrory where the snapshot is stored
	ClientPort               int    // Port at which clients will connect
	AutoPurgeSnapRetainCount int    // Number of snapshots to retain in dataDir
	AutoPurgePurgeInterval   int    // Purge task internal in hours (0 to disable auto purge)
	Servers                  []ServerConfigServer
}

func (sc ServerConfig) Marshall(w io.Writer) error {
	if sc.DataDir == "" {
		return ErrMissingServerConfigField("dataDir")
	}
	fmt.Fprintf(w, "dataDir=%s\n", sc.DataDir)
	if sc.TickTime <= 0 {
		sc.TickTime = DefaultServerTickTime
	}
	fmt.Fprintf(w, "tickTime=%d\n", sc.TickTime)
	if sc.InitLimit <= 0 {
		sc.InitLimit = DefaultServerInitLimit
	}
	fmt.Fprintf(w, "initLimit=%d\n", sc.InitLimit)
	if sc.SyncLimit <= 0 {
		sc.SyncLimit = DefaultServerSyncLimit
	}
	fmt.Fprintf(w, "syncLimit=%d\n", sc.SyncLimit)
	if sc.ClientPort <= 0 {
		sc.ClientPort = DefaultPort
	}
	fmt.Fprintf(w, "clientPort=%d\n", sc.ClientPort)
	if sc.AutoPurgePurgeInterval > 0 {
		if sc.AutoPurgeSnapRetainCount <= 0 {
			sc.AutoPurgeSnapRetainCount = DefaultServerAutoPurgeSnapRetainCount
		}
		fmt.Fprintf(w, "autopurge.snapRetainCount=%d\n", sc.AutoPurgeSnapRetainCount)
		fmt.Fprintf(w, "autopurge.purgeInterval=%d\n", sc.AutoPurgePurgeInterval)
	}
	if len(sc.Servers) > 0 {
		for _, srv := range sc.Servers {
			if srv.PeerPort <= 0 {
				srv.PeerPort = DefaultPeerPort
			}
			if srv.LeaderElectionPort <= 0 {
				srv.LeaderElectionPort = DefaultLeaderElectionPort
			}
			fmt.Fprintf(w, "server.%d=%s:%d:%d\n", srv.ID, srv.Host, srv.PeerPort, srv.LeaderElectionPort)
		}
	}
	return nil
}

var jarSearchPaths = []string{
	"zookeeper-*/contrib/fatjar/zookeeper-*-fatjar.jar",
	"../zookeeper-*/contrib/fatjar/zookeeper-*-fatjar.jar",
	"/usr/share/java/zookeeper-*.jar",
	"/usr/local/zookeeper-*/contrib/fatjar/zookeeper-*-fatjar.jar",
	"/usr/local/Cellar/zookeeper/*/libexec/contrib/fatjar/zookeeper-*-fatjar.jar",
}

func findZookeeperFatJar() string {
	var paths []string
	zkPath := os.Getenv("ZOOKEEPER_PATH")
	if zkPath == "" {
		paths = jarSearchPaths
	} else {
		paths = []string{filepath.Join(zkPath, "contrib/fatjar/zookeeper-*-fatjar.jar")}
	}
	for _, path := range paths {
		matches, _ := filepath.Glob(path)
		// TODO: could sort by version and pick latest
		if len(matches) > 0 {
			return matches[0]
		}
	}
	return ""
}

type serverState int

const (
	serverStateNew = iota
	serverStateStopped
	serverStateStarted
	serverStateInconsistent
)

var serverStates = [...]string{
	"init",
	"stopped",
	"started",
	"inconsistent",
}

func (state serverState) String() string {
	return serverStates[state]
}

type Server struct {
	JarPath        string
	ConfigPath     string
	Stdout, Stderr io.Writer
	Address        string
	cmd            *exec.Cmd
	state          serverState
}

func (srv *Server) Start() (err error) {
	DefaultLogger.Printf("starting %s [state:%v]", srv.Address, srv.state)
	defer func() {
		if err == nil {
			srv.state = serverStateStarted
			DefaultLogger.Printf("started %s successfully", srv.Address)
		} else {
			srv.state = serverStateInconsistent
			DefaultLogger.Printf("start %s failed: %v", srv.Address, err)
		}
	}()
	if srv.JarPath == "" {
		srv.JarPath = findZookeeperFatJar()
		if srv.JarPath == "" {
			return fmt.Errorf("zk: unable to find server jar")
		}
	}
	srv.cmd = exec.Command("java", "-jar", srv.JarPath, "server", srv.ConfigPath)
	srv.cmd.Stdout = srv.Stdout
	srv.cmd.Stderr = srv.Stderr
	if err = srv.cmd.Start(); err == nil {
		var i int
		for ; i < maxStartStopPolls; i++ {
			if ok := FLWRuok([]string{srv.Address}, time.Second); ok[0] {
				return nil
			}
			time.Sleep(startStopPollInterval)
		}
		err = newRetryError(fmt.Sprintf("starting %v", srv), startStopPollInterval, i)
	}
	return
}

func (srv *Server) Stop() (err error) {
	DefaultLogger.Printf("stopping %s [state:%v]", srv.Address, srv.state)
	defer func() {
		if err == nil {
			srv.state = serverStateStopped
			DefaultLogger.Printf("stopped %s successfully", srv.Address)
		} else {
			if err.Error() == "os: process already finished" {
				srv.state = serverStateStopped
			} else {
				srv.state = serverStateInconsistent
			}
			DefaultLogger.Printf("stop %s failed: %v", srv.Address, err)
		}
	}()
	errOk := srv.cmd.Process.Signal(os.Kill)
	if errOk != nil {
		DefaultLogger.Printf("error signaling kill while stopping %s: %v", srv.Address, errOk)
	}
	if errOk = srv.cmd.Wait(); errOk.Error() != "signal: killed" {
		DefaultLogger.Printf("unexpected error from wait while stopping %s: %v", srv.Address, errOk)
	}
	var i int
	for ; i < maxStartStopPolls; i++ {
		if ok := FLWRuok([]string{srv.Address}, time.Second); !ok[0] {
			return nil
		}
		time.Sleep(startStopPollInterval)
	}
	err = newRetryError(fmt.Sprintf("stopping %v", srv), startStopPollInterval, i)
	return
}
