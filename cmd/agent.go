package cmd

import (
	"net"

	"github.com/dalonghahaha/avenger/components/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc/reflection"

	"Asgard/applications"
	"Asgard/client"
	"Asgard/rpc"
	"Asgard/server"
)

var (
	masterClient rpc.MasterClient
	agentIP      string
	agentPort    string
)

func init() {
	agentCommonCmd.PersistentFlags().StringP("conf", "c", "conf", "config path")
	rootCmd.AddCommand(agentCommonCmd)
}

var agentCommonCmd = &cobra.Command{
	Use:    "agent",
	Short:  "run as agent",
	PreRun: PreRun,
	Run: func(cmd *cobra.Command, args []string) {
		client.InitMasterClient()
		go StartAgent()
		go StartAgentRpcServer()
		NotityKill(StopAgent)
	},
}

func StartAgent() {
	agentIP = viper.GetString("agent.rpc.ip")
	agentPort = viper.GetString("agent.rpc.port")
	if agentIP == "" && agentPort == "" {
		panic("agent config error")
	}
	err := client.AgentRegister(agentIP, agentPort)
	if err != nil {
		panic(err)
	}
	err = AppsRegister()
	if err != nil {
		panic(err)
	}
	err = JobsRegister()
	if err != nil {
		panic(err)
	}
	applications.AppStartAll(false)
	applications.JobStartAll(false)
	applications.MoniterStart()
}

func StopAgent() {
	applications.AppStopAll()
	applications.JobStopAll()
}

func StartAgentRpcServer() {
	port := viper.GetString("agent.rpc.port")
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		logger.Error("failed to listen:", err)
		panic(err)
	}
	s := server.DefaultServer()
	rpc.RegisterGuardServer(s, &server.GuardServer{})
	rpc.RegisterCronServer(s, &server.CronServer{})
	reflection.Register(s)
	logger.Info("agent rpc started at ", port)
	err = s.Serve(listen)
	if err != nil {
		logger.Error("failed to serve:", err)
		panic(err)
	}
}

func AppsRegister() error {
	apps, err := client.GetAppList(agentIP, agentPort)
	if err != nil {
		return err
	}
	for _, info := range apps {
		logger.Debug("app register: ", info.GetName())
		config := map[string]interface{}{
			"id":           info.GetId(),
			"name":         info.GetName(),
			"dir":          info.GetDir(),
			"program":      info.GetProgram(),
			"args":         info.GetArgs(),
			"stdout":       info.GetStdOut(),
			"stderr":       info.GetStdErr(),
			"auto_restart": info.GetAutoRestart(),
			"is_monitor":   info.GetIsMonitor(),
		}
		app, err := applications.AppRegister(info.GetId(), config)
		if err != nil {
			return err
		}
		app.MonitorReport = func(monitor *applications.Monitor) {
			client.AppMonitorReport(app, monitor)
		}
		app.ArchiveReport = func(command *applications.Command) {
			client.AppArchiveReport(app, command)
		}
	}
	return nil
}

func JobsRegister() error {
	jobs, err := client.GetJobList(agentIP, agentPort)
	if err != nil {
		return err
	}
	for _, info := range jobs {
		logger.Debug("app register: ", info.GetName())
		config := map[string]interface{}{
			"id":         info.GetId(),
			"name":       info.GetName(),
			"dir":        info.GetDir(),
			"program":    info.GetProgram(),
			"args":       info.GetArgs(),
			"stdout":     info.GetStdOut(),
			"stderr":     info.GetStdErr(),
			"spec":       info.GetSpec(),
			"timeout":    info.GetTimeout(),
			"is_monitor": info.GetIsMonitor(),
		}
		job, err := applications.JobRegister(info.GetId(), config)
		if err != nil {
			return err
		}
		job.MonitorReport = func(monitor *applications.Monitor) {
			client.JobMonitorReport(job, monitor)
		}
		job.ArchiveReport = func(command *applications.Command) {
			client.JobArchiveReport(job, command)
		}
	}
	return nil
}
