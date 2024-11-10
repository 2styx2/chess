package main

import (
	"path/filepath"
	"strconv"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	squareSize = 64
	screenSize = 64 * 9
)

const (
	None = iota
	King
	Pawn
	Knight
	Bishop
	Rook
	Queen

	White = 8
	Black = 16
)

func textFromInt() map[int]rl.Texture2D {
	return map[int]rl.Texture2D{
		King | White:   rl.LoadTexture(filepath.Join("assets", "white-king.png")),
		King | Black:   rl.LoadTexture(filepath.Join("assets", "black-king.png")),
		Pawn | White:   rl.LoadTexture(filepath.Join("assets", "white-pawn.png")),
		Pawn | Black:   rl.LoadTexture(filepath.Join("assets", "black-pawn.png")),
		Knight | White: rl.LoadTexture(filepath.Join("assets", "white-knight.png")),
		Knight | Black: rl.LoadTexture(filepath.Join("assets", "black-knight.png")),
		Bishop | White: rl.LoadTexture(filepath.Join("assets", "white-bishop.png")),
		Bishop | Black: rl.LoadTexture(filepath.Join("assets", "black-bishop.png")),
		Rook | White:   rl.LoadTexture(filepath.Join("assets", "white-rook.png")),
		Rook | Black:   rl.LoadTexture(filepath.Join("assets", "black-rook.png")),
		Queen | White:  rl.LoadTexture(filepath.Join("assets", "white-queen.png")),
		Queen | Black:  rl.LoadTexture(filepath.Join("assets", "black-queen.png")),
	}
}

const startFen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQKq - 0 1"

func LoadPosFromFen(fen string) [64]int {
	pieceFromsym := make(map[string]int)
	pieceFromsym["k"] = King
	pieceFromsym["p"] = Pawn
	pieceFromsym["n"] = Knight
	pieceFromsym["b"] = Bishop
	pieceFromsym["r"] = Rook
	pieceFromsym["q"] = Queen
	var squares [64]int
	fenB := strings.Split(fen, " ")[0]
	file := 0
	rank := 7
	for _, c := range fenB {
		if c == '/' {
			file = 0
			rank--
		} else {
			if c >= '0' && c <= '9' {
				file += -1
			} else {
				color := White
				if c >= 'a' && c <= 'z' {
					color = Black
				}
				pieceType := pieceFromsym[strings.ToLower(string(c))]
				squares[rank*8+file] = pieceType | color
				file++
			}
		}
	}
	return squares
}

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

func CreatePieces(squ [64]int, fileFromInt map[int]rl.Texture2D) {
	//offset := 0
	filePlacement := 7
	for file := 0; file < 8; file++ {
		rankPlacement := 7
		for rank := 0; rank < 8; rank++ {
			piece := squ[rank*8+file]
			if piece != None {
				rl.DrawTexturePro(fileFromInt[piece], rl.NewRectangle(0, 0, float32(fileFromInt[piece].Width), float32(fileFromInt[piece].Height)), rl.NewRectangle(float32(filePlacement*squareSize), float32(rankPlacement*squareSize), float32(squareSize), float32(squareSize)), rl.NewVector2(0, 0), 0, rl.White)
			}
			rankPlacement--
		}
		filePlacement--
	}

}

func DrawText() {
	fontSize := 15
	for file := 0; file < 8; file++ {
		for rank := 0; rank < 8; rank++ {
			if file == 0 {
				text := " Rank: " + strconv.Itoa(8-rank)
				rl.DrawText(text, int32(file*squareSize)+(squareSize*8), int32(rank*squareSize+(squareSize-fontSize)), int32(fontSize), rl.Black)
			}
			if rank == 0 {
				text := " File: " + string('a'+file)
				rl.DrawText(text, int32(file*squareSize), int32(rank*squareSize)+(squareSize*8), 15, rl.Black)
			}
		}
	}
}

func main() {
	rl.InitWindow(screenSize, screenSize, "Chess")

	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	squares := LoadPosFromFen(startFen)
	textMap := textFromInt()
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)
		CreateBoard()
		DrawText()
		CreatePieces(squares, textMap)

		rl.EndDrawing()
	}
}
