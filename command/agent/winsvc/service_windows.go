// +build windows

package winsvc

import (
	wsvc "golang.org/x/sys/windows/svc"
)

type serviceWindows struct{}

func init() {
	inService, err := wsvc.IsWindowsService()
	if err != nil {
		panic(err)
	}

	// Cannot run as a service when running interactively
	if !inService {
		return
	}

	go wsvc.Run("", serviceWindows{})
}

// Execute implements the Windows service Handler type. It will be
// called at the start of the service, and the service will exit
// once Execute completes.
func (serviceWindows) Execute(args []string, r <-chan wsvc.ChangeRequest, s chan<- wsvc.Status) (svcSpecificEC bool, exitCode uint32) {
	const accCommands = wsvc.AcceptStop | wsvc.AcceptShutdown
	s <- wsvc.Status{State: wsvc.StartPending}
	s <- wsvc.Status{State: wsvc.Running, Accepts: accCommands}
	for {
		c := <-r
		switch c.Cmd {
		case wsvc.Interrogate:
			s <- c.CurrentStatus
		case wsvc.Stop, wsvc.Shutdown:
			s <- wsvc.Status{State: wsvc.StopPending}
			chanGraceExit <- 1
			return false, 0
		}
	}

	return false, 0
}
