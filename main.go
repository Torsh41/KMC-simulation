package main

import (
    // "os"
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
    func (ctx *Context, c *Vec2) func() {
        return EventMove(ctx, c, Vec2{0, -1})
    },
    func (ctx *Context, c *Vec2) func() {
        return EventMove(ctx, c, Vec2{1, 0})
    },
    func (ctx *Context, c *Vec2) func() {
        return EventMove(ctx, c, Vec2{0, 1})
    },
    func (ctx *Context, c *Vec2) func() {
        return EventMove(ctx, c, Vec2{-1, 0})
    },
    // Cloning
    func (ctx *Context, c *Vec2) func() {
        return EventClone(ctx, *c, Vec2{0, -1})
    },
    func (ctx *Context, c *Vec2) func() {
        return EventClone(ctx, *c, Vec2{1, 0})
    },
    func (ctx *Context, c *Vec2) func() {
        return EventClone(ctx, *c, Vec2{0, 1})
    },
    func (ctx *Context, c *Vec2) func() {
        return EventClone(ctx, *c, Vec2{-1, 0})
    },
    // Diying
    func (ctx *Context, c *Vec2) func() {
        return func() { EventDie(ctx, *c) }
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
const GRID_SYMBOL_WHITE = '█' // "\u2592"
const GRID_SYMBOL_BLACK = '▒' // "\u2593"
const GRID_SYMBOL_UNKNOWN = '?'
const GRID_VALUE_PROMISE = 2
const GRID_VALUE_POINT = 1
const GRID_VALUE_EMPTY = 0



type Vec2 struct {
    x, y int
}

type Matrix struct {
    Size Vec2
    _data []rune
}

type Context struct {
    Points []Vec2
    Time float64
}

/// EVENT METHODS



/// QUEUE METHODS


/// CONVENIENCE METHODS

func Abs(val int) int {
    if val < 0 {
        val *= -1
    }
    return val
}


func NewMatrix(size Vec2) (m Matrix, err error) {
    if size.x <= 0 || size.y <= 0 {
        return m, errors.New("Invalid size")
    }

    m.Size = size
    m._data = make([]rune, size.x * size.y)
    return m, nil
}


func NewContext() (ctx Context) {
    ctx.Time = 0.0
    ctx.Points = make([]Vec2, 0)
    return ctx
}


func (m Matrix)Get(i, j int) (ret rune, err error) {
    if (i < 0 || i >= m.Size.x || j < 0 || j >= m.Size.y) {
        return 0, errors.New("Invalid index")
    }

    ret = m._data[m.Size.x * j + i]
    return ret, nil
}


func (m *Matrix)Set(i, j int, val rune) (err error) {
    if (i < 0 || i >= m.Size.x || j < 0 || j >= m.Size.y) {
        return errors.New("Invalid index")
    }

    m._data[m.Size.x * j + i] = val
    return nil
}


func (ctx *Context)AddPoint(i, j int) {
    ctx.Points = append(ctx.Points, Vec2{i, j})
}


func (ctx *Context)RemovePoint(p Vec2) (err error) {
    for i, v := range ctx.Points {
        if v.x == p.x && v.y == p.y {
            ctx.Points = append(ctx.Points[:i], ctx.Points[i+1:]...)
            return nil
        }
    }

    return errors.New("Error removing point: not found")
}


func (ctx *Context)MovePoint(src, dst Vec2) (err error) {
    for i, v := range ctx.Points {
        if v.x == src.x && v.y == src.y {
            ctx.Points[i] = dst
            return nil
        }
    }

    return errors.New("Error moving point: not found")
}


func (ctx Context)FindPointIdx(p Vec2) (ret int, err error) {
    for idx, v := range ctx.Points {
        if v.x == p.x && v.y == p.y {
            return idx, nil
        }
    }
    return 0, errors.New("Error moving point: invalid coordinate")
}


func (ctx Context)CountAdjacentPoints(p Vec2) (count int, directions [4]bool) {
    count = 0
    for _, v := range ctx.Points {
        if Abs(v.x - p.x) + Abs(v.y - p.y) == 1 {
            count += 1
            directions[EventMap_UP] = directions[EventMap_UP] || (p.y + 1 == v.y)
            directions[EventMap_DOWN] = directions[EventMap_DOWN] || (p.y - 1 == v.y)
            directions[EventMap_RIGHT] = directions[EventMap_RIGHT] || (p.x + 1 == v.x)
            directions[EventMap_LEFT] = directions[EventMap_LEFT] || (p.x - 1 == v.x)
        }
    }
    return count, directions
}


/// PRINTING METHODS


func (m *Matrix)Clear() {
    for i := 0; i < len(m._data); i++ {
        m._data[i] = GRID_SYMBOL_BLACK
    }
}


func (m *Matrix)DrawPoint(p Vec2) {
    m.Set(p.x, p.y, GRID_SYMBOL_WHITE)
}


func (m Matrix)Print() {
    for j := 0; j < m.Size.y; j++ {
        s := string(m._data[m.Size.x * j:m.Size.x * (j + 1)])
        fmt.Println(s)
    }
}


// func (m Matrix)PrintColor() {
//     for j := 0; j < m.Size.y; j++ {
//         for i := 0; i < m.Size.x; i++ {
//             var s string
//             if m._data[m.Size.x * j + i] == GRID_VALUE_POINT {
//                 s = GRID_SYMBOL_WHITE
//             } else if m._data[m.Size.x * j + i] == GRID_VALUE_EMPTY {
//                 s = GRID_SYMBOL_BLACK
//             } else if m._data[m.Size.x * j + i] == GRID_VALUE_PROMISE {
//                 s = GRID_SYMBOL_UNKNOWN
//             }
//             fmt.Print(s)
//         }
//         fmt.Println()
//     }
// }


/// THE ALGORYTHM PART


func UpdateEventSpeed(ctx Context, m Matrix, p Vec2, eventSpeed []int) {
    // Check boundary conditions
    if p.y == 0 {
        eventSpeed[EventMap_MOVE + EventMap_UP] = 0
        eventSpeed[EventMap_CLONE + EventMap_UP] = 0
    } else if p.y == m.Size.y - 1 {
        eventSpeed[EventMap_MOVE + EventMap_DOWN] = 0
        eventSpeed[EventMap_CLONE + EventMap_DOWN] = 0
    }
    if p.x == 0 {
        eventSpeed[EventMap_MOVE + EventMap_LEFT] = 0
        eventSpeed[EventMap_CLONE + EventMap_LEFT] = 0
    } else if p.x == m.Size.x - 1 {
        eventSpeed[EventMap_MOVE + EventMap_RIGHT] = 0
        eventSpeed[EventMap_CLONE + EventMap_RIGHT] = 0
    }

    // Check adjacent points
    adjacentCount, directions := ctx.CountAdjacentPoints(p)
    for direction, cond := range directions {
        if cond == true {
            eventSpeed[EventMap_MOVE + direction] = 0
            eventSpeed[EventMap_CLONE + direction] = 0
        }
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


func (ctx *Context)IterationAdvance(m Matrix) {
    // Select dt
    K := 0
    for _, Ki := range EventSpeedDefault {
        K += Ki
    }
    dt := - math.Log(rand.Float64()) / float64(K * m.Size.x * m.Size.y)
    ctx.Time += dt

    // Select cell
    // cell_coord := Vec2{rand.Intn(m.Size.x), rand.Intn(m.Size.y)}
    // cell_idx, err := ctx.FindPointIdx(cell_coord)
    cell_idx := rand.Intn(len(ctx.Points))
    // cell_coord := ctx.Points[cell_idx]

    // Filter speed for impossible events
    eventSpeed := make([]int, len(EventSpeedDefault))
    copy(eventSpeed, EventSpeedDefault)
    UpdateEventSpeed(*ctx, m, ctx.Points[cell_idx], eventSpeed)

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

    // Call the event
    callback := EventCallbackMap[idxSel](ctx, &ctx.Points[cell_idx])
    callback()
}


func (ctx Context)FindIdenticalPoints() {
    for i := 0; i < len(ctx.Points); i++ {
        for j := i + 1; j < len(ctx.Points); j++ {
            if ctx.Points[i].x == ctx.Points[j].x && ctx.Points[i].y == ctx.Points[j].y {
                fmt.Println("Found a bad clone: ", ctx.Points[i])
            }
        }
    }
}


func main() {
    rand.Seed(time.Now().UTC().UnixNano())
    m, _ := NewMatrix(Vec2{20, 20})
    ctx := NewContext()

    // Initial conditions
    ctx.AddPoint(5, 5)
    ctx.AddPoint(5, 15)
    ctx.AddPoint(15, 5)
    ctx.AddPoint(15, 15)
    // ctx.AddPoint(15, 14)
    // ctx.AddPoint(15, 16)
    // ctx.AddPoint(14, 15)

    for ctx.Time < 100.0 {
        ctx.IterationAdvance(m)
        // / View evolution in real time
        // fmt.Print(ctx.Time, ctx.Points)
        m.Clear()
        for _, p := range ctx.Points {
            m.DrawPoint(p)
        }
        m.Print()
        fmt.Print("\n\n\n\n\n\n\n\n\n")
        time.Sleep(time.Second / 60)
    }

    /// Sanity Validation
    // ctx.FindIdenticalPoints()
    // fmt.Println(len(ctx.Points))
    fmt.Println(ctx.Points)
}
