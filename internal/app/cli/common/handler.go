package common

import (
	"asvsoft/internal/app/config"
	"asvsoft/internal/pkg/communication"
	"asvsoft/internal/pkg/logger"
	"asvsoft/internal/pkg/proto"
	serialport "asvsoft/internal/pkg/serial-port"
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	CtrlCfgPath string
)

type module struct {
	rcvr *communication.Receiver
	sncr *communication.Syncer
}

// ModuleHandler инициализирует sender и syncer, запускает sender.
func ModuleHandler(cfg *config.ModuleConfig, mode RunMode) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		ctx := config.WrapContext(cmd.Context(), cfg)

		sndr, sncr, err := Init(ctx, mode)
		if err != nil {
			return err
		}

		err = sncr.Sync()
		if err != nil {
			return fmt.Errorf("cannot sync: %v", err)
		}

		err = sndr.Start(ctx)
		if err != nil {
			return err
		}

		return nil
	}
}

// ControllerHandler ...
func ControllerHandler(moduleID proto.ModuleID) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		log := logger.Wrap(log.StandardLogger(), "[main]")

		ctrlCfg, err := config.NewControllerConfig(CtrlCfgPath)
		if err != nil {
			return fmt.Errorf("failed to get controller config: %w", err)
		}

		moduleCfg := ctrlCfg.Modules
		modules := make(map[string]module, len(moduleCfg))

		for name, connectionCfg := range moduleCfg {
			if !connectionCfg.Enabled {
				continue
			}

			srcPort, err := serialport.New(connectionCfg.Listener)
			if err != nil {
				return fmt.Errorf("cannot open serial port %s: %w", connectionCfg.Listener, err)
			}

			srcPort.SetLogger(logger.Wrap(logrus.StandardLogger(), fmt.Sprintf("[%s]", name)))

			log.Debugf("successfull create serail port: %s", connectionCfg.Listener)

			s := communication.NewSyncer(moduleID).WithReadWriter(srcPort)

			r := communication.NewReceiver(srcPort)
			defer func() {
				err = r.Close()
				if err != nil {
					log.Errorf("cannot close receiver: %v", err)
				}
			}()

			modules[name] = module{rcvr: r, sncr: s}
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		closeCount := 0
		closeChannel := make(chan struct{}, len(modules))

		for moduleName, module := range modules {
			go receiving(ctx, moduleName, module, closeChannel)
		}

		quitChannel := make(chan os.Signal, 2)
		signal.Notify(quitChannel, os.Kill, os.Interrupt)

		for {
			select {
			case signal := <-quitChannel:
				log.Infof("%s signal called, cancel operations", signal.String())
				cancel()
			case <-closeChannel:
				closeCount++
				if closeCount == len(modules) {
					return nil
				}
			default:
			}
		}
	}
}

func receiving(
	ctx context.Context,
	moduleName string,
	module module,
	closeChannel chan struct{},
) {
	log := logger.Wrap(log.StandardLogger(), fmt.Sprintf("[%s]", moduleName))

	log.Infof("starting receive message...")

	for {
		select {
		case <-ctx.Done():
			log.Infof("stop receiving, context done")
			closeChannel <- struct{}{}

			return
		default:
			msg, err := module.rcvr.Receive()
			if err != nil {
				log.Errorf("receive failed: %v", err)
				continue
			}

			log.Infof("received message: %v", msg)

			if msg.MsgID == proto.SyncRequest {
				resp, err := module.sncr.ProcessSyncRequest(msg)
				if err != nil {
					log.Errorf("failed to process sync request: %v", msg)
				}

				log.Infof("sent sync response: %+v", resp)

				continue
			}
		}
	}
}
