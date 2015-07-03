package memory

import "testing"
import "fmt"

// Helpers
func checkEqualBytes(val, expectedVal byte, t *testing.T){
    if val != expectedVal {
        t.Error("Expected 0x", fmt.Sprintf("%x", expectedVal), ", got 0x", fmt.Sprintf("%x", val))
    }
}

// Actual tests
func TestBiosAccess(t *testing.T){
    // At startup, the BIOS data should be mapped at 0x0000 - 0x00FF
    // (It shadows ROM bank 0, and is therefore removed later)
    mmu := MakeMmu()

    firstBiosByte := mmu.ReadByte(0x00)
    checkEqualBytes(firstBiosByte, 0x31, t)

    fifthBiosByte := mmu.ReadByte(0x04)
    checkEqualBytes(fifthBiosByte, 0x21, t)

    lastBiosByte := mmu.ReadByte(0xff)
    checkEqualBytes(lastBiosByte, 0x50, t)
}
