module github.com/snippetor/bingo

require github.com/smallnest/rpcx v0.0.0-20180823121341-26b8cb57c531 // indirect

replace (
	golang.org/x/crypto v0.0.0-20180820150726-614d502a4dac => github.com/golang/crypto v0.0.0-20180820150726-614d502a4dac
	golang.org/x/net v0.0.0-20180821023952-922f4815f713 => github.com/golang/net v0.0.0-20180821023952-922f4815f713
	golang.org/x/text v0.3.0 => github.com/golang/text v0.3.0
)
