package main

import (
	"fmt"
	"math/rand/v2"
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
	Queen
	Rook
	Bishop
	Knight
	Pawn

	White = 8
	Black = 16
)

const (
	white = iota
	black
)

const (
	a8 = iota
	b8
	c8
	d8
	e8
	f8
	g8
	h8
	a7
	b7
	c7
	d7
	e7
	f7
	g7
	h7
	a6
	b6
	c6
	d6
	e6
	f6
	g6
	h6
	a5
	b5
	c5
	d5
	e5
	f5
	g5
	h5
	a4
	b4
	c4
	d4
	e4
	f4
	g4
	h4
	a3
	b3
	c3
	d3
	e3
	f3
	g3
	h3
	a2
	b2
	c2
	d2
	e2
	f2
	g2
	h2
	a1
	b1
	c1
	d1
	e1
	f1
	g1
	h1
)

var squareCords = [...]string{
	"a8", "b8", "c8", "d8", "e8", "f8", "g8", "h8",
	"a7", "b7", "c7", "d7", "e7", "f7", "g7", "h7",
	"a6", "b6", "c6", "d6", "e6", "f6", "g6", "h6",
	"a5", "b5", "c5", "d5", "e5", "f5", "g5", "h5",
	"a4", "b4", "c4", "d4", "e4", "f4", "g4", "h4",
	"a3", "b3", "c3", "d3", "e3", "f3", "g3", "h3",
	"a2", "b2", "c2", "d2", "e2", "f2", "g2", "h2",
	"a1", "b1", "c1", "d1", "e1", "f1", "g1", "h1",
}

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

func genMagicNumber() uint64 {
	return uint64(rand.Int()) & uint64(rand.Int()) & uint64(rand.Int())
}

func getBit(bitboard uint64, square int) bool {
	return bitboard&(1<<uint(square)) != 0
}

func setBit(bitboard *uint64, square int) {
	*bitboard |= 1 << uint(square)
}

func clearBit(bitboard *uint64, square int) {
	if getBit(*bitboard, square) {
		*bitboard ^= 1 << uint(square)
	}
}

func countBits(bitboard uint64) int {
	count := 0
	for bitboard != 0 {
		count++
		bitboard &= bitboard - 1
	}
	return count
}

func getLSB(bitboard uint64) int {
	if bitboard == 0 {
		return -1
	}
	return countBits((bitboard & -bitboard) - 1)
}

var fileMasks = [8]uint64{
	0x0101010101010101,
	0x0202020202020202,
	0x0404040404040404,
	0x0808080808080808,
	0x1010101010101010,
	0x2020202020202020,
	0x4040404040404040,
	0x8080808080808080,
}

var bishopRelBits = [64]uint64{
	6, 5, 5, 5, 5, 5, 5, 6,
	5, 5, 5, 5, 5, 5, 5, 5,
	5, 5, 7, 7, 7, 7, 5, 5,
	5, 5, 7, 9, 9, 7, 5, 5,
	5, 5, 7, 9, 9, 7, 5, 5,
	5, 5, 7, 7, 7, 7, 5, 5,
	5, 5, 5, 5, 5, 5, 5, 5,
	6, 5, 5, 5, 5, 5, 5, 6,
}

var rookRelBits = [64]uint64{
	12, 11, 11, 11, 11, 11, 11, 12,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	12, 11, 11, 11, 11, 11, 11, 12,
}

var bishopMagics = [64]uint64{
	0x40040844404084,
	0x2004208a004208,
	0x10190041080202,
	0x108060845042010,
	0x581104180800210,
	0x2112080446200010,
	0x1080820820060210,
	0x3c0808410220200,
	0x4050404440404,
	0x21001420088,
	0x24d0080801082102,
	0x1020a0a020400,
	0x40308200402,
	0x4011002100800,
	0x401484104104005,
	0x801010402020200,
	0x400210c3880100,
	0x404022024108200,
	0x810018200204102,
	0x4002801a02003,
	0x85040820080400,
	0x810102c808880400,
	0xe900410884800,
	0x8002020480840102,
	0x220200865090201,
	0x2010100a02021202,
	0x152048408022401,
	0x20080002081110,
	0x4001001021004000,
	0x800040400a011002,
	0xe4004081011002,
	0x1c004001012080,
	0x8004200962a00220,
	0x8422100208500202,
	0x2000402200300c08,
	0x8646020080080080,
	0x80020a0200100808,
	0x2010004880111000,
	0x623000a080011400,
	0x42008c0340209202,
	0x209188240001000,
	0x400408a884001800,
	0x110400a6080400,
	0x1840060a44020800,
	0x90080104000041,
	0x201011000808101,
	0x1a2208080504f080,
	0x8012020600211212,
	0x500861011240000,
	0x180806108200800,
	0x4000020e01040044,
	0x300000261044000a,
	0x802241102020002,
	0x20906061210001,
	0x5a84841004010310,
	0x4010801011c04,
	0xa010109502200,
	0x4a02012000,
	0x500201010098b028,
	0x8040002811040900,
	0x28000010020204,
	0x6000020202d0240,
	0x8918844842082200,
	0x4010011029020020,
}
var rookMagics = [64]uint64{
	0x200c06040000408,
	0x10620920000211,
	0x140041400ac0,
	0x13440a0001a04229,
	0x208000001000000,
	0x2060010010020401,
	0x1008920c00602040,
	0x4002009104002008,
	0x2200000000000703,
	0x800004080200814,
	0x420200003020c401,
	0x23400004c1021380,
	0x2204400800200100,
	0x40181000000840,
	0x85000000010090,
	0x10880009200108,
	0x2102090810018008,
	0x48800240100000,
	0x420000000c108000,
	0x4e04000400180040,
	0x42015140024000,
	0x2110010000004a00,
	0x84040c90090000,
	0x408200910008,
	0x840200409,
	0x80010041000020,
	0x608404080408,
	0x2100048080000580,
	0x210420100101,
	0x100240800001000,
	0x220000008000044,
	0x20000000400,
	0x2004004400000280,
	0x5100002050000059,
	0x100080411004440,
	0x4000001843124210,
	0xa0420800248000,
	0x800000036020000,
	0x210080040a012,
	0x4028004410002640,
	0x80000120240000,
	0x1000084001406808,
	0x20000206c0060000,
	0x200809000c0000,
	0x408180000250a01,
	0x404041000050800,
	0x1002200020,
	0x2010020400004,
	0x2200288011c00844,
	0x4002100000,
	0x40004002020001a0,
	0x100070008001046,
	0x548400000040020,
	0x4002400000b0000,
	0x1000810130c10250,
	0x50044801040201,
	0x80400814008,
	0x6008108080122090,
	0x10490ca0000200,
	0x1102800801,
	0x812040009000020,
	0x6002020000010100,
	0x3604008030800,
	0xa0000098208408,
}

var bishopMasks = [64]uint64{}
var rookMasks = [64]uint64{}
var bishopAttacks = [64][512]uint64{}
var rookAttacks = [64][4096]uint64{}

// Attacks
var pawnAttacks = [2][64]uint64{}

var knightAttacks = [64]uint64{}

var kingAttacks = [64]uint64{}

func maskPawnAttacks(square int, side int) uint64 {
	var bitboard uint64 = 0
	var attacks uint64 = 0
	setBit(&bitboard, square)
	if side == white {
		attacks |= bitboard >> 7 &^ fileMasks[0]
		attacks |= bitboard >> 9 &^ fileMasks[7]
	} else {
		attacks |= bitboard << 7 &^ fileMasks[7]
		attacks |= bitboard << 9 &^ fileMasks[0]
	}
	return attacks
}

func maskKnightAttacks(square int) uint64 {
	var bitboard uint64 = 0
	var attacks uint64 = 0
	setBit(&bitboard, square)
	attacks |= bitboard >> 17 &^ fileMasks[7]
	attacks |= bitboard >> 15 &^ fileMasks[0]
	attacks |= bitboard >> 10 & ((^fileMasks[7]) & (^fileMasks[6]))
	attacks |= bitboard >> 6 & ((^fileMasks[0]) & (^fileMasks[1]))
	attacks |= bitboard << 17 &^ fileMasks[0]
	attacks |= bitboard << 15 &^ fileMasks[7]
	attacks |= bitboard << 10 & ((^fileMasks[0]) & (^fileMasks[1]))
	attacks |= bitboard << 6 & ((^fileMasks[7]) & (^fileMasks[6]))
	return attacks
}

func maskKingAttacks(square int) uint64 {
	var bitboard uint64 = 0
	var attacks uint64 = 0
	setBit(&bitboard, square)
	attacks |= bitboard >> 7 &^ fileMasks[0]
	attacks |= bitboard >> 8
	attacks |= bitboard >> 9 &^ fileMasks[7]
	attacks |= bitboard >> 1 &^ fileMasks[7]
	attacks |= bitboard << 7 &^ fileMasks[7]
	attacks |= bitboard << 8
	attacks |= bitboard << 9 &^ fileMasks[0]
	attacks |= bitboard << 1 &^ fileMasks[0]
	return attacks
}

func maskBishopAttacks(square int) uint64 {
	var bitboard uint64 = 0
	var attacks uint64 = 0
	setBit(&bitboard, square)
	tr := square / 8
	tf := square % 8
	for r, f := tr+1, tf+1; r <= 7 && f <= 7; r, f = r+1, f+1 {
		attacks |= uint64((1 << uint(r*8+f)))
	}
	for r, f := tr-1, tf-1; r >= 0 && f >= 0; r, f = r-1, f-1 {
		attacks |= uint64((1 << uint(r*8+f)))
	}
	for r, f := tr+1, tf-1; r <= 7 && f >= 0; r, f = r+1, f-1 {
		attacks |= uint64((1 << uint(r*8+f)))
	}
	for r, f := tr-1, tf+1; r >= 0 && f <= 7; r, f = r-1, f+1 {
		attacks |= uint64((1 << uint(r*8+f)))
	}
	return attacks
}

func bishopAttacksNow(square int, block uint64) uint64 {
	var bitboard uint64 = 0
	var attacks uint64 = 0
	setBit(&bitboard, square)
	tr := square / 8
	tf := square % 8
	for r, f := tr+1, tf+1; r <= 7 && f <= 7; r, f = r+1, f+1 {
		attacks |= uint64((1 << uint(r*8+f)))
		if (1<<uint(r*8+f))&block != 0 {
			break
		}
	}
	for r, f := tr-1, tf-1; r >= 0 && f >= 0; r, f = r-1, f-1 {
		attacks |= uint64((1 << uint(r*8+f)))
		if (1<<uint(r*8+f))&block != 0 {
			break
		}
	}
	for r, f := tr+1, tf-1; r <= 7 && f >= 0; r, f = r+1, f-1 {
		attacks |= uint64((1 << uint(r*8+f)))
		if (1<<uint(r*8+f))&block != 0 {
			break
		}
	}
	for r, f := tr-1, tf+1; r >= 0 && f <= 7; r, f = r-1, f+1 {
		attacks |= uint64((1 << uint(r*8+f)))
		if (1<<uint(r*8+f))&block != 0 {
			break
		}
	}
	return attacks
}

func maskRookAttacks(square int) uint64 {
	var bitboard uint64 = 0
	var attacks uint64 = 0
	setBit(&bitboard, square)
	tr := square / 8
	tf := square % 8
	for r := tr + 1; r <= 6; r++ {
		attacks |= uint64((1 << uint(r*8+tf)))
	}
	for r := tr - 1; r >= 1; r-- {
		attacks |= uint64((1 << uint(r*8+tf)))
	}
	for f := tf + 1; f <= 6; f++ {
		attacks |= uint64((1 << uint(tr*8+f)))
	}
	for f := tf - 1; f >= 1; f-- {
		attacks |= uint64((1 << uint(tr*8+f)))
	}
	return attacks
}

func rookAttacksNow(square int, block uint64) uint64 {
	var bitboard uint64 = 0
	var attacks uint64 = 0
	setBit(&bitboard, square)
	tr := square / 8
	tf := square % 8
	for r := tr + 1; r <= 7; r++ {
		attacks |= uint64((1 << uint(r*8+tf)))
		if (1<<uint(r*8+tf))&block != 0 {
			break
		}
	}
	for r := tr - 1; r >= 0; r-- {
		attacks |= uint64((1 << uint(r*8+tf)))
		if (1<<uint(r*8+tf))&block != 0 {
			break
		}
	}
	for f := tf + 1; f <= 7; f++ {
		attacks |= uint64((1 << uint(tr*8+f)))
		if (1<<uint(tr*8+f))&block != 0 {
			break
		}
	}
	for f := tf - 1; f >= 0; f-- {
		attacks |= uint64((1 << uint(tr*8+f)))
		if (1<<uint(tr*8+f))&block != 0 {
			break
		}
	}
	return attacks
}

func initLeaperAttacks() {
	for square := 0; square < 64; square++ {
		pawnAttacks[white][square] = maskPawnAttacks(square, white)
		pawnAttacks[black][square] = maskPawnAttacks(square, black)
		knightAttacks[square] = maskKnightAttacks(square)
		kingAttacks[square] = maskKingAttacks(square)
	}
}

func setOccupancy(index int, bitIntMask int, attackMask *uint64) uint64 {
	var occupancy uint64 = 0
	for count := 0; count < bitIntMask; count++ {
		square := getLSB(*attackMask)
		clearBit(attackMask, square)
		if index&(1<<uint(count)) != 0 {
			occupancy |= 1 << uint(square)
		}
	}
	return occupancy
}

func CreatePiecesFromBit(bitboard uint64, fileFromInt map[int]rl.Texture2D) {
	for file := 0; file < 8; file++ {
		for rank := 0; rank < 8; rank++ {
			square := file*8 + rank
			if getBit(bitboard, square) {
				rl.DrawTexturePro(fileFromInt[King|White], rl.NewRectangle(0, 0, float32(fileFromInt[King|White].Width), float32(fileFromInt[King|White].Height)), rl.NewRectangle(float32(rank*squareSize), float32(file*squareSize), float32(squareSize), float32(squareSize)), rl.NewVector2(0, 0), 0, rl.White)
			}
		}
	}
}

func printBitboard(bitboard uint64) {
	for rank := 0; rank < 8; rank++ {
		for file := 0; file < 8; file++ {
			square := rank*8 + file
			if getBit(bitboard, square) {
				print("1")
			} else {
				print("0")
			}
		}
		println()
	}
}

// Find magic number
func findMagic(square int, relBits int, isBishop bool) uint64 {
	var occupancy [4096]uint64
	var attacks [4096]uint64
	var usedAttacks [4096]bool
	var mask uint64
	if isBishop {
		mask = maskBishopAttacks(square)
	} else {
		mask = maskRookAttacks(square)
	}
	var occupancyIndicies int = 1 << relBits
	for index := 0; index < occupancyIndicies; index++ {
		occupancy[index] = setOccupancy(index, relBits, &mask)
		if isBishop {
			attacks[index] = bishopAttacksNow(square, occupancy[index])
		} else {
			attacks[index] = rookAttacksNow(square, occupancy[index])
		}
	}
	for i := 0; i < 1000000; i++ {
		var magic uint64 = genMagicNumber()
		var fail bool = false
		for index := 0; index < occupancyIndicies; index++ {
			var magicIndex uint64 = (occupancy[index] * magic) >> uint(64-relBits)
			if usedAttacks[magicIndex] && attacks[magicIndex] != attacks[index] {
				fail = true
				break
			}
			usedAttacks[magicIndex] = true
		}
		if !fail {
			return magic
		}
	}
	return 0
}

func initSliders(isBishop bool) {
	for square := 0; square < 64; square++ {
		if isBishop {
			bishopMasks[square] = maskBishopAttacks(square)
		} else {
			rookMasks[square] = maskRookAttacks(square)
		}
		var attackMask uint64
		if isBishop {
			attackMask = bishopMasks[square]
		} else {
			attackMask = rookMasks[square]
		}
		var occupancyIndicies int = 1 << bishopRelBits[square]
		for index := 0; index < occupancyIndicies; index++ {
			var occupancy uint64 = setOccupancy(index, countBits(attackMask), &attackMask)
			if isBishop {
				bishopAttacks[square][index] = bishopAttacksNow(square, occupancy)
			} else {
				rookAttacks[square][index] = rookAttacksNow(square, occupancy)
			}
		}

	}

}

func getBishopAttacks(square int, occupancy uint64) uint64 {
	occupancy &= bishopMasks[square]
	occupancy *= bishopMagics[square]
	occupancy >>= 64 - bishopRelBits[square]
	return bishopAttacks[square][int(occupancy)]
}

func getRookAttacks(square int, occupancy uint64) uint64 {
	occupancy &= rookMasks[square]
	occupancy *= rookMagics[square]
	occupancy >>= uint(64 - rookRelBits[square])
	return rookAttacks[square][int(occupancy)]
}

func initMagicNumbers() {
	for square := 0; square < 64; square++ {
		bishopMagics[square] = findMagic(square, int(bishopRelBits[square]), true)
	}
	fmt.Println()
	for square := 0; square < 64; square++ {
		rookMagics[square] = findMagic(square, int(rookRelBits[square]), false)
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

func initAll() {
	initLeaperAttacks()
	initSliders(true)
	initSliders(false)
	//initMagicNumbers()
}

func main() {
	rl.InitWindow(screenSize, screenSize, "Chess")

	defer rl.CloseWindow()

	rl.SetTargetFPS(60)
	initAll()

	//squares := LoadPosFromFen(startFen)
	textMap := textFromInt()

	var occupancy uint64 = 0
	setBit(&occupancy, c5)
	setBit(&occupancy, d4)
	occupancy = getBishopAttacks(c5, occupancy)

	for rank := 0; rank < 8; rank++ {
		for file := 0; file < 8; file++ {
			square := rank*8 + file
			if getBit(occupancy, square) {
				print("1")
			} else {
				print("0")
			}
		}
		println()
	}

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)
		CreateBoard()
		DrawText()
		CreatePiecesFromBit(occupancy, textMap)

		//CreatePieces(squares, textMap)
		bitboardText := strconv.FormatUint(occupancy, 10)
		rl.DrawText(bitboardText, 500, 550, 15, rl.Black)
		rl.EndDrawing()
	}
}
