package store

import "context"

type Store struct {
	Users  UserStore
	Hotels HotelStore
	Rooms  RoomStore
}
type Dropper interface {
	Drop(context context.Context) error
}
