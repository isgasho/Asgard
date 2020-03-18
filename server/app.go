package server

import (
	"fmt"

	"Asgard/applications"
	"Asgard/client"
	"Asgard/rpc"
)

func GetAppList() []*rpc.App {
	list := []*rpc.App{}
	for _, app := range applications.APPs {
		list = append(list, rpc.BuildApp(app))
	}
	return list
}

func GetApp(id int64) *rpc.App {
	if app, ok := applications.APPs[id]; ok {
		return rpc.BuildApp(app)
	}
	return nil
}

func GetAppByName(name string) *rpc.App {
	for _, app := range applications.APPs {
		if name == app.Name {
			return rpc.BuildApp(app)
		}
	}
	return nil
}

func AddApp(id int64, appRequest *rpc.App) error {
	_, ok := applications.APPs[id]
	if ok {
		return nil
	}
	app, err := applications.AppRegister(id, rpc.BuildAppConfig(appRequest))
	if err != nil {
		return err
	}
	app.MonitorReport = func(monitor *applications.Monitor) {
		client.AppMonitorReport(rpc.BuildAppMonitor(app, monitor))
	}
	app.ArchiveReport = func(command *applications.Command) {
		client.AppArchiveReport(rpc.BuildAppArchive(app, command))
	}
	ok = applications.AppStartByID(id)
	if !ok {
		return fmt.Errorf("app %d start failed", id)
	}
	return nil
}

func UpdateApp(id int64, appRequest *rpc.App) error {
	if _, ok := applications.APPs[id]; ok {
		if err := DeleteApp(id); err != nil {
			return err
		}
		return AddApp(id, appRequest)
	} else {
		return fmt.Errorf("no app %d", id)
	}
}

func DeleteApp(id int64) error {
	if app, ok := applications.APPs[id]; ok {
		app.AutoRestart = false
		if ok := applications.AppStopByID(id); !ok {
			return fmt.Errorf("app %d stop failed", id)
		} else {
			delete(applications.APPs, id)
			return nil
		}
	} else {
		return nil
	}
}

func DeleteAppByName(name string) error {
	for _, app := range applications.APPs {
		if name == app.Name {
			return DeleteApp(app.ID)
		}
	}
	return nil
}
