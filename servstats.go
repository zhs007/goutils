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
)

type ServStatsMsgNode struct {
	Start time.Time
	End   time.Time
}

type ServStatsMsg struct {
	Name         string     `json:"name,omitempty"`
	TotalTime    float64    `json:"totalTime,omitempty"`
	TotalTimes   int        `json:"totalTimes,omitempty"`
	MaxTime      float64    `json:"maxTime,omitempty"`
	MinTime      float64    `json:"minTime,omitempty"`
	MaxParallels int        `json:"maxParallels,omitempty"`
	Nodes        []float64  `json:"nodes,omitempty"`
	pool         sync.Pool  `json:"-"`
	lock         sync.Mutex `json:"-"`
	LastMsgNums  int        `json:"lastMsgNums"`
}

func newServStatsMsg(name string, poolSize int) *ServStatsMsg {
	msg := &ServStatsMsg{
		Name:    name,
		MaxTime: 0,
		MinTime: math.MaxFloat64,
	}

	msg.pool = sync.Pool{
		New: func() interface{} {
			return &ServStatsMsgNode{}
		},
	}

	return msg
}

func (msg *ServStatsMsg) startMsg() *ServStatsMsgNode {
	msg.lock.Lock()
	msg.LastMsgNums++
	if msg.LastMsgNums > msg.MaxParallels {
		msg.MaxParallels = msg.LastMsgNums
	}
	msg.lock.Unlock()

	n := msg.pool.Get().(*ServStatsMsgNode)
	n.Start = time.Now()

	return n
}

func (msg *ServStatsMsg) sort() {
	msg.lock.Lock()

	sort.Slice(msg.Nodes, func(i, j int) bool {
		return msg.Nodes[i] > msg.Nodes[j]
	})

	msg.lock.Unlock()
}

func (msg *ServStatsMsg) endMsg(node *ServStatsMsgNode, maxNodes int) {
	node.End = time.Now()

	dt := node.End.Sub(node.Start).Seconds()

	msg.lock.Lock()

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

	msg.LastMsgNums--

	msg.lock.Unlock()

	msg.pool.Put(node)
}

type ServStats struct {
	MapMsgs      map[string]*ServStatsMsg `json:"mapMsgs,omitempty"`
	MaxNodes     int                      `json:"-"`
	ChanState    chan int                 `json:"-"`
	TickerOutput *time.Ticker             `json:"-"`
	PathOutput   string                   `json:"-"`
	poolSize     int                      `json:"-"`
	prefixFN     string                   `json:"-"`
}

func NewServStats(maxNodes int, chanSize int, outputTimer time.Duration, pathOutput string, poolSize int, prefixFN string) *ServStats {
	stats := &ServStats{
		MapMsgs:      make(map[string]*ServStatsMsg),
		MaxNodes:     maxNodes,
		ChanState:    make(chan int),
		TickerOutput: time.NewTicker(outputTimer),
		PathOutput:   pathOutput,
		poolSize:     poolSize,
		prefixFN:     prefixFN,
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

	for _, v := range stats.MapMsgs {
		v.sort()
	}

	b, err := json.Marshal(stats)
	if err != nil {
		Warn("ServStats.output:Marshal",
			Err(err))

		return
	}

	err = os.WriteFile(path.Join(stats.PathOutput, fmt.Sprintf("%v.%v.json", stats.prefixFN, time.Now().Unix())), b, 0644)
	if err != nil {
		Warn("ServStats.output:WriteFile",
			Err(err))

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
