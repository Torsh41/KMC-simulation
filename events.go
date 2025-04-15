package main

import (
    "log"
)


type EventCallback func (ctx *Context, c *Vec2) func()



func EventMove(ctx *Context, c *Vec2, d Vec2) func() {
    src := *c
    dst := Vec2{src.x + d.x, src.y + d.y}
    // (*c).x += d.x
    // (*c).y += d.y
    // (*ctx).Set(dst.x, dst.y, GRID_VALUE_PROMISE)
    return func() {
        // log.Print("Removing a point: ", future)
        // (*ctx).Set(future.x, future.y, GRID_VALUE_EMPTY)
        // (*ctx).Set(future.x + d.x, future.y + d.y, GRID_VALUE_POINT)
        err := ctx.MovePoint(src, dst)
        if err != nil {
            log.Println(src, err)
        }
    }
}


func EventClone(ctx *Context, c Vec2, d Vec2) func() {
    dst := Vec2{c.x + d.x, c.y + d.y}
    // (*ctx).Set(dst.x, dst.y, GRID_VALUE_PROMISE)
    return func() {
        // (*ctx).Set(dst.x, dst.y, GRID_VALUE_POINT)
        (*ctx).AddPoint(dst.x, dst.y)
    }
}


func EventDie(ctx *Context, c Vec2) {
    ctx.RemovePoint(c)
}


