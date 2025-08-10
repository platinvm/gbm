package tags

import (
    "encoding/binary"
    "strings"
)

func ParseOrder(tag string) (binary.ByteOrder, bool) {
    for _, p := range strings.Split(tag, ",") {
        p = strings.TrimSpace(p)
        if p == "LE" {
            return binary.LittleEndian, true
        }
        if p == "BE" {
            return binary.BigEndian, true
        }
    }
    return nil, false
}

func ParseLen(tag string) (string, bool) {
    for _, p := range strings.Split(tag, ",") {
        p = strings.TrimSpace(p)
        if strings.HasPrefix(p, "len=") {
            return strings.TrimPrefix(p, "len="), true
        }
    }
    return "", false
}

func ParseEnc(tag string) (utf8Enc bool) {
    for _, p := range strings.Split(tag, ",") {
        p = strings.TrimSpace(p)
        if strings.HasPrefix(p, "enc=") {
            val := strings.TrimPrefix(p, "enc=")
            if strings.EqualFold(val, "utf-8") || strings.EqualFold(val, "utf8") {
                return true
            }
            return false
        }
    }
    return true // default UTF-8
}
