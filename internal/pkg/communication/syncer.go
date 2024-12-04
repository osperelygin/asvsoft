package communication

import (
	"asvsoft/internal/pkg/proto"
	"fmt"
	"io"

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
		return nil
	}

	var req proto.Message

	b, err := req.Marshal(nil, s.moduleID, proto.SyncRequest)
	if err != nil {
		return fmt.Errorf("cannot marshal msg: %w", err)
	}

	_, err = s.rw.Write(b)
	if err != nil {
		return fmt.Errorf("cannot write measures: %w", err)
	}

	log.Debugf("raw sync request: %+v", b)
	log.Infof("sync request: %+v", req)

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
	log.Infof("sync response: %+v", resp)

	if resp.MsgID != proto.SyncResponse {
		return fmt.Errorf("unexpected msgID: %#X", resp.MsgID)
	}

	startStamp := resp.Payload.(uint32)
	proto.SetStartStamp(startStamp)

	log.Debugf("synced system time, start stamp: %d (ms)", startStamp*1000)

	return nil
}

func (s *Syncer) Serve() error {
	if s.rw == nil {
		return nil
	}

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
	log.Infof("sync request: %+v", req)

	if req.MsgID != proto.SyncRequest {
		return fmt.Errorf("unexpected msgID: %#X", req.MsgID)
	}

	var resp proto.Message

	b, err := resp.Marshal(proto.GetStartStamp(), s.moduleID, proto.SyncResponse)
	if err != nil {
		return fmt.Errorf("cannot marshal resp: %w", err)
	}

	_, err = s.rw.Write(b)
	if err != nil {
		return fmt.Errorf("cannot write resp: %w", err)
	}

	log.Debugf("raw sync response: %+v", b)
	log.Infof("sync response: %+v", resp)

	return nil
}
