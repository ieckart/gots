/*
MIT License

Copyright 2016 Comcast Cable Communications Management, LLC

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package pes

const (
	// MaxPts is the max value allowed for a PTS time, 2^33 - 1
	MaxPts = 8589934591
	// UpperPtsRolloverThreshold is the threshold for a rollover on the upper end, maxPts = 30 min
	UpperPtsRolloverThreshold = 8427934591
	// LowerPtsRolloverThreshold is the threshold for a rollover on the lower end, 30 min
	LowerPtsRolloverThreshold = 162000000
)

// After checks if this PTS is after the other PTS
func (p PTS) After(other PTS) bool {
	switch {
	case other == PtsPositiveInfinity:
		return false
	case other == PtsNegativeInfinity:
		return true
	case p.RolledOver(other):
		return true
	case other.RolledOver(p):
		return false
	default:
		return p > other
	}
}

// GreaterOrEqual returns true if the method reciever is >= the provided PTS
func (p PTS) GreaterOrEqual(other PTS) bool {
	if p == other {
		return true
	}

	return p.After(other)
}

// RolledOver checks if this PTS just rollover compared to the other PTS
func (p PTS) RolledOver(other PTS) bool {
	if p < LowerPtsRolloverThreshold && other > UpperPtsRolloverThreshold {
		return true
	}
	return false
}

// DurationFrom returns the difference between the two pts times. This number is always positive.
func (p PTS) DurationFrom(from PTS) uint64 {
	switch {
	case p.RolledOver(from):
		return uint64((PTS_MAX + 1 - from) + p)
	case from.RolledOver(p):
		return uint64((PTS_MAX + 1 - p) + from)
	case p < from:
		return uint64(from - p)
	default:
		return uint64(p - from)
	}
}

// Add adds the two PTS times together and returns a new PTS
func (p PTS) Add(x PTS) PTS {
	result := p + x
	if result > PTS_MAX {
		result = result - PTS_MAX
	}
	return PTS(result)
}
