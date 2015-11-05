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
	big "github.com/locusf/gmp"
	"testing"
)

func TestBabelian(t *testing.T) {
	byt := make([]byte, 10)
	_, err := rand.Read(byt)
	if err != nil {
		t.Log("Error: ", err)
	}
	z := big.NewInt(0)
	z = z.SetBytes(byt)
	z = z.Abs(z)
	zipbytes := ToBabelianAddressCompressed(z.Bytes())
	ret := FromBabelianAddressCompressed(zipbytes)
	if !bytes.Equal(z.Bytes(), ret) {
		t.Fail()
	}
}
func BenchmarkToBabelian(b *testing.B) {
	byt := make([]byte, 1000000)
	rand.Read(byt)
	z := big.NewInt(0).SetBytes(byt)
	zipbytes := ToBabelianAddressCompressed(z.Bytes())
	ret := FromBabelianAddressCompressed(zipbytes)
	if !bytes.Equal(ret, z.Bytes()) {
		b.Fail()
	}
}
