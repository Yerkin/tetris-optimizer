package main

import (
	"bufio"
	"fmt"
	"os"
)

// XY Our new struction
type XY struct {
	x int
	y int
}

var tetro [][]XY
var squareMaxSize int
var ans [][]rune
var ban, ban2 [][]int
var curByte []byte
var line []int
var cur int
var dxy []XY

func check(inX, inY int, pos []XY) bool {
	for _, i := range pos {
		if inX+i.x >= fin || inY+i.y >= fin || inY+i.y < 0 || inX+i.x < 0 || ban[inX+i.x][inY+i.y] == 2 {
			return false
		}
	}
	return true
}
func accept(inX, inY, color int, pos []XY) {
	for _, i := range pos {
		ans[inX+i.x][inY+i.y] = rune('A' + color)
		ban[inX+i.x][inY+i.y] = 2
	}

}

var lol int

func reject(inX, inY int, pos []XY) {
	for _, i := range pos {
		ans[inX+i.x][inY+i.y] = '.'
		ban[inX+i.x][inY+i.y] = 1
	}
}

var maxX, maxY int
var mn, sz int
var res [][]rune
var (
	st = map[string]bool{}
)

var fin int
var SawIt []bool

// Next Will help us optimize solution
var Next []int

// BruteForce tetrominos
func calc(start int, freeCells int, step []int) {

	if start == fin*fin || freeCells-1 > fin*fin-sz*4 {
		return
	}
	sX := start / fin
	sY := start % fin

	for z := 0; z < sz; z++ {
		if step[z] == 0 && !SawIt[z] && check(sX, sY, tetro[z]) {

			step[z] = 1

			accept(sX, sY, z, tetro[z])
			if Next[z] != -1 {
				SawIt[Next[z]] = false
			}

			cur--
			if cur == 0 {
				for i := 0; i < fin; i++ {
					for j := 0; j < fin; j++ {
						fmt.Print(string(ans[i][j]))
					}
					fmt.Println()
				}
				os.Exit(0)
				return
			}
			calc(start+1, freeCells, step)
			cur++

			reject(sX, sY, tetro[z])
			if Next[z] != -1 {
				SawIt[Next[z]] = true
			}
			step[z] = 0
		}
	}
	ok := 0
	if ban[sX][sY] != 2 {
		ok = 1
	}
	calc(start+1, freeCells+ok, step)

}

func main() {

	// now := time.Now()
	// defer func() {
	// 	fmt.Println(time.Since(now))
	// }()
	dxy = append(dxy, XY{x: 0, y: 1})
	dxy = append(dxy, XY{x: 0, y: -1})
	dxy = append(dxy, XY{x: 1, y: 0})
	dxy = append(dxy, XY{x: -1, y: 0})
	data := []string{}
	arg := os.Args[1:]
	if len(arg) != 1 {
		fmt.Println("chose file for test")
		os.Exit(0)
	}
	file, _ := os.Open(arg[0])
	reader := bufio.NewReader(file)

	cnt := 0
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		data = append(data, line)
		cnt++
	}
	// Checking our tetrominos and transform
	err := false
	for i := 0; i < cnt; i += 5 {
		minX := cnt
		minY := cnt
		tetro = append(tetro, []XY{})
		pointsCount := 0
		for j := i; j < i+4; j++ {
			for k, c := range data[j] {
				if c == '#' {
					if j < minX {
						minY = k
						minX = j
					} else if j == minX {
						if k < minY {
							minY = k
						}
					}
					tetro[cur] = append(tetro[cur], XY{x: j, y: k})
				}
				if c == '.' {
					pointsCount++
				}
			}
		}
		if i+4 < cnt && len(data[i+4]) != 1 || len(tetro[cur]) != 4 || pointsCount != 12 {
			err = true
			break
		}
		for i := range tetro[cur] {
			tetro[cur][i].x -= minX
			tetro[cur][i].y -= minY
		}
		con := 0
		for j := 0; j < 4; j++ {
			for k := 0; k < 4; k++ {
				for _, c := range dxy {
					if tetro[cur][j].x+c.x == tetro[cur][k].x && tetro[cur][j].y+c.y == tetro[cur][k].y {
						con++
					}
				}
			}
		}
		if con < 6 {
			err = true
			break
		}

		for j := 0; j < 4; j++ {
			for k := j + 1; k < 4; k++ {
				if tetro[cur][j].x > tetro[cur][k].x || tetro[cur][j].x == tetro[cur][k].x && tetro[cur][j].y > tetro[cur][k].y {
					tetro[cur][j], tetro[cur][k] = tetro[cur][k], tetro[cur][j]
				}
			}
		}
		cur++
	}

	// Optimizing solution with making array of used tetrominos
	Next = make([]int, cur)
	SawIt = make([]bool, cur)
	for i := 0; i < cur; i++ {
		Next[i] = -1
		SawIt[i] = false
	}
	for i := 0; i < cur; i++ {
		for j := i + 1; j < cur; j++ {
			ok := true
			for k := 0; k < 4; k++ {
				if tetro[i][k].x != tetro[j][k].x || tetro[i][k].y != tetro[j][k].y {
					ok = false
					break
				}

			}
			if ok {
				Next[i] = j
				SawIt[j] = true
				break
			}
		}
	}
	//  #
	// ##
	//  #
	// If error ocuured
	if err {
		fmt.Println("Error")
		return
	}

	for i := 0; i < cur; i++ {
		line = append(line, 0)
	}
	for squareMaxSize*squareMaxSize < cur*6 {
		squareMaxSize++
	}
	squarMinSize := 0
	for squarMinSize*squarMinSize < cur*4 {
		squarMinSize++
	}

	if squareMaxSize < 4 {
		squareMaxSize = 4
	}

	// Bruteforce from small square to one possible square
	for i := 0; i < squareMaxSize; i++ {
		ban = append(ban, []int{})
		ans = append(ans, []rune{})
		for j := 0; j < squareMaxSize; j++ {
			ban[i] = append(ban[i], 1)
			ans[i] = append(ans[i], '.')
			curByte = append(curByte, '0')
		}
	}
	l := cur
	sz = cur
	for k := squarMinSize; k <= squareMaxSize; k++ {
		cur = l
		fin = k
		for i := 0; i < len(ban); i++ {
			for j := 0; j < len(ban[i]); j++ {
				ban[i][j] = 1
			}
		}

		for i := 0; i < len(ans); i++ {
			for j := 0; j < len(ans[i]); j++ {
				ans[i][j] = '.'
			}
		}
		for i := 0; i < len(curByte); i++ {
			curByte[i] = '0'
		}
		sz = len(line)
		ban[0][0] = 0
		calc(0, 0, line)
	}
}
