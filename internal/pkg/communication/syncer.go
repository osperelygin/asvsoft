package communication

import (
	"asvsoft/internal/pkg/logger"
	"asvsoft/internal/pkg/proto"
	"asvsoft/internal/pkg/utils"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	DefaultSyncerRetries = 10
	DefaultSyncerSleep   = 500 * time.Millisecond
)

func NewSyncer(moduleID proto.ModuleID) *Syncer {
	return &Syncer{
		moduleID: moduleID,
		retries:  DefaultSyncerRetries,
		sleep:    DefaultSyncerSleep,
	}
}

// Syncer осуществляется синхронизацию системного времени между модулем и контроллером управления.
// Системное время - время в мс, отсчитываемое от момента старта работы контроллера управления
// или полученное из переменной окружения START_STAMP.
type Syncer struct {
	moduleID proto.ModuleID
	rw       io.ReadWriter
	retries  int
	sleep    time.Duration
}

func (s *Syncer) WithReadWriter(rw io.ReadWriter) *Syncer {
	s.rw = rw
	return s
}

func (s *Syncer) WithRetries(retries int) *Syncer {
	s.retries = retries
	return s
}

func (s *Syncer) WithSleep(sleep time.Duration) *Syncer {
	s.sleep = sleep
	return s
}

// SyncSystemTime осуществляет синхронизацию системного времении. При установленом s.rw отправляет
// запрос синхронизации и назначает систменое время из полученного ответа. При неустановленном
// s.rw пытается получить начало отсчета системношго времи из START_STAMP, если переменная
// не установлена, то назчает в качестве начало отсчета время запуска утилиты на модуле.
func (s *Syncer) SyncSystemTime() error {
	setStartStampFn := func(startStamp uint32, srcStartStamp string) {
		proto.SetStartStamp(startStamp)
		log.Infof("time was synced from %s, start stamp: %d", srcStartStamp, startStamp)
	}

	if s.rw == nil {
		startStampStr := os.Getenv("START_STAMP")
		if startStampStr == "" {
			setStartStampFn(proto.GetStartStamp(), "cli_start_time")
			return nil
		}

		startStampInt, err := strconv.ParseInt(startStampStr, 10, 64)
		if err != nil {
			return fmt.Errorf("bad start stamp from env var: %w", err)
		}

		setStartStampFn(uint32(startStampInt), "env_var")

		return nil
	}

	b, err := proto.NewMessage(s.moduleID, proto.SyncRequest, nil).Marshal()
	if err != nil {
		return fmt.Errorf("cannot marshal msg: %w", err)
	}

	var startStamp uint32

	err = utils.RunWithRetries(func() error {
		_, err := s.rw.Write(b)
		if err != nil {
			return fmt.Errorf("cannot write measures: %v", err)
		}

		rawResp, err := proto.Read(s.rw)
		if err != nil {
			return fmt.Errorf("cannot read response: %v", err)
		}

		var resp proto.Message

		err = resp.Unmarshal(rawResp)
		if err != nil {
			return fmt.Errorf("unmarshal msg failed: %v", err)
		}

		if resp.MsgID != proto.SyncResponse {
			return fmt.Errorf("unexpected msgID: %#X", resp.MsgID)
		}

		d, ok := resp.Payload.(*proto.SyncData)
		if !ok {
			return fmt.Errorf("unexpected payload type")
		}

		startStamp = uint32(*d)

		return nil
	}, logger.Wrap(log.StandardLogger(), "[syncer]"), s.retries, s.sleep)
	if err != nil {
		return fmt.Errorf("failed to get start time from remote: %w", err)
	}

	setStartStampFn(startStamp, "remote")

	return nil
}

func (s *Syncer) ProcessSyncRequest(req proto.Message) (*proto.Message, error) {
	startStamp := proto.SyncData(proto.GetStartStamp())

	resp := proto.NewMessage(s.moduleID, proto.SyncResponse, &startStamp)

	if req.MsgID != proto.SyncRequest {
		return resp, fmt.Errorf("unexpected msgID: %#X", req.MsgID)
	}

	b, err := resp.Marshal()
	if err != nil {
		return resp, fmt.Errorf("cannot marshal resp: %w", err)
	}

	time.Sleep(time.Second)

	log.Traceln("writing sync response...")

	_, err = s.rw.Write(b)
	if err != nil {
		return resp, fmt.Errorf("cannot write resp: %w", err)
	}

	log.Debugf("raw sync response: %+v", b)

	return resp, nil
}
