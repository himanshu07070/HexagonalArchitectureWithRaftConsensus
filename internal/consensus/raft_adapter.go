package consensus

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"go.etcd.io/raft/v3"
	"go.etcd.io/raft/v3/raftpb"
)

type RaftConsensusAdapter struct {
	node    raft.Node
	storage *raft.MemoryStorage
}

func NewRaftConsensusAdapter(nodeID int, peers []raft.Peer) *RaftConsensusAdapter {
	config := raft.Config{
		ID:              uint64(nodeID),
		ElectionTick:    10,
		HeartbeatTick:   1,
		Storage:         raft.NewMemoryStorage(),
		MaxSizePerMsg:   4096,
		MaxInflightMsgs: 256,
	}

	node := raft.StartNode(&config, peers)
	adapter := &RaftConsensusAdapter{
		node:    node,
		storage: config.Storage.(*raft.MemoryStorage),
	}
	// Start the Raft node
	go adapter.serveRaft()

	return adapter
}

func (r *RaftConsensusAdapter) UpdateDatabase(fileName string, fileSize int64) {
	// Create a Raft log entry with the update
	entry := raftpb.Entry{
		Term: r.node.Status().Term,
		Type: raftpb.EntryNormal,
		Data: []byte(fmt.Sprintf("%s:%d", fileName, fileSize)),
	}
	entryBytes, _ := json.Marshal(entry)
	// Propose the entry to the Raft node
	err := r.node.Propose(context.TODO(), entryBytes)
	if err != nil {
		log.Printf("Error proposing entry to Raft node: %v", err)
	}
}

func (r *RaftConsensusAdapter) serveRaft() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			r.node.Tick()
		case rd := <-r.node.Ready():
			// Handle Raft Ready events
			r.storage.Append(rd.Entries)

			for _, entry := range rd.CommittedEntries {
				if raftpb.EntryType(*entry.Type.Enum()) == raftpb.EntryNormal {
					// Extract the fileName and fileSize from the entry's data and update the database
					data := string(entry.Data)
					var fileName string
					var fileSize int64
					fmt.Sscanf(data, "%s:%d", &fileName, &fileSize)

					log.Printf("Updating database with fileName: %s, fileSize: %d", fileName, fileSize)
				}
			}

			// Advance the Raft node
			r.node.Advance()
		}
	}
}
