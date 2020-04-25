package main

import (
	"encoding/binary"
	"errors"
	"github.com/fogleman/nes/nes"
	"io"
)

const iNESFileMagic = 0x1a53454e

type iNESFileHeader struct {
	Magic    uint32  // iNES magic number
	NumPRG   byte    // number of PRG-ROM banks (16KB each)
	NumCHR   byte    // number of CHR-ROM banks (8KB each)
	Control1 byte    // control bits
	Control2 byte    // control bits
	NumRAM   byte    // PRG-RAM size (x 8KB)
	_        [7]byte // unused padding
}

func LoadNES(reader io.Reader) (*nes.Cartridge, error) {
	// open file
	// read file header
	header := iNESFileHeader{}
	if err := binary.Read(reader, binary.LittleEndian, &header); err != nil {
		return nil, err
	}

	// verify header magic number
	if header.Magic != iNESFileMagic {
		return nil, errors.New("invalid .nes file")
	}

	// mapper type
	mapper1 := header.Control1 >> 4
	mapper2 := header.Control2 >> 4
	mapper := mapper1 | mapper2<<4

	// mirroring type
	mirror1 := header.Control1 & 1
	mirror2 := (header.Control1 >> 3) & 1
	mirror := mirror1 | mirror2<<1

	// battery-backed RAM
	battery := (header.Control1 >> 1) & 1

	// read trainer if present (unused)
	if header.Control1&4 == 4 {
		trainer := make([]byte, 512)
		if _, err := io.ReadFull(reader, trainer); err != nil {
			return nil, err
		}
	}

	// read prg-rom bank(s)
	prg := make([]byte, int(header.NumPRG)*16384)
	if _, err := io.ReadFull(reader, prg); err != nil {
		return nil, err
	}

	// read chr-rom bank(s)
	chr := make([]byte, int(header.NumCHR)*8192)
	if _, err := io.ReadFull(reader, chr); err != nil {
		return nil, err
	}

	// provide chr-rom/ram if not in file
	if header.NumCHR == 0 {
		chr = make([]byte, 8192)
	}

	// success
	return nes.NewCartridge(prg, chr, mapper, mirror, battery), nil
}
