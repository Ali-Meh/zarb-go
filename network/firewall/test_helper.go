package firewall

import (
	"fmt"

	"github.com/libp2p/go-libp2p-core/peer"
)

func getFaultyPeerSelector(faultystore, faultyload bool) PeerSelector {
	return &faultyPeerSelector{
		faultyLoad:  faultyload,
		faultyStore: faultystore,
	}
}

func getFaultyFirewall(faultystore, faultyload bool) *Firewall {
	fw, _ := NewFirewall(getFaultyPeerSelector(faultystore, faultyload), nil)
	return fw
}

func getGoodFirewall() *Firewall {
	fw, _ := NewFirewall(nil, nil)
	return fw

}

type faultyPeerSelector struct {
	faultyStore bool
	faultyLoad  bool
}

func (fps *faultyPeerSelector) Store() error {
	if fps.faultyStore {
		return fmt.Errorf("couldn't save")
	}
	return nil
}
func (fps *faultyPeerSelector) Load() error {
	if fps.faultyLoad {
		return fmt.Errorf("couldn't load")
	}
	return nil
}
func (fps *faultyPeerSelector) CanConnect(peer peer.AddrInfo) bool { return false }
