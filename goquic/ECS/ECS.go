package ECS

import (
	"github.com/deeean/go-vector/vector2"
)

// Component defines a piece of data attached to an entity
type Component interface{}

// Position component stores the x and y coordinates of an entity
type Transform struct {
	Position vector2.Vector2
	Velocity vector2.Vector2
}

// Entity is a unique identifier with attached components
type Entity struct {
	EntityId   int
	Components map[string]Component
}

// System defines a group of entities with specific components and logic to operate on them
type System interface {
	Update(entities []Entity)
}

// Create player entity
type Player struct {
	Entity
	Transform
	MaxVel    float64
	Speed     float64
	JumpForce float64
	Grounded  bool
}

// MovementSystem updates entities' positions based on their velocities
type MovementSystem struct{}

func (ms MovementSystem) Update(players []*Player) {
	for _, player := range players {
		newPosition := player.Transform.Position.Add(&player.Transform.Velocity)
		player.Transform.Position = *newPosition
		// fmt.Printf("Player %d new position: %+v\n", player.ID, player.Transform.Position)
	}
}

// func CreatePlayerEntity(id int, pos vector2.Vector2, vel vector2.Vector2) Entity {
// 	// Initialize the Position component
// 	positionComponent := Position{
// 		Position: pos,
// 		Velocity: vel,
// 	}

// 	// Create the entity and attach the components
// 	player := Entity{
// 		ID: id,
// 		Components: map[string]Component{
// 			"Position": positionComponent,
// 		},
// 	}

// 	return player
// }
