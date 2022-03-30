package goutils

import (
	"fmt"
	"math"
	"os"
	"path"
	"sort"
	"sync"
	"time"

	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

type SenderStatsNode struct {
	Name       string `json:"name,omitempty"`
	TotalBytes int64  `json:"totalBytes,omitempty"`
	TotalTimes int    `json:"totalTimes,omitempty"`
	MaxBytes   int    `json:"maxBytes,omitempty"`
	MinBytes   int    `json:"minBytes,omitempty"`
	Nodes      []int  `json:"nodes,omitempty"`
}

func newSenderStatsNode(name string) *SenderStatsNode {
	msg := &SenderStatsNode{
		Name:     name,
		MaxBytes: 0,
		MinBytes: math.MaxInt32,
	}

	return msg
}

func (node *SenderStatsNode) push(bytes int, maxNodes int) {
	if bytes > node.MaxBytes {
		node.MaxBytes = bytes
	}

	if bytes < node.MinBytes {
		node.MinBytes = bytes
	}

	node.TotalTimes++
	node.TotalBytes += int64(bytes)
	node.Nodes = append(node.Nodes, bytes)

	if len(node.Nodes) > maxNodes*2 {
		sort.Slice(node.Nodes, func(i, j int) bool {
			return node.Nodes[i] > node.Nodes[j]
		})

		node.Nodes = node.Nodes[:maxNodes]
	}
}

func (node *SenderStatsNode) sort() {
	sort.Slice(node.Nodes, func(i, j int) bool {
		return node.Nodes[i] > node.Nodes[j]
	})
}

type SenderStats struct {
	MapNodes     map[string]*SenderStatsNode `json:"mapMsgs,omitempty"`
	MaxNodes     int                         `json:"-"`
	ChanState    chan int                    `json:"-"`
	TickerOutput *time.Ticker                `json:"-"`
	PathOutput   string                      `json:"-"`
	prefixFN     string                      `json:"-"`
	lock         sync.Mutex                  `json:"-"`
}

func NewSenderStats(maxNodes int, chanSize int, outputTimer time.Duration, pathOutput string, prefixFN string) *SenderStats {
	stats := &SenderStats{
		MapNodes:     make(map[string]*SenderStatsNode),
		MaxNodes:     maxNodes,
		ChanState:    make(chan int),
		TickerOutput: time.NewTicker(outputTimer),
		PathOutput:   pathOutput,
		prefixFN:     prefixFN,
	}

	return stats
}

func (stats *SenderStats) Start() {
	go stats.mainLoop()
}

func (stats *SenderStats) Stop() {
	stats.TickerOutput.Stop()
	stats.ChanState <- 0
}

func (stats *SenderStats) Push(name string, bytes int) {
	stats.lock.Lock()
	defer stats.lock.Unlock()

	node, isok := stats.MapNodes[name]
	if isok {
		node.push(bytes, stats.MaxNodes)
	} else {
		stats.MapNodes[name] = newSenderStatsNode(name)
		stats.MapNodes[name].push(bytes, stats.MaxNodes)
	}
}

func (stats *SenderStats) Output() {
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	stats.lock.Lock()
	for _, v := range stats.MapNodes {
		v.sort()
	}

	b, err := json.Marshal(stats)
	if err != nil {
		stats.lock.Unlock()

		Warn("SenderStats.output:Marshal",
			zap.Error(err))

		return
	}
	stats.lock.Unlock()

	err = os.WriteFile(path.Join(stats.PathOutput, fmt.Sprintf("%v.%v.json", stats.prefixFN, time.Now().Unix())), b, 0644)
	if err != nil {
		Warn("SenderStats.output:WriteFile",
			zap.Error(err))

		return
	}
}

func (stats *SenderStats) mainLoop() {
	for {
		select {
		case <-stats.ChanState:
			Info("SenderStats:mainLoop:ChanState")

			return
		case <-stats.TickerOutput.C:
			stats.Output()
		}
	}
}
