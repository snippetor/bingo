module github.com/snippetor/bingo

require (
	github.com/RoaringBitmap/roaring v0.4.16 // indirect
	github.com/anacrolix/tagflag v0.0.0-20180803105420-3a8ff5428f76 // indirect
	github.com/dustin/go-humanize v0.0.0-20180713052910-9f541cc9db5d // indirect
	github.com/smallnest/rpcx v0.0.0-20180827063508-e9348955605f // indirect
)

replace (
	golang.org/x/crypto v0.0.0-20180820150726-614d502a4dac => github.com/golang/crypto v0.0.0-20180820150726-614d502a4dac
	golang.org/x/net v0.0.0-20180821023952-922f4815f713 => github.com/golang/net v0.0.0-20180821023952-922f4815f713
	golang.org/x/text v0.3.0 => github.com/golang/text v0.3.0
)
