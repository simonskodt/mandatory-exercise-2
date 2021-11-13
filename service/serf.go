package service

import (
	"fmt"
	"github.com/hashicorp/serf/serf"
	"log"
	"strconv"
	"utils"
)

type SerfService struct {
	bindPort int
	clusterPort int
	cluster *serf.Serf
	eventCh chan serf.Event
	logger *utils.Logger
}

func (ss *SerfService) Init(address string, ch chan string) {
	err := ss.joinCluster(address)
	if err != nil {
		ss.logger.ErrorFatalf("Could not join cluster with Serf. :: %v", err)
	}

	go func() {
		for {
			select {
			case e := <-ss.eventCh:
				if e.EventType() == serf.EventMemberJoin {
					m := e.(serf.MemberEvent)

					for _, member := range m.Members {
						ch <- fmt.Sprintf("[%v] %v:%v\n", member.Name, member.Addr.String(), member.Port)
					}
				}
			}
		}
	}()
}

func (ss *SerfService) joinCluster(address string) error {
	_, err := ss.cluster.Join([]string{address + ":" + strconv.Itoa(ss.clusterPort)}, true)
	if err != nil {
		return err
	}

	return nil
}

func createConfig(nodeName string, nodeAddress string, serfPort int, eventCh chan serf.Event) *serf.Config {
	conf := serf.DefaultConfig()
	conf.Init()
	//conf.LogOutput = os.Stdin
	//conf.MemberlistConfig.LogOutput = os.Stdin
	conf.MemberlistConfig.AdvertiseAddr = nodeAddress
	conf.MemberlistConfig.AdvertisePort = serfPort
	conf.MemberlistConfig.BindAddr = nodeAddress
	conf.MemberlistConfig.BindPort = serfPort
	conf.EventCh = eventCh
	conf.MemberlistConfig.ProtocolVersion = 3
	conf.NodeName = nodeName

	return conf
}

func NewSerfService(nodeName string, nodeAddress string, serfPort int, clusterPort int) *SerfService {
	eventCh := make(chan serf.Event, 64)

	// Setup cluster configuration
	config := createConfig(nodeName, nodeAddress, serfPort, eventCh)

	// Create cluster
	cluster, err := serf.Create(config)
	if err != nil {
		log.Println(err)
		return nil
	}

	return &SerfService{
		bindPort: serfPort,
		clusterPort: clusterPort,
		cluster: cluster,
		eventCh: eventCh,
	}
}
