package godude

import (
	"crypto/elliptic"
	"crypto/subtle"
	"math/big"
	"math/rand"
	"testing"
	"time"
)

func TestBytesCT(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	target := make([]byte, 100)
	_, _ = r.Read(target)
	tv := Measure(func(t []byte) int {
		return subtle.ConstantTimeCompare(target, t)
	}, func() []byte {
		t := make([]byte, 100)
		_, _ = r.Read(t)
		return t
	}, func() []byte {
		t := make([]byte, 100)
		copy(t, target)
		return t
	}, 1000000)

	t.Logf("%f", tv)
	if tv > 10 {
		t.Errorf("appears not constant time")
	}
}

func TestBytesBad(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	target := make([]byte, 100)
	_, _ = r.Read(target)
	tv := Measure(func(t []byte) int {
		for i, b := range target {
			if b != t[i] {
				return 1
			}
		}
		return 0
	}, func() []byte {
		t := make([]byte, 100)
		_, _ = r.Read(t)
		return t
	}, func() []byte {
		t := make([]byte, 100)
		copy(t, target)
		return t
	}, 1000000)

	t.Logf("%f", tv)
	if tv <= 10 {
		t.Errorf("appears to be constant time")
	}
}

func TestP256(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	p256 := elliptic.P256()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	tv := Measure(func(t []byte) *big.Int {
		x, _ := p256.ScalarBaseMult(t)
		return x
	}, func() []byte {
		b := make([]byte, 32)
		_, _ = r.Read(b)
		return b
	}, func() []byte {
		t := make([]byte, 1)
		t[0] = 1
		return t
	}, 100000)
	t.Logf("%f", tv)
	if tv > 10 {
		t.Errorf("appears not constant time")
	}
}

func TestP521(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	p521 := elliptic.P521()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	tv := Measure(func(t []byte) *big.Int {
		x, _ := p521.ScalarBaseMult(t)
		return x
	}, func() []byte {
		b := make([]byte, 32)
		_, _ = r.Read(b)
		return b
	}, func() []byte {
		t := make([]byte, 1)
		t[0] = 1
		return t
	}, 10000)
	t.Logf("%f", tv)
	if tv <= 10 {
		t.Errorf("appears to be constant time")
	}
}
