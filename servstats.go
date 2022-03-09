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

type ServStatsMsgNode struct {
	// ID    int
	Start time.Time
	End   time.Time
}

type ServStatsMsg struct {
	Name         string    `json:"name,omitempty"`
	TotalTime    float64   `json:"totalTime,omitempty"`
	TotalTimes   int       `json:"totalTimes,omitempty"`
	MaxTime      float64   `json:"maxTime,omitempty"`
	MinTime      float64   `json:"minTime,omitempty"`
	MaxParallels int       `json:"maxParallels,omitempty"`
	Nodes        []float64 `json:"nodes,omitempty"`
	pool         sync.Pool `json:"-"`
	// NoEndNodes   map[int]*ServStatsMsgNode `json:"noEndNodes,omitempty"`
	// pool         []*ServStatsMsgNode       `json:"-"`
	// lock         sync.Mutex                `json:"-"`
	// poolSize     int                       `json:"-"`
	// curID        int                       `json:"-"`
}

func newServStatsMsg(name string, poolSize int) *ServStatsMsg {
	msg := &ServStatsMsg{
		Name:    name,
		MinTime: math.MaxFloat64,
		// NoEndNodes: make(map[int]*ServStatsMsgNode),
		// poolSize:   poolSize,
	}

	msg.pool = sync.Pool{
		New: func() interface{} {
			return &ServStatsMsgNode{
				// ID: msg.curID,
			}
		},
	}

	return msg
}

// func (msg *ServStatsMsg) _newNode() *ServStatsMsgNode {
// 	if len(msg.pool) <= 0 {
// 		for i := 0; i < msg.poolSize; i++ {
// 			msg.curID++

// 			msg.pool = append(msg.pool, &ServStatsMsgNode{
// 				ID: msg.curID,
// 			})
// 		}
// 	}

// 	node := msg.pool[len(msg.pool)-1]
// 	msg.pool = msg.pool[:(len(msg.pool) - 1)]
// 	return node
// }

func (msg *ServStatsMsg) startMsg() *ServStatsMsgNode {
	// msg.lock.Lock()
	// n := msg._newNode()

	// n.Start = time.Now()

	// msg.NoEndNodes[n.ID] = n

	// if len(msg.NoEndNodes) > msg.MaxParallels {
	// 	msg.MaxParallels = len(msg.NoEndNodes)
	// }

	// msg.lock.Unlock()

	n := msg.pool.Get().(*ServStatsMsgNode)
	n.Start = time.Now()

	return n
}

func (msg *ServStatsMsg) endMsg(node *ServStatsMsgNode, maxNodes int) {
	// msg.lock.Lock()
	// _, isok := msg.NoEndNodes[node.ID]
	// if !isok {
	// 	Warn("ServStatsMsg:EndMsg",
	// 		zap.Error(ErrNoMsgCtx))

	// 	msg.lock.Unlock()

	// 	return
	// }

	node.End = time.Now()

	dt := node.End.Sub(node.Start).Seconds()
	if dt > msg.MaxTime {
		msg.MaxTime = dt
	}

	if dt < msg.MinTime {
		msg.MinTime = dt
	}

	msg.TotalTime += dt
	msg.TotalTimes++

	msg.Nodes = append(msg.Nodes, dt)
	if len(msg.Nodes) > maxNodes*2 {
		sort.Slice(msg.Nodes, func(i, j int) bool {
			return msg.Nodes[i] > msg.Nodes[j]
		})

		msg.Nodes = msg.Nodes[:maxNodes]
	}

	msg.pool.Put(node)

	// delete(msg.NoEndNodes, node.ID)
	// msg.pool = append(msg.pool, node)

	// msg.lock.Unlock()
}

type ServStats struct {
	MapMsgs      map[string]*ServStatsMsg `json:"mapMsgs,omitempty"`
	MaxNodes     int                      `json:"-"`
	ChanState    chan int                 `json:"-"`
	TickerOutput *time.Ticker             `json:"-"`
	PathOutput   string                   `json:"-"`
	poolSize     int                      `json:"-"`
}

func NewServStats(maxNodes int, chanSize int, outputTimer time.Duration, pathOutput string, poolSize int) *ServStats {
	stats := &ServStats{
		MapMsgs:      make(map[string]*ServStatsMsg),
		MaxNodes:     maxNodes,
		ChanState:    make(chan int),
		TickerOutput: time.NewTicker(outputTimer),
		PathOutput:   pathOutput,
		poolSize:     poolSize,
	}

	return stats
}

func (stats *ServStats) Start() {
	go stats.mainLoop()
}

func (stats *ServStats) Stop() {
	stats.TickerOutput.Stop()
	stats.ChanState <- 0
}

func (stats *ServStats) StartMsg(msgname string) *ServStatsMsgNode {
	msg, isok := stats.MapMsgs[msgname]
	if isok {
		return msg.startMsg()
	}

	return nil
}

func (stats *ServStats) EndMsg(msgname string, node *ServStatsMsgNode) {
	msg, isok := stats.MapMsgs[msgname]
	if isok {
		msg.endMsg(node, stats.MaxNodes)
	}
}

func (stats *ServStats) RegMsg(name string) {
	stats.MapMsgs[name] = newServStatsMsg(name, stats.poolSize)
}

func (stats *ServStats) Output() {
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	// for _, v := range stats.MapMsgs {
	// 	v.lock.Lock()
	// }

	b, err := json.Marshal(stats)
	if err != nil {
		Warn("ServStats.output:Marshal",
			zap.Error(err))

		// for _, v := range stats.MapMsgs {
		// 	v.lock.Unlock()
		// }

		return
	}

	// for _, v := range stats.MapMsgs {
	// 	v.lock.Unlock()
	// }

	err = os.WriteFile(path.Join(stats.PathOutput, fmt.Sprintf("%v.json", time.Now().Unix())), b, 0644)
	if err != nil {
		Warn("ServStats.output:WriteFile",
			zap.Error(err))

		return
	}
}

func (stats *ServStats) mainLoop() {
	for {
		select {
		case <-stats.ChanState:
			Info("ServStats:mainLoop:ChanState")

			return
		case <-stats.TickerOutput.C:
			stats.Output()
		}
	}
}
