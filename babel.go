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
	big "github.com/ncw/gmp"
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

type Page struct {
	Hex, Wall, Shelf, Volume, Page *big.Int
}

func ToBabelianAddressCompressed(input []byte) []byte {
	nonbabel := bytes.NewBuffer(input)
	var pages []Page
	var blocks [][]byte
	var pagenum = int(float64(nonbabel.Len())/float64(3239)) + 1
	for i := 0; i < pagenum; i++ {
		blocks = append(blocks, nonbabel.Next(3239))
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
		multed := big.NewInt(0).Mul(loc_int, loc_mult)
		key := big.NewInt(0).Sub(hex, multed)
		ret.Write(key.Bytes())
	}
	return ret.Bytes()
}
