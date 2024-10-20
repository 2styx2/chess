package main

import rl "github.com/gen2brain/raylib-go/raylib"

const (
	squareSize = 100
	screenSize = 800
)

func CreateBoard() {
	for file := 0; file < 8; file++ {
		for rank := 0; rank < 8; rank++ {
			if (file+rank)%2 == 0 {
				rl.DrawRectangle(int32(file*squareSize), int32(rank*squareSize), squareSize, squareSize, rl.LightGray)
			} else {
				rl.DrawRectangle(int32(file*squareSize), int32(rank*squareSize), squareSize, squareSize, rl.DarkGreen)
			}
		}
	}
}

func main() {
	rl.InitWindow(screenSize, screenSize, "raylib [core] example - basic window")

	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)
		CreateBoard()
		rl.DrawText("text string", 0, 800, 100, rl.Black)

		rl.EndDrawing()
	}
}
