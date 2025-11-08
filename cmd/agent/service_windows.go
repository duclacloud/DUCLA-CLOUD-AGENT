// +build windows

package main

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

const serviceName = "DuclaAgent"
const serviceDesc = "Ducla Cloud Agent for distributed task execution"

type windowsService struct{}

func (m *windowsService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown
	changes <- svc.Status{State: svc.StartPending}
	
	// Start service
	go runAgent()
	
	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
	
	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				changes <- svc.Status{State: svc.StopPending}
				stopAgent()
				return
			}
		}
	}
}

func installService() error {
	exepath, err := os.Executable()
	if err != nil {
		return err
	}
	
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	
	s, err := m.OpenService(serviceName)
	if err == nil {
		s.Close()
		return fmt.Errorf("service %s already exists", serviceName)
	}
	
	s, err = m.CreateService(serviceName, exepath, mgr.Config{
		DisplayName: serviceDesc,
		StartType:   mgr.StartAutomatic,
	})
	if err != nil {
		return err
	}
	defer s.Close()
	
	return nil
}

func uninstallService() error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	
	s, err := m.OpenService(serviceName)
	if err != nil {
		return fmt.Errorf("service %s is not installed", serviceName)
	}
	defer s.Close()
	
	err = s.Delete()
	if err != nil {
		return err
	}
	
	return nil
}
