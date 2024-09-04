package pow

import (
	"bytes"
	"encoding/hex"
	"math/rand"
	"testing"

	"github.com/zilong-dai/karlsen-miner/consensus/model/externalapi"
	"github.com/zilong-dai/karlsen-miner/consensus/utils/hashes"
)

func BenchmarkMatrix_HeavyHash(b *testing.B) {
	input := []byte("BenchmarkMatrix_HeavyHash")
	writer := hashes.NewPoWHashWriter()
	writer.InfallibleWrite(input)
	hash := writer.Finalize()
	matrix := generateMatrix(hash)
	for i := 0; i < b.N; i++ {
		hash = matrix.HeavyHash(hash)
	}
}

func BenchmarkMatrix_Generate(b *testing.B) {
	r := rand.New(rand.NewSource(0))
	h := [32]byte{}
	r.Read(h[:])
	hash := externalapi.NewDomainHashFromByteArray(&h)
	for i := 0; i < b.N; i++ {
		generateMatrix(hash)
	}
}

func BenchmarkMatrix_Rank(b *testing.B) {
	r := rand.New(rand.NewSource(0))
	h := [32]byte{}
	r.Read(h[:])
	hash := externalapi.NewDomainHashFromByteArray(&h)
	matrix := generateMatrix(hash)

	for i := 0; i < b.N; i++ {
		matrix.computeRank()
	}
}

func TestMatrix_Rank(t *testing.T) {
	var mat matrix
	if mat.computeRank() != 0 {
		t.Fatalf("The zero matrix should have rank 0, instead got: %d", mat.computeRank())
	}

	r := rand.New(rand.NewSource(0))
	h := [32]byte{}
	r.Read(h[:])
	hash := externalapi.NewDomainHashFromByteArray(&h)
	mat = *generateMatrix(hash)

	if mat.computeRank() != 64 {
		t.Fatalf("generateMatrix() should always return full rank matrix, instead got: %d", mat.computeRank())
	}

	for i := range mat {
		mat[0][i] = mat[1][i] * 2
	}
	if mat.computeRank() != 63 {
		t.Fatalf("Made a linear depenency between 2 rows, rank should be 63, instead got: %d", mat.computeRank())
	}

	for i := range mat {
		mat[33][i] = mat[32][i] * 3
	}
	if mat.computeRank() != 62 {
		t.Fatalf("Made anoter linear depenency between 2 rows, rank should be 62, instead got: %d", mat.computeRank())
	}
}

func TestGenerateMatrix(t *testing.T) {
	hashBytes := [32]byte{42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42}
	hash := externalapi.NewDomainHashFromByteArray(&hashBytes)
	expectedMatrix := matrix{
		{4, 5, 4, 5, 4, 5, 4, 5, 4, 5, 4, 5, 4, 5, 4, 5, 15, 3, 15, 3, 15, 3, 15, 3, 15, 3, 15, 3, 15, 3, 15, 3, 2, 10, 2, 10, 2, 10, 2, 10, 2, 10, 2, 10, 2, 10, 2, 10, 14, 1, 2, 2, 14, 10, 4, 12, 4, 12, 10, 10, 10, 10, 10, 10},
		{9, 11, 1, 11, 1, 11, 9, 11, 9, 11, 9, 3, 12, 13, 11, 5, 15, 15, 5, 0, 6, 8, 1, 8, 6, 11, 15, 5, 3, 6, 7, 3, 2, 15, 14, 3, 7, 11, 14, 7, 3, 6, 14, 12, 3, 9, 5, 1, 1, 0, 8, 4, 10, 15, 9, 10, 6, 13, 1, 1, 7, 4, 4, 6},
		{2, 6, 0, 8, 11, 15, 4, 0, 5, 2, 7, 13, 15, 3, 11, 12, 6, 2, 1, 8, 13, 4, 11, 4, 10, 14, 13, 2, 6, 15, 10, 6, 6, 5, 6, 9, 3, 3, 3, 1, 9, 12, 12, 15, 6, 0, 1, 5, 7, 13, 14, 1, 10, 10, 5, 14, 4, 0, 12, 13, 2, 15, 8, 4},
		{8, 6, 5, 1, 0, 6, 4, 8, 13, 0, 8, 12, 7, 2, 4, 3, 10, 5, 9, 3, 12, 13, 2, 4, 13, 14, 7, 7, 9, 12, 10, 8, 11, 6, 14, 3, 12, 8, 8, 0, 2, 10, 0, 9, 1, 9, 7, 8, 5, 2, 9, 13, 15, 6, 13, 10, 1, 9, 1, 10, 6, 2, 10, 9},
		{4, 2, 6, 14, 4, 2, 5, 7, 15, 6, 0, 4, 11, 9, 12, 0, 3, 2, 0, 4, 10, 5, 12, 3, 3, 4, 10, 1, 0, 13, 3, 12, 15, 0, 7, 10, 2, 2, 15, 0, 2, 15, 8, 2, 15, 12, 10, 6, 6, 2, 13, 3, 8, 14, 3, 13, 10, 5, 4, 5, 1, 6, 5, 10},
		{0, 3, 13, 12, 11, 4, 11, 13, 1, 12, 4, 11, 15, 14, 13, 4, 7, 1, 3, 0, 10, 3, 8, 8, 1, 2, 5, 14, 4, 5, 14, 1, 1, 3, 3, 1, 5, 15, 7, 5, 11, 8, 8, 12, 10, 5, 7, 9, 2, 10, 13, 11, 4, 2, 12, 15, 10, 6, 6, 0, 6, 6, 3, 12},
		{9, 12, 3, 3, 5, 8, 12, 13, 7, 4, 5, 11, 4, 0, 7, 2, 2, 15, 12, 14, 12, 5, 4, 2, 8, 8, 8, 13, 6, 1, 1, 5, 0, 15, 12, 13, 8, 5, 0, 4, 13, 1, 6, 1, 12, 14, 1, 0, 13, 12, 10, 10, 1, 4, 13, 13, 8, 4, 15, 13, 6, 6, 14, 10},
		{14, 15, 8, 0, 7, 2, 5, 10, 5, 3, 12, 0, 11, 3, 4, 2, 8, 11, 6, 14, 14, 3, 3, 12, 3, 7, 6, 2, 6, 12, 15, 1, 1, 13, 0, 6, 9, 9, 7, 7, 13, 4, 4, 2, 15, 5, 2, 15, 13, 13, 10, 6, 9, 15, 2, 9, 6, 10, 6, 14, 14, 3, 5, 11},
		{6, 4, 7, 8, 11, 0, 13, 11, 0, 7, 0, 0, 13, 6, 3, 11, 15, 14, 10, 2, 7, 8, 13, 14, 8, 15, 10, 8, 14, 6, 10, 14, 3, 11, 5, 11, 13, 5, 3, 12, 3, 0, 2, 0, 6, 14, 4, 12, 4, 4, 8, 15, 7, 8, 12, 11, 3, 9, 5, 13, 10, 14, 13, 4},
		{10, 0, 0, 15, 1, 4, 13, 3, 15, 10, 2, 5, 11, 2, 9, 14, 7, 3, 2, 8, 6, 15, 0, 12, 1, 4, 1, 9, 3, 0, 15, 8, 9, 13, 0, 7, 9, 10, 6, 14, 3, 7, 9, 7, 4, 0, 11, 8, 4, 6, 5, 8, 8, 0, 5, 14, 7, 12, 12, 2, 5, 6, 5, 6},
		{12, 0, 0, 14, 8, 3, 0, 3, 13, 10, 5, 13, 5, 7, 2, 4, 13, 11, 3, 1, 11, 2, 14, 5, 10, 5, 5, 9, 12, 15, 12, 8, 1, 0, 11, 13, 8, 1, 1, 11, 10, 0, 11, 15, 13, 9, 12, 14, 5, 4, 5, 14, 2, 7, 2, 1, 4, 12, 11, 11, 9, 12, 11, 15},
		{3, 15, 9, 8, 13, 12, 15, 7, 8, 7, 14, 6, 10, 3, 0, 5, 2, 2, 6, 6, 3, 2, 5, 12, 11, 2, 10, 11, 13, 3, 9, 7, 7, 6, 8, 15, 14, 14, 11, 11, 9, 7, 1, 3, 8, 5, 11, 11, 1, 2, 15, 8, 13, 8, 11, 4, 1, 5, 3, 12, 5, 3, 7, 7},
		{13, 13, 2, 14, 4, 3, 15, 2, 0, 15, 1, 5, 4, 1, 5, 1, 4, 14, 5, 1, 11, 13, 15, 1, 3, 3, 5, 13, 14, 1, 0, 4, 6, 1, 15, 7, 7, 0, 15, 8, 15, 3, 14, 7, 7, 8, 12, 10, 2, 14, 9, 2, 11, 11, 7, 10, 4, 3, 12, 13, 4, 13, 0, 14},
		{12, 14, 15, 15, 2, 0, 0, 13, 4, 6, 4, 2, 14, 11, 5, 6, 14, 8, 14, 7, 13, 15, 6, 15, 7, 9, 1, 0, 11, 9, 9, 0, 2, 12, 8, 8, 14, 11, 7, 5, 3, 0, 11, 12, 9, 2, 8, 9, 0, 0, 9, 8, 9, 8, 2, 14, 12, 2, 0, 14, 13, 8, 4, 10},
		{7, 10, 1, 15, 12, 14, 7, 4, 7, 13, 4, 8, 13, 12, 1, 7, 10, 6, 5, 14, 14, 3, 14, 4, 11, 14, 6, 12, 15, 12, 15, 12, 4, 5, 9, 8, 7, 7, 3, 0, 5, 7, 3, 8, 4, 4, 7, 5, 6, 12, 13, 0, 12, 10, 2, 5, 14, 9, 6, 4, 13, 13, 14, 5},
		{14, 5, 8, 3, 4, 15, 13, 14, 14, 10, 7, 14, 15, 2, 11, 14, 13, 13, 12, 10, 6, 9, 5, 5, 6, 13, 15, 13, 7, 0, 15, 11, 4, 12, 15, 7, 7, 4, 3, 11, 8, 14, 5, 10, 2, 4, 4, 12, 3, 6, 1, 9, 15, 1, 1, 13, 7, 5, 0, 14, 15, 7, 8, 6},
		{1, 2, 10, 5, 2, 13, 1, 11, 15, 10, 4, 9, 9, 12, 14, 13, 3, 5, 0, 3, 7, 11, 10, 3, 12, 5, 10, 2, 13, 7, 1, 7, 13, 8, 2, 8, 3, 14, 10, 3, 5, 12, 0, 9, 3, 9, 11, 2, 10, 9, 0, 6, 4, 0, 1, 14, 11, 0, 8, 6, 1, 15, 3, 10},
		{13, 9, 0, 5, 8, 7, 12, 15, 10, 10, 5, 1, 1, 7, 6, 1, 14, 5, 15, 2, 3, 5, 3, 5, 7, 3, 7, 7, 1, 4, 3, 14, 5, 0, 12, 0, 12, 10, 10, 6, 12, 6, 3, 5, 5, 11, 10, 1, 11, 3, 13, 3, 9, 11, 1, 7, 14, 14, 0, 8, 15, 5, 2, 7},
		{8, 5, 11, 6, 15, 0, 1, 13, 1, 6, 7, 15, 4, 3, 14, 12, 9, 3, 11, 6, 4, 12, 1, 11, 6, 12, 5, 11, 1, 12, 2, 3, 1, 2, 11, 12, 0, 5, 11, 5, 3, 13, 11, 3, 11, 14, 10, 8, 3, 9, 4, 8, 13, 11, 9, 11, 2, 4, 12, 3, 0, 14, 7, 11},
		{10, 11, 4, 10, 7, 8, 3, 14, 15, 8, 15, 6, 9, 8, 5, 6, 12, 1, 15, 6, 5, 5, 14, 13, 2, 12, 14, 6, 5, 5, 14, 9, 1, 10, 11, 14, 8, 6, 14, 11, 1, 15, 6, 11, 11, 8, 1, 2, 8, 5, 4, 15, 6, 8, 0, 8, 0, 11, 0, 1, 0, 7, 8, 15},
		{0, 15, 5, 0, 11, 4, 4, 2, 0, 4, 8, 12, 2, 2, 0, 8, 1, 2, 6, 5, 6, 12, 3, 1, 12, 1, 6, 10, 2, 5, 0, 2, 0, 11, 8, 6, 13, 4, 14, 4, 15, 5, 8, 11, 9, 6, 2, 6, 9, 1, 4, 2, 14, 10, 4, 4, 1, 1, 11, 8, 6, 11, 11, 9},
		{7, 3, 6, 5, 9, 1, 11, 0, 15, 13, 13, 13, 4, 14, 14, 12, 3, 7, 9, 3, 1, 6, 5, 9, 7, 6, 2, 11, 10, 4, 11, 14, 10, 13, 11, 8, 11, 8, 1, 15, 5, 0, 10, 5, 6, 0, 5, 15, 11, 6, 6, 4, 10, 11, 8, 12, 0, 10, 11, 11, 11, 1, 13, 6},
		{7, 15, 0, 0, 11, 5, 7, 13, 3, 7, 3, 2, 5, 12, 6, 11, 14, 4, 9, 8, 9, 9, 13, 0, 15, 2, 13, 2, 15, 6, 15, 1, 1, 7, 4, 0, 10, 1, 8, 14, 0, 10, 12, 4, 5, 13, 9, 0, 7, 12, 13, 11, 11, 8, 8, 15, 2, 15, 4, 4, 9, 3, 10, 7},
		{0, 9, 3, 5, 14, 6, 7, 14, 7, 2, 13, 7, 3, 15, 9, 15, 2, 8, 0, 4, 6, 0, 15, 6, 2, 1, 14, 8, 5, 8, 2, 4, 2, 11, 9, 2, 15, 13, 11, 12, 8, 15, 3, 13, 2, 2, 10, 13, 1, 8, 7, 15, 13, 6, 7, 7, 4, 3, 14, 7, 0, 9, 15, 11},
		{8, 13, 7, 7, 8, 8, 7, 8, 1, 4, 10, 1, 12, 4, 14, 11, 7, 12, 15, 0, 10, 15, 9, 2, 14, 2, 14, 2, 4, 5, 13, 3, 2, 10, 0, 15, 7, 6, 8, 11, 7, 6, 10, 10, 4, 7, 10, 6, 6, 14, 10, 4, 14, 6, 12, 2, 8, 1, 9, 13, 3, 4, 3, 14},
		{10, 10, 6, 3, 8, 5, 10, 7, 11, 10, 9, 4, 8, 14, 9, 10, 0, 9, 8, 14, 11, 15, 8, 13, 13, 7, 13, 13, 13, 9, 12, 11, 6, 3, 9, 6, 0, 0, 6, 6, 11, 6, 4, 8, 1, 5, 1, 7, 9, 6, 13, 4, 3, 8, 8, 11, 9, 10, 6, 11, 12, 13, 14, 14},
		{14, 10, 0, 15, 14, 4, 3, 0, 12, 4, 0, 14, 11, 9, 0, 6, 4, 6, 0, 9, 8, 14, 4, 4, 6, 8, 2, 8, 10, 3, 8, 0, 1, 1, 15, 4, 2, 4, 13, 9, 9, 4, 0, 5, 5, 1, 2, 5, 11, 6, 2, 1, 7, 8, 10, 10, 1, 5, 8, 6, 7, 0, 4, 14},
		{0, 15, 10, 11, 13, 12, 7, 7, 4, 0, 9, 5, 2, 8, 0, 10, 6, 6, 7, 5, 6, 7, 9, 0, 1, 4, 8, 14, 10, 3, 5, 5, 11, 5, 1, 10, 6, 10, 0, 14, 1, 15, 11, 12, 8, 2, 7, 8, 4, 0, 3, 11, 9, 15, 3, 5, 15, 15, 14, 15, 3, 4, 5, 14},
		{5, 12, 12, 8, 0, 0, 14, 1, 4, 15, 3, 2, 2, 6, 1, 10, 7, 10, 14, 5, 14, 0, 8, 5, 9, 0, 12, 8, 9, 10, 3, 12, 3, 2, 0, 0, 12, 12, 7, 13, 2, 6, 4, 7, 10, 10, 14, 1, 11, 6, 10, 3, 12, 2, 1, 10, 7, 13, 10, 12, 14, 11, 14, 8},
		{9, 5, 3, 12, 4, 3, 10, 14, 7, 5, 11, 12, 2, 13, 9, 8, 5, 2, 6, 2, 4, 9, 10, 10, 4, 3, 4, 0, 11, 1, 10, 9, 4, 10, 4, 5, 8, 11, 1, 7, 13, 7, 6, 6, 3, 12, 0, 0, 15, 6, 12, 12, 13, 7, 14, 14, 11, 15, 7, 14, 12, 6, 15, 2},
		{15, 2, 0, 12, 15, 14, 8, 14, 7, 14, 0, 3, 3, 11, 12, 2, 3, 14, 13, 5, 12, 9, 6, 11, 7, 4, 5, 1, 7, 12, 0, 11, 1, 5, 6, 6, 8, 6, 12, 2, 12, 3, 10, 3, 4, 10, 3, 3, 3, 10, 10, 14, 3, 13, 15, 0, 7, 6, 15, 6, 13, 7, 4, 11},
		{11, 15, 5, 14, 0, 1, 1, 14, 2, 3, 15, 14, 4, 3, 11, 1, 6, 6, 0, 12, 3, 5, 15, 6, 3, 11, 13, 11, 7, 7, 8, 11, 5, 9, 10, 10, 9, 14, 7, 1, 7, 2, 8, 6, 6, 5, 1, 9, 6, 5, 8, 14, 2, 14, 2, 9, 3, 3, 4, 15, 13, 5, 2, 7},
		{7, 8, 13, 9, 15, 8, 11, 7, 1, 9, 15, 12, 6, 9, 3, 1, 10, 10, 11, 0, 0, 8, 14, 5, 11, 12, 14, 4, 3, 9, 12, 9, 14, 0, 0, 9, 12, 4, 1, 13, 3, 6, 3, 4, 13, 10, 2, 9, 3, 7, 7, 10, 7, 10, 10, 3, 5, 15, 8, 9, 11, 7, 1, 14},
		{5, 5, 9, 1, 15, 3, 3, 11, 6, 11, 13, 13, 4, 12, 7, 12, 4, 8, 14, 13, 7, 12, 13, 8, 10, 2, 1, 12, 11, 7, 0, 8, 10, 9, 15, 1, 3, 9, 10, 0, 9, 1, 14, 1, 1, 9, 2, 2, 8, 9, 5, 6, 3, 2, 15, 9, 15, 6, 3, 11, 14, 4, 0, 4},
		{9, 2, 10, 2, 0, 9, 6, 13, 13, 0, 13, 14, 3, 12, 1, 15, 9, 3, 12, 2, 5, 15, 6, 6, 15, 11, 7, 11, 0, 4, 0, 11, 10, 12, 7, 9, 3, 0, 2, 2, 13, 13, 9, 6, 9, 2, 6, 4, 3, 6, 5, 10, 10, 9, 7, 2, 4, 9, 13, 11, 2, 13, 6, 8},
		{13, 15, 9, 8, 6, 2, 3, 2, 2, 12, 5, 3, 8, 6, 11, 6, 15, 7, 10, 3, 15, 8, 7, 5, 3, 8, 4, 2, 11, 1, 0, 4, 1, 1, 6, 1, 13, 6, 5, 1, 2, 6, 7, 10, 4, 3, 10, 6, 2, 0, 7, 13, 15, 1, 13, 0, 12, 10, 15, 6, 2, 4, 14, 3},
		{5, 11, 14, 4, 0, 7, 12, 4, 4, 14, 12, 3, 4, 10, 7, 14, 6, 4, 14, 7, 0, 12, 5, 9, 15, 6, 15, 6, 3, 12, 0, 10, 11, 7, 1, 14, 13, 5, 1, 14, 5, 15, 12, 1, 9, 13, 9, 13, 14, 5, 10, 11, 12, 10, 15, 11, 9, 13, 2, 14, 9, 12, 2, 11},
		{2, 12, 5, 7, 1, 5, 2, 11, 8, 4, 15, 6, 9, 14, 5, 1, 15, 4, 3, 1, 11, 4, 2, 1, 4, 5, 4, 4, 7, 3, 3, 12, 4, 3, 2, 15, 13, 1, 14, 15, 1, 4, 6, 11, 13, 15, 6, 12, 12, 13, 6, 8, 10, 0, 10, 12, 1, 10, 3, 2, 9, 8, 2, 8},
		{10, 12, 12, 6, 8, 5, 4, 4, 5, 3, 6, 7, 15, 5, 10, 3, 8, 15, 14, 5, 6, 2, 14, 4, 1, 7, 1, 3, 12, 3, 12, 4, 10, 15, 6, 6, 0, 6, 6, 8, 6, 9, 5, 7, 5, 1, 9, 2, 4, 9, 0, 8, 1, 1, 14, 3, 7, 14, 8, 9, 0, 4, 11, 7},
		{13, 11, 14, 7, 0, 4, 0, 10, 12, 11, 10, 8, 6, 12, 13, 15, 9, 2, 14, 9, 3, 0, 12, 14, 11, 15, 4, 7, 15, 14, 4, 8, 15, 12, 9, 14, 7, 7, 9, 13, 14, 14, 4, 9, 13, 8, 1, 13, 6, 3, 12, 7, 0, 15, 6, 15, 7, 2, 3, 0, 9, 5, 13, 0},
		{3, 8, 12, 11, 5, 9, 9, 14, 8, 14, 14, 5, 9, 9, 12, 10, 3, 12, 13, 0, 0, 0, 6, 7, 12, 4, 2, 3, 8, 8, 9, 15, 11, 1, 12, 13, 10, 15, 11, 1, 2, 13, 10, 1, 7, 2, 7, 11, 8, 15, 7, 6, 4, 6, 5, 11, 11, 15, 2, 1, 11, 1, 1, 8},
		{10, 7, 7, 1, 4, 13, 9, 10, 2, 2, 3, 7, 12, 8, 5, 5, 5, 5, 3, 1, 5, 6, 8, 2, 8, 11, 5, 0, 4, 12, 12, 6, 7, 9, 14, 10, 11, 8, 0, 9, 11, 4, 14, 7, 7, 8, 2, 15, 12, 7, 4, 4, 13, 2, 0, 3, 14, 0, 1, 5, 2, 15, 7, 11},
		{3, 8, 10, 4, 1, 7, 3, 13, 5, 14, 0, 9, 3, 1, 0, 11, 2, 15, 4, 9, 6, 5, 14, 0, 2, 8, 1, 14, 7, 6, 1, 5, 5, 7, 2, 0, 5, 3, 4, 15, 13, 10, 9, 13, 13, 12, 5, 11, 11, 14, 13, 10, 8, 14, 0, 8, 1, 7, 2, 10, 12, 12, 1, 11},
		{11, 14, 4, 13, 3, 11, 10, 6, 15, 2, 5, 10, 14, 4, 13, 3, 12, 7, 12, 10, 4, 0, 0, 1, 14, 6, 1, 2, 2, 12, 9, 2, 3, 11, 1, 4, 10, 4, 4, 7, 7, 12, 4, 3, 12, 11, 9, 3, 15, 13, 6, 13, 7, 11, 5, 12, 5, 13, 15, 12, 0, 13, 12, 9},
		{8, 7, 2, 2, 5, 3, 10, 15, 10, 8, 1, 0, 4, 5, 7, 6, 15, 13, 2, 14, 6, 2, 9, 5, 9, 0, 5, 12, 8, 6, 4, 12, 6, 8, 14, 15, 7, 15, 11, 2, 2, 12, 7, 9, 7, 11, 15, 7, 0, 4, 5, 13, 7, 2, 5, 9, 0, 5, 7, 6, 7, 12, 4, 1},
		{11, 4, 2, 13, 6, 10, 9, 4, 12, 9, 9, 6, 4, 2, 14, 14, 9, 5, 5, 15, 15, 9, 8, 11, 4, 2, 8, 11, 14, 3, 8, 10, 14, 9, 6, 6, 4, 7, 11, 2, 3, 7, 5, 1, 14, 2, 9, 4, 0, 1, 10, 7, 6, 7, 1, 3, 13, 7, 3, 2, 12, 3, 6, 6},
		{11, 1, 14, 3, 14, 3, 9, 9, 0, 11, 14, 6, 14, 7, 14, 8, 4, 2, 5, 6, 13, 3, 4, 10, 8, 8, 10, 11, 5, 1, 15, 15, 7, 0, 4, 14, 15, 13, 14, 13, 3, 2, 1, 6, 0, 6, 6, 4, 15, 6, 0, 12, 5, 11, 1, 7, 3, 3, 13, 12, 12, 6, 3, 2},
		{7, 2, 10, 14, 14, 13, 4, 14, 10, 6, 0, 2, 7, 7, 2, 5, 14, 1, 5, 14, 15, 1, 2, 9, 2, 13, 1, 3, 6, 1, 3, 13, 10, 6, 11, 13, 1, 7, 13, 15, 2, 11, 9, 6, 13, 7, 9, 2, 3, 13, 10, 10, 6, 2, 5, 9, 1, 3, 0, 3, 1, 5, 3, 12},
		{11, 14, 4, 2, 10, 11, 15, 5, 9, 7, 8, 11, 10, 9, 5, 7, 14, 3, 12, 2, 7, 15, 12, 15, 4, 15, 12, 9, 2, 6, 6, 6, 8, 5, 0, 7, 14, 15, 14, 14, 3, 12, 7, 12, 2, 4, 1, 7, 1, 3, 4, 7, 1, 9, 11, 15, 15, 3, 7, 1, 10, 9, 14, 14},
		{4, 13, 11, 1, 9, 6, 5, 1, 11, 6, 6, 8, 3, 9, 8, 15, 13, 12, 3, 13, 5, 9, 10, 5, 12, 1, 15, 14, 12, 1, 10, 11, 5, 7, 3, 12, 9, 12, 0, 2, 2, 3, 14, 4, 2, 13, 1, 15, 11, 8, 3, 13, 0, 10, 5, 4, 6, 0, 14, 8, 1, 0, 6, 15},
		{15, 2, 0, 5, 2, 14, 9, 0, 10, 5, 12, 8, 5, 6, 0, 1, 9, 4, 4, 1, 4, 6, 14, 5, 3, 0, 2, 2, 14, 9, 7, 0, 2, 15, 12, 0, 10, 12, 9, 12, 15, 1, 9, 4, 15, 3, 0, 13, 0, 6, 5, 0, 2, 6, 11, 9, 13, 15, 6, 3, 5, 4, 0, 8},
		{4, 14, 8, 14, 13, 4, 4, 10, 6, 12, 15, 11, 7, 2, 15, 6, 9, 9, 1, 11, 13, 2, 7, 10, 4, 4, 5, 12, 14, 15, 8, 5, 6, 1, 11, 15, 4, 11, 5, 2, 5, 7, 3, 4, 5, 7, 3, 8, 10, 13, 7, 5, 6, 5, 10, 1, 12, 13, 3, 6, 2, 8, 7, 15},
		{3, 15, 4, 9, 14, 12, 6, 1, 7, 0, 7, 15, 10, 6, 5, 5, 15, 5, 9, 4, 7, 6, 14, 2, 1, 4, 10, 3, 12, 1, 7, 1, 0, 10, 2, 11, 14, 13, 7, 10, 5, 11, 5, 11, 15, 5, 0, 3, 15, 1, 2, 14, 13, 13, 10, 9, 15, 12, 10, 5, 2, 10, 0, 6},
		{4, 6, 5, 13, 11, 10, 15, 4, 2, 15, 13, 6, 7, 7, 4, 0, 4, 6, 7, 4, 9, 1, 6, 7, 6, 1, 4, 2, 0, 11, 6, 3, 14, 5, 9, 2, 2, 10, 1, 2, 13, 14, 4, 11, 4, 7, 12, 9, 8, 2, 2, 9, 5, 7, 9, 12, 8, 15, 0, 9, 12, 11, 1, 12},
		{11, 12, 11, 9, 8, 15, 4, 12, 13, 10, 6, 6, 6, 12, 3, 0, 6, 15, 15, 10, 6, 12, 5, 7, 10, 2, 7, 1, 6, 12, 9, 11, 11, 14, 1, 12, 15, 0, 6, 2, 12, 15, 4, 15, 14, 8, 3, 4, 15, 4, 13, 3, 14, 1, 3, 7, 6, 13, 9, 1, 0, 12, 4, 14},
		{12, 11, 13, 10, 10, 10, 3, 7, 12, 3, 13, 9, 6, 0, 12, 10, 4, 11, 5, 4, 11, 5, 7, 14, 6, 10, 12, 12, 13, 15, 12, 1, 13, 15, 15, 7, 1, 2, 8, 6, 1, 12, 12, 0, 4, 3, 3, 3, 7, 8, 9, 10, 7, 7, 0, 0, 11, 13, 15, 4, 9, 5, 10, 9},
		{6, 12, 3, 0, 9, 11, 6, 4, 9, 9, 1, 5, 9, 14, 3, 7, 15, 3, 5, 0, 5, 11, 7, 6, 13, 5, 10, 2, 12, 10, 2, 6, 0, 1, 1, 13, 9, 3, 11, 7, 8, 2, 10, 9, 13, 6, 6, 4, 12, 0, 3, 10, 9, 4, 15, 11, 14, 1, 9, 3, 0, 14, 6, 1},
		{11, 14, 10, 10, 11, 6, 4, 7, 10, 0, 7, 9, 3, 2, 13, 13, 9, 9, 2, 3, 3, 14, 10, 4, 14, 1, 10, 7, 14, 4, 9, 15, 3, 11, 5, 10, 7, 8, 3, 0, 1, 2, 2, 3, 12, 9, 6, 2, 11, 15, 3, 9, 3, 6, 8, 0, 4, 5, 7, 3, 0, 14, 7, 9},
		{4, 11, 13, 12, 6, 2, 3, 15, 15, 3, 5, 1, 0, 5, 10, 2, 5, 3, 7, 10, 15, 0, 5, 3, 2, 10, 12, 10, 8, 3, 9, 15, 5, 3, 7, 13, 5, 7, 13, 12, 5, 10, 2, 9, 10, 1, 9, 4, 14, 1, 10, 13, 1, 2, 2, 12, 5, 3, 14, 7, 7, 8, 13, 13},
		{10, 12, 11, 10, 0, 15, 4, 3, 0, 8, 3, 0, 15, 0, 3, 10, 10, 9, 15, 3, 13, 3, 8, 3, 8, 2, 14, 7, 1, 6, 13, 8, 2, 2, 12, 3, 3, 0, 10, 12, 0, 1, 1, 7, 5, 0, 13, 10, 7, 13, 9, 9, 13, 7, 0, 1, 0, 2, 14, 2, 13, 0, 8, 3},
		{11, 3, 11, 10, 12, 15, 11, 6, 14, 8, 8, 5, 7, 11, 3, 1, 13, 7, 13, 4, 15, 7, 2, 3, 8, 7, 3, 8, 9, 15, 10, 15, 9, 0, 5, 4, 1, 7, 13, 8, 2, 7, 1, 10, 1, 12, 12, 1, 7, 12, 13, 5, 14, 10, 9, 15, 12, 2, 10, 3, 10, 3, 9, 12},
		{9, 8, 11, 0, 5, 6, 1, 5, 9, 1, 0, 12, 12, 0, 12, 11, 2, 8, 4, 0, 1, 7, 7, 5, 1, 14, 1, 9, 13, 7, 2, 12, 8, 9, 12, 13, 1, 11, 5, 3, 12, 14, 15, 4, 9, 8, 12, 7, 11, 1, 3, 9, 11, 5, 7, 14, 4, 6, 12, 3, 4, 12, 7, 9},
		{10, 12, 2, 14, 14, 1, 11, 8, 3, 7, 13, 7, 2, 1, 14, 13, 7, 6, 15, 8, 15, 12, 13, 10, 11, 15, 4, 2, 6, 13, 12, 3, 2, 10, 15, 14, 10, 11, 8, 14, 9, 3, 12, 9, 15, 2, 14, 14, 5, 13, 7, 6, 2, 1, 1, 4, 1, 0, 13, 10, 1, 0, 2, 9},
		{10, 5, 11, 14, 12, 1, 12, 7, 12, 8, 10, 5, 6, 10, 0, 7, 5, 6, 11, 11, 13, 12, 0, 13, 0, 6, 11, 0, 14, 4, 2, 1, 12, 7, 1, 10, 7, 15, 5, 3, 14, 15, 1, 3, 1, 2, 10, 4, 11, 8, 2, 11, 2, 5, 5, 4, 15, 5, 10, 3, 1, 7, 2, 14},
	}
	mat := generateMatrix(hash)

	if *mat != expectedMatrix {
		t.Fatal("The generated matrix doesn't match the test vector")
	}
}

func TestMatrix_HeavyHash(t *testing.T) {
	//original expect
	//expected, err := hex.DecodeString("87689f379943eaf9b7475ca95325687772bfcc68fc7899caeb4409ec4590c325")
	//karlsenhashv1 expect
	expected, err := hex.DecodeString("f32d27ec12b058b602252dc79dbd5c32958ef5c153731dfc8220d25b58ffa691")

	if err != nil {
		t.Fatal(err)
	}
	input := []byte{0xC1, 0xEC, 0xFD, 0xFC}
	writer := hashes.NewPoWHashWriter()
	writer.InfallibleWrite(input)
	hashed := testMatrix.HeavyHash(writer.Finalize())

	if !bytes.Equal(expected, hashed.ByteSlice()) {
		t.Fatalf("expected: %x == %s", expected, hashed)
	}

}

var testMatrix = matrix{
	{13, 2, 14, 13, 2, 15, 14, 3, 10, 4, 1, 8, 4, 3, 8, 15, 15, 15, 15, 15, 2, 11, 15, 15, 15, 1, 7, 12, 12, 4, 2, 0, 6, 1, 14, 10, 12, 14, 15, 8, 10, 12, 0, 5, 13, 3, 14, 10, 10, 6, 12, 11, 11, 7, 6, 6, 10, 2, 2, 4, 11, 12, 0, 5},
	{4, 13, 0, 2, 1, 15, 13, 13, 11, 2, 5, 12, 15, 7, 0, 10, 7, 2, 6, 3, 12, 0, 12, 0, 2, 6, 7, 7, 7, 7, 10, 12, 11, 14, 12, 12, 4, 11, 10, 0, 10, 11, 2, 10, 1, 7, 7, 12, 15, 9, 5, 14, 9, 12, 3, 0, 12, 13, 4, 13, 8, 15, 11, 6},
	{14, 6, 15, 9, 8, 2, 2, 12, 2, 3, 4, 12, 13, 15, 4, 5, 13, 4, 3, 0, 14, 3, 5, 14, 3, 13, 4, 15, 9, 12, 7, 15, 5, 1, 13, 12, 9, 9, 8, 11, 14, 11, 4, 10, 12, 6, 12, 8, 6, 3, 9, 8, 1, 6, 0, 5, 8, 9, 12, 5, 14, 15, 2, 2},
	{9, 6, 7, 6, 0, 11, 5, 6, 2, 14, 12, 6, 4, 13, 8, 9, 2, 1, 9, 7, 4, 5, 10, 8, 11, 11, 11, 15, 7, 11, 1, 14, 3, 8, 14, 8, 2, 8, 13, 7, 8, 8, 15, 7, 1, 13, 7, 9, 1, 7, 15, 15, 0, 0, 12, 15, 13, 5, 13, 10, 1, 5, 6, 13},
	{4, 0, 12, 10, 6, 11, 14, 2, 2, 15, 4, 1, 2, 4, 2, 12, 13, 1, 9, 10, 8, 0, 2, 10, 13, 8, 9, 7, 5, 3, 8, 2, 6, 6, 1, 12, 3, 0, 1, 4, 2, 8, 3, 13, 6, 15, 0, 13, 14, 4, 15, 0, 7, 3, 7, 8, 5, 14, 14, 5, 5, 0, 1, 2},
	{12, 14, 6, 3, 3, 4, 6, 7, 1, 3, 2, 7, 15, 15, 15, 10, 9, 12, 0, 6, 3, 8, 5, 0, 13, 5, 0, 6, 0, 14, 2, 12, 10, 4, 11, 2, 10, 7, 7, 6, 8, 11, 4, 4, 11, 9, 3, 12, 10, 5, 2, 6, 5, 5, 10, 13, 12, 10, 1, 6, 14, 7, 12, 4},
	{7, 14, 6, 7, 7, 12, 4, 1, 8, 6, 8, 13, 13, 5, 12, 14, 10, 8, 6, 2, 12, 3, 8, 15, 5, 15, 15, 3, 14, 0, 8, 6, 9, 12, 9, 7, 3, 8, 4, 0, 7, 14, 3, 3, 13, 14, 3, 7, 3, 2, 2, 3, 3, 12, 6, 7, 4, 1, 14, 10, 6, 10, 2, 9},
	{14, 11, 15, 5, 7, 10, 1, 11, 4, 2, 6, 2, 9, 7, 4, 0, 9, 12, 11, 2, 3, 13, 1, 5, 4, 10, 5, 6, 6, 12, 8, 1, 1, 15, 4, 2, 12, 12, 0, 4, 14, 3, 11, 1, 7, 5, 9, 4, 3, 15, 7, 3, 15, 9, 8, 3, 8, 3, 3, 6, 7, 6, 9, 2},
	{10, 4, 6, 10, 5, 2, 15, 12, 0, 14, 14, 15, 14, 0, 12, 9, 1, 12, 4, 5, 5, 2, 10, 4, 2, 13, 11, 3, 1, 8, 10, 0, 7, 0, 12, 4, 11, 1, 14, 6, 14, 5, 5, 11, 11, 1, 3, 8, 0, 6, 11, 11, 8, 4, 7, 6, 14, 4, 9, 14, 9, 7, 13, 9},
	{12, 7, 9, 8, 2, 3, 3, 5, 14, 8, 0, 9, 7, 4, 2, 15, 15, 3, 11, 11, 8, 5, 7, 5, 0, 15, 10, 8, 0, 13, 1, 14, 8, 10, 1, 4, 13, 1, 13, 3, 11, 11, 2, 3, 10, 6, 8, 14, 15, 2, 10, 10, 12, 7, 7, 6, 6, 3, 13, 8, 1, 14, 2, 1},
	{2, 11, 6, 9, 13, 3, 12, 6, 0, 4, 6, 13, 8, 14, 6, 9, 10, 2, 10, 8, 4, 13, 6, 5, 0, 13, 15, 4, 2, 2, 1, 7, 5, 3, 3, 13, 7, 3, 5, 9, 15, 14, 14, 6, 0, 15, 11, 2, 4, 15, 6, 9, 8, 9, 15, 2, 6, 9, 15, 8, 4, 4, 11, 1},
	{10, 11, 8, 3, 11, 13, 10, 2, 2, 5, 2, 14, 15, 10, 2, 11, 0, 1, 8, 2, 14, 1, 10, 0, 3, 7, 5, 10, 7, 8, 15, 7, 2, 5, 13, 4, 10, 3, 6, 2, 3, 9, 6, 11, 7, 14, 1, 11, 9, 3, 3, 7, 6, 0, 9, 11, 4, 10, 4, 1, 9, 7, 4, 15},
	{13, 8, 15, 14, 11, 12, 5, 3, 9, 14, 1, 5, 14, 13, 14, 5, 13, 5, 4, 10, 9, 9, 0, 0, 6, 12, 5, 7, 2, 7, 2, 6, 6, 6, 1, 12, 9, 15, 7, 11, 11, 10, 11, 1, 10, 10, 0, 8, 1, 4, 5, 5, 8, 10, 10, 15, 6, 8, 13, 11, 11, 3, 15, 5},
	{8, 11, 5, 10, 1, 10, 9, 1, 12, 7, 6, 11, 1, 1, 4, 1, 2, 8, 4, 4, 7, 7, 8, 2, 7, 1, 14, 1, 8, 15, 15, 12, 10, 4, 15, 11, 3, 6, 10, 7, 4, 0, 10, 9, 11, 7, 1, 14, 4, 14, 3, 14, 10, 4, 13, 12, 5, 3, 12, 7, 10, 8, 0, 3},
	{9, 11, 6, 15, 14, 10, 0, 4, 7, 7, 6, 0, 7, 7, 12, 15, 5, 4, 12, 3, 7, 3, 0, 12, 2, 7, 11, 6, 7, 3, 2, 8, 5, 11, 9, 4, 3, 8, 11, 12, 3, 5, 14, 12, 4, 13, 12, 0, 3, 14, 4, 9, 1, 1, 9, 14, 10, 14, 8, 15, 6, 14, 10, 15},
	{10, 14, 10, 0, 10, 12, 15, 0, 3, 9, 11, 10, 3, 5, 1, 1, 9, 1, 7, 15, 7, 8, 10, 10, 12, 11, 5, 1, 10, 3, 6, 6, 13, 0, 13, 1, 4, 5, 9, 4, 9, 15, 8, 4, 13, 13, 4, 5, 5, 11, 1, 13, 15, 3, 10, 15, 7, 11, 10, 15, 8, 12, 10, 3},
	{8, 5, 11, 3, 8, 13, 15, 15, 3, 12, 1, 13, 1, 7, 1, 5, 6, 13, 7, 8, 5, 1, 12, 3, 10, 7, 12, 6, 14, 12, 15, 5, 3, 12, 2, 15, 11, 13, 1, 13, 8, 5, 8, 0, 13, 15, 7, 13, 6, 13, 10, 1, 11, 0, 8, 9, 5, 11, 2, 9, 9, 10, 4, 15},
	{0, 4, 12, 14, 3, 1, 7, 5, 11, 13, 5, 3, 11, 12, 6, 8, 10, 15, 11, 8, 7, 10, 0, 2, 5, 15, 6, 10, 4, 2, 3, 1, 13, 7, 6, 12, 14, 7, 6, 14, 12, 10, 6, 14, 12, 0, 12, 11, 6, 9, 3, 1, 12, 15, 15, 3, 5, 5, 10, 11, 7, 15, 13, 3},
	{12, 14, 2, 14, 13, 6, 15, 7, 8, 8, 14, 13, 9, 2, 2, 10, 3, 15, 6, 10, 11, 7, 13, 0, 12, 1, 5, 8, 8, 12, 1, 11, 1, 3, 2, 4, 10, 7, 7, 7, 3, 10, 7, 2, 2, 3, 0, 1, 13, 5, 8, 2, 14, 0, 11, 13, 9, 3, 13, 2, 14, 2, 15, 4},
	{0, 0, 13, 6, 9, 12, 15, 7, 8, 0, 7, 4, 12, 15, 3, 2, 7, 1, 14, 4, 9, 3, 13, 12, 11, 12, 9, 9, 3, 7, 10, 9, 1, 9, 10, 2, 10, 14, 11, 0, 14, 4, 15, 12, 12, 9, 9, 8, 14, 1, 9, 14, 0, 6, 1, 0, 13, 9, 7, 6, 13, 2, 3, 9},
	{8, 0, 10, 13, 0, 7, 9, 7, 5, 1, 0, 3, 7, 10, 3, 15, 1, 15, 3, 11, 2, 6, 3, 10, 0, 10, 10, 3, 4, 15, 8, 6, 11, 11, 7, 5, 8, 5, 7, 15, 1, 11, 7, 13, 13, 6, 13, 13, 4, 2, 3, 15, 9, 5, 10, 6, 6, 6, 3, 11, 15, 13, 1, 15},
	{1, 1, 2, 10, 2, 2, 9, 5, 9, 2, 0, 1, 14, 2, 11, 6, 11, 6, 1, 0, 13, 7, 14, 1, 15, 14, 13, 7, 12, 11, 8, 11, 2, 11, 6, 10, 2, 3, 0, 0, 15, 0, 4, 6, 4, 12, 5, 5, 7, 14, 10, 6, 0, 3, 13, 0, 8, 1, 13, 10, 5, 1, 7, 5},
	{0, 5, 2, 12, 10, 2, 5, 1, 14, 0, 1, 4, 15, 11, 8, 7, 11, 14, 15, 6, 4, 1, 6, 6, 7, 13, 12, 5, 13, 2, 1, 6, 2, 13, 5, 15, 0, 8, 8, 6, 5, 5, 2, 0, 3, 13, 14, 2, 10, 5, 7, 6, 14, 5, 1, 4, 11, 2, 11, 1, 8, 15, 2, 4},
	{9, 9, 4, 5, 2, 5, 3, 12, 14, 5, 1, 3, 3, 0, 0, 6, 7, 14, 0, 15, 14, 11, 3, 10, 1, 9, 4, 14, 7, 14, 1, 0, 15, 11, 5, 9, 4, 0, 0, 10, 4, 4, 0, 7, 8, 15, 12, 8, 10, 8, 1, 2, 1, 11, 12, 14, 14, 14, 8, 10, 1, 5, 13, 10},
	{5, 10, 4, 4, 11, 10, 0, 6, 0, 12, 10, 5, 9, 11, 8, 10, 11, 3, 11, 14, 12, 9, 4, 6, 11, 12, 8, 7, 6, 14, 0, 6, 12, 4, 5, 3, 9, 0, 11, 6, 1, 3, 2, 12, 8, 9, 7, 12, 14, 7, 12, 6, 11, 13, 0, 2, 1, 3, 1, 8, 12, 2, 15, 15},
	{10, 11, 2, 3, 11, 10, 1, 7, 1, 10, 10, 14, 5, 13, 10, 3, 11, 15, 9, 14, 11, 11, 3, 15, 11, 6, 15, 13, 13, 1, 1, 10, 5, 1, 5, 11, 10, 3, 9, 12, 12, 1, 5, 6, 3, 3, 1, 1, 12, 8, 3, 15, 6, 2, 8, 14, 3, 4, 10, 9, 7, 13, 2, 6},
	{12, 0, 1, 0, 4, 3, 3, 6, 8, 3, 1, 13, 6, 12, 1, 1, 1, 4, 12, 4, 4, 9, 9, 14, 15, 3, 6, 4, 11, 1, 12, 5, 6, 0, 10, 9, 1, 8, 14, 5, 2, 8, 4, 15, 12, 13, 7, 14, 12, 2, 6, 9, 4, 13, 0, 15, 10, 10, 6, 12, 7, 12, 9, 10},
	{0, 8, 5, 11, 12, 12, 11, 7, 2, 9, 2, 15, 1, 1, 0, 0, 6, 5, 10, 1, 11, 12, 8, 7, 1, 7, 10, 4, 2, 8, 2, 5, 1, 1, 2, 9, 2, 0, 3, 7, 5, 1, 5, 5, 3, 1, 4, 3, 14, 8, 11, 7, 8, 0, 2, 13, 3, 15, 1, 13, 14, 15, 11, 13},
	{8, 13, 5, 14, 2, 9, 9, 13, 15, 8, 2, 14, 4, 2, 6, 0, 1, 13, 10, 13, 6, 12, 15, 11, 6, 11, 9, 9, 2, 9, 6, 14, 2, 9, 12, 1, 13, 9, 5, 11, 10, 4, 4, 5, 8, 9, 13, 10, 9, 0, 5, 15, 4, 12, 7, 10, 6, 5, 5, 15, 8, 8, 11, 14},
	{6, 9, 6, 7, 1, 15, 0, 1, 4, 15, 5, 3, 10, 9, 15, 9, 14, 12, 7, 6, 3, 0, 12, 8, 12, 2, 11, 8, 11, 8, 1, 10, 10, 7, 7, 5, 3, 5, 1, 2, 13, 11, 2, 5, 2, 10, 10, 1, 14, 14, 8, 1, 11, 1, 2, 6, 15, 10, 8, 7, 10, 7, 0, 3},
	{12, 6, 11, 1, 1, 7, 8, 1, 5, 5, 8, 4, 6, 5, 6, 4, 2, 8, 4, 1, 0, 0, 14, 2, 10, 14, 14, 11, 2, 9, 14, 15, 12, 14, 9, 3, 7, 14, 4, 7, 12, 9, 3, 5, 1, 0, 12, 9, 10, 5, 11, 12, 10, 10, 6, 14, 6, 13, 13, 5, 5, 10, 13, 10},
	{12, 6, 13, 0, 8, 0, 10, 6, 15, 15, 7, 3, 0, 10, 13, 14, 10, 13, 5, 13, 15, 14, 3, 4, 10, 10, 9, 6, 6, 15, 2, 7, 0, 10, 6, 14, 2, 9, 11, 7, 5, 5, 13, 14, 11, 15, 9, 4, 2, 0, 15, 5, 4, 14, 14, 1, 3, 4, 5, 8, 1, 1, 10, 12},
	{2, 5, 0, 4, 11, 5, 5, 6, 10, 4, 6, 7, 10, 3, 0, 14, 14, 0, 12, 15, 11, 12, 13, 7, 6, 3, 9, 1, 9, 8, 8, 8, 4, 10, 3, 1, 7, 10, 3, 2, 12, 6, 15, 14, 0, 6, 8, 10, 1, 9, 12, 12, 15, 7, 1, 11, 15, 13, 0, 4, 10, 0, 12, 11},
	{8, 12, 14, 15, 14, 15, 10, 0, 2, 14, 3, 1, 2, 6, 0, 2, 1, 7, 9, 0, 15, 13, 5, 14, 6, 8, 15, 4, 15, 6, 10, 6, 15, 3, 12, 8, 5, 4, 10, 5, 3, 0, 4, 13, 10, 9, 8, 4, 6, 3, 9, 6, 12, 11, 9, 13, 8, 10, 9, 9, 8, 12, 1, 2},
	{11, 10, 15, 15, 5, 14, 15, 7, 5, 9, 14, 14, 7, 11, 6, 6, 3, 8, 2, 3, 4, 14, 11, 1, 12, 15, 11, 6, 0, 0, 13, 7, 14, 3, 12, 14, 0, 15, 6, 1, 11, 2, 11, 8, 3, 13, 4, 12, 10, 13, 7, 14, 9, 13, 3, 10, 2, 14, 13, 4, 12, 13, 14, 10},
	{1, 11, 2, 12, 1, 10, 7, 12, 3, 3, 14, 9, 1, 10, 0, 11, 8, 10, 12, 12, 4, 12, 2, 11, 5, 0, 3, 15, 8, 2, 14, 3, 10, 2, 1, 13, 6, 14, 0, 0, 8, 11, 6, 13, 15, 10, 12, 7, 7, 11, 14, 9, 2, 7, 6, 8, 14, 9, 14, 10, 11, 9, 9, 12},
	{5, 10, 14, 2, 1, 4, 11, 5, 10, 2, 13, 9, 6, 12, 11, 5, 13, 4, 5, 14, 8, 7, 15, 9, 8, 4, 5, 2, 9, 11, 5, 3, 12, 2, 6, 1, 7, 4, 11, 4, 15, 0, 5, 2, 13, 11, 11, 2, 15, 10, 0, 12, 5, 8, 10, 1, 4, 11, 3, 13, 11, 7, 9, 14},
	{9, 8, 10, 5, 0, 2, 5, 8, 7, 3, 3, 6, 11, 1, 13, 15, 4, 4, 11, 6, 2, 6, 13, 11, 2, 6, 9, 4, 5, 13, 12, 2, 8, 7, 7, 12, 14, 15, 5, 12, 7, 0, 15, 15, 0, 5, 15, 0, 3, 9, 10, 15, 9, 11, 10, 10, 5, 3, 9, 3, 12, 13, 0, 13},
	{1, 11, 15, 0, 10, 5, 3, 5, 6, 7, 1, 11, 4, 11, 4, 2, 5, 12, 2, 5, 5, 6, 1, 5, 14, 9, 1, 5, 14, 12, 6, 10, 0, 8, 5, 11, 11, 11, 12, 10, 8, 10, 10, 1, 14, 1, 0, 8, 4, 7, 0, 11, 3, 1, 11, 12, 11, 8, 14, 15, 9, 3, 1, 14},
	{14, 11, 12, 12, 4, 6, 8, 14, 15, 1, 11, 2, 13, 3, 6, 2, 7, 1, 8, 1, 4, 9, 11, 15, 8, 1, 10, 13, 4, 13, 2, 7, 7, 10, 5, 2, 12, 12, 12, 3, 10, 8, 2, 11, 0, 3, 8, 9, 4, 2, 15, 7, 15, 6, 4, 6, 12, 7, 14, 9, 9, 8, 14, 12},
	{15, 4, 8, 12, 11, 11, 9, 5, 0, 0, 7, 6, 10, 5, 8, 2, 5, 6, 14, 11, 13, 0, 13, 15, 5, 4, 9, 15, 13, 12, 14, 15, 10, 2, 3, 6, 10, 14, 1, 8, 6, 7, 10, 1, 14, 9, 12, 13, 7, 2, 12, 10, 6, 11, 15, 1, 15, 11, 13, 0, 6, 13, 7, 15},
	{3, 3, 12, 5, 14, 9, 14, 14, 8, 0, 9, 1, 2, 2, 14, 11, 7, 1, 3, 1, 14, 15, 12, 8, 14, 2, 4, 13, 10, 5, 10, 8, 1, 7, 6, 5, 4, 2, 11, 5, 4, 13, 14, 6, 13, 15, 6, 6, 7, 12, 11, 5, 13, 10, 9, 13, 9, 14, 5, 6, 7, 14, 11, 7},
	{14, 12, 11, 5, 0, 5, 10, 5, 7, 1, 7, 11, 1, 0, 13, 6, 5, 14, 3, 0, 5, 14, 6, 7, 8, 5, 8, 6, 6, 3, 6, 1, 8, 3, 10, 7, 15, 6, 11, 6, 6, 7, 13, 2, 2, 0, 0, 11, 1, 15, 2, 14, 5, 1, 4, 8, 0, 1, 8, 0, 1, 1, 2, 2},
	{10, 13, 13, 3, 15, 14, 9, 12, 15, 15, 8, 5, 8, 10, 5, 9, 6, 6, 7, 15, 1, 0, 14, 9, 1, 11, 6, 11, 13, 4, 6, 14, 9, 12, 13, 8, 14, 6, 14, 2, 3, 15, 4, 4, 14, 4, 9, 12, 8, 0, 9, 11, 13, 10, 8, 14, 3, 5, 7, 11, 6, 7, 15, 2},
	{9, 9, 11, 6, 11, 0, 5, 4, 8, 10, 8, 11, 2, 12, 8, 7, 11, 13, 6, 1, 13, 13, 11, 4, 5, 7, 7, 9, 6, 4, 12, 0, 11, 8, 6, 12, 11, 4, 15, 11, 12, 8, 11, 11, 1, 3, 6, 14, 9, 6, 7, 5, 0, 10, 3, 15, 13, 7, 0, 1, 13, 15, 1, 14},
	{10, 6, 8, 7, 3, 6, 9, 15, 1, 3, 10, 14, 9, 0, 0, 10, 0, 15, 2, 0, 0, 0, 6, 0, 13, 9, 9, 1, 8, 6, 13, 2, 1, 9, 14, 9, 1, 4, 8, 4, 2, 0, 8, 5, 0, 11, 12, 15, 13, 1, 14, 14, 15, 7, 8, 4, 4, 12, 1, 12, 8, 3, 9, 5},
	{12, 11, 1, 4, 10, 14, 8, 12, 2, 4, 15, 2, 9, 7, 7, 11, 15, 12, 10, 11, 7, 4, 13, 0, 8, 6, 8, 8, 10, 5, 5, 13, 3, 7, 9, 13, 13, 14, 6, 8, 1, 5, 7, 12, 4, 4, 6, 9, 13, 1, 6, 1, 6, 14, 5, 8, 2, 10, 4, 10, 1, 9, 6, 15},
	{4, 13, 4, 9, 6, 11, 1, 8, 7, 11, 11, 1, 3, 10, 12, 11, 1, 10, 6, 10, 0, 7, 3, 0, 0, 6, 3, 9, 2, 1, 4, 8, 2, 10, 2, 15, 9, 15, 14, 14, 15, 14, 3, 2, 7, 6, 6, 10, 8, 8, 4, 11, 1, 13, 6, 0, 2, 10, 0, 11, 15, 14, 6, 9},
	{15, 0, 12, 13, 0, 9, 10, 4, 11, 5, 10, 0, 8, 7, 3, 2, 12, 6, 3, 8, 5, 15, 14, 2, 13, 13, 6, 11, 5, 6, 9, 10, 14, 5, 14, 4, 9, 7, 5, 11, 13, 2, 7, 1, 14, 9, 0, 7, 8, 12, 11, 15, 2, 1, 5, 11, 3, 7, 5, 1, 6, 3, 8, 6},
	{0, 3, 8, 1, 4, 6, 3, 1, 3, 8, 2, 0, 15, 15, 14, 15, 13, 10, 11, 9, 2, 11, 5, 12, 3, 3, 0, 1, 5, 3, 11, 6, 10, 11, 8, 5, 7, 15, 4, 12, 8, 8, 12, 12, 12, 1, 9, 4, 11, 6, 10, 11, 1, 12, 8, 12, 5, 6, 1, 14, 2, 10, 3, 0},
	{10, 13, 6, 9, 11, 1, 4, 10, 0, 13, 8, 7, 4, 12, 15, 5, 14, 12, 6, 9, 0, 0, 10, 5, 13, 10, 15, 3, 0, 8, 7, 0, 9, 8, 10, 6, 11, 8, 10, 13, 11, 7, 5, 5, 9, 13, 1, 15, 0, 5, 15, 5, 4, 7, 9, 9, 15, 8, 2, 6, 3, 8, 5, 8},
	{14, 0, 6, 2, 4, 12, 2, 13, 6, 10, 5, 2, 2, 1, 6, 11, 1, 6, 9, 13, 0, 13, 9, 3, 12, 4, 3, 8, 7, 0, 9, 12, 0, 1, 7, 10, 10, 7, 3, 9, 13, 5, 15, 4, 13, 0, 8, 5, 4, 14, 11, 3, 3, 13, 15, 9, 9, 12, 9, 5, 2, 0, 1, 14},
	{4, 14, 13, 0, 14, 15, 11, 10, 11, 1, 3, 3, 9, 1, 12, 8, 6, 5, 15, 11, 1, 7, 5, 3, 8, 13, 0, 13, 11, 5, 8, 1, 8, 6, 13, 4, 13, 7, 12, 6, 5, 5, 7, 0, 12, 1, 1, 8, 1, 6, 4, 2, 8, 8, 15, 11, 11, 11, 4, 4, 4, 7, 13, 12},
	{14, 15, 10, 0, 4, 3, 1, 9, 13, 7, 9, 9, 15, 5, 0, 3, 9, 6, 4, 7, 13, 11, 3, 2, 7, 1, 6, 8, 13, 7, 10, 4, 3, 9, 5, 9, 2, 6, 10, 7, 9, 13, 2, 14, 2, 14, 7, 2, 14, 2, 8, 8, 0, 9, 0, 9, 12, 6, 7, 7, 6, 8, 12, 13},
	{5, 15, 8, 12, 11, 3, 13, 4, 5, 14, 10, 4, 15, 15, 1, 10, 9, 14, 6, 6, 4, 12, 4, 9, 12, 2, 15, 13, 2, 5, 12, 2, 3, 2, 15, 11, 12, 2, 6, 2, 11, 6, 7, 9, 12, 10, 5, 1, 1, 5, 9, 6, 14, 11, 3, 11, 6, 10, 11, 11, 0, 12, 15, 1},
	{12, 6, 8, 10, 2, 5, 7, 9, 8, 14, 15, 15, 13, 10, 15, 3, 10, 10, 6, 10, 14, 10, 7, 5, 3, 7, 6, 12, 11, 12, 8, 9, 12, 9, 15, 15, 15, 7, 8, 3, 15, 14, 1, 12, 0, 0, 4, 0, 9, 10, 8, 7, 14, 10, 8, 14, 6, 2, 8, 1, 11, 10, 0, 1},
	{12, 1, 2, 12, 7, 10, 4, 11, 5, 14, 10, 2, 2, 9, 4, 13, 3, 14, 3, 15, 5, 0, 14, 7, 7, 15, 6, 5, 2, 8, 15, 9, 6, 6, 13, 10, 9, 8, 6, 3, 14, 7, 12, 9, 7, 8, 13, 12, 14, 13, 6, 0, 5, 1, 9, 12, 14, 0, 11, 11, 6, 3, 11, 7},
	{15, 4, 8, 12, 8, 11, 4, 15, 1, 6, 2, 13, 1, 7, 7, 12, 0, 8, 14, 14, 10, 14, 0, 12, 0, 3, 3, 11, 7, 4, 2, 13, 0, 0, 11, 2, 5, 8, 12, 11, 6, 5, 6, 0, 0, 4, 0, 0, 1, 9, 9, 11, 3, 2, 13, 4, 13, 9, 15, 4, 7, 8, 3, 2},
	{3, 13, 8, 8, 12, 10, 5, 4, 7, 13, 10, 13, 14, 3, 2, 12, 11, 0, 9, 5, 6, 4, 14, 4, 6, 9, 2, 5, 10, 3, 9, 10, 5, 0, 12, 5, 15, 5, 15, 15, 2, 12, 3, 11, 0, 15, 9, 14, 1, 5, 6, 6, 14, 5, 8, 0, 5, 9, 3, 7, 7, 12, 15, 1},
	{1, 11, 7, 4, 13, 3, 0, 8, 11, 9, 15, 1, 4, 12, 2, 12, 10, 4, 14, 3, 9, 14, 14, 2, 3, 11, 12, 4, 5, 10, 6, 15, 2, 13, 13, 9, 9, 1, 11, 12, 12, 14, 1, 5, 15, 1, 7, 14, 12, 10, 11, 13, 13, 5, 2, 4, 7, 7, 9, 4, 14, 15, 13, 10},
	{14, 15, 9, 14, 9, 5, 13, 2, 0, 0, 14, 8, 6, 2, 0, 7, 11, 10, 2, 13, 2, 14, 9, 6, 4, 11, 5, 14, 6, 1, 6, 14, 6, 3, 9, 5, 2, 9, 3, 11, 1, 14, 5, 4, 12, 5, 3, 5, 11, 3, 11, 6, 13, 7, 13, 7, 4, 9, 4, 13, 8, 3, 5, 11},
	{13, 12, 12, 13, 8, 2, 4, 2, 10, 6, 3, 5, 7, 7, 6, 13, 8, 6, 15, 4, 12, 7, 15, 4, 3, 9, 8, 15, 0, 3, 12, 1, 9, 8, 13, 10, 15, 4, 14, 1, 6, 15, 0, 4, 8, 9, 3, 1, 3, 15, 5, 5, 1, 11, 11, 10, 11, 10, 8, 8, 5, 4, 13, 0},
	{8, 4, 15, 9, 14, 9, 5, 8, 8, 10, 5, 15, 9, 8, 12, 5, 11, 10, 2, 12, 13, 1, 0, 2, 6, 13, 11, 9, 12, 0, 5, 0, 11, 5, 14, 12, 3, 4, 2, 10, 3, 12, 5, 15, 4, 8, 14, 1, 0, 13, 9, 5, 2, 4, 13, 8, 2, 5, 8, 9, 15, 3, 5, 5},
	{0, 3, 3, 4, 6, 5, 5, 1, 3, 2, 14, 5, 10, 7, 15, 11, 7, 13, 15, 4, 0, 12, 9, 15, 12, 0, 3, 1, 14, 1, 12, 9, 13, 8, 9, 15, 12, 3, 5, 11, 3, 11, 4, 1, 9, 4, 13, 7, 4, 10, 6, 14, 13, 0, 9, 11, 15, 15, 3, 3, 13, 15, 10, 15},
}
