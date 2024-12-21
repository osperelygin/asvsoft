package communication

import (
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

type Syncer struct {
	moduleID proto.ModuleID
	rw       io.ReadWriter
}

func (s *Syncer) WithReadWriter(rw io.ReadWriter) *Syncer {
	s.rw = rw
	return s
}

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

		log.Tracef("start stamp set to %d", startStamp)
		proto.SetStartStamp(uint32(startStamp))

		return nil
	}

	var req proto.Message

	b, err := req.Marshal(nil, s.moduleID, proto.SyncRequest)
	if err != nil {
		return fmt.Errorf("cannot marshal msg: %w", err)
	}

	log.Traceln("writing sync request...")

	_, err = s.rw.Write(b)
	if err != nil {
		return fmt.Errorf("cannot write measures: %w", err)
	}

	log.Debugf("raw sync request: %+v", b)
	log.Debugf("sync request: %+v", req)

	time.Sleep(500 * time.Millisecond)

	log.Traceln("reading sync response...")

	rawResp, err := proto.Read(s.rw)
	if err != nil {
		return fmt.Errorf("cannot read response: %w", err)
	}

	var resp proto.Message

	err = resp.Unmarshal(rawResp)
	if err != nil {
		return fmt.Errorf("unmarshal msg failed: %v", err)
	}

	log.Debugf("raw sync response: %+v", rawResp)
	log.Debugf("sync response: %+v", resp)

	if resp.MsgID != proto.SyncResponse {
		return fmt.Errorf("unexpected msgID: %#X", resp.MsgID)
	}

	startStamp := resp.Payload.(uint32)
	proto.SetStartStamp(startStamp)

	log.Infof("time was synced, start stamp: %d", startStamp)

	return nil
}

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

	if req.MsgID != proto.SyncRequest {
		return fmt.Errorf("unexpected msgID: %#X", req.MsgID)
	}

	var resp proto.Message

	b, err := resp.Marshal(proto.GetStartStamp(), s.moduleID, proto.SyncResponse)
	if err != nil {
		return fmt.Errorf("cannot marshal resp: %w", err)
	}

	time.Sleep(time.Second)

	log.Traceln("writing sync response...")

	_, err = s.rw.Write(b)
	if err != nil {
		return fmt.Errorf("cannot write resp: %w", err)
	}

	log.Debugf("raw sync response: %+v", b)
	log.Debugf("sync response: %+v", resp)

	return nil
}
