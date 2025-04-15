package main



type EventQueue struct {
    Events []func()
    DeltaT []float64
}


func NewEventQueue() (q EventQueue) {
    q.Events = nil
    q.DeltaT = nil
    return q
}


func (q *EventQueue)PushBack(f func(), dt float64) {
    q.Events = append(q.Events, f) 
    q.DeltaT = append(q.DeltaT, dt) 
}


// TODO: these don't work anymore
func (q *EventQueue)Pop() {
    if len(q.Events) == 0 {
        return
    }
    q.Events[0]()
    q.Events = q.Events[1:]
}


func (q *EventQueue)Flush() {
    for _, ev := range q.Events {
        ev()
    }
    q.Events = nil
}
