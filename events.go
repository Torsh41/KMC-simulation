package main

import (
    "log"
)


type EventCallback func (g *Grid, c *Coordinate) func()



func EventMove(g *Grid, c *Coordinate, d Coordinate) func() {
    src := *c
    dst := Coordinate{src.x + d.x, src.y + d.y}
    // (*c).x += d.x
    // (*c).y += d.y
    // (*g).Set(dst.x, dst.y, GRID_VALUE_PROMISE)
    return func() {
        // log.Print("Removing a point: ", future)
        // (*g).Set(future.x, future.y, GRID_VALUE_EMPTY)
        // (*g).Set(future.x + d.x, future.y + d.y, GRID_VALUE_POINT)
        err := g.MovePoint(src, dst)
        if err != nil {
            log.Println(src, err)
        }
    }
}


func EventClone(g *Grid, c Coordinate, d Coordinate) func() {
    dst := Coordinate{c.x + d.x, c.y + d.y}
    // (*g).Set(dst.x, dst.y, GRID_VALUE_PROMISE)
    return func() {
        (*g).Set(dst.x, dst.y, GRID_VALUE_POINT)
        err := (*g).AddPoint(dst.x, dst.y)
        if err != nil {
            log.Println(err)
        }
    }
}


func EventDie(g *Grid, c Coordinate) {
    g.RemovePoint(c)
}


