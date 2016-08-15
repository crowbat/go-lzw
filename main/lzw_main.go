package main

import (
    "github.com/crowbat/lzw"
    "flag"
    "log"
)

func main() {
    var infile = flag.String("i", "in.txt", "input file")
    var outfile = flag.String("o", "out.lzw", "output file")
    var direction = flag.Int("t", 1, "1 for encode, 2 for decode")
    flag.Parse()
    if *direction == 1 {
        lzw.Encode(*infile, *outfile)
    } else if *direction == 2 {
        lzw.Decode(*infile, *outfile)
    } else {
        log.Fatal("invalid argument")
    }
}