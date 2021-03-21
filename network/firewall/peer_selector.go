package firewall

import (
	"github.com/libp2p/go-libp2p-core/peer"
)

//PeerSelector filters wether or not a peer is allowd and trusted to be connected or not
//its main firewall interface on p2p layer to be implemented
type PeerSelector interface {
	//check to see if Peer is allowed to connect
	CanConnect(peer.AddrInfo) bool
	//Store Selection related Configuration
	Store() error
	//Load Selection related Configuration
	Load() error
}

type DefaultSelector struct {
	trustedNodes *peer.Set
}

func NewSelector( /* , logger *logger.Logger */ ) *DefaultSelector {
	return &DefaultSelector{
		trustedNodes: peer.NewSet(),
	}
}

func (ps *DefaultSelector) CanConnect(peer peer.AddrInfo) bool {
	// logger.Trace("Checking if %s Can Connect", peer.ID.String())
	if ps.trustedNodes.Size() > 0 {
		// logger.Info("id donna", len(fw.trustedNodes.Peers()))
		return ps.trustedNodes.Contains(peer.ID)
	}
	return true
}
func (ps *DefaultSelector) Store() error {
	return nil

}

func (ps *DefaultSelector) Load() error {
	return nil
}
