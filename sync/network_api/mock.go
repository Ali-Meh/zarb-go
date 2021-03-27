package network_api

import (
	"fmt"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/sync/firewall"
	"github.com/zarbchain/zarb-go/sync/message"
	"github.com/zarbchain/zarb-go/sync/message/payload"
)

type errorFunc func() error

type MockNetworkAPI struct {
	ch             chan *message.Message
	id             peer.ID
	Firewall       *firewall.Firewall
	ParsFn         ParsMessageFn
	OtherAPI       *MockNetworkAPI
	DropPeerFn     errorFunc
	AllowFromValue bool
}

func MockingNetworkAPI(id peer.ID, opts ...errorFunc) *MockNetworkAPI {
	return &MockNetworkAPI{
		ch:         make(chan *message.Message, 100),
		id:         id,
		DropPeerFn: func() error { return nil },
	}
}
func (mock *MockNetworkAPI) Start() error {
	return nil
}
func (mock *MockNetworkAPI) Stop() {
}
func (mock *MockNetworkAPI) JoinDownloadTopic() error {
	return nil
}
func (mock *MockNetworkAPI) LeaveDownloadTopic() {}
func (mock *MockNetworkAPI) PublishMessage(msg *message.Message) error {
	mock.ch <- msg
	return nil
}
func (mock *MockNetworkAPI) SelfID() peer.ID {
	return mock.id
}

func (mock *MockNetworkAPI) DropPeer(peerId peer.ID) error {
	return mock.DropPeerFn()
}

func (mock *MockNetworkAPI) AllowFrom(peerId peer.ID) bool {
	return mock.AllowFromValue
}

func (mock *MockNetworkAPI) CheckAndParsMessage(msg *message.Message, id peer.ID) bool {
	d, _ := msg.Encode()
	msg2 := mock.Firewall.ParsMessage(d, id)
	if msg2 != nil {
		mock.ParsFn(msg2, mock.id)
		return true
	}
	return false
}

func (mock *MockNetworkAPI) sendMessageToOtherPeer(m *message.Message) {
	mock.OtherAPI.CheckAndParsMessage(m, mock.id)
}

func (mock *MockNetworkAPI) ShouldPublishMessageWithThisType(t *testing.T, payloadType payload.PayloadType) *message.Message {
	timeout := time.NewTimer(2 * time.Second)

	for {
		select {
		case <-timeout.C:
			require.NoError(t, fmt.Errorf("Timeout"))
			return nil
		case msg := <-mock.ch:
			logger.Info("shouldPublishMessageWithThisType", "id", mock.id, "msg", msg, "type", payloadType.String())
			mock.sendMessageToOtherPeer(msg)

			if msg.PayloadType() == payloadType {
				return msg
			}
		}
	}
}

func (mock *MockNetworkAPI) ShouldNotPublishMessageWithThisType(t *testing.T, payloadType payload.PayloadType) {
	timeout := time.NewTimer(300 * time.Millisecond)

	for {
		select {
		case <-timeout.C:
			return
		case msg := <-mock.ch:
			logger.Info("shouldNotPublishMessageWithThisType", "id", mock.id, "msg", msg, "type", payloadType.String())
			mock.sendMessageToOtherPeer(msg)

			assert.NotEqual(t, msg.PayloadType(), payloadType)
		}
	}
}
