package main

import (
    "bufio"
    "encoding/hex"
    "flag"
    "fmt"
    "image"
    "image/color"
    "image/png"
    "log"
    "os"
    "strings"
)

func main() {
    inPath := flag.String("input", "", "Input PNG file (1-bit or any grayscale)")
    outPath := flag.String("output", "", "Output ZPL file (default stdout)")
    invert := flag.Bool("invert", false, "Invert black/white")
    flag.Parse()

    if *inPath == "" {
        log.Fatalf("usage: png2zpl -input file.png [-output out.zpl] [-invert]")
    }

    f, err := os.Open(*inPath)
    if err != nil {
        log.Fatalf("open input: %v", err)
    }
    defer f.Close()

    img, err := png.Decode(f)
    if err != nil {
        log.Fatalf("decode png: %v", err)
    }

    w := img.Bounds().Dx()
    h := img.Bounds().Dy()

    bytesPerRow := (w + 7) / 8
    totalBytes := bytesPerRow * h
    buf := make([]byte, totalBytes)

    switch im := img.(type) {
    case *image.Gray:
        fillFromGray(im, buf, w, h, bytesPerRow, *invert)
    case *image.Paletted:
        fillFromPaletted(im, buf, w, h, bytesPerRow, *invert)
    default:
        fillFromGeneric(img, buf, w, h, bytesPerRow, *invert)
    }

    hexData := make([]byte, hex.EncodedLen(len(buf)))
    hex.Encode(hexData, buf)

    compressed := compressZPL(string(hexData))

    zpl := fmt.Sprintf("^XA\n^FO0,0\n^GFA,%d,%d,%d,%s\n^XZ\n", totalBytes, totalBytes, bytesPerRow, compressed)

    if *outPath == "" {
        w := bufio.NewWriter(os.Stdout)
        _, _ = w.WriteString(zpl)
        w.Flush()
    } else {
        of, err := os.Create(*outPath)
        if err != nil {
            log.Fatalf("create output: %v", err)
        }
        defer of.Close()
        w := bufio.NewWriter(of)
        _, _ = w.WriteString(zpl)
        w.Flush()
    }
}

func isBlackColor(c color.Color) bool {
    r, g, b, _ := c.RGBA()
    lum := (299*int(r) + 587*int(g) + 114*int(b)) / 1000 / 257
    return lum < 128
}

func fillFromGray(img *image.Gray, buf []byte, w, h, bytesPerRow int, invert bool) {
    stride := img.Stride
    for y := 0; y < h; y++ {
        offset := y * bytesPerRow
        rowBase := y * stride
        bit := 7
        bidx := offset
        var cur byte = 0
        for x := 0; x < w; x++ {
            v := img.Pix[rowBase+x]
            black := v < 128
            if invert {
                black = !black
            }
            if black {
                cur |= 1 << uint(bit)
            }
            bit--
            if bit < 0 {
                buf[bidx] = cur
                bidx++
                bit = 7
                cur = 0
            }
        }
        if bit != 7 {
            buf[bidx] = cur
        }
    }
}

func fillFromPaletted(img *image.Paletted, buf []byte, w, h, bytesPerRow int, invert bool) {
    stride := img.Stride
    pal := img.Palette
    palIsBlack := make([]bool, len(pal))
    for i, c := range pal {
        palIsBlack[i] = isBlackColor(c)
    }

    for y := 0; y < h; y++ {
        offset := y * bytesPerRow
        rowBase := y * stride
        bit := 7
        bidx := offset
        var cur byte = 0
        for x := 0; x < w; x++ {
            idx := int(img.Pix[rowBase+x])
            black := palIsBlack[idx]
            if invert {
                black = !black
            }
            if black {
                cur |= 1 << uint(bit)
            }
            bit--
            if bit < 0 {
                buf[bidx] = cur
                bidx++
                bit = 7
                cur = 0
            }
        }
        if bit != 7 {
            buf[bidx] = cur
        }
    }
}

func fillFromGeneric(img image.Image, buf []byte, w, h, bytesPerRow int, invert bool) {
    for y := 0; y < h; y++ {
        offset := y * bytesPerRow
        bit := 7
        bidx := offset
        var cur byte = 0
        for x := 0; x < w; x++ {
            c := img.At(x, y)
            black := isBlackColor(c)
            if invert {
                black = !black
            }
            if black {
                cur |= 1 << uint(bit)
            }
            bit--
            if bit < 0 {
                buf[bidx] = cur
                bidx++
                bit = 7
                cur = 0
            }
        }
        if bit != 7 {
            buf[bidx] = cur
        }
    }
}

func compressZPL(data string) string {
    if len(data) == 0 {
        return data
    }
    var out strings.Builder
    n := len(data)
    i := 0
    for i < n {
        ch := data[i]
        j := i + 1
        for j < n && data[j] == ch {
            j++
        }
        run := j - i
        if run >= 3 {
            repeat := run
            var count strings.Builder
            // 400-blocks -> 'z'
            if repeat > 400 {
                zcount := repeat / 400
                for k := 0; k < zcount; k++ {
                    count.WriteByte('z')
                }
                repeat = repeat % 400
            }
            // 20-blocks -> 'f' + floor(repeat/20)
            if repeat > 19 {
                times := repeat / 20
                count.WriteByte(byte('f' + byte(times)))
                repeat = repeat % 20
            }
            // remainder 1..19 -> 'F' + remainder
            if repeat > 0 {
                count.WriteByte(byte('F' + byte(repeat)))
            }

            out.WriteString(count.String())
            out.WriteByte(ch)
        } else {
            // copy raw if run < 3 (same behavior as your PHP regex)
            out.WriteString(data[i:j])
        }
        i = j
    }
    return out.String()
}

