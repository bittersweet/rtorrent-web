package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatRatio(t *testing.T) {
	torrent := Torrent{
		RatioRaw: 13505,
	}

	fmt.Println(torrent)
	assert.Equal(t, "13.51", torrent.FormatRatio())
}

func TestFormatRatioLowNumber(t *testing.T) {
	torrent := Torrent{
		RatioRaw: 296,
	}

	fmt.Println(torrent)
	assert.Equal(t, "0.30", torrent.FormatRatio())
}
