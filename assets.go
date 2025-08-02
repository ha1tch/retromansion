package main

import (
	"fmt"
	"path/filepath"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (fg *FilmationGame) LoadSprites() error {
	fmt.Println("Loading sprites from asset files...")

	if fg.AssetPath == "" {
		fg.AssetPath = "./game_assets/sprites"
	}

	floorNames := []string{"floor_stone", "floor_wood", "floor_grass", "floor_sand"}
	for i, name := range floorNames {
		path := filepath.Join(fg.AssetPath, "tiles", "floors", name+".png")
		texture := rl.LoadTexture(path)
		if texture.ID == 0 {
			return fmt.Errorf("failed to load floor texture: %s", path)
		}
		fg.Sprites.FloorTiles[i] = texture
		fmt.Printf("  Loaded: %s\n", name)
	}

	wallNames := []string{"wall_stone", "wall_brick", "wall_wood", "wall_metal"}
	for i, name := range wallNames {
		path := filepath.Join(fg.AssetPath, "tiles", "walls", name+".png")
		texture := rl.LoadTexture(path)
		if texture.ID == 0 {
			return fmt.Errorf("failed to load wall texture: %s", path)
		}
		fg.Sprites.WallTiles[i] = texture
		fmt.Printf("  Loaded: %s\n", name)
	}

	specialSprites := []struct {
		name    string
		texture *rl.Texture2D
	}{
		{"pillar", &fg.Sprites.PillarTile},
		{"stairs", &fg.Sprites.StairsTile},
		{"door_closed", &fg.Sprites.DoorTiles[0]},
		{"door_open", &fg.Sprites.DoorTiles[1]},
		{"ceiling", &fg.Sprites.CeilingTile},
	}

	for _, special := range specialSprites {
		path := filepath.Join(fg.AssetPath, "tiles", "special", special.name+".png")
		texture := rl.LoadTexture(path)
		if texture.ID == 0 {
			return fmt.Errorf("failed to load special texture: %s", path)
		}
		*special.texture = texture
		fmt.Printf("  Loaded: %s\n", special.name)
	}

	playerDirections := []string{"down", "left", "up", "right"}
	for i, direction := range playerDirections {
		path := filepath.Join(fg.AssetPath, "entities", "player", "player_"+direction+".png")
		texture := rl.LoadTexture(path)
		if texture.ID == 0 {
			return fmt.Errorf("failed to load player texture: %s", path)
		}
		fg.Sprites.PlayerSprites[i] = texture
		fmt.Printf("  Loaded: player_%s\n", direction)
	}

	itemNames := []string{"key", "gem", "potion", "sword", "shield", "food"}
	for i, name := range itemNames {
		path := filepath.Join(fg.AssetPath, "entities", "items", "item_"+name+".png")
		texture := rl.LoadTexture(path)
		if texture.ID == 0 {
			return fmt.Errorf("failed to load item texture: %s", path)
		}
		fg.Sprites.ItemSprites[i] = texture
		fmt.Printf("  Loaded: item_%s\n", name)
	}

	enemyNames := []string{"goblin", "orc", "troll", "skeleton"}
	for i, name := range enemyNames {
		path := filepath.Join(fg.AssetPath, "entities", "enemies", "enemy_"+name+".png")
		texture := rl.LoadTexture(path)
		if texture.ID == 0 {
			return fmt.Errorf("failed to load enemy texture: %s", path)
		}
		fg.Sprites.EnemySprites[i] = texture
		fmt.Printf("  Loaded: enemy_%s\n", name)
	}

	fmt.Println("All sprites loaded successfully!")
	return nil
}

func (fg *FilmationGame) LoadMusic() error {
	fmt.Println("Loading background music...")
	
	// Use consistent path structure like sprites do
	musicPath := filepath.Join("./game_assets", "music", "retromansion.wav")
	fg.BackgroundMusic = rl.LoadMusicStream(musicPath)
	
	if fg.BackgroundMusic.Stream.SampleRate == 0 {
		return fmt.Errorf("failed to load music: %s", musicPath)
	}
	
	fmt.Printf("  Loaded: %s\n", musicPath)
	return nil
}

func (fg *FilmationGame) StartMusic() {
	if fg.BackgroundMusic.Stream.SampleRate > 0 {
		rl.PlayMusicStream(fg.BackgroundMusic)
		fg.BackgroundMusic.Looping = true
		fmt.Println("Background music started")
	}
}

func (fg *FilmationGame) UpdateMusic() {
	if fg.BackgroundMusic.Stream.SampleRate > 0 {
		rl.UpdateMusicStream(fg.BackgroundMusic)
	}
}

func (fg *FilmationGame) CleanupSprites() {
	fmt.Println("Unloading sprites...")

	for i := 0; i < 4; i++ {
		rl.UnloadTexture(fg.Sprites.FloorTiles[i])
		rl.UnloadTexture(fg.Sprites.WallTiles[i])
		rl.UnloadTexture(fg.Sprites.PlayerSprites[i])
		rl.UnloadTexture(fg.Sprites.EnemySprites[i])
	}

	for i := 0; i < 6; i++ {
		rl.UnloadTexture(fg.Sprites.ItemSprites[i])
	}

	rl.UnloadTexture(fg.Sprites.PillarTile)
	rl.UnloadTexture(fg.Sprites.StairsTile)
	rl.UnloadTexture(fg.Sprites.DoorTiles[0])
	rl.UnloadTexture(fg.Sprites.DoorTiles[1])
	rl.UnloadTexture(fg.Sprites.CeilingTile)

	fmt.Println("All sprites unloaded")
}

func (fg *FilmationGame) CleanupAudio() {
	if fg.BackgroundMusic.Stream.SampleRate > 0 {
		rl.UnloadMusicStream(fg.BackgroundMusic)
		fmt.Println("Background music unloaded")
	}
}