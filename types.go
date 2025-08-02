package main

import rl "github.com/gen2brain/raylib-go/raylib"

type Point3D struct {
	X, Y, Z float32
}

type Point2D struct {
	X, Y float32
}

type BoundingBox3D struct {
	Min, Max Point3D
}

type Direction int

const (
	DirDown Direction = iota
	DirLeft
	DirUp
	DirRight
)

type TileType int

const (
	TileEmpty TileType = iota
	TileStoneFloor
	TileWoodFloor
	TileGrassFloor
	TileSandFloor
	TileStoneWall
	TileBrickWall
	TileWoodWall
	TileMetalWall
	TilePillar
	TileStairs
	TileDoor
	TileCeiling
)

type Tile3D struct {
	Type     TileType
	Position Point3D
	Solid    bool
	Height   float32
}

type EntityType int

const (
	EntityPlayer EntityType = iota
	EntityItem
	EntityEnemy
	EntityNPC
	EntityEffect
)

type GameEntity struct {
	ID        int
	Type      EntityType
	Position  Point3D
	Bounds    BoundingBox3D
	SpriteID  int
	Direction Direction
	Color     rl.Color
	Active    bool
	Health    int
	MaxHealth int
	Frame     int
	AnimSpeed float32
	
	TargetPosition Point3D
	IsMoving       bool
	MoveSpeed      float32
	MoveTimer      float32
}

type SpriteCache struct {
	FloorTiles  [4]rl.Texture2D
	WallTiles   [4]rl.Texture2D
	PillarTile  rl.Texture2D
	StairsTile  rl.Texture2D
	DoorTiles   [2]rl.Texture2D
	CeilingTile rl.Texture2D

	PlayerSprites [4]rl.Texture2D
	ItemSprites   [6]rl.Texture2D
	EnemySprites  [4]rl.Texture2D
}

type World3D struct {
	Width, Height, Depth int
	Tiles                [][][]Tile3D
	Entities             []GameEntity
	PlayerSpawn          Point3D
}

type FilmationGame struct {
	Sprites SpriteCache
	World   World3D
	Player  *GameEntity
	Camera  Point3D
	Scale   float32
	ScreenW int32
	ScreenH int32

	Rooms RoomManager

	AssetPath string

	// Music support - this is the key addition
	BackgroundMusic rl.Music

	RenderOrder []RenderItem

	GameTime   float32
	InputDelay float32
	AnimTime   float32

	ItemsCollected int
	EnemiesKilled  int
}

type RenderItem struct {
	Position Point3D
	Depth    float32
	Type     string
	TileData *Tile3D
	EntityID int
}

func BoundingBoxesIntersect(a, b BoundingBox3D) bool {
	return a.Min.X <= b.Max.X && a.Max.X >= b.Min.X &&
		a.Min.Y <= b.Max.Y && a.Max.Y >= b.Min.Y &&
		a.Min.Z <= b.Max.Z && a.Max.Z >= b.Min.Z
}

func PointsNearlyEqual(a, b Point3D, tolerance float32) bool {
	dx := a.X - b.X
	dy := a.Y - b.Y  
	dz := a.Z - b.Z
	if dx < 0 { dx = -dx }
	if dy < 0 { dy = -dy }
	if dz < 0 { dz = -dz }
	return dx < tolerance && dy < tolerance && dz < tolerance
}

func LerpPoint3D(a, b Point3D, t float32) Point3D {
	if t <= 0 { return a }
	if t >= 1 { return b }
	return Point3D{
		X: a.X + (b.X - a.X) * t,
		Y: a.Y + (b.Y - a.Y) * t,
		Z: a.Z + (b.Z - a.Z) * t,
	}
}