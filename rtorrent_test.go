package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatRatio(t *testing.T) {
	torrent := Torrent{
		RatioRaw: 13505,
	}

	assert.Equal(t, "13.51", torrent.FormatRatio())
}

func TestFormatRatioLowNumber(t *testing.T) {
	torrent := Torrent{
		RatioRaw: 296,
	}

	assert.Equal(t, "0.30", torrent.FormatRatio())
}

func TestCompletedPercentageZero(t *testing.T) {
	torrent := Torrent{
		BytesDoneRaw: 0,
		SizeBytesRaw: 14450390875,
	}

	assert.Equal(t, "0%", torrent.CalculateCompletedPercentage())
}

func TestCompletedPercentageSomething(t *testing.T) {
	torrent := Torrent{
		BytesDoneRaw: 3845519195,
		SizeBytesRaw: 14450390875,
	}

	assert.Equal(t, "27%", torrent.CalculateCompletedPercentage())
}
