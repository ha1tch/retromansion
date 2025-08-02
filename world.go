package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (fg *FilmationGame) BuildWorld() {
	fmt.Println("Building 3D tile-based world...")

	world := World3D{
		Width:       12,
		Height:      8,
		Depth:       12,
		PlayerSpawn: Point3D{X: 6, Y: 1, Z: 6},
	}

	world.Tiles = make([][][]Tile3D, world.Width)
	for x := range world.Tiles {
		world.Tiles[x] = make([][]Tile3D, world.Height)
		for y := range world.Tiles[x] {
			world.Tiles[x][y] = make([]Tile3D, world.Depth)
		}
	}

	for x := 0; x < world.Width; x++ {
		for z := 0; z < world.Depth; z++ {
			world.Tiles[x][0][z] = Tile3D{
				Type:     TileType((x + z) % 4),
				Position: Point3D{X: float32(x), Y: 0, Z: float32(z)},
				Solid:    false,
				Height:   1.0,
			}
		}
	}

	for x := 0; x < world.Width; x++ {
		for z := 0; z < world.Depth; z++ {
			if x == 0 || x == world.Width-1 || z == 0 || z == world.Depth-1 {
				world.Tiles[x][1][z] = Tile3D{
					Type:     TileStoneWall,
					Position: Point3D{X: float32(x), Y: 1, Z: float32(z)},
					Solid:    true,
					Height:   1.0,
				}
			}
		}
	}

	for x := 4; x <= 7; x++ {
		for z := 4; z <= 7; z++ {
			if x == 4 || x == 7 || z == 4 || z == 7 {
				world.Tiles[x][1][z] = Tile3D{
					Type:     TileBrickWall,
					Position: Point3D{X: float32(x), Y: 1, Z: float32(z)},
					Solid:    true,
					Height:   1.0,
				}
			}
		}
	}

	world.Tiles[5][1][5] = Tile3D{
		Type:     TilePillar,
		Position: Point3D{X: 5, Y: 1, Z: 5},
		Solid:    true,
		Height:   1.0,
	}
	world.Tiles[6][1][6] = Tile3D{
		Type:     TilePillar,
		Position: Point3D{X: 6, Y: 1, Z: 6},
		Solid:    true,
		Height:   1.0,
	}

	world.Tiles[4][1][5] = Tile3D{
		Type:     TileDoor,
		Position: Point3D{X: 4, Y: 1, Z: 5},
		Solid:    false,
		Height:   1.0,
	}
	world.Tiles[7][1][6] = Tile3D{
		Type:     TileDoor,
		Position: Point3D{X: 7, Y: 1, Z: 6},
		Solid:    false,
		Height:   1.0,
	}

	positions := []Point3D{
		{X: 2, Y: 1, Z: 3}, {X: 2, Y: 1, Z: 4},
		{X: 9, Y: 1, Z: 2}, {X: 10, Y: 1, Z: 2},
		{X: 3, Y: 1, Z: 9}, {X: 8, Y: 1, Z: 8},
	}

	for i, pos := range positions {
		world.Tiles[int(pos.X)][int(pos.Y)][int(pos.Z)] = Tile3D{
			Type:     TileType(int(TileStoneWall) + i%4),
			Position: pos,
			Solid:    true,
			Height:   1.0,
		}
	}

	for x := 4; x <= 7; x++ {
		for z := 4; z <= 7; z++ {
			world.Tiles[x][2][z] = Tile3D{
				Type:     TileCeiling,
				Position: Point3D{X: float32(x), Y: 2, Z: float32(z)},
				Solid:    false,
				Height:   1.0,
			}
		}
	}

	entityID := 0

	itemPositions := []Point3D{
		{X: 2, Y: 1, Z: 2}, {X: 9, Y: 1, Z: 3},
		{X: 3, Y: 1, Z: 8}, {X: 8, Y: 1, Z: 9},
		{X: 1, Y: 1, Z: 6}, {X: 10, Y: 1, Z: 5},
	}

	for i, pos := range itemPositions {
		entity := GameEntity{
			ID:       entityID,
			Type:     EntityItem,
			Position: pos,
			Bounds: BoundingBox3D{
				Min: Point3D{X: pos.X - 0.3, Y: pos.Y - 0.3, Z: pos.Z - 0.3},
				Max: Point3D{X: pos.X + 0.3, Y: pos.Y + 0.3, Z: pos.Z + 0.3},
			},
			SpriteID:  i % 6,
			Active:    true,
			Color:     rl.White,
			Health:    1,
			MaxHealth: 1,
		}
		world.Entities = append(world.Entities, entity)
		entityID++
	}

	enemyPositions := []Point3D{
		{X: 3, Y: 1, Z: 2}, {X: 8, Y: 1, Z: 3},
		{X: 2, Y: 1, Z: 8}, {X: 9, Y: 1, Z: 9},
	}

	for i, pos := range enemyPositions {
		entity := GameEntity{
			ID:        entityID,
			Type:      EntityEnemy,
			Position:  pos,
			Direction: DirDown,
			Bounds: BoundingBox3D{
				Min: Point3D{X: pos.X - 0.4, Y: pos.Y - 0.4, Z: pos.Z - 0.4},
				Max: Point3D{X: pos.X + 0.4, Y: pos.Y + 0.4, Z: pos.Z + 0.4},
			},
			SpriteID:  i % 4,
			Active:    true,
			Color:     rl.White,
			Health:    3,
			MaxHealth: 3,
		}
		world.Entities = append(world.Entities, entity)
		entityID++
	}

	player := GameEntity{
		ID:        entityID,
		Type:      EntityPlayer,
		Position:  world.PlayerSpawn,
		Direction: DirDown,
		Bounds: BoundingBox3D{
			Min: Point3D{X: world.PlayerSpawn.X - 0.4, Y: world.PlayerSpawn.Y - 0.4, Z: world.PlayerSpawn.Z - 0.4},
			Max: Point3D{X: world.PlayerSpawn.X + 0.4, Y: world.PlayerSpawn.Y + 0.4, Z: world.PlayerSpawn.Z + 0.4},
		},
		Active:    true,
		Color:     rl.White,
		Health:    10,
		MaxHealth: 10,
	}

	world.Entities = append(world.Entities, player)
	fg.Player = &world.Entities[len(world.Entities)-1]

	fg.World = world

	fmt.Printf("Built 3D world: %dx%dx%d with %d entities\n", world.Width, world.Height, world.Depth, len(world.Entities))
}

func (fg *FilmationGame) IsPositionSolid(pos Point3D) bool {
	x, y, z := int(pos.X), int(pos.Y), int(pos.Z)

	if x < 0 || x >= fg.World.Width || y < 0 || y >= fg.World.Height || z < 0 || z >= fg.World.Depth {
		return true
	}

	tile := &fg.World.Tiles[x][y][z]
	if tile.Type != TileEmpty && tile.Solid {
		return true
	}

	checkBounds := BoundingBox3D{
		Min: Point3D{X: pos.X - 0.4, Y: pos.Y - 0.4, Z: pos.Z - 0.4},
		Max: Point3D{X: pos.X + 0.4, Y: pos.Y + 0.4, Z: pos.Z + 0.4},
	}

	for i := range fg.World.Entities {
		entity := &fg.World.Entities[i]
		if entity.Active && entity.Type == EntityEnemy {
			if BoundingBoxesIntersect(checkBounds, entity.Bounds) {
				return true
			}
		}
	}

	return false
}