package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
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
	for i := 0; i < k; i++ {
		n, _ := rand.Int(rand.Reader, max)
		r[i] = *n
	}

	return RandParams{*p, r}
}

func eval(poly []big.Int, x int64, p big.Int) *big.Int {

	acc := big.NewInt(0)
	// Poly in reverse, see: https://rosettacode.org/wiki/Horner%27s_rule_for_polynomial_evaluation
	for i := len(poly) - 1; i >= 0; i-- {
		acc.Mul(acc, big.NewInt(x))
		acc.Add(acc, &poly[i])
		acc.Mod(acc, &p)
	}
	return acc
}

func get_points(n int, params RandParams) []big.Int {
	ps := make([]big.Int, n)
	for i := 1; i <= n; i++ {
		// evaluate poly for each
		ps[i-1] = *eval(params.Rands, int64(i), params.P)
	}
	return ps
}

func main() {
	//reader := bufio.NewReader(os.Stdin)

	fmt.Println("What is the secret?")
	fmt.Println("What is the number of parts?")
	fmt.Println("What is the threshold?")
	rands := generate_rands(5)
	fmt.Println(rands)

	points := get_points(3, rands)
	fmt.Println(points)
}
