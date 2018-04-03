package ParseID3

import (
    "testing"
    "os"
    "fmt"
)

func TestMain(t *testing.T) {
    f, _ := os.OpenFile("C:\\Users\\user\\Desktop\\s\\Aaum0ELD1P2fL5m7sDrROI0TJoEF1p1a.mp3", os.O_RDONLY, 0777)
    a, _ := os.Create("C:\\Users\\user\\Desktop\\a.jpg")
    defer func() {
        f.Close()
    } ()
    id32, _ := ParseID3v2(f)
    for k, v := range FrameIDMap {
        if COVER_PHOTO_ID != k {
            c, ce := id32.GetFrameContent(k)
            if ce != nil {
                fmt.Printf("%s - %-20s %-20s\n", k, v, "")
            } else {
                fmt.Printf("%s - %-20s %-20s\n", k, v, string(c))
            }
        } else {
            binaryBytes, typeBytes := id32.GetCover()
            fmt.Printf("%s - %-20s %-20s\n", k, v, string(typeBytes))
            a.Write(binaryBytes)
        }
    }
}
