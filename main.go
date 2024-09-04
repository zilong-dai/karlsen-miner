package main

import (
	"encoding/binary"
	"fmt"

	"github.com/pkg/errors"
	"github.com/zilong-dai/karlsen-miner/consensus/utils/hashes"
	"github.com/zilong-dai/karlsen-miner/consensus/utils/pow"
	"github.com/zilong-dai/karlsen-miner/consensus/utils/serialization"
)

func main() {
	timestamp := int64(1702373833378)
	nonce := uint64(0x85607312505c273c)
	prePowHash := [4]uint64{1783683135831606672,11442366678174958974,1054617894611757111,6012950097848553811}
	prePowHashBytes := make([]byte, 32)
	for i, num := range prePowHash {
		binary.LittleEndian.PutUint64(prePowHashBytes[i*8:i*8+8], num)
	}

	writer := hashes.NewPoWHashWriter()
	writer.InfallibleWrite(prePowHashBytes)
	err := serialization.WriteElement(writer, timestamp)
	if err != nil {
		panic(errors.Wrap(err, "this should never happen. Hash digest should never return an error"))
	}

	zeroes := [32]byte{}
	writer.InfallibleWrite(zeroes[:])
	err = serialization.WriteElement(writer, nonce)
	if err != nil {
		panic(errors.Wrap(err, "this should never happen. Hash digest should never return an error"))
	}

	powHash := writer.Finalize()

	context := pow.GetContext(false)

	middleHash := pow.Fishhash(context, powHash)
	writer2 := hashes.NewPoWHashWriter()
	writer2.InfallibleWrite(middleHash.ByteSlice())
	finalHash := writer2.Finalize()

	fmt.Println("finalHash:", finalHash)

}
