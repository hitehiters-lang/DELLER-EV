package schema

import (
	"github.com/google/uuid"
)

type MotoSpecifications struct {
	Id                 uuid.UUID `bun:"id,pk,autoincrement"`
	Gear               string    `bun:"gear"`
	PowerRatedMax      string    `bun:"power_rated_max"`
	TorqueMoment       string    `bun:"torque_moment"`
	MaxSpeed           string    `bun:"max_speed"`
	Battery            string    `bun:"battery"`
	PowerReserve       string    `bun:"power_reserve"`
	Weight             string    `bun:"weight"`
	LoadCap            string    `bun:"load_cap"`
	DimensionsLWH      string    `bun:"dimensions_l_w_h"`
	SeatHeight         string    `bun:"seat_height"`
	Tires              string    `bun:"tires"`
	Break              string    `bun:"break"`
	Clearance          string    `bun:"clearance"`
	Charging           string    `bun:"charging"`
	MoistureProtection string    `bun:"moisture_protection"`
}
