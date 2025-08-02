package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (fg *FilmationGame) WorldToScreen(p Point3D) Point2D {
	rotatedX := p.X
	rotatedZ := p.Z
	rotatedY := p.Y

	worldCenterX := float32(fg.World.Width-1) / 2.0
	worldCenterZ := float32(fg.World.Depth-1) / 2.0
	worldCenterY := float32(fg.World.Height-1) / 2.0

	rotatedCenterX := worldCenterX
	rotatedCenterZ := worldCenterZ
	rotatedCenterY := worldCenterY

	tileWidth := fg.Scale
	tileHeight := fg.Scale * 0.5

	projectedX := (rotatedX - rotatedZ) * tileWidth * 0.5
	projectedY := (rotatedX+rotatedZ)*tileHeight*0.5 - rotatedY*tileHeight

	projectedCenterX := (rotatedCenterX - rotatedCenterZ) * tileWidth * 0.5
	projectedCenterY := (rotatedCenterX+rotatedCenterZ)*tileHeight*0.5 - rotatedCenterY*tileHeight

	verticalOffset := float32(50)
	screenX := projectedX - projectedCenterX + float32(fg.ScreenW)/2
	screenY := projectedY - projectedCenterY + float32(fg.ScreenH)/2 - verticalOffset

	return Point2D{X: screenX, Y: screenY}
}

func (fg *FilmationGame) CalculateRenderOrder() {
	fg.RenderOrder = nil

	for x := 0; x < fg.World.Width; x++ {
		for y := 0; y < fg.World.Height; y++ {
			for z := 0; z < fg.World.Depth; z++ {
				tile := &fg.World.Tiles[x][y][z]
				if tile.Type != TileEmpty {
					depth := tile.Position.X + tile.Position.Z + tile.Position.Y*2

					renderItem := RenderItem{
						Position: tile.Position,
						Depth:    depth,
						Type:     "tile",
						TileData: tile,
						EntityID: -1,
					}
					fg.RenderOrder = append(fg.RenderOrder, renderItem)
				}
			}
		}
	}

	for i, entity := range fg.World.Entities {
		if entity.Active {
			depth := entity.Position.X + entity.Position.Z + entity.Position.Y*2

			renderItem := RenderItem{
				Position: entity.Position,
				Depth:    depth,
				Type:     "entity",
				TileData: nil,
				EntityID: i,
			}
			fg.RenderOrder = append(fg.RenderOrder, renderItem)
		}
	}

	for i := 0; i < len(fg.RenderOrder); i++ {
		for j := i + 1; j < len(fg.RenderOrder); j++ {
			if fg.RenderOrder[i].Depth > fg.RenderOrder[j].Depth {
				fg.RenderOrder[i], fg.RenderOrder[j] = fg.RenderOrder[j], fg.RenderOrder[i]
			}
		}
	}
}

func (fg *FilmationGame) RenderTile(tile *Tile3D) {
	screenPos := fg.WorldToScreen(tile.Position)

	var texture rl.Texture2D

	switch tile.Type {
	case TileStoneFloor, TileWoodFloor, TileGrassFloor, TileSandFloor:
		texture = fg.Sprites.FloorTiles[int(tile.Type)]
	case TileStoneWall, TileBrickWall, TileWoodWall, TileMetalWall:
		texture = fg.Sprites.WallTiles[int(tile.Type)-int(TileStoneWall)]
	case TilePillar:
		texture = fg.Sprites.PillarTile
	case TileStairs:
		texture = fg.Sprites.StairsTile
	case TileDoor:
		texture = fg.Sprites.DoorTiles[0]
	case TileCeiling:
		texture = fg.Sprites.CeilingTile
	default:
		return
	}

	renderX := screenPos.X - float32(texture.Width)/2
	renderY := screenPos.Y - float32(texture.Height)/2

	switch tile.Type {
	case TileStoneFloor, TileWoodFloor, TileGrassFloor, TileSandFloor:
		renderY += 16
	case TileStoneWall, TileBrickWall, TileWoodWall, TileMetalWall:
		renderY -= 20
	case TilePillar:
		renderY -= 16
	case TileStairs:
		renderY -= 12
	case TileDoor:
		renderY -= 20
	case TileCeiling:
		renderY -= 40
	}

	rl.DrawTexture(texture, int32(renderX), int32(renderY), rl.White)
}

func (fg *FilmationGame) RenderEntity(entityID int) {
	entity := &fg.World.Entities[entityID]
	if !entity.Active {
		return
	}

	screenPos := fg.WorldToScreen(entity.Position)

	var texture rl.Texture2D

	switch entity.Type {
	case EntityPlayer:
		texture = fg.Sprites.PlayerSprites[entity.Direction]
	case EntityItem:
		texture = fg.Sprites.ItemSprites[entity.SpriteID]
	case EntityEnemy:
		texture = fg.Sprites.EnemySprites[entity.SpriteID]
	default:
		return
	}

	renderX := screenPos.X - float32(texture.Width)/2
	renderY := screenPos.Y - float32(texture.Height)/2

	switch entity.Type {
	case EntityPlayer, EntityEnemy:
		renderY -= 8
	case EntityItem:
		renderY -= 4
	}

	color := entity.Color
	if entity.Type == EntityEnemy && entity.Health < entity.MaxHealth {
		color = rl.Color{R: 255, G: 150, B: 150, A: 255}
	}

	rl.DrawTexture(texture, int32(renderX), int32(renderY), color)

	if entity.Type == EntityEnemy && entity.MaxHealth > 0 {
		barWidth := float32(16)
		healthPercent := float32(entity.Health) / float32(entity.MaxHealth)
		rl.DrawRectangle(int32(screenPos.X-barWidth/2), int32(screenPos.Y-25), int32(barWidth), 2, rl.Color{R: 100, G: 100, B: 100, A: 200})
		rl.DrawRectangle(int32(screenPos.X-barWidth/2), int32(screenPos.Y-25), int32(barWidth*healthPercent), 2, rl.Color{R: 255, G: 0, B: 0, A: 255})
	}
}

func (fg *FilmationGame) Render() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.Color{R: 20, G: 25, B: 35, A: 255})

	for _, item := range fg.RenderOrder {
		if item.Type == "tile" {
			fg.RenderTile(item.TileData)
		} else if item.Type == "entity" {
			fg.RenderEntity(item.EntityID)
		}
	}

	rl.DrawText("RETROMANSION", 10, 10, 20, rl.White)
	rl.DrawText("WASD: MOVE | SPACE: ATTACK", 10, 45, 10, rl.LightGray)

	rl.DrawText(fmt.Sprintf("HEALTH: %d/%d", fg.Player.Health, fg.Player.MaxHealth), 10, 60, 10, rl.Color{R: 255, G: 100, B: 100, A: 255})
	rl.DrawText(fmt.Sprintf("ITEMS: %d | Enemies: %d", fg.ItemsCollected, fg.EnemiesKilled), 10, 75, 10, rl.White)
	rl.DrawText(fmt.Sprintf("POS: (%.0f,%.0f,%.0f)", fg.Player.Position.X, fg.Player.Position.Y, fg.Player.Position.Z), 10, 90, 10, rl.Gray)

	rl.DrawText(fmt.Sprintf("WORLD: %dx%dx%d | RENDERED: %d", fg.World.Width, fg.World.Height, fg.World.Depth, len(fg.RenderOrder)), 10, fg.ScreenH-30, 10, rl.DarkGray)

	if fg.Player.Health <= 0 {
		rl.DrawRectangle(0, 0, fg.ScreenW, fg.ScreenH, rl.Color{R: 0, G: 0, B: 0, A: 180})
		rl.DrawText("GAME OVER", fg.ScreenW/2-120, fg.ScreenH/2-20, 36, rl.Color{R: 255, G: 0, B: 0, A: 255})
	}

	rl.EndDrawing()
}
