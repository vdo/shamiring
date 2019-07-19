package main

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"strings"
)

// s = secret ==== a0
// n = parts ===>  we need to obtain n points in the poly
// k = threshold ===> we need to obtain (k-1) random numbers

type RandParams struct {
	P     big.Int   //finite field prime
	Rands []big.Int // random coeficients
}

// f(x) = s + n1*x + n2*x^2 + n3*x^3 ...
func generate_rands(k int) RandParams {
	// Generate prime number for finite field
	p, _ := rand.Prime(rand.Reader, 256)

	r := make([]big.Int, k)
	max := new(big.Int)
	max.Exp(big.NewInt(2), big.NewInt(130), nil).Sub(max, big.NewInt(1))
	for i := 1; i < k; i++ {
		n, _ := rand.Int(rand.Reader, max)
		r[i] = *n
	}
	// r[0] is the secret!
	return RandParams{*p, r}
}

func eval(poly []big.Int, x int, p big.Int) *big.Int {
	// Poly in reverse, see: https://rosettacode.org/wiki/Horner%27s_rule_for_polynomial_evaluation
	acc := big.NewInt(int64(0))
	for i := len(poly) - 1; i >= 0; i-- {
		acc.Mul(acc, big.NewInt(int64(x)))
		acc.Add(acc, &poly[i])
		acc.Mod(acc, &p)
	}
	return acc
}

func get_points(n int, params RandParams) []big.Int {
	ps := make([]big.Int, n)
	for i := 1; i <= n; i++ {
		// evaluate poly for each
		ps[i-1] = *eval(params.Rands, i, params.P)
	}
	return ps
}

// Lagrange Interpolation, kudos to arnaucube
func lagrange(xs, ys []*big.Int, p *big.Int) *big.Int {
	resultN := big.NewInt(int64(0))
	resultD := big.NewInt(int64(0))
	for i := 0; i < len(xs); i++ {
		lagNum := big.NewInt(int64(1))
		lagDen := big.NewInt(int64(1))
		for j := 0; j < len(xs); j++ {
			if xs[i] != xs[j] {
				currLagNum := xs[j]
				currLagDen := new(big.Int).Sub(xs[j], xs[i])
				lagNum = new(big.Int).Mul(lagNum, currLagNum)
				lagDen = new(big.Int).Mul(lagDen, currLagDen)
			}
		}
		numerator := new(big.Int).Mul(xs[i], lagNum)
		quo := new(big.Int).Quo(numerator, lagDen)
		if quo.Int64() != 0 {
			resultN = resultN.Add(resultN, quo)
		} else {
			resultNMULlagDen := new(big.Int).Mul(resultN, lagDen)
			resultN = new(big.Int).Add(resultNMULlagDen, numerator)
			resultD = resultD.Add(resultD, lagDen)
		}
	}
	var modinvMul *big.Int
	if resultD.Int64() != 0 {
		modinv := new(big.Int).ModInverse(resultD, p)
		modinvMul = new(big.Int).Mul(resultN, modinv)
	} else {
		modinvMul = resultN
	}
	r := new(big.Int).Mod(modinvMul, p)
	return r
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("What is the secret?")
	secret, err := reader.ReadString('\n')
	check(err)
	if len(secret) <= 1 {
		fmt.Println("Please enter a secret!")
		return
	}
	secret = strings.TrimSuffix(secret, "\n")

	// Convert to big.Int
	z := new(big.Int)
	z.SetBytes([]byte(secret))
	fmt.Println(z)
	// Convert back to string!
	bytes := z.Bytes()
	fmt.Println(string(bytes))

	fmt.Println("How many shares?")
	shares, err := reader.ReadString('\n')
	check(err)
	if len(shares) <= 1 {
		fmt.Println("Please enter a number!")
		return
	}

	fmt.Println("Threshold?")
	threshold, err := reader.ReadString('\n')
	check(err)
	if len(threshold) <= 1 {
		fmt.Println("Please enter a number!")
		return
	}
	rands := generate_rands(5)
	fmt.Println(rands)

	points := get_points(3, rands)
	fmt.Println(points)
}
