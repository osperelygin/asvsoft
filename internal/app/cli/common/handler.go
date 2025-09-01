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
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type ModuleOptions struct {
	SendMode proto.MessageID
}

// ModuleHandler инициализирует sender и syncer, запускает sender.
func ModuleHandler(cfg *config.ModuleConfig, mode RunMode, opts ...ModuleOptions) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		ctx := config.WrapContext(cmd.Context(), cfg)

		sndr, sncr, err := Init(ctx, mode, opts...)
		if err != nil {
			return err
		}

		err = sncr.SyncSystemTime()
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

type module struct {
	rcvr *communication.Receiver
	sncr *communication.Syncer
}

// ControllerHandler ...
func ControllerHandler(moduleID proto.ModuleID, ctrlCfgPath *string) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		log := logger.Wrap(logrus.StandardLogger(), "[main]")

		ctrlCfg, err := config.NewControllerConfig(*ctrlCfgPath)
		if err != nil {
			return fmt.Errorf("failed to get controller config: %w", err)
		}

		moduleCfg := ctrlCfg.Modules
		modules := make(map[string]module, len(moduleCfg))

		for name, connCfg := range moduleCfg {
			if !connCfg.Enabled {
				continue
			}

			srcPort, err := serialport.New(connCfg.Listener.Short())
			if err != nil {
				return fmt.Errorf("cannot open serial port %s: %w", connCfg.Listener, err)
			}

			srcPort.SetLogger(logger.Wrap(logrus.StandardLogger(), fmt.Sprintf("[%s]", name)))

			log.Debugf("successfull create serail port: %s", connCfg.Listener)

			sncr := communication.NewSyncer(moduleID).WithReadWriter(srcPort)

			rcvr := communication.NewReceiver(srcPort, moduleID).
				WithSync(connCfg.Listener.Sync).
				WithChunkSize(connCfg.Listener.ChunkSize).
				WithRetriesLimit(connCfg.Listener.RetriesLimit)

			defer func() {
				err = rcvr.Close()
				if err != nil {
					log.Errorf("cannot close receiver: %v", err)
				}
			}()

			modules[name] = module{rcvr: rcvr, sncr: sncr}
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		closeCount := 0
		closeChannel := make(chan struct{}, len(modules))

		for moduleName, module := range modules {
			go receiving(ctx, moduleName, module, closeChannel)
		}

		quitChannel := make(chan os.Signal, 2)
		signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)

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
	log := logger.Wrap(logrus.StandardLogger(), fmt.Sprintf("[%s]", moduleName))

	log.Infof("starting receive message...")

	module.rcvr.WithLogger(log)

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

			// TODO: использовать общий подход к обработке сообщения каждого модуля
			if msg.ModuleID == proto.CameraModuleID && msg.MsgID == proto.WritingModeB {
				err = handleCameraRegistratorMsg(log, msg)
				if err != nil {
					log.Errorf("failed to handle camera message: %v", err)
					continue
				}
			}
		}
	}
}

func handleCameraRegistratorMsg(log logger.Logger, msg proto.Message) error {
	payload, ok := msg.Payload.(*proto.CameraData)
	if !ok {
		return fmt.Errorf("unexecpted message payload type")
	}

	fileName := fmt.Sprintf("camera_%d.jpeg", msg.SystemTime)

	err := os.WriteFile(fileName, payload.RawImagePart, 0666) // nolint:gosec
	if err != nil {
		return fmt.Errorf("failed to write image to file: %w", err)
	}

	log.Infof("successfully saved recieved camera image to %s", fileName)

	return nil
}
