package sb

import "errors"

var (
	ErrRepeatedShot = errors.New("repeated shot")
	ErrGameEnded    = errors.New("game ended")
)

type info struct {
	Destroy bool `json:"destroy"`
	Knock   bool `json:"knock"`
	End     bool `json:"end"`
}

type state struct {
	ShipCount int `json:"ship_count"`
	Destroyed int `json:"destroyed"`
	Knocked   int `json:"knocked"`
	ShotCount int `json:"shot_count"`
}

type SeaBattle struct {
	active bool
	field  *field
	ships  []*Ship
}

func Create(size int) (*SeaBattle, error) {
	f, err := NewField(size, size)
	if err != nil {
		return nil, err
	}
	return &SeaBattle{
		field: f,
	}, nil
}

func (sb *SeaBattle) Setup(ships []*Ship) error {
	for _, ship := range ships {
		if err := sb.field.add(ship); err != nil {
			return err
		}
	}

	sb.ships = ships
	sb.active = true

	return nil
}

func (sb *SeaBattle) Shot(p Point) (*info, error) {
	if !sb.active {
		return nil, ErrGameEnded
	}

	if err := sb.field.check(p); err != nil {
		return nil, err
	}

	c, err := sb.field.cell(p)
	if err != nil {
		return nil, err
	}

	// Повторный выстрел
	if c.hit {
		return nil, ErrRepeatedShot
	}
	c.hit = true

	i := info{}

	// Промах
	if c.ship == nil {
		return &i, nil
	}

	s := c.ship
	i.Knock = s.Knock(p)
	i.Destroy = s.Destroyed()

	var destroyed int
	for _, ship := range sb.ships {
		if ship.Destroyed() {
			destroyed++
		}
	}
	if destroyed == len(sb.ships) {
		i.End = true
		sb.active = false
	}

	return &i, nil
}

func (sb *SeaBattle) State() *state {
	var destroyed, knocked int
	for _, ship := range sb.ships {
		if ship.Destroyed() {
			destroyed++
		}
		if ship.Knocked() {
			knocked++
		}
	}

	var shots int
	for x := 0; x < len(*sb.field); x++ {
		for y := 0; y < len((*sb.field)[0]); y++ {
			if (*sb.field)[x][y].hit {
				shots++
			}
		}
	}

	return &state{
		ShipCount: len(sb.ships),
		Destroyed: destroyed,
		Knocked:   knocked,
		ShotCount: shots,
	}
}
