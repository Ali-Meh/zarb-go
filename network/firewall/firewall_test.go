package firewall

import (
	"context"
	"fmt"
	"testing"

	"github.com/libp2p/go-libp2p-core/network"
	mocknet "github.com/libp2p/go-libp2p/p2p/net/mock"
	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/logger"
)

func TestMain(m *testing.M) {
	logger.InitLogger(&logger.Config{})
	m.Run()
}
func TestFirewall_Connected(t *testing.T) {
	ctx := context.Background()
	fw, _ := NewFirewall(nil, nil)
	net, _ := mocknet.WithNPeers(ctx, 5)
	host := net.Host(net.Peers()[0])
	allowedPeer, _ := net.GenPeer()
	forbidenPeer, _ := net.GenPeer()

	//add firewall to host
	host.Network().Notify(fw)
	fw.ps.(*DefaultSelector).trustedNodes.TryAdd(forbidenPeer.ID())

	net.LinkAll()
	//check if allowed peer can connect
	err := net.ConnectAllButSelf()
	assert.Nil(t, err)
	err = host.Connect(ctx, allowedPeer.Peerstore().PeerInfo(allowedPeer.ID()))
	assert.Nil(t, err)
	assert.Equal(t, network.Connected, host.Network().Connectedness(allowedPeer.ID()))

	//check if forbiden peer can connect
	err = host.Connect(ctx, forbidenPeer.Peerstore().PeerInfo(forbidenPeer.ID()))
	fmt.Println(err)
	assert.Equal(t, network.NotConnected, host.Network().Connectedness(forbidenPeer.ID()))
}
func TestFirewall_Close(t *testing.T) {
	tests := []struct {
		name    string
		fw      *Firewall
		wantErr bool
	}{
		{
			name:    "faulty firewall save want error",
			fw:      getFaultyFirewall(true, false),
			wantErr: true,
		},
		{
			name:    "faulty firewall load want error",
			fw:      getFaultyFirewall(false, true),
			wantErr: false,
		},
		{
			name:    "good firewall should save without error",
			fw:      getGoodFirewall(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.fw.Stop(); (err != nil) != tt.wantErr {
				t.Errorf("Firewall.Stop() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
