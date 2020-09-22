package kademlia

import (
	"crypto/md5"
	"fmt"

	log "github.com/sirupsen/logrus"
)

type Node struct {
	RT      *RoutingTable
	network Network
	content map[string]string
}

// InitNode initializes the Kademlia Node
// with a Routing Table and a Network
func (kademlia *Node) InitNode(id *NodeID) {
	kademlia.network = NewNetwork(kademlia)
	ip := kademlia.network.ip
	go kademlia.network.Listen(ip, "8080")

	me := NewContact(id, ip+":8080")
	rendezvousNode := NewContact(NewNodeID("00000000000000000000000000000000FFFFFFFF"), "10.0.8.3:8080")
	kademlia.RT = NewRoutingTable(me)
	kademlia.RT.AddContact(rendezvousNode)

	kademlia.content = make(map[string]string)
}

func (kademlia *Node) NodeLookup(target *Contact) {
	table := kademlia.RT.FindClosestContacts(target.ID, BucketSize)
	fmt.Println("table = ", table)
}

func (kademlia *Node) FindValue(hash string) {
	md5 := md5.Sum([]byte(hash))
	var content = kademlia.content[string(md5[:])]
	if content == "" {
		fmt.Println("Content not found!")
	} else {
		fmt.Println("content = ", content)
	}
}

func (kademlia *Node) StoreValue(data string) {
	md5 := md5.Sum([]byte(data))
	kademlia.content[string(md5[:])] = data
}

func (kademlia *Node) Ping() {
	target := &kademlia.RT.FindClosestContacts(kademlia.RT.me.ID, BucketSize)[0]
	rpc, err := kademlia.network.SendPingMessage(target, &kademlia.RT.me)

	if err != nil {
		log.Warn(err)
		kademlia.RT.RemoveContact(*target)
	} else if *rpc.Type == "OK" {
		kademlia.RT.AddContact(*target)
	}
}

// SearchStore looks for a value in the node's store. Returns the value
// if found else nil.
func (kademlia *Node) SearchStore(key string) *string {
	value, exists := kademlia.content[key]
	if exists {
		return nil
	}
	return &value
}