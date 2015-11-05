/*
   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/
package babel

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"fmt"
	big "github.com/locusf/gmp"
	"math/rand"
	"strings"
	"time"
)

var length_of_page = big.NewInt(3239)
var seededRandom = rand.New(rand.NewSource(time.Now().UnixNano()))
var loc_mult = big.NewInt(0).Exp(big.NewInt(30), length_of_page, big.NewInt(0))
var tobabelreplacer = strings.NewReplacer(
	"0", "a",
	"1", "b",
	"2", "c",
	"3", "d",
	"4", "e",
	"5", "f",
	"6", "g",
	"7", "h",
	"8", "i",
	"9", "j",
	"a", "k",
	"b", "l",
	"c", "m",
	"d", "n",
	"e", "o",
	"f", "p",
	"g", "r",
	"h", "s",
	"i", "t",
	"j", "u",
	"k", "v",
	"l", "w",
	"m", "x",
	"n", "y",
	"o", "z",
	"p", ",",
	"q", " ",
	"r", ".",
)

func substr(s []byte, from, length int) []byte {
	//create array like string view
	wb := make([][]byte, length)
	wb = bytes.Split(s, []byte(""))

	//miss nil pointer error
	to := from + length

	if to > len(wb) {
		to = len(wb)
	}

	if from > len(wb) {
		from = len(wb)
	}

	return bytes.Join(wb[from:to], []byte(""))
}

type Page struct {
	Hex, Wall, Shelf, Volume, Page *big.Int
}

func ToBabelianAddressCompressed(input []byte) []byte {
	inp := big.NewInt(0).SetBytes(input)
	nonbabel := inp.Bytes()
	var pages []Page
	var blocks [][]byte
	for i := 0; i < ((len(nonbabel) / 3239) + 1); i++ {
		blocks = append(blocks, substr(nonbabel, i, i+3239))
	}
	for _, subabel := range blocks {
		wall := big.NewInt(0).Rand(seededRandom, big.NewInt(4))
		shelf := big.NewInt(0).Rand(seededRandom, big.NewInt(5))
		volume := big.NewInt(0).Rand(seededRandom, big.NewInt(410))
		page := big.NewInt(0).Rand(seededRandom, big.NewInt(410))
		loc_int := big.NewInt(0).Add(page, big.NewInt(0).Add(volume, big.NewInt(0).Add(shelf, wall)))
		x := big.NewInt(0).SetBytes(subabel)
		multed := big.NewInt(0).Mul(loc_int, loc_mult)
		added := big.NewInt(0).Add(x, multed)
		pages = append(pages, Page{added, wall, shelf, volume, page})
	}
	var network bytes.Buffer
	var output bytes.Buffer
	enc := gob.NewEncoder(&network)
	err := enc.Encode(pages)
	if err != nil {
		fmt.Println(err)
	}
	zipper := gzip.NewWriter(&output)
	_, err = zipper.Write(network.Bytes())
	zipper.Close()
	if err != nil {
		fmt.Println("Error: ", err)
	}
	return output.Bytes()
}
func FromBabelianAddressCompressed(input []byte) []byte {
	var ret bytes.Buffer
	inputBuf := bytes.NewBuffer(input)
	zipper, err := gzip.NewReader(inputBuf)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	defer zipper.Close()
	dec := gob.NewDecoder(zipper)
	var pages []Page
	err = dec.Decode(&pages)
	if err != nil {
		fmt.Println("Decode error: ", err)
	}
	for _, pagestr := range pages {
		hex := pagestr.Hex
		wall := pagestr.Wall
		shelf := pagestr.Shelf
		volume := pagestr.Volume
		page := pagestr.Page
		loc_int := big.NewInt(0).Add(page, big.NewInt(0).Add(volume, big.NewInt(0).Add(shelf, wall)))
		key := big.NewInt(0).Sub(hex, big.NewInt(0).Mul(loc_int, loc_mult))
		ret.Write(key.Bytes())
	}
	return ret.Bytes()
}
