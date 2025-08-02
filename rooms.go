package main

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Room struct {
	ID          int
	Name        string
	World       World3D
	Connections []RoomConnection
}

type RoomConnection struct {
	Position    Point3D
	ToRoomID    int
	ToPosition  Point3D
	Direction   Direction
	RequiresKey bool
	Active      bool
}

type RoomManager struct {
	Rooms       map[int]*Room
	CurrentRoom int
	PlayerPos   Point3D
}

func (fg *FilmationGame) InitRoomSystem() {
	fg.Rooms = RoomManager{
		Rooms:       make(map[int]*Room),
		CurrentRoom: 1,
	}
}

func (fg *FilmationGame) CreateRoom(id int, name string, width, height, depth int) *Room {
	room := &Room{
		ID:   id,
		Name: name,
		World: World3D{
			Width:  width,
			Height: height,
			Depth:  depth,
		},
		Connections: []RoomConnection{},
	}

	room.World.Tiles = make([][][]Tile3D, room.World.Width)
	for x := range room.World.Tiles {
		room.World.Tiles[x] = make([][]Tile3D, room.World.Height)
		for y := range room.World.Tiles[x] {
			room.World.Tiles[x][y] = make([]Tile3D, room.World.Depth)
		}
	}

	fg.Rooms.Rooms[id] = room
	return room
}

func (fg *FilmationGame) BuildBasicRoom(room *Room, floorType TileType) {
	for x := 0; x < room.World.Width; x++ {
		for z := 0; z < room.World.Depth; z++ {
			room.World.Tiles[x][0][z] = Tile3D{
				Type:     floorType,
				Position: Point3D{X: float32(x), Y: 0, Z: float32(z)},
				Solid:    false,
				Height:   1.0,
			}
		}
	}

	for x := 0; x < room.World.Width; x++ {
		for z := 0; z < room.World.Depth; z++ {
			if x == 0 || x == room.World.Width-1 || z == 0 || z == room.World.Depth-1 {
				room.World.Tiles[x][1][z] = Tile3D{
					Type:     TileStoneWall,
					Position: Point3D{X: float32(x), Y: 1, Z: float32(z)},
					Solid:    true,
					Height:   1.0,
				}
			}
		}
	}
}

func (fg *FilmationGame) AddRoomConnection(fromRoomID, toRoomID int, fromPos, toPos Point3D, direction Direction, requiresKey bool) {
	room := fg.Rooms.Rooms[fromRoomID]
	if room == nil {
		return
	}

	connection := RoomConnection{
		Position:    fromPos,
		ToRoomID:    toRoomID,
		ToPosition:  toPos,
		Direction:   direction,
		RequiresKey: requiresKey,
		Active:      true,
	}
	room.Connections = append(room.Connections, connection)

	x, y, z := int(fromPos.X), int(fromPos.Y), int(fromPos.Z)
	if x >= 0 && x < room.World.Width && y >= 0 && y < room.World.Height && z >= 0 && z < room.World.Depth {
		room.World.Tiles[x][y][z] = Tile3D{
			Type:     TileDoor,
			Position: fromPos,
			Solid:    false,
			Height:   1.0,
		}
	}
}

func (fg *FilmationGame) CheckRoomTransitions() {
	currentRoom := fg.Rooms.Rooms[fg.Rooms.CurrentRoom]
	if currentRoom == nil {
		return
	}

	playerGridPos := Point3D{
		X: float32(int(fg.Player.Position.X + 0.5)),
		Y: float32(int(fg.Player.Position.Y + 0.5)),
		Z: float32(int(fg.Player.Position.Z + 0.5)),
	}

	for _, connection := range currentRoom.Connections {
		if !connection.Active {
			continue
		}

		if playerGridPos.X == connection.Position.X &&
			playerGridPos.Y == connection.Position.Y &&
			playerGridPos.Z == connection.Position.Z {

			if connection.RequiresKey && fg.ItemsCollected < 1 {
				fmt.Println("A key is required to enter this room!")
				return
			}

			fg.TransitionToRoom(connection.ToRoomID, connection.ToPosition, connection.Direction)
			break
		}
	}
}

func (fg *FilmationGame) TransitionToRoom(roomID int, newPos Point3D, direction Direction) {
	newRoom := fg.Rooms.Rooms[roomID]
	if newRoom == nil {
		fmt.Printf("Error: Room %d not found\n", roomID)
		return
	}

	fmt.Printf("Transitioning to room: %s\n", newRoom.Name)

	fg.Rooms.CurrentRoom = roomID
	fg.World = newRoom.World

	fg.Player.Position = newPos
	fg.Player.TargetPosition = newPos
	fg.Player.IsMoving = false
	fg.Player.Direction = direction
	fg.UpdatePlayerBounds()

	fg.CalculateRenderOrder()
}

func (fg *FilmationGame) SetupPlayerInRoom(roomID int, position Point3D) {
	room := fg.Rooms.Rooms[roomID]
	if room == nil {
		return
	}

	playerEntity := GameEntity{
		ID:        999,
		Type:      EntityPlayer,
		Position:  position,
		Direction: DirDown,
		Bounds: BoundingBox3D{
			Min: Point3D{X: position.X - 0.4, Y: position.Y - 0.4, Z: position.Z - 0.4},
			Max: Point3D{X: position.X + 0.4, Y: position.Y + 0.4, Z: position.Z + 0.4},
		},
		Active:    true,
		Color:     rl.White,
		Health:    10,
		MaxHealth: 10,
		
		TargetPosition: position,
		IsMoving:       false,
		MoveSpeed:      4.0,
		MoveTimer:      0.0,
	}

	room.World.Entities = append(room.World.Entities, playerEntity)

	fg.Player = &room.World.Entities[len(room.World.Entities)-1]

	fmt.Printf("Player setup: Health=%d, Position=(%.1f,%.1f,%.1f)\n",
		fg.Player.Health, fg.Player.Position.X, fg.Player.Position.Y, fg.Player.Position.Z)
}

func (fg *FilmationGame) BuildSampleRooms() {
	fmt.Println("Building room-based world...")

	fg.InitRoomSystem()

	room1 := fg.CreateRoom(1, "Starting Chamber", 10, 3, 10)
	fg.BuildBasicRoom(room1, TileStoneFloor)

	room1.World.Tiles[2][1][2] = Tile3D{
		Type:     TilePillar,
		Position: Point3D{X: 2, Y: 1, Z: 2},
		Solid:    true,
		Height:   1.0,
	}
	room1.World.Tiles[7][1][7] = Tile3D{
		Type:     TilePillar,
		Position: Point3D{X: 7, Y: 1, Z: 7},
		Solid:    true,
		Height:   1.0,
	}

	room2 := fg.CreateRoom(2, "Treasure Vault", 6, 3, 6)
	fg.BuildBasicRoom(room2, TileGrassFloor)

	room3 := fg.CreateRoom(3, "Ancient Library", 10, 3, 6)
	fg.BuildBasicRoom(room3, TileWoodFloor)

	for x := 3; x <= 6; x++ {
		room3.World.Tiles[x][1][3] = Tile3D{
			Type:     TileBrickWall,
			Position: Point3D{X: float32(x), Y: 1, Z: 3},
			Solid:    true,
			Height:   1.0,
		}
	}

	fg.AddRoomConnection(1, 2, Point3D{X: 9, Y: 1, Z: 5}, Point3D{X: 1, Y: 1, Z: 3}, DirRight, false)
	fg.AddRoomConnection(2, 1, Point3D{X: 0, Y: 1, Z: 3}, Point3D{X: 8, Y: 1, Z: 5}, DirLeft, false)

	fg.AddRoomConnection(1, 3, Point3D{X: 5, Y: 1, Z: 0}, Point3D{X: 5, Y: 1, Z: 5}, DirUp, true)
	fg.AddRoomConnection(3, 1, Point3D{X: 5, Y: 1, Z: 5}, Point3D{X: 5, Y: 1, Z: 1}, DirDown, false)

	fg.AddRoomEntities()

	startPos := Point3D{X: 5, Y: 1, Z: 8}
	fg.SetupPlayerInRoom(1, startPos)

	fg.World = room1.World

	fmt.Printf("Built %d rooms, player health: %d\n", len(fg.Rooms.Rooms), fg.Player.Health)
}

func (fg *FilmationGame) AddRoomEntities() {
	entityID := 0

	room1 := fg.Rooms.Rooms[1]

	keyEntity := GameEntity{
		ID:       entityID,
		Type:     EntityItem,
		Position: Point3D{X: 1, Y: 1, Z: 1},
		Bounds: BoundingBox3D{
			Min: Point3D{X: 0.7, Y: 0.7, Z: 0.7},
			Max: Point3D{X: 1.3, Y: 1.3, Z: 1.3},
		},
		SpriteID:  0,
		Active:    true,
		Color:     rl.White,
		Health:    1,
		MaxHealth: 1,
		TargetPosition: Point3D{X: 1, Y: 1, Z: 1},
		IsMoving:       false,
		MoveSpeed:      0,
		MoveTimer:      0,
	}
	room1.World.Entities = append(room1.World.Entities, keyEntity)
	entityID++

	enemyPos := Point3D{X: 1, Y: 1, Z: 2}
	enemyEntity := GameEntity{
		ID:        entityID,
		Type:      EntityEnemy,
		Position:  enemyPos,
		Direction: DirDown,
		Bounds: BoundingBox3D{
			Min: Point3D{X: 0.6, Y: 0.6, Z: 1.6},
			Max: Point3D{X: 1.4, Y: 1.4, Z: 2.4},
		},
		SpriteID:  0,
		Active:    true,
		Color:     rl.White,
		Health:    3,
		MaxHealth: 3,
		TargetPosition: enemyPos,
		IsMoving:       false,
		MoveSpeed:      2.0,
		MoveTimer:      1.0,
	}
	room1.World.Entities = append(room1.World.Entities, enemyEntity)
	entityID++

	room2 := fg.Rooms.Rooms[2]
	gemPos := Point3D{X: 3, Y: 1, Z: 3}
	gemEntity := GameEntity{
		ID:       entityID,
		Type:     EntityItem,
		Position: gemPos,
		Bounds: BoundingBox3D{
			Min: Point3D{X: 2.7, Y: 0.7, Z: 2.7},
			Max: Point3D{X: 3.3, Y: 1.3, Z: 3.3},
		},
		SpriteID:  1,
		Active:    true,
		Color:     rl.White,
		Health:    1,
		MaxHealth: 1,
		TargetPosition: gemPos,
		IsMoving:       false,
		MoveSpeed:      0,
		MoveTimer:      0,
	}
	room2.World.Entities = append(room2.World.Entities, gemEntity)
	entityID++

	room3 := fg.Rooms.Rooms[3]
	potionPos := Point3D{X: 8, Y: 1, Z: 2}
	potionEntity := GameEntity{
		ID:       entityID,
		Type:     EntityItem,
		Position: potionPos,
		Bounds: BoundingBox3D{
			Min: Point3D{X: 7.7, Y: 0.7, Z: 1.7},
			Max: Point3D{X: 8.3, Y: 1.3, Z: 2.3},
		},
		SpriteID:  2,
		Active:    true,
		Color:     rl.White,
		Health:    1,
		MaxHealth: 1,
		TargetPosition: potionPos,
		IsMoving:       false,
		MoveSpeed:      0,
		MoveTimer:      0,
	}
	room3.World.Entities = append(room3.World.Entities, potionEntity)
	entityID++
}

func (fg *FilmationGame) UpdateEnemies() {
	deltaTime := rl.GetFrameTime()
	
	for i := range fg.World.Entities {
		entity := &fg.World.Entities[i]
		if entity.Type == EntityEnemy && entity.Active {
			entity.MoveTimer -= deltaTime
			
			if !entity.IsMoving && entity.MoveTimer <= 0 {
				fmt.Printf("Enemy %d considering move from (%.1f,%.1f,%.1f) toward player at (%.1f,%.1f,%.1f)\n",
					i, entity.Position.X, entity.Position.Y, entity.Position.Z,
					fg.Player.Position.X, fg.Player.Position.Y, fg.Player.Position.Z)

				newPos := entity.Position

				dx := fg.Player.Position.X - entity.Position.X
				dz := fg.Player.Position.Z - entity.Position.Z

				if math.Abs(float64(dx)) > math.Abs(float64(dz)) {
					if dx > 0 {
						newPos.X += 1.0
						entity.Direction = DirRight
					} else {
						newPos.X -= 1.0
						entity.Direction = DirLeft
					}
				} else {
					if dz > 0 {
						newPos.Z += 1.0
						entity.Direction = DirDown
					} else {
						newPos.Z -= 1.0
						entity.Direction = DirUp
					}
				}

				if !fg.IsPositionSolid(newPos) {
					entity.TargetPosition = newPos
					entity.IsMoving = true
					fmt.Printf("Enemy %d starting move to (%.1f,%.1f,%.1f)\n", i, newPos.X, newPos.Y, newPos.Z)
				}
				
				entity.MoveTimer = 3.0
			}
			
			if !entity.IsMoving && BoundingBoxesIntersect(entity.Bounds, fg.Player.Bounds) {
				fg.Player.Health--
				fmt.Printf("PLAYER HIT! Health now: %d\n", fg.Player.Health)
			}
		}
	}
}

func abs(x float32) float32 {
	if x < 0 {
		return -x
	}
	return x
}