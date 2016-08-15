package lzw

import (
    "bufio"
    "fmt"
    "os"
    "log"
    "io"
    "strconv"
    "github.com/crowbat/bits"
)

func Init_reverse_table() (reverse_table map[int]string) {
    r := make(map[int]string)
    for byte_counter:=0; byte_counter <= 255; byte_counter++ {           //leave index 0 open
        if byte_counter < 16 {
            r[byte_counter+1] = "0" + strconv.FormatInt(int64(byte_counter),16)
        } else {
            r[byte_counter+1] = strconv.FormatInt(int64(byte_counter),16)
        }
    }
    return r
}

func Init_forward_table() (lzw_table map[string] int) {
    l := make(map[string]int)
    for byte_counter:=0; byte_counter <= 255; byte_counter++ {           //leave index 0 open
        if byte_counter < 16 {
            l["0" + strconv.FormatInt(int64(byte_counter),16)] = byte_counter+1
        } else {
            l[strconv.FormatInt(int64(byte_counter),16)] = byte_counter+1
        }
    }
    return l
}

func Encode(infilename string, outfilename string) {
    lzw_table := Init_forward_table()
    curr_sequence := []byte{}
    table_index := 257              //change
    matches := make([]int, 1000)
    can_write := true

    infile, err := os.Open(infilename)
    if err != nil {
        log.Fatal(err)
    }
    defer infile.Close()

    outfile, err := os.Create(outfilename)
    if err != nil {
        log.Fatal(err)
    }
    defer outfile.Close()

    inreader := bufio.NewReader(infile)
    outwriter := bufio.NewWriter(outfile)
    outBitWriter := bits.BitWriter{BufioWriter: outwriter}
    first_byte, err := inreader.ReadByte()
    if err != nil {
        log.Fatal("error reading first byte")
    }
    curr_sequence = append(curr_sequence, []byte{first_byte}...)
    seq_string := bytearray_to_string(curr_sequence)
    counter := 1
    for {
        next_byte, err := inreader.ReadByte()
        if err == io.EOF {
            outBitWriter.WriteUint(uint(lzw_table[bytearray_to_string(curr_sequence)]), 12)
            fmt.Printf("reached EOF\n")
            break
        }
        counter++
        if err != nil {
            log.Fatal("err")
            return
        }
        prev_seq_string := seq_string
        curr_sequence = append(curr_sequence, []byte{next_byte}...)
        seq_string = bytearray_to_string(curr_sequence)
        _, ok := lzw_table[seq_string]
        if !ok {
            if can_write {
                lzw_table[seq_string] = table_index
                table_index++
            }
            outBitWriter.WriteUint(uint(lzw_table[prev_seq_string]), 12)            //change
            curr_sequence = []byte{next_byte}
            seq_string = bytearray_to_string(curr_sequence)
            matches = append(matches[1:len(matches)], 0)
        } else {
            matches = append(matches[1:len(matches)], 1)
        }
        if table_index >= 4096 {        //change
            can_write = false
            sum := 0
            for _,n := range(matches) {
                sum += n
            }
            if float64(sum)/1000 < .1 {
                can_write = true
                lzw_table = Init_forward_table()
                table_index = 257
                outBitWriter.WriteUint(0,12)
            }
        }
    }
    outBitWriter.FinishByte()
    fmt.Printf("read %v bytes\n", counter)
}

func Decode(infilename string, outfilename string) {
    reverse_table := Init_reverse_table()
    curr_string := ""
    table_index := 257                  //change
    can_write := true

    infile, err := os.Open(infilename)
    if err != nil {
        log.Fatal(err)
    }
    defer infile.Close()

    outfile, err := os.Create(outfilename)
    if err != nil {
        log.Fatal(err)
    }
    defer outfile.Close()

    inreader := bufio.NewReader(infile)
    inBitReader := bits.BitReader{BufioReader: inreader}
    outwriter := bufio.NewWriter(outfile)

    first_index,_ := inBitReader.ReadBits(12)
    curr_string = reverse_table[first_index]
    curr_bytes := string_to_bytearray(curr_string)
    for i:=0;i<len(curr_bytes);i++ {
        outwriter.WriteByte(curr_bytes[i])
        outwriter.Flush()
    }
    for {
        next_index, err := inBitReader.ReadBits(12)
        if err != nil {
            return
        }
        if next_index == 0 {
            can_write = true
            reverse_table = Init_reverse_table()
            table_index = 257
            first_index,_ := inBitReader.ReadBits(12)
            curr_string = reverse_table[first_index]
            curr_bytes := string_to_bytearray(curr_string)
            for i:=0;i<len(curr_bytes);i++ {
                outwriter.WriteByte(curr_bytes[i])
                outwriter.Flush()
            }
        } else {
            next_string, ok := reverse_table[next_index]
            if ok {
                next_bytes := string_to_bytearray(next_string)
                for i:=0;i<len(next_bytes);i++ {
                    outwriter.WriteByte(next_bytes[i])
                    outwriter.Flush()
                }
                curr_string += next_string[0:2]
                if can_write {
                    reverse_table[table_index] = curr_string
                    table_index++
                }
                curr_string = next_string
            } else {
                curr_string += curr_string[0:2]
                if can_write {
                    reverse_table[next_index] = curr_string
                    table_index = next_index + 1
                }
                curr_bytes = string_to_bytearray(curr_string)
                for i:=0;i<len(curr_bytes);i++ {
                    outwriter.WriteByte(curr_bytes[i])
                    outwriter.Flush()
                }
            }
        }
    }
}

func bytearray_to_string(barray []byte) (output string) {
    output = ""
    for i:=0;i<len(barray);i++ {
        if int64(barray[i])<16 {
            output += "0"
            output += strconv.FormatUint(uint64(barray[i]),16)
        } else {
            output += strconv.FormatUint(uint64(barray[i]),16)
        }
    }
    return output
}

func string_to_bytearray(s string) (output []byte) {
    output = make([]byte, len(s)/2)
    for i:=0;i<len(s)-1;i=i+2 {
        nextint,_ := strconv.ParseUint(s[i:i+2],16,8)
        output[i/2] = byte(nextint)
    }
    return output
}