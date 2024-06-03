package ECS

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Component defines a piece of data attached to an entity
type Component interface{}

// Position component stores the x and y coordinates of an entity
type Transform struct {
	Position rl.Vector2
	Velocity rl.Vector2
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
	MaxVel  float32
	Speed   float32
	Texture rl.Texture2D
}

// MovementSystem updates entities' positions based on their velocities
type MovementSystem struct{}

func (ms MovementSystem) Update(players []*Player) {
	for _, player := range players {
		newPosition := rl.Vector2Add(player.Transform.Position, player.Transform.Velocity)
		// newPosition := player.Transform.Position.Add(&player.Transform.Velocity)
		player.Transform.Position = newPosition
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
