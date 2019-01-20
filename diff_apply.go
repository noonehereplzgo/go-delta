// -----------------------------------------------------------------------------
// (c) balarabe@protonmail.com                                      License: MIT
// :v: 2019-01-20 06:32:50 F48889                       go-delta/[diff_apply.go]
// -----------------------------------------------------------------------------

package delta

import (
	"bytes"
	"fmt"
)

// Apply __
func (ob *Diff) Apply(source []byte) ([]byte, error) {
	if DebugTiming {
		tmr.Start("Diff.Apply")
		defer tmr.Stop("Diff.Apply")
	}
	if len(source) != ob.sourceSize {
		return nil, mod.Error(fmt.Sprintf(
			"Size of source [%d] does not match expected [%d]",
			len(source), ob.sourceSize))
	}
	if !bytes.Equal(makeHash(source), ob.sourceHash) {
		return nil, mod.Error("Diff. can not be applied to source")
	}
	var buf = bytes.NewBuffer(make([]byte, 0, ob.targetSize))
	for i, pt := range ob.parts {
		var data []byte
		switch {
		case pt.sourceLoc == -1:
			data = pt.data
		case pt.sourceLoc < 0 || pt.sourceLoc >= ob.sourceSize:
			return nil, mod.Error("part", i, "sourceLoc:", pt.sourceLoc,
				"out of range 0 -", ob.sourceSize-1)
		case pt.sourceLoc+pt.size > ob.sourceSize:
			return nil, mod.Error("part", i, "sourceLoc:", pt.sourceLoc,
				"+ size:", pt.size, "extends beyond", ob.sourceSize)
		default:
			data = source[pt.sourceLoc : pt.sourceLoc+pt.size]
		}
		var n, err = buf.Write(data)
		if err != nil {
			return nil, mod.Error(err)
		}
		if n != pt.size {
			return nil, mod.Error("Wrote", n, "bytes instead of", pt.size)
		}
	}
	var ret = buf.Bytes()
	if !bytes.Equal(makeHash(ret), ob.targetHash) {
		return nil, mod.Error("Result does not match target hash.")
	}
	return buf.Bytes(), nil
} //                                                                       Apply

//end
