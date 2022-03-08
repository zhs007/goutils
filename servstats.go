package goutils

import (
	"context"
	"fmt"
	"math"
	"os"
	"path"
	"sort"
	"time"

	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

type ServStatsMsgNode struct {
	Start     time.Time `json:"-"`
	End       time.Time `json:"-"`
	SecondOff float64   `json:"secondOff,omitempty"`
}

type ServStatsMsg struct {
	Name         string                                `json:"name,omitempty"`
	TotalTime    float64                               `json:"totalTime,omitempty"`
	TotalTimes   int                                   `json:"totalTimes,omitempty"`
	MaxTime      float64                               `json:"maxTime,omitempty"`
	MinTime      float64                               `json:"minTime,omitempty"`
	MaxParallels int                                   `json:"maxParallels,omitempty"`
	Nodes        []*ServStatsMsgNode                   `json:"nodes,omitempty"`
	NoEndNodes   map[context.Context]*ServStatsMsgNode `json:"noEndNodes,omitempty"`
}

func newServStatsMsg(name string) *ServStatsMsg {
	return &ServStatsMsg{
		Name:       name,
		MinTime:    math.MaxFloat64,
		NoEndNodes: make(map[context.Context]*ServStatsMsgNode),
	}
}

func (msg *ServStatsMsg) startMsg(ctx context.Context) {
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

func (msg *ServStatsMsg) endMsg(ctx context.Context, maxNodes int) {
	_, isok := msg.NoEndNodes[ctx]
	if !isok {
		Warn("ServStatsMsg:EndMsg",
			zap.Error(ErrNoMsgCtx))
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
	MapMsgs     map[string]*ServStatsMsg `json:"mapMsgs,omitempty"`
	MaxNodes    int                      `json:"-"`
	ChanStart   chan context.Context     `json:"-"`
	ChanEnd     chan context.Context     `json:"-"`
	ChanState   chan int                 `json:"-"`
	TimerOutput *time.Timer              `json:"-"`
	PathOutput  string                   `json:"-"`
}

func NewServStats(maxNodes int, chanSize int, outputTimer time.Duration, pathOutput string) *ServStats {
	return &ServStats{
		MapMsgs:     make(map[string]*ServStatsMsg),
		MaxNodes:    maxNodes,
		ChanStart:   make(chan context.Context, chanSize),
		ChanEnd:     make(chan context.Context, chanSize),
		ChanState:   make(chan int),
		TimerOutput: time.NewTimer(outputTimer),
	}
}

func (stats *ServStats) Start() {
	go stats.mainLoop()
}

func (stats *ServStats) Stop() {
	stats.TimerOutput.Stop()
	stats.ChanState <- 0
}

func (stats *ServStats) StartMsg(ctx context.Context) {
	stats.ChanStart <- ctx
}

func (stats *ServStats) EndMsg(ctx context.Context) {
	stats.ChanEnd <- ctx
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
		case ctx := <-stats.ChanStart:
			name, isok := ctx.Value("msgname").(string)
			if !isok {
				Error("ServStats:mainLoop:ChanStart",
					zap.Error(ErrNoMsgName))
			}

			msg, isok := stats.MapMsgs[name]
			if !isok {
				Error("ServStats:mainLoop:ChanStart",
					zap.String("msgname", name),
					zap.Error(ErrInvalidMsgName))
			}

			msg.startMsg(ctx)
		case ctx := <-stats.ChanEnd:
			name, isok := ctx.Value("msgname").(string)
			if !isok {
				Error("ServStats:mainLoop:ChanEnd",
					zap.Error(ErrNoMsgName))
			}

			msg, isok := stats.MapMsgs[name]
			if !isok {
				Error("ServStats:mainLoop:ChanEnd",
					zap.String("msgname", name),
					zap.Error(ErrInvalidMsgName))
			}

			msg.endMsg(ctx, stats.MaxNodes)
		case <-stats.ChanState:
			Info("ServStats:mainLoop:ChanState")

			return
		case <-stats.TimerOutput.C:
			stats.output()
		}
	}
}
