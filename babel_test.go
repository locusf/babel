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
	"crypto/rand"
	"testing"
)

var lenghts = []int{300, 1000, 3200, 3238, 3239, 3250, 3300, 3400, 3500, 10000,
				100000, 1000000, 32000000}

func TestBabelian(t *testing.T) {
	for _, blen := range lenghts {
		byt := make([]byte, blen)
		t.Log("Testing length", blen)
		_, err := rand.Read(byt)
		if err != nil {
			t.Log("Error: ", err)
		}
		zipbytes := ToBabelianAddressCompressed(byt)
		ret := FromBabelianAddressCompressed(zipbytes)
		if !bytes.Equal(byt, ret) {
			t.Log("Failed for lenght:", blen, "array lengths are", blen, len(ret))
			t.Fail()
		} else {
			t.Log("Passed for length: ", blen)
		}
	}
}
func BenchmarkToBabelian(b *testing.B) {
	blen := 1000000
	byt := make([]byte, 1000000)
	rand.Read(byt)
	zipbytes := ToBabelianAddressCompressed(byt)
	ret := FromBabelianAddressCompressed(zipbytes)
	if !bytes.Equal(ret, byt) {
		b.Log("Failed for lenght:", blen, "array lengths are", blen, len(ret))
		b.Fail()
	}
}
