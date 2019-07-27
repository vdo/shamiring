package main

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"strings"
)

type RandParams struct {
	P     *big.Int   //finite field prime
	Rands []*big.Int // random coeficients
}

func generate_rands(k int) RandParams {
	// Generate prime number for finite field
	p, _ := rand.Prime(rand.Reader, 256)

	r := make([]*big.Int, k)
	max := new(big.Int)
	max.Exp(big.NewInt(2), big.NewInt(130), nil).Sub(max, big.NewInt(1))
	for i := 1; i < k; i++ {
		n, _ := rand.Int(rand.Reader, max)
		r[i] = n
	}
	// r[0] is the secret!
	return RandParams{p, r}
}

func eval(poly []*big.Int, x int, p *big.Int) *big.Int {
	// Poly in reverse, see: https://rosettacode.org/wiki/Horner%27s_rule_for_polynomial_evaluation
	acc := big.NewInt(0)
	for i := len(poly) - 1; i >= 0; i-- {
		acc.Mul(acc, big.NewInt(int64(x)))
		acc.Add(acc, poly[i])
		acc.Mod(acc, p)
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

func modInverse(number *big.Int, p *big.Int) *big.Int {
	copy := big.NewInt(int64(0)).Set(number)
	copy = copy.Mod(copy, p)
	pcopy := big.NewInt(int64(0)).Set(p)
	x := big.NewInt(int64(0))
	y := big.NewInt(int64(0))
	copy.GCD(x, y, pcopy, copy)
	result := big.NewInt(int64(0)).Set(p)
	result = result.Add(result, y)
	result = result.Mod(result, p)
	return result
}

func interpolate(xs []*big.Int, ys []*big.Int, p *big.Int) (value *big.Int) {
	value = big.NewInt(int64(0))
	for i := 0; i < len(ys); i++ {
		num := big.NewInt(int64(1))
		den := big.NewInt(int64(1))
		for j := 0; j < len(xs); j++ {
			if i != j {
				top := big.NewInt(int64(0))
				top = top.Mul(xs[j], big.NewInt(int64(-1)))
				bottom := new(big.Int).Sub(xs[j], xs[i])
				num = num.Mul(num, top)
				num = num.Mod(num, p)
				den = den.Mul(den, bottom)
				den = den.Mod(den, p)
			}
		}
		acc := big.NewInt(int64(0)).Set(ys[i])
		acc = acc.Mul(acc, num)
		acc = acc.Mul(acc, modInverse(den, p))
		value = value.Add(value, acc)
		value = value.Mod(value, p)
	}
	return
}

func check(e error) {
	if e != nil {
		fmt.Println("ERROR: ", e)
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
	ptest := big.NewInt(int64(91994388364979))
	//ptest, _ := rand.Prime(rand.Reader, 1024)
	var ytest []*big.Int
	var xtest []*big.Int

	xtest = append(xtest, big.NewInt(int64(1)))
	xtest = append(xtest, big.NewInt(int64(3)))
	xtest = append(xtest, big.NewInt(int64(5)))
	ytest = append(ytest, big.NewInt(int64(9285275391624)))
	ytest = append(ytest, big.NewInt(int64(53079135586964)))
	ytest = append(ytest, big.NewInt(int64(37709686632597)))
	res := interpolate(xtest, ytest, ptest)
	fmt.Println("First example:", res)

	// rands := generate_rands(5)
	// fmt.Println(rands)

	// points := get_points(3, rands)
	// fmt.Println(points)

}
