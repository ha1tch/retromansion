package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (fg *FilmationGame) HandleInput() {
	if fg.InputDelay > 0 {
		fg.InputDelay -= rl.GetFrameTime()
		return
	}

	if fg.Player.IsMoving && !PointsNearlyEqual(fg.Player.Position, fg.Player.TargetPosition, 0.1) {
		return
	}

	moved := false
	newTargetPos := fg.Player.Position

	if fg.Player.IsMoving {
		fg.Player.Position = fg.Player.TargetPosition
		fg.Player.IsMoving = false
		fg.UpdatePlayerBounds()
	}

	if rl.IsKeyPressed(rl.KeyLeft) || rl.IsKeyPressed(rl.KeyA) {
		newTargetPos.X -= 1.0
		fg.Player.Direction = DirLeft
		moved = true
	} else if rl.IsKeyPressed(rl.KeyRight) || rl.IsKeyPressed(rl.KeyD) {
		newTargetPos.X += 1.0
		fg.Player.Direction = DirRight
		moved = true
	}

	if rl.IsKeyPressed(rl.KeyUp) || rl.IsKeyPressed(rl.KeyW) {
		newTargetPos.Z -= 1.0
		fg.Player.Direction = DirUp
		moved = true
	} else if rl.IsKeyPressed(rl.KeyDown) || rl.IsKeyPressed(rl.KeyS) {
		newTargetPos.Z += 1.0
		fg.Player.Direction = DirDown
		moved = true
	}

	if rl.IsKeyPressed(rl.KeySpace) {
		fg.PlayerAttack()
		fg.InputDelay = 0.2
		return
	}

	if moved {
		fg.InputDelay = 0.05

		if !fg.IsPositionSolid(newTargetPos) {
			fg.Player.TargetPosition = newTargetPos
			fg.Player.IsMoving = true
		}
	}
}

func (fg *FilmationGame) UpdatePlayerBounds() {
	pos := fg.Player.Position
	fg.Player.Bounds = BoundingBox3D{
		Min: Point3D{X: pos.X - 0.4, Y: pos.Y - 0.4, Z: pos.Z - 0.4},
		Max: Point3D{X: pos.X + 0.4, Y: pos.Y + 0.4, Z: pos.Z + 0.4},
	}
}

func (fg *FilmationGame) UpdateMovement() {
	deltaTime := rl.GetFrameTime()

	for i := range fg.World.Entities {
		entity := &fg.World.Entities[i]
		if !entity.Active || !entity.IsMoving {
			continue
		}

		moveSpeed := entity.MoveSpeed
		if moveSpeed <= 0 {
			moveSpeed = 4.0
		}

		moveDistance := moveSpeed * deltaTime

		dx := entity.TargetPosition.X - entity.Position.X
		dy := entity.TargetPosition.Y - entity.Position.Y
		dz := entity.TargetPosition.Z - entity.Position.Z

		totalDistance := float32(0)
		if dx != 0 {
			totalDistance += dx * dx
		}
		if dy != 0 {
			totalDistance += dy * dy
		}
		if dz != 0 {
			totalDistance += dz * dz
		}

		if totalDistance > 0 {
			totalDistance = float32(sqrt(float64(totalDistance)))
		}

		if moveDistance >= totalDistance {
			entity.Position = entity.TargetPosition
			entity.IsMoving = false

			if entity.Type == EntityPlayer {
				fg.UpdatePlayerBounds()
				fg.CheckInteractions()
			} else {
				fg.UpdateEntityBounds(entity)
			}
		} else {
			if totalDistance > 0 {
				moveRatio := moveDistance / totalDistance
				entity.Position.X += dx * moveRatio
				entity.Position.Y += dy * moveRatio
				entity.Position.Z += dz * moveRatio

				if entity.Type == EntityPlayer {
					fg.UpdatePlayerBounds()
				} else {
					fg.UpdateEntityBounds(entity)
				}
			}
		}
	}
}

func (fg *FilmationGame) UpdateEntityBounds(entity *GameEntity) {
	pos := entity.Position
	entity.Bounds = BoundingBox3D{
		Min: Point3D{X: pos.X - 0.4, Y: pos.Y - 0.4, Z: pos.Z - 0.4},
		Max: Point3D{X: pos.X + 0.4, Y: pos.Y + 0.4, Z: pos.Z + 0.4},
	}
}

func (fg *FilmationGame) CheckInteractions() {
	for i := range fg.World.Entities {
		entity := &fg.World.Entities[i]
		if entity.Type == EntityItem && entity.Active {
			if BoundingBoxesIntersect(fg.Player.Bounds, entity.Bounds) {
				entity.Active = false
				fg.ItemsCollected++
				fmt.Printf("Picked up %s!\n", fg.getItemName(entity.SpriteID))
			}
		}
	}

	fg.CheckRoomTransitions()
}

func (fg *FilmationGame) PlayerAttack() {
	attackPos := fg.Player.Position

	switch fg.Player.Direction {
	case DirLeft:
		attackPos.X -= 1.0
	case DirRight:
		attackPos.X += 1.0
	case DirUp:
		attackPos.Z -= 1.0
	case DirDown:
		attackPos.Z += 1.0
	}

	attackBounds := BoundingBox3D{
		Min: Point3D{X: attackPos.X - 0.5, Y: attackPos.Y - 0.5, Z: attackPos.Z - 0.5},
		Max: Point3D{X: attackPos.X + 0.5, Y: attackPos.Y + 0.5, Z: attackPos.Z + 0.5},
	}

	for i := range fg.World.Entities {
		entity := &fg.World.Entities[i]
		if entity.Type == EntityEnemy && entity.Active {
			if BoundingBoxesIntersect(attackBounds, entity.Bounds) {
				entity.Health--
				if entity.Health <= 0 {
					entity.Active = false
					fg.EnemiesKilled++
					fmt.Printf("Enemy defeated! Total: %d\n", fg.EnemiesKilled)
				}
				break
			}
		}
	}
}

func (fg *FilmationGame) getItemName(spriteID int) string {
	names := []string{"Key", "Gem", "Potion", "Sword", "Shield", "Apple"}
	if spriteID >= 0 && spriteID < len(names) {
		return names[spriteID]
	}
	return "Item"
}

func sqrt(x float64) float64 {
	if x == 0 {
		return 0
	}
	z := x
	for i := 0; i < 10; i++ {
		z = (z + x/z) / 2
	}
	return z
}

func (fg *FilmationGame) Update() {
	fg.GameTime += rl.GetFrameTime()
	fg.AnimTime += rl.GetFrameTime()

	fg.HandleInput()
	fg.UpdateMovement()
	fg.UpdateEnemies()
	fg.CalculateRenderOrder()
}

func main() {
	const screenWidth = 800
	const screenHeight = 600

	rl.InitWindow(screenWidth, screenHeight, "RETROMANSION")
	rl.SetTargetFPS(60)

	// Initialize audio system for music support
	rl.InitAudioDevice()
	defer rl.CloseAudioDevice()

	game := FilmationGame{
		Scale:     54,
		ScreenW:   screenWidth,
		ScreenH:   screenHeight,
		AssetPath: "./game_assets/sprites",
	}

	err := game.LoadSprites()
	if err != nil {
		fmt.Printf("Failed to load sprites: %v\n", err)
		rl.CloseWindow()
		return
	}

	// Load and start background music
	err = game.LoadMusic()
	if err != nil {
		fmt.Printf("Failed to load music: %v\n", err)
		// Continue without music - it's optional
	} else {
		game.StartMusic()
	}

	game.BuildSampleRooms()
	game.CalculateRenderOrder()

	fmt.Println("World ready!")

	for !rl.WindowShouldClose() {
		// Update music stream each frame
		game.UpdateMusic()
		
		game.Update()
		game.Render()
	}

	game.CleanupSprites()
	rl.CloseWindow()
}