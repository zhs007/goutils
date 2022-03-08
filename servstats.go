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

type ServStatsContext struct {
	ID      int
	MsgName string
	State   int
}

// type ServStatsKey int

// const (
// 	ServStatsMsgName ServStatsKey = iota
// )

type ServStatsMsgNode struct {
	Start     time.Time `json:"-"`
	End       time.Time `json:"-"`
	SecondOff float64   `json:"secondOff,omitempty"`
}

type ServStatsMsg struct {
	Name         string                                  `json:"name,omitempty"`
	TotalTime    float64                                 `json:"totalTime,omitempty"`
	TotalTimes   int                                     `json:"totalTimes,omitempty"`
	MaxTime      float64                                 `json:"maxTime,omitempty"`
	MinTime      float64                                 `json:"minTime,omitempty"`
	MaxParallels int                                     `json:"maxParallels,omitempty"`
	Nodes        []*ServStatsMsgNode                     `json:"nodes,omitempty"`
	NoEndNodes   map[*ServStatsContext]*ServStatsMsgNode `json:"noEndNodes,omitempty"`
}

func newServStatsMsg(name string) *ServStatsMsg {
	return &ServStatsMsg{
		Name:       name,
		MinTime:    math.MaxFloat64,
		NoEndNodes: make(map[*ServStatsContext]*ServStatsMsgNode),
	}
}

func (msg *ServStatsMsg) startMsg(ctx *ServStatsContext) {
	lastnode, isok := msg.NoEndNodes[ctx]
	if isok {
		Warn("ServStatsMsg:StartMsg",
			JSON("node", lastnode),
			zap.Error(ErrDuplicateMsgCtx))
	}

	msg.NoEndNodes[ctx] = &ServStatsMsgNode{
		Start: time.Now(),
	}

	if len(msg.NoEndNodes) > msg.MaxParallels {
		msg.MaxParallels = len(msg.NoEndNodes)
	}
}

func (msg *ServStatsMsg) endMsg(ctx *ServStatsContext, maxNodes int) {
	_, isok := msg.NoEndNodes[ctx]
	if !isok {
		Warn("ServStatsMsg:EndMsg",
			zap.Error(ErrNoMsgCtx))

		return
	}

	msg.NoEndNodes[ctx].End = time.Now()

	dt := msg.NoEndNodes[ctx].End.Sub(msg.NoEndNodes[ctx].Start).Seconds()
	if dt > msg.MaxTime {
		msg.MaxTime = dt
	}

	if dt < msg.MinTime {
		msg.MinTime = dt
	}

	msg.NoEndNodes[ctx].SecondOff = dt

	msg.TotalTime += dt
	msg.TotalTimes++

	msg.Nodes = append(msg.Nodes, msg.NoEndNodes[ctx])
	if len(msg.Nodes) > maxNodes*2 {
		sort.Slice(msg.Nodes, func(i, j int) bool {
			return msg.Nodes[i].SecondOff > msg.Nodes[j].SecondOff
		})

		msg.Nodes = msg.Nodes[:maxNodes]
	}

	delete(msg.NoEndNodes, ctx)
}

type ServStats struct {
	MapMsgs       map[string]*ServStatsMsg `json:"mapMsgs,omitempty"`
	MaxNodes      int                      `json:"-"`
	ChanCtx       chan *ServStatsContext   `json:"-"`
	ChanState     chan int                 `json:"-"`
	TimerOutput   *time.Timer              `json:"-"`
	PathOutput    string                   `json:"-"`
	TotalPoolSize int                      `json:"totalPoolSize,omitempty"`
	LastPoolSize  int                      `json:"lastPoolSize,omitempty"`
	PoolSize      int                      `json:"poolSize,omitempty"`
	poolContext   []*ServStatsContext      `json:"-"`
	lock          sync.Mutex               `json:"-"`
}

func NewServStats(maxNodes int, chanSize int, outputTimer time.Duration, pathOutput string, poolSize int) *ServStats {
	stats := &ServStats{
		MapMsgs:     make(map[string]*ServStatsMsg),
		MaxNodes:    maxNodes,
		ChanCtx:     make(chan *ServStatsContext, chanSize),
		ChanState:   make(chan int),
		TimerOutput: time.NewTimer(outputTimer),
		PoolSize:    poolSize,
	}

	stats.newPool()

	return stats
}

func (stats *ServStats) newPool() {
	for i := 0; i < stats.PoolSize; i++ {
		stats.poolContext = append(stats.poolContext, &ServStatsContext{
			ID: len(stats.poolContext) + 1,
		})
	}

	stats.TotalPoolSize += stats.PoolSize
}

func (stats *ServStats) Start() {
	go stats.mainLoop()
}

func (stats *ServStats) Stop() {
	stats.TimerOutput.Stop()
	stats.ChanState <- 0
}

func (stats *ServStats) StartMsg(msgname string) *ServStatsContext {
	stats.lock.Lock()
	if len(stats.poolContext) <= 0 {
		stats.newPool()
	}

	ctx := stats.poolContext[len(stats.poolContext)-1]
	stats.poolContext = stats.poolContext[:(len(stats.poolContext) - 1)]
	stats.lock.Unlock()

	ctx.MsgName = msgname
	ctx.State = 0

	stats.ChanCtx <- ctx

	return ctx
}

func (stats *ServStats) EndMsg(ctx *ServStatsContext) {
	ctx.State = 1
	stats.ChanCtx <- ctx
}

func (stats *ServStats) RegMsg(name string) {
	stats.MapMsgs[name] = newServStatsMsg(name)
}

func (stats *ServStats) output() {
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	b, err := json.Marshal(stats)
	if err != nil {
		Warn("ServStats.output:Marshal",
			zap.Error(err))

		return
	}

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
		case ctx := <-stats.ChanCtx:
			msg, isok := stats.MapMsgs[ctx.MsgName]
			if !isok {
				Error("ServStats:mainLoop:ChanStart",
					zap.String("msgname", ctx.MsgName),
					zap.Error(ErrInvalidMsgName))
			} else {
				if ctx.State == 0 {
					msg.startMsg(ctx)
				} else {
					msg.endMsg(ctx, stats.MaxNodes)

					stats.lock.Lock()
					stats.poolContext = append(stats.poolContext, ctx)
					stats.lock.Unlock()
				}
			}
		// case ctx := <-stats.ChanEnd:
		// 	// name, isok := ctx.Value(ServStatsMsgName).(string)
		// 	// if !isok {
		// 	// 	Error("ServStats:mainLoop:ChanEnd",
		// 	// 		zap.Error(ErrNoMsgName))
		// 	// }

		// 	msg, isok := stats.MapMsgs[ctx.MsgName]
		// 	if !isok {
		// 		Error("ServStats:mainLoop:ChanEnd",
		// 			zap.String("msgname", ctx.MsgName),
		// 			zap.Error(ErrInvalidMsgName))
		// 	}

		// 	msg.endMsg(ctx, stats.MaxNodes)

		// 	stats.lock.Lock()
		// 	stats.poolContext = append(stats.poolContext, ctx)
		// 	stats.lock.Unlock()
		case <-stats.ChanState:
			Info("ServStats:mainLoop:ChanState")

			return
		case <-stats.TimerOutput.C:
			stats.output()
		}
	}
}
