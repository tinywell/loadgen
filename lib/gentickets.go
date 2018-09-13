package lib

import (
	"fmt"
	"tinywell/loadgen/model"
)

// GenTickets 载荷发生器票池
type GenTickets struct {
	pool   chan struct{}
	active bool
	total  int32
}

// NewTickets 生成新的票池并进行初始化
func NewTickets(capacity int32) (model.Tickets, error) {
	if capacity <= 0 {
		err := fmt.Errorf("want capacity >0,actuly %d", capacity)
		return nil, err
	}
	ticket := &GenTickets{
		pool:   make(chan struct{}, capacity),
		active: false,
		total:  capacity,
	}
	ticket.initTicket()
	return ticket, nil
}

func (ticket *GenTickets) initTicket() error {
	for i := 0; i < int(ticket.total); i++ {
		ticket.pool <- struct{}{}
	}
	ticket.active = true
	return nil
}

func (ticket *GenTickets) Tack() {
	<-ticket.pool
}

func (ticket *GenTickets) Return() {
	ticket.pool <- struct{}{}
}

func (ticket *GenTickets) Active() bool {
	return ticket.active
}

func (ticket *GenTickets) Total() int32 {
	return ticket.total
}

func (ticket *GenTickets) Remain() int32 {
	return int32(len(ticket.pool))
}
