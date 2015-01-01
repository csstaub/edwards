package edwards

import (
	"io"
	"math/big"
	"sync"
)

var initonce sync.Once
var e222 Curve
var e382 Curve
var e448 Curve
var e521 Curve
var c1174 Curve
var c41417 Curve

var zero = big.NewInt(0)
var one = big.NewInt(1)

// Curve represents a Twisted Edwards curve.
type Curve struct {
	P       *big.Int
	A, D    *big.Int
	N       *big.Int
	Gx, Gy  *big.Int
	BitSize int
}

// A point on the curve represented using projective coordinates.
type point struct {
	X, Y, Z *big.Int
}

func initAll() {
	initE222()
	initE382()
	initE448()
	initE521()
	initC1174()
	initC41417()
}

func initE222() {
	e222.P = new(big.Int).Exp(big.NewInt(2), big.NewInt(222), nil)
	e222.P.Sub(e222.P, big.NewInt(117))

	e222.A = big.NewInt(1)
	e222.D = big.NewInt(160102)

	e222.Gx = new(big.Int)
	e222.Gx.SetString("2705691079882681090389589001251962954446177367541711474502428610129", 10)
	e222.Gy = big.NewInt(28)

	e222.N = new(big.Int)
	e222.N.SetString("1684996666696914987166688442938726735569737456760058294185521417407", 10)

	e222.BitSize = 222
}

func initE382() {
	e382.P = new(big.Int).Exp(big.NewInt(2), big.NewInt(382), nil)
	e382.P.Sub(e382.P, big.NewInt(105))

	e382.A = big.NewInt(1)
	e382.D = big.NewInt(-67254)

	e382.Gx = new(big.Int)
	e382.Gx.SetString("3914921414754292646847594472454013487047137431784830634731377862923477302047857640522480241298429278603678181725699", 10)
	e382.Gy = big.NewInt(17)

	e382.N = new(big.Int)
	e382.N.SetString("2462625387274654950767440006258975862817483704404090416745738034557663054564649171262659326683244604346084081047321", 10)

	e382.BitSize = 382
}

func initE448() {
	two := big.NewInt(2)

	p224 := new(big.Int).Exp(two, big.NewInt(224), nil)

	e448.P = new(big.Int).Exp(two, big.NewInt(448), nil)
	e448.P.Sub(e448.P, p224)
	e448.P.Sub(e448.P, one)

	e448.A = big.NewInt(1)
	e448.D = big.NewInt(-39081)

	e448.Gx = new(big.Int)
	e448.Gx.SetString("117812161263436946737282484343310064665180535357016373416879082147939404277809514858788439644911793978499419995990477371552926308078495", 10)
	e448.Gy = big.NewInt(19)

	e448.N = new(big.Int)
	e448.N.SetString("181709681073901722637330951972001133588410340171829515070372549795146003961539585716195755291692375963310293709091662304773755859649779", 10)

	e448.BitSize = 448
}

func initE521() {
	e521.P = new(big.Int).Exp(big.NewInt(2), big.NewInt(521), nil)
	e521.P.Sub(e521.P, one)

	e521.A = big.NewInt(1)
	e521.D = big.NewInt(-376014)

	e521.Gx = new(big.Int)
	e521.Gx.SetString("1571054894184995387535939749894317568645297350402905821437625181152304994381188529632591196067604100772673927915114267193389905003276673749012051148356041324", 10)
	e521.Gy = big.NewInt(12)

	e521.N = new(big.Int)
	e521.N.SetString("1716199415032652428745475199770348304317358825035826352348615864796385795849413675475876651663657849636693659065234142604319282948702542317993421293670108523", 10)

	e521.BitSize = 521
}

func initC1174() {
	c1174.P = new(big.Int).Exp(big.NewInt(2), big.NewInt(251), nil)
	c1174.P.Sub(c1174.P, big.NewInt(9))

	c1174.A = big.NewInt(1)
	c1174.D = big.NewInt(-1174)

	c1174.Gx = new(big.Int)
	c1174.Gx.SetString("1582619097725911541954547006453739763381091388846394833492296309729998839514", 10)
	c1174.Gy = new(big.Int)
	c1174.Gy.SetString("3037538013604154504764115728651437646519513534305223422754827055689195992590", 10)

	c1174.N = new(big.Int)
	c1174.N.SetString("904625697166532776746648320380374280092339035279495474023489261773642975601", 10)

	c1174.BitSize = 251
}

func initC41417() {
	c41417.P = new(big.Int).Exp(big.NewInt(2), big.NewInt(414), nil)
	c41417.P.Sub(c41417.P, big.NewInt(17))

	c41417.A = big.NewInt(1)
	c41417.D = big.NewInt(3617)

	c41417.Gx = new(big.Int)
	c41417.Gx.SetString("17319886477121189177719202498822615443556957307604340815256226171904769976866975908866528699294134494857887698432266169206165", 10)
	c41417.Gy = big.NewInt(34)

	c41417.N = new(big.Int)
	c41417.N.SetString("5288447750321988791615322464262168318627237463714249754277190328831105466135348245791335989419337099796002495788978276839289", 10)

	c41417.BitSize = 414
}

func (crv Curve) IsOnCurve(x, y *big.Int) bool {
	// a*x^2 + y^2 = 1 + d*x^2*y^2
	xx := new(big.Int).Mul(x, x)
	yy := new(big.Int).Mul(y, y)

	l := new(big.Int).Mul(crv.A, xx)
	l.Add(l, yy)
	l.Mod(l, crv.P)

	r := new(big.Int).Mul(crv.D, xx)
	r.Mul(r, yy)
	r.Add(r, big.NewInt(1))
	r.Mod(r, crv.P)

	return l.Cmp(r) == 0
}

func (crv Curve) GenerateKey(rand io.Reader) (priv []byte, x, y *big.Int, err error) {
	priv = make([]byte, (crv.BitSize+7)>>3)

	for {
		_, err = io.ReadFull(rand, priv)
		if err != nil {
			return
		}

		priv[0] &= 0xFF >> uint((8-crv.BitSize)%8)
		x, y = crv.ScalarBaseMult(priv)

		// We don't want to return the neutral element
		if x.Cmp(zero) != 0 && y.Cmp(one) != 0 {
			return
		}
	}
}

func (crv Curve) Add(Lx, Ly, Rx, Ry *big.Int) (x, y *big.Int) {
	return crv.affineFromProjective(crv.projectiveAdd(projectiveFromAffine(Lx, Ly), projectiveFromAffine(Rx, Ry)))
}

func (crv Curve) ScalarBaseMult(k []byte) (x, y *big.Int) {
	return crv.ScalarMult(crv.Gx, crv.Gy, k)
}

func (crv Curve) ScalarMult(Bx, By *big.Int, k []byte) (x, y *big.Int) {
	r0 := projectiveFromAffine(zero, one)
	r1 := projectiveFromAffine(Bx, By)

	for _, b := range k {
		for n := 0; n < 8; n++ {
			if b&0x80 == 0x80 {
				r0 = crv.projectiveAdd(r0, r1)
				r1 = crv.projectiveAdd(r1, r1)
			} else {
				r1 = crv.projectiveAdd(r0, r1)
				r0 = crv.projectiveAdd(r0, r0)
			}
			b <<= 1
		}
	}

	return crv.affineFromProjective(r0)
}

func (crv Curve) projectiveAdd(l, r point) point {
	a := new(big.Int).Mul(l.Z, r.Z)
	a.Mod(a, crv.P)
	c := new(big.Int).Mul(l.X, r.X)
	c.Mod(c, crv.P)
	d := new(big.Int).Mul(l.Y, r.Y)
	d.Mod(d, crv.P)

	b := new(big.Int).Mul(a, a)
	b.Mod(b, crv.P)

	e := new(big.Int)
	e.Mul(c, d)
	e.Mul(e, crv.D)
	e.Mod(e, crv.P)

	f := new(big.Int).Sub(b, e)
	g := new(big.Int).Add(b, e)

	x1y1 := new(big.Int).Add(l.X, l.Y)
	x1y1.Mod(x1y1, crv.P)
	x2y2 := new(big.Int).Add(r.X, r.Y)
	x2y2.Mod(x2y2, crv.P)

	x3 := new(big.Int)
	x3.Mul(x1y1, x2y2)
	x3.Sub(x3, c)
	x3.Sub(x3, d)
	x3.Mul(x3, f)
	x3.Mul(x3, a)
	x3.Mod(x3, crv.P)

	ac := new(big.Int).Mul(crv.A, c)
	ac.Mod(ac, crv.P)

	y3 := new(big.Int)
	y3.Sub(d, ac)
	y3.Mul(y3, g)
	y3.Mul(y3, a)
	y3.Mod(y3, crv.P)

	z3 := new(big.Int).Mul(f, g)
	z3.Mod(z3, crv.P)

	return point{
		X: x3,
		Y: y3,
		Z: z3,
	}
}

func projectiveFromAffine(x, y *big.Int) point {
	return point{
		X: x,
		Y: y,
		Z: big.NewInt(1),
	}
}

func (crv Curve) affineFromProjective(p point) (x, y *big.Int) {
	x = new(big.Int)
	y = new(big.Int)

	if p.Z.Sign() == 0 {
		return
	}

	zinv := new(big.Int).ModInverse(p.Z, crv.P)

	x.Mul(p.X, zinv)
	x.Mod(x, crv.P)
	y.Mul(p.Y, zinv)
	y.Mod(y, crv.P)

	return
}

func E222() Curve {
	initonce.Do(initAll)
	return e222
}

func E382() Curve {
	initonce.Do(initAll)
	return e382
}

func Ed448Goldilocks() Curve {
	initonce.Do(initAll)
	return e448
}

func E521() Curve {
	initonce.Do(initAll)
	return e521
}

func Curve1174() Curve {
	initonce.Do(initAll)
	return c1174
}

func Curve41417() Curve {
	initonce.Do(initAll)
	return c41417
}
