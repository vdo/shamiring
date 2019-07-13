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

// Python code to evaluate a poly at point x, then modulo 'prime'
// def _eval_at(poly, x, prime):
//     '''evaluates polynomial (coefficient tuple) at x, used to generate a
//     shamir pool in make_random_shares below.
//     '''
//     accum = 0
//     for coeff in reversed(poly):
//         accum *= x
//         accum += coeff
//         accum %= prime
//     return accum

func eval(poly []big.Int, x int, p big.Int) big.Int {

	acc := big.NewInt(0)
	for _, x := range poly {
		acc = acc * x

	}

}

func get_points(n int, r []big.Int, s big.Int, k int) []big.Int {
	points := make([]big.Int, n)
	for i := 0; i < n; i++ {
		// evaluate poly

	}
}

func main() {
	//reader := bufio.NewReader(os.Stdin)

	fmt.Println("What is the secret?")
	fmt.Println("What is the number of parts?")
	fmt.Println("What is the threshold?")
	rands := generate_rands(5)
	fmt.Println(rands)

	points
}
