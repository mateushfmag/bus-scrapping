package main

type SeatType string

// const (
// 	Comfort SeatType = "COMFORT"
// 	Common  SeatType = "COMMON"
// )

type Itinerary struct {
	Date string
	Hour string
	Location string
}

type BusSummary struct {
	Price float32
	Seat string
	Origin Itinerary
	Destination Itinerary
	Available bool
}