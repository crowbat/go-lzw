# go-lzw

A fun project to try and implement the lzw compression algorithm using Google Go.  Currently I have an encoder that can process files and generate a "compressed" file, and the decoder is able to take this file and recover the original file.  Unfortunately the size of the "compressed" file is actually usually larger than that of the original file so I'm still trying to figure out why this is.

To run the encoder, cd into main/ and run  
`./lzw_main -i <input file> -o <encoded output> -t 1`  

To run the decoder, run  
`./lzw_main -i <encoded file> -o <output file> -t 2`

Next steps:
1. Figure out how to produce smaller compressed files  
2. Allow user to choose size of encoding/decoding table (currently hardcoded at 12 bits).