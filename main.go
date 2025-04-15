package main

import (
    "fmt"
    "math"
    "time"
    "math/rand"
    "errors"
)


const EVENT_COUNT = 9

// All posible events
var EventSpeedDefault = []int{
    // move u/r/d/l
    1, 1, 1, 1,
    // clone u/r/d/l
    3, 3, 3, 3,
    // die
    2,
}

var EventCallbackMap = []EventCallback {
    // Moving
    func (g *Grid, c *Coordinate) func() {
        return EventMove(g, c, Coordinate{0, -1})
    },
    func (g *Grid, c *Coordinate) func() {
        return EventMove(g, c, Coordinate{1, 0})
    },
    func (g *Grid, c *Coordinate) func() {
        return EventMove(g, c, Coordinate{0, 1})
    },
    func (g *Grid, c *Coordinate) func() {
        return EventMove(g, c, Coordinate{-1, 0})
    },
    // Cloning
    func (g *Grid, c *Coordinate) func() {
        return EventClone(g, *c, Coordinate{0, -1})
    },
    func (g *Grid, c *Coordinate) func() {
        return EventClone(g, *c, Coordinate{1, 0})
    },
    func (g *Grid, c *Coordinate) func() {
        return EventClone(g, *c, Coordinate{0, 1})
    },
    func (g *Grid, c *Coordinate) func() {
        return EventClone(g, *c, Coordinate{-1, 0})
    },
    // Diying
    func (g *Grid, c *Coordinate) func() {
        return func() { EventDie(g, *c) }
    },
}

// Indices for the events
const EventMap_MOVE     = 0
const EventMap_CLONE    = 4
const EventMap_DIE      = 8

const EventMap_UP       = 0
const EventMap_RIGHT    = 1
const EventMap_DOWN     = 2
const EventMap_LEFT     = 3

// Const values for Grid
const GRID_SYMBOL_WHITE = "█" // "\u2592"
const GRID_SYMBOL_BLACK = "▒" // "\u2593"
const GRID_SYMBOL_UNKNOWN = "?"
const GRID_VALUE_PROMISE = 2
const GRID_VALUE_POINT = 1
const GRID_VALUE_EMPTY = 0



type Coordinate struct {
    x, y int
}

type Grid struct {
    Size Coordinate
    Points []Coordinate
    _data []int
    Time float64
}


/// EVENT METHODS



/// QUEUE METHODS


/// CONVENIENCE METHODS


func NewGrid(size Coordinate) (g Grid, err error) {
    if size.x <= 0 || size.y <= 0 {
        return g, errors.New("Invalid size")
    }

    g.Time = 0.0
    g.Size = size
    g._data = make([]int, size.x * size.y)
    g.Points = make([]Coordinate, 0)
    return g, nil
}


func (g Grid)Get(i, j int) (ret int, err error) {
    if (i < 0 || i >= g.Size.x || j < 0 || j >= g.Size.y) {
        return 0, errors.New("Invalid index")
    }

    ret = g._data[g.Size.x * j + i]
    return ret, nil
}


func (g *Grid)Set(i, j int, val int) (err error) {
    if (i < 0 || i >= g.Size.x || j < 0 || j >= g.Size.y) {
        return errors.New("Invalid index")
    }

    g._data[g.Size.x * j + i] = val
    return nil
}


func (g *Grid)AddPoint(i, j int) (err error) {
    if (i < 0 || i >= g.Size.x || j < 0 || j >= g.Size.y) {
        return errors.New("Invalid index")
    }

    g.Points = append(g.Points, Coordinate{i, j})
    g.Set(i, j, GRID_VALUE_POINT)
    return nil
}


func (g *Grid)RemovePoint(p Coordinate) (err error) {
    for i, v := range g.Points {
        if v.x == p.x && v.y == p.y {
            g.Points = append(g.Points[:i], g.Points[i+1:]...)
            g.Set(v.x, v.y, GRID_VALUE_EMPTY)
            return nil
        }
    }

    return errors.New("Error removing point: not found")
}


func (g *Grid)MovePoint(src Coordinate, dst Coordinate) (err error) {
    if (dst.x < 0 || dst.x >= g.Size.x ||
        dst.y < 0 || dst.y >= g.Size.y){
        return errors.New("Error moving point: invalid dst coordinate")
    }

    for i, v := range g.Points {
        if v.x == src.x && v.y == src.y {
            g.Set(src.x, src.y, GRID_VALUE_EMPTY)
            g.Set(dst.x, dst.y, GRID_VALUE_POINT)
            g.Points[i] = dst
            return nil
        }
    }

    return errors.New("Error moving point: not found")
}


func (g Grid)FindPointIdx(p Coordinate) (ret int, err error) {
    for idx, v := range g.Points {
        if v.x == p.x && v.y == p.y {
            return idx, nil
        }
    }
    return 0, errors.New("Error moving point: invalid coordinate")
}


/// PRINTING METHODS


func (g Grid)Print() {
    for j := 0; j < g.Size.y; j++ {
        for i := 0; i < g.Size.x; i++ {
            fmt.Print(g._data[g.Size.x * j + i], " ")
        }
        fmt.Println()
    }
}


func (g Grid)PrintColor() {
    for j := 0; j < g.Size.y; j++ {
        for i := 0; i < g.Size.x; i++ {
            var s string
            if g._data[g.Size.x * j + i] == GRID_VALUE_POINT {
                s = GRID_SYMBOL_WHITE
            } else if g._data[g.Size.x * j + i] == GRID_VALUE_EMPTY {
                s = GRID_SYMBOL_BLACK
            } else if g._data[g.Size.x * j + i] == GRID_VALUE_PROMISE {
                s = GRID_SYMBOL_UNKNOWN
            }
            fmt.Print(s)
        }
        fmt.Println()
    }
}


/// THE ALGORYTHM PART


func UpdateEventSpeed(g Grid, p Coordinate, eventSpeed []int) {
    // Check boundary conditions
    if p.y == 0 {
        eventSpeed[EventMap_MOVE + EventMap_UP] = 0
        eventSpeed[EventMap_CLONE + EventMap_UP] = 0
    } else if p.y == g.Size.y - 1 {
        eventSpeed[EventMap_MOVE + EventMap_DOWN] = 0
        eventSpeed[EventMap_CLONE + EventMap_DOWN] = 0
    }
    if p.x == 0 {
        eventSpeed[EventMap_MOVE + EventMap_LEFT] = 0
        eventSpeed[EventMap_CLONE + EventMap_LEFT] = 0
    } else if p.x == g.Size.x - 1 {
        eventSpeed[EventMap_MOVE + EventMap_RIGHT] = 0
        eventSpeed[EventMap_CLONE + EventMap_RIGHT] = 0
    }

    // Check adjacent points
    var v int
    var err error
    adjacentCount := 0
    checkAdjacent := func(v int, direction int) (ret int) {
        if v == GRID_VALUE_POINT /* || v == GRID_VALUE_PROMISE */ {
            eventSpeed[EventMap_MOVE + direction] = 0
            eventSpeed[EventMap_CLONE + direction] = 0
            return 1
        }
        return 0
    }

    v, err = g.Get(p.x, p.y - 1)
    if err == nil {
        adjacentCount += checkAdjacent(v, EventMap_UP)
    }
    v, err = g.Get(p.x, p.y + 1)
    if err == nil {
        adjacentCount += checkAdjacent(v, EventMap_DOWN)
    }
    v, err = g.Get(p.x - 1, p.y)
    if err == nil {
        adjacentCount += checkAdjacent(v, EventMap_LEFT)
    }
    v, err = g.Get(p.x + 1, p.y)
    if err == nil {
        adjacentCount += checkAdjacent(v, EventMap_RIGHT)
    }

    // Special case, unable to clone
    if adjacentCount == 0 {
        eventSpeed[EventMap_CLONE + EventMap_UP] = 0
        eventSpeed[EventMap_CLONE + EventMap_RIGHT] = 0
        eventSpeed[EventMap_CLONE + EventMap_DOWN] = 0
        eventSpeed[EventMap_CLONE + EventMap_LEFT] = 0
    }

    // Check "die" condition
    if adjacentCount < 3 {
        eventSpeed[EventMap_DIE] = 0
    }
}


func (g *Grid)IterationAdvance() {
    // Select dt
    K := 0
    for _, Ki := range EventSpeedDefault {
        K += Ki
    }
    dt := - math.Log(rand.Float64()) / float64(K * g.Size.x * g.Size.y)
    g.Time += dt

    // Select cell
    cell_coord := Coordinate{rand.Intn(g.Size.x), rand.Intn(g.Size.y)}
    // cell_coord := Coordinate{15, 15}
    cell_val, _ := g.Get(cell_coord.x, cell_coord.y)
    if cell_val == GRID_VALUE_EMPTY {
        return
    }
    cell_idx, _ := g.FindPointIdx(cell_coord)

    // Filter speed for impossible events
    eventSpeed := make([]int, len(EventSpeedDefault))
    copy(eventSpeed, EventSpeedDefault)
    UpdateEventSpeed(*g, g.Points[cell_idx], eventSpeed)


    // Select the event
    K = 0
    for _, Ki := range eventSpeed {
        K += Ki
    }
    evAccumulate := 0.0
    idxSel := -1
    r := rand.Float64()
    // fmt.Println(eventSpeed, r)
    for i := 0; i < len(eventSpeed); i++ {
        // fmt.Println("idx: ", i, "; spd: ", eventSpeed[i])
        dspeed := float64(eventSpeed[i]) / float64(K)



        // fmt.Println(dspeed, evAccumulate, evAccumulate + dspeed)
        if (evAccumulate <= r && r < evAccumulate + dspeed) {
            // check, if Event is possible at all
            if eventSpeed[i] != 0.0 {
                idxSel = i
            }
            break;
        }
        evAccumulate += dspeed
    }

    // Event is impossible
    if idxSel == -1 {
        return;
    }

    // fmt.Println(g.Time, len(g.Points))

    // fmt.Println("r: ", r)
    // fmt.Println("cell_coord: ", cell_coord)
    // fmt.Println("cell: ", g.Points[cell_idx])
    // fmt.Println(g.Points)

    // Call the event
    callback := EventCallbackMap[idxSel](g, &g.Points[cell_idx])
    callback()
}


func (g Grid)FindIdenticalPoints() {
    for i := 0; i < len(g.Points); i++ {
        for j := i + 1; j < len(g.Points); j++ {
            if g.Points[i].x == g.Points[j].x && g.Points[i].y == g.Points[j].y {
                fmt.Println("Found a bad clone: ", g.Points[i])
            }
        }
    }
}


func main() {
    rand.Seed(time.Now().UTC().UnixNano())
    g, _ := NewGrid(Coordinate{20, 20})

    // Initial conditions
    g.AddPoint(5, 5)
    g.AddPoint(5, 15)
    g.AddPoint(15, 5)
    g.AddPoint(15, 15)
    // g.AddPoint(15, 14)
    // g.AddPoint(15, 16)
    // g.AddPoint(14, 15)

    for g.Time < 100.0 {
        g.IterationAdvance()
        // / View evolution in real time
        // fmt.Print(g.Time, g.Points)
        g.PrintColor()
        fmt.Print("\n\n\n\n\n\n\n\n\n")
        // time.Sleep(time.Second / 60)
    }

    /// Sanity Validation
    // g.FindIdenticalPoints()
    // fmt.Println(len(g.Points))
    // fmt.Println(g.Points)
    g.PrintColor()
}
