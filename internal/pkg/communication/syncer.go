package communication

import (
	"asvsoft/internal/pkg/logger"
	"asvsoft/internal/pkg/proto"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

func NewSyncer(moduleID proto.ModuleID) *Syncer {
	return &Syncer{moduleID: moduleID}
}

// Syncer осуществляется синхронизацию системного времени между модулем и контроллером управления.
// Системное время - время в мс, отсчитываемое от момента старта работы контроллера управления
// или полученное из переменной окружения START_STAMP.
type Syncer struct {
	moduleID proto.ModuleID
	rw       io.ReadWriter
}

func (s *Syncer) WithReadWriter(rw io.ReadWriter) *Syncer {
	s.rw = rw
	return s
}

// Sync осуществляет синхронизацию системного времении. При установленом s.rw отправляет
// запрос синхронизации и назначает систменое время из полученного ответа. При неустановленном
// s.rw пытается получить начало отсчета системношго времи из START_STAMP, если переменная
// не установлена, то назчает в качестве начало отсчета время запуска утилиты на модуле.
func (s *Syncer) Sync() error {
	if s.rw == nil {
		log.Traceln("rw == nil: try getting START_STAMP env var")

		startStampStr := os.Getenv("START_STAMP")
		if startStampStr == "" {
			log.Tracef("START_STAMP is not set, use CLI start stamp: %d", proto.GetStartStamp())
			return nil
		}

		startStamp, err := strconv.ParseInt(startStampStr, 10, 64)
		if err != nil {
			return fmt.Errorf("cannot parse start stamp %q: %w", startStampStr, err)
		}

		proto.SetStartStamp(uint32(startStamp))

		log.Tracef("start stamp set to %d", startStamp)

		return nil
	}

	var req proto.Message

	b, err := req.Marshal(nil, s.moduleID, proto.SyncRequest)
	if err != nil {
		return fmt.Errorf("cannot marshal msg: %w", err)
	}

	const syncRequestRetries = 10

	for retry := range syncRequestRetries {
		log := logger.Wrap(
			log.StandardLogger(),
			fmt.Sprintf("[retry #%d]", retry),
		)

		log.Tracef("writing sync request...")

		_, err = s.rw.Write(b)
		if err != nil {
			log.Errorf("cannot write measures: %v", err)
			continue
		}

		log.Debugf("raw sync request: %+v", b)
		log.Debugf("sync request: %+v", req)

		time.Sleep(500 * time.Millisecond)

		log.Tracef("reading sync response...")

		var rawResp []byte

		rawResp, err = proto.Read(s.rw)
		if err != nil {
			log.Errorf("cannot read response: %v", err)
			continue
		}

		var resp proto.Message

		err = resp.Unmarshal(rawResp)
		if err != nil {
			log.Errorf("unmarshal msg failed: %v", err)
			continue
		}

		log.Debugf("raw sync response: %+v", rawResp)
		log.Debugf("sync response: %+v", resp)

		if resp.MsgID != proto.SyncResponse {
			log.Errorf("unexpected msgID: %#X", resp.MsgID)
			continue
		}

		startStamp := resp.Payload.(uint32)
		proto.SetStartStamp(startStamp)

		log.Infof("time was synced, start stamp: %d", startStamp)

		break
	}

	if err != nil {
		return fmt.Errorf("failed to sync: %w", err)
	}

	return nil
}

// Serve ожидает запрос синхронизации, возвращает в ответе начало отсчета системного времени.
// Запускается только на контроллере управления.
func (s *Syncer) Serve() error {
	if s.rw == nil {
		log.Traceln("skipping serve: rw == nil")
		return nil
	}

	log.Traceln("reading sync request...")

	rawReq, err := proto.Read(s.rw)
	if err != nil {
		return fmt.Errorf("cannot read req: %w", err)
	}

	var req proto.Message

	err = req.Unmarshal(rawReq)
	if err != nil {
		return fmt.Errorf("unmarshal req failed: %v", err)
	}

	log.Debugf("raw sync request: %+v", rawReq)
	log.Debugf("sync request: %+v", req)

	_, err = s.ProcessSyncRequest(req)

	return err
}

func (s *Syncer) ProcessSyncRequest(req proto.Message) (proto.Message, error) {
	var resp proto.Message

	if req.MsgID != proto.SyncRequest {
		return resp, fmt.Errorf("unexpected msgID: %#X", req.MsgID)
	}

	b, err := resp.Marshal(proto.GetStartStamp(), s.moduleID, proto.SyncResponse)
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
