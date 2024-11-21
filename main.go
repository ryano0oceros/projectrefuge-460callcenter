package main

import (
	"container/heap"
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type Event struct {
	timestamp int
	eventType string // "CallArrival", "CallCompletion", "CallAbandonment"
	callID    int
	agentID   int
}

type EventQueue []*Event

func (eq EventQueue) Len() int           { return len(eq) }
func (eq EventQueue) Less(i, j int) bool { return eq[i].timestamp < eq[j].timestamp }
func (eq EventQueue) Swap(i, j int)      { eq[i], eq[j] = eq[j], eq[i] }

func (eq *EventQueue) Push(x interface{}) {
	*eq = append(*eq, x.(*Event))
}

func (eq *EventQueue) Pop() interface{} {
	old := *eq
	n := len(old)
	event := old[n-1]
	*eq = old[0 : n-1]
	return event
}

type CallCenter struct {
	numAgents       int
	availableAgents []bool
	eventQueue      EventQueue
	currentTime     int
	callCounter     int
	callQueue       []*Event
	logs            []string
	maxWaitTime     int
	totalCalls      int
	abandonedCalls  int
	busyTime        []int
	lambda          float64
	averageCallTime float64
}

func NewCallCenter(numAgents int, maxWaitTime int, lambda, averageCallTime float64) *CallCenter {
	return &CallCenter{
		numAgents:       numAgents,
		availableAgents: make([]bool, numAgents),
		eventQueue:      EventQueue{},
		currentTime:     0,
		callCounter:     0,
		callQueue:       []*Event{},
		logs:            []string{},
		maxWaitTime:     maxWaitTime,
		busyTime:        make([]int, numAgents),
		lambda:          lambda,
		averageCallTime: averageCallTime,
	}
}

func (cc *CallCenter) ScheduleEvent(timestamp int, eventType string, callID, agentID int) {
	event := &Event{
		timestamp: timestamp,
		eventType: eventType,
		callID:    callID,
		agentID:   agentID,
	}
	heap.Push(&cc.eventQueue, event)
}

func (cc *CallCenter) ProcessNextEvent() {
	event := heap.Pop(&cc.eventQueue).(*Event)
	cc.currentTime = event.timestamp
	switch event.eventType {
	case "CallArrival":
		cc.handleCallArrival(event)
	case "CallCompletion":
		cc.handleCallCompletion(event)
	case "CallAbandonment":
		cc.handleCallAbandonment(event)
	}
}

func (cc *CallCenter) handleCallArrival(event *Event) {
	cc.totalCalls++
	if cc.assignAgentToCall(event) {
		return
	}
	cc.callQueue = append(cc.callQueue, event)
	cc.ScheduleEvent(cc.currentTime+cc.maxWaitTime, "CallAbandonment", event.callID, -1)
}

func (cc *CallCenter) assignAgentToCall(event *Event) bool {
	for i := 0; i < cc.numAgents; i++ {
		if !cc.availableAgents[i] {
			cc.availableAgents[i] = true
			callDuration := exponential(1.0 / cc.averageCallTime)
			cc.ScheduleEvent(cc.currentTime+callDuration, "CallCompletion", event.callID, i)
			cc.busyTime[i] += callDuration
			return true
		}
	}
	return false
}

func (cc *CallCenter) handleCallCompletion(event *Event) {
	cc.availableAgents[event.agentID] = false
	if len(cc.callQueue) > 0 {
		nextCall := cc.callQueue[0]
		cc.callQueue = cc.callQueue[1:]
		cc.assignAgentToCall(nextCall)
	}
}

func (cc *CallCenter) handleCallAbandonment(event *Event) {
	for i, call := range cc.callQueue {
		if call.callID == event.callID {
			cc.callQueue = append(cc.callQueue[:i], cc.callQueue[i+1:]...)
			cc.abandonedCalls++
			break
		}
	}
}

func exponential(lambda float64) int {
	return int(-math.Log(1.0-rand.Float64()) / lambda)
}

func (cc *CallCenter) RunSimulation(simulationTime int) (int, int, float64) {
	heap.Init(&cc.eventQueue)
	t := 0
	for t < simulationTime {
		interArrivalTime := exponential(cc.lambda)
		t += interArrivalTime
		cc.callCounter++
		cc.ScheduleEvent(t, "CallArrival", cc.callCounter, -1)
	}
	for len(cc.eventQueue) > 0 && cc.currentTime < simulationTime {
		cc.ProcessNextEvent()
	}
	totalBusyTime := 0
	for _, bt := range cc.busyTime {
		totalBusyTime += bt
	}
	utilization := float64(totalBusyTime) / float64(cc.numAgents*simulationTime)
	return cc.totalCalls, cc.abandonedCalls, utilization
}

func main() {
	numAgents := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	simulationTime := 1440
	maxWaitTimes := []int{5, 10, 15}
	lambdas := []float64{0.5, 1, 1.5, 2}
	averageCallTimes := []float64{3, 5, 7, 9}

	rand.Seed(time.Now().UnixNano())

	file, err := os.Create("simulation_results.csv")
	if err != nil {
		log.Fatalf("Error creating CSV file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{"NumAgents", "SimulationTime", "MaxWaitTime", "Lambda", "AverageCallTime", "TotalCalls", "AbandonedCalls", "Utilization"}
	writer.Write(headers)

	for _, agents := range numAgents {
		for _, maxWait := range maxWaitTimes {
			for _, lambda := range lambdas {
				for _, avgCallTime := range averageCallTimes {
					callCenter := NewCallCenter(agents, maxWait, lambda, avgCallTime)
					totalCalls, abandonedCalls, utilization := callCenter.RunSimulation(simulationTime)
					writer.Write([]string{
						strconv.Itoa(agents),
						strconv.Itoa(simulationTime),
						strconv.Itoa(maxWait),
						fmt.Sprintf("%.1f", lambda),
						fmt.Sprintf("%.1f", avgCallTime),
						strconv.Itoa(totalCalls),
						strconv.Itoa(abandonedCalls),
						fmt.Sprintf("%.2f", utilization*100),
					})
				}
			}
		}
	}

	fmt.Println("Simulation completed. Results written to simulation_results.csv")
}
