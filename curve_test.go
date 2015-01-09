package edwards

import (
	"crypto/rand"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func testCurves() []Curve {
	return []Curve{E222(), E382(), E521(), Ed448Goldilocks(), Curve1174(), Curve41417()}
}

func TestScalarMultZero(t *testing.T) {
	for _, crv := range testCurves() {
		_, x0, y0, err := crv.GenerateKey(rand.Reader)
		if err != nil {
			panic(err)
		}
		assert.True(t, crv.IsOnCurve(x0, y0), "Generated P must be on curve")

		rx, ry := crv.ScalarMult(x0, y0, []byte{0, 0})

		assert.Equal(t, zero, rx, "(0*P).X != 0")
		assert.Equal(t, one, ry, "(0*P).Y != 1")
	}
}

func TestScalarMultOne(t *testing.T) {
	for _, crv := range testCurves() {
		_, x0, y0, err := crv.GenerateKey(rand.Reader)
		if err != nil {
			panic(err)
		}
		assert.True(t, crv.IsOnCurve(x0, y0), "Generated P must be on curve")

		rx, ry := crv.ScalarMult(x0, y0, []byte{0, 1})

		assert.Equal(t, rx, x0, "(1*P).X != P.X")
		assert.Equal(t, ry, y0, "(1*P).Y != P.Y")
	}
}

func TestScalarMultTwo(t *testing.T) {
	for _, crv := range testCurves() {
		_, x0, y0, err := crv.GenerateKey(rand.Reader)
		if err != nil {
			panic(err)
		}
		assert.True(t, crv.IsOnCurve(x0, y0), "Generated P must be on curve")

		rx, ry := crv.ScalarMult(x0, y0, []byte{2})
		assert.True(t, crv.IsOnCurve(rx, ry), "2*P must be on curve")

		sx, sy := crv.Add(x0, y0, x0, y0)
		assert.True(t, crv.IsOnCurve(sx, sy), "P+P must be on curve")

		assert.Equal(t, rx, sx, "(2*P).X != (P+P).X")
		assert.Equal(t, ry, sy, "(2*P).Y != (P+P).Y")
	}
}

func TestScalarMultCommutative(t *testing.T) {
	for _, crv := range testCurves() {
		priv0, x0, y0, err := crv.GenerateKey(rand.Reader)
		if err != nil {
			panic(err)
		}
		assert.True(t, crv.IsOnCurve(x0, y0), "Generated P must be on curve")

		priv1, x1, y1, err := crv.GenerateKey(rand.Reader)
		if err != nil {
			panic(err)
		}
		assert.True(t, crv.IsOnCurve(x1, y1), "Generated P must be on curve")

		rx, ry := crv.ScalarMult(x0, y0, priv1)
		assert.True(t, crv.IsOnCurve(rx, ry), "priv1*P0 must be on curve")

		sx, sy := crv.ScalarMult(x1, y1, priv0)
		assert.True(t, crv.IsOnCurve(sx, sy), "priv0*P1 must be on curve")

		assert.Equal(t, rx, sx, "(priv0*P1).X != (priv1*P0).X")
		assert.Equal(t, ry, sy, "(priv0*P1).Y != (priv1*P0).Y")
	}
}

func TestAdditionCommutative(t *testing.T) {
	for _, crv := range testCurves() {
		_, x0, y0, err := crv.GenerateKey(rand.Reader)
		if err != nil {
			panic(err)
		}
		assert.True(t, crv.IsOnCurve(x0, y0), "Generated P must be on curve")

		_, x1, y1, err := crv.GenerateKey(rand.Reader)
		if err != nil {
			panic(err)
		}
		assert.True(t, crv.IsOnCurve(x1, y1), "Generated P must be on curve")

		rx, ry := crv.Add(x0, y0, x1, y1)
		assert.True(t, crv.IsOnCurve(rx, ry), "P0+P1 must be on curve")

		sx, sy := crv.Add(x1, y1, x0, y0)
		assert.True(t, crv.IsOnCurve(sx, sy), "P1+P0 must be on curve")

		assert.Equal(t, rx, sx, "(P0+P1).X != (P1+P0).X")
		assert.Equal(t, ry, sy, "(P0+P1).Y != (P1+P0).Y")
	}
}

func TestBasePointOnCurve(t *testing.T) {
	for _, crv := range testCurves() {
		assert.True(t, crv.IsOnCurve(crv.Gx, crv.Gy), "Base point must be on curve")
	}
}

func TestBasePointOrderCorrect(t *testing.T) {
	for _, crv := range testCurves() {
		rx, ry := crv.ScalarBaseMult(crv.N.Bytes())

		assert.Equal(t, zero, rx, "(N*G).X != 0")
		assert.Equal(t, one, ry, "(N*G).Y != 1")
	}
}

func TestAdditiveInverse(t *testing.T) {
	for _, crv := range testCurves() {
		_, x0, y0, err := crv.GenerateKey(rand.Reader)
		if err != nil {
			panic(err)
		}
		assert.True(t, crv.IsOnCurve(x0, y0), "Generated P must be on curve")

		x1 := new(big.Int).Neg(x0)

		rx, ry := crv.Add(x0, y0, x1, y0)

		assert.Equal(t, zero, rx, "(P+(-P)).X != 0")
		assert.Equal(t, one, ry, "(P+(-P)).Y != 1")
	}
}

func TestNeutralAddition(t *testing.T) {
	for _, crv := range testCurves() {
		x, y := crv.Add(zero, one, zero, one)
		assert.Equal(t, zero, x, "(E+E).X != 0")
		assert.Equal(t, one, y, "(E+E).Y != 1")
	}
}

func TestNeutralScalarMult(t *testing.T) {
	for _, crv := range testCurves() {
		x, y := crv.ScalarMult(zero, one, []byte{0})
		assert.Equal(t, zero, x, "(0*E).X != 0")
		assert.Equal(t, one, y, "(0*E).Y != 1")

		x, y = crv.ScalarMult(zero, one, []byte{1})
		assert.Equal(t, zero, x, "(1*E).X != 0")
		assert.Equal(t, one, y, "(1*E).Y != 1")

		x, y = crv.ScalarMult(zero, one, []byte{2})
		assert.Equal(t, zero, x, "(2*E).X != 0")
		assert.Equal(t, one, y, "(2*E).Y != 1")
	}
}

func TestKnownPoints(t *testing.T) {
	for _, crv := range testCurves() {
		x, y := crv.ScalarMult(zero, big.NewInt(-1), []byte{2})
		assert.Equal(t, zero, x, "(2*(0, -1)).X != 0")
		assert.Equal(t, one, y, "(2*(0, -1)).Y != 1")

		x, y = crv.ScalarMult(one, zero, []byte{4})
		assert.Equal(t, zero, x, "(2*(1, 0)).X != 0")
		assert.Equal(t, one, y, "(2*(1, 0)).Y != 1")

		x, y = crv.ScalarMult(big.NewInt(-1), zero, []byte{4})
		assert.Equal(t, zero, x, "(4*(-1, 0)).X != 0")
		assert.Equal(t, one, y, "(4*(-1, 0)).Y != 1")
	}
}
