package firewall

import (
	"os"
	"testing"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/util"
)

func TestDefaultSelector_CanConnect(t *testing.T) {
	allowed_peer := util.RandomPeerID()
	selector := NewSelector()
	selector.trustedNodes.Add(allowed_peer)

	tests := []struct {
		name string
		fw   *DefaultSelector
		peer peer.AddrInfo
		want bool
	}{
		{
			name: "TrustedPeer empty",
			fw:   NewSelector(),
			want: true,
			peer: peer.AddrInfo{ID: allowed_peer},
		},
		{
			name: "TrustedPeer Contain Our Id",
			fw:   selector,
			want: true,
			peer: peer.AddrInfo{ID: allowed_peer},
		},
		{
			name: "TrustedPeer Doesn't Contain Our Id",
			fw:   selector,
			want: false,
			peer: peer.AddrInfo{ID: peer.ID("QmcCqNeNoUZZRmwLXZoCAumAEPCWHXodFUeVQQaQ6huNTh")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fw.CanConnect(tt.peer); got != tt.want {
				t.Errorf("%s => got %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestDefaultSelector_Store_Load(t *testing.T) {
	os.Remove(configFile)
	ps := NewSelector()
	err := ps.Load()
	assert.Nil(t, err, "Loading had issue")
	assert.Equal(t, 0, ps.trustedNodes.Size())
	ps.trustedNodes.Add(util.RandomPeerID())
	assert.Equal(t, 1, ps.trustedNodes.Size())
	err = ps.Store()
	assert.Nil(t, err, "Storing config Errored")
	ps2 := NewSelector()
	err = ps2.Load()
	assert.Nil(t, err, "2nd Loading had issue")
	assert.Equal(t, 1, ps2.trustedNodes.Size())
	assert.Equal(t, ps, ps2)
}
