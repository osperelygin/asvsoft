package communication

import (
	"asvsoft/internal/pkg/logger"
	"asvsoft/internal/pkg/proto"
	"asvsoft/internal/pkg/utils"
	"fmt"
	"io"
)

type Receiver struct {
	rwc          io.ReadWriteCloser
	moduleID     proto.ModuleID
	sync         bool
	chunkSize    int
	retriesLimit int
	log          logger.Logger
}

func NewReceiver(rwc io.ReadWriteCloser, moduleID proto.ModuleID) *Receiver {
	return &Receiver{
		rwc:          rwc,
		moduleID:     moduleID,
		log:          logger.DummyLogger{},
		chunkSize:    DefaultChunkSize,
		retriesLimit: DefaultRetriesLimit,
	}
}

func (r *Receiver) WithChunkSize(chunkSize int) *Receiver {
	r.chunkSize = chunkSize
	return r
}

func (r *Receiver) WithRetriesLimit(retriesLimit int) *Receiver {
	r.retriesLimit = retriesLimit
	return r
}

func (r *Receiver) WithSync(sync bool) *Receiver {
	r.sync = sync
	return r
}

func (r *Receiver) WithLogger(log logger.Logger) *Receiver {
	r.log = log
	return r
}

// Receive читает данный из r.rc и распаковывает пакет в сообщение.
func (r *Receiver) Receive() (proto.Message, error) {
	msg, err := r.receive()
	if err != nil {
		return msg, err
	}

	// TODO: унести в отдельный метод
	if msg.ModuleID == proto.CameraModuleID && msg.MsgID == proto.WritingModeB {
		payload, ok := msg.Payload.(*proto.CameraData)
		if !ok {
			return msg, fmt.Errorf("failed to handle chunked message: unexpected type")
		}

		rawImage := make([]byte, 0, int(payload.TotalChunckes)*r.chunkSize)
		rawImage = append(rawImage, payload.RawImagePart...)

		for payload.CurrentChunck < payload.TotalChunckes {
			msg, err := r.receive()
			if err != nil {
				return msg, err
			}

			r.log.Debugf("received chunked message: %s", msg)

			payload, ok = msg.Payload.(*proto.CameraData)
			if !ok {
				return msg, fmt.Errorf("failed to handle chunked message: unexpected type")
			}

			rawImage = append(rawImage, payload.RawImagePart...)
		}

		msg.Payload = &proto.CameraData{RawImagePart: rawImage}
		msg.CheckSum = 0
		msg.PayloadSize = 0
	}

	return msg, err
}

func (r *Receiver) receive() (proto.Message, error) {
	var (
		msg proto.Message
		err error
	)

	err = utils.RunWithRetries(func() error {
		rawData, err := proto.Read(r.rwc)
		if err != nil {
			return fmt.Errorf("read msg failed: %v", err)
		}

		r.log.Debugf("raw received msg: %+v", rawData)

		err = msg.Unmarshal(rawData)
		if err != nil {
			_ = r.sendMsg(proto.ResponseFail)

			return fmt.Errorf("unmarshal msg failed: %v", err)
		}

		if msg.MsgID != proto.SyncRequest {
			_ = r.sendMsg(proto.ResponseOK)
		}

		return nil
	}, r.log, r.retriesLimit, 0)

	return msg, err
}

func (r *Receiver) sendMsg(msgID proto.MessageID) error {
	if !r.sync {
		return nil
	}

	rawResp, err := proto.NewMessage(r.moduleID, msgID, nil).Marshal()
	if err != nil {
		err = fmt.Errorf("failed to marshal response: %w", err)
		r.log.Errorf("%v", err)

		return err
	}

	_, err = r.rwc.Write(rawResp)
	if err != nil {
		err = fmt.Errorf("failed to write response: %w", err)
		r.log.Errorf("%v", err)

		return err
	}

	r.log.Debugf("successfully sent %d msg", msgID)

	return nil
}

func (r *Receiver) Close() error {
	return r.rwc.Close()
}
