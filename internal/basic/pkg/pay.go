package pkg

import (
	"fmt"
	"github.com/smartwalle/alipay/v3"
	"github.com/zeromicro/go-zero/core/logx"
)

func Pay(Num, OutTradeNo, TotalAmount string) (payURL string) {
	var err error
	var appId = "9021000144698155"
	var privateKey = "MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDD0unrVKzeBO8A0Gswj7yO5+h9GJqBe7LlWtHLmFro4uTLiGzypTMHEW3l3z+TyMVyZVUTt8lhxOIl+ghavCRZlbY+ONXfdxzXJFMPqOHX+nZ6eXPQm5H5thJIT3Xh34TEpDhOPcIpmLEJN03iRhVQeMVNSkyUSHgejI3/3PADe5IAFzhNUCfbwedlLv/IPoHQZXlF2xCveapFoxCoV89eT304gyxC96Xlbx81Xk1nr2r6rstW/tXjnZ0+Ckl8mNTe036CrXGlbzywazD4UE8j3wuExOy9JsCpo83+ME7LVTqF3H0ytnZhZjHS5cPTWFoiS2zp8oNZ3MIyfdpi8UuzAgMBAAECggEBAIe+jzT4z5RgKyWPVJ6nJPiTPhBpm8EidJKU8FWH2Y0Sq7ODcLKLXeIKbPoqVbImPOjd4an3fvvtNS9KMbxkK3aGV3yufWOY+D8TCRkT4uqKztJ7mpMDJJ/LXMUPgBIBldGSXQ2vtgaLuD1BPxYZqvDLp0m6tXcc/Vd/63dwOljOeimoz65V4gvISr71VdmYb1cNsuORcRahUlqHWxtxj1FLzU9U7R8xZuTLXmaDc1P525C9GHPjrzOX73aXCWJWWnO2xGlWHI2knsdX+koHIHKepS3PirbxH0Es83b8ZlWzs9GVXLUD6t4VJ0BPHLzHs2XPRBUXnBR1NpgyqpwFHVECgYEA9yi4LGMRr/d9CLIu2b6Voc81cfJ6R9s9RMQ5A/lFTGC4NwaoojMG3c6yyLchDrlRFKQkmaK5iLuxqNMDCtmSG7vSvPst7/55jIGYG9Jmh4VW33ToTEsrCJN37Vrqg4uRvsgIQKvIbR8p0c4H2r5k2Zntx8YY96pvLgorrksTyucCgYEAytQb1eE4GTxxu6/yKGkU9g1DZvfAGmmsyZt2CB1vvqwRbV/bWSNeAsKAaHg16K6taSfd3iotdYHFjTcGenwIgrxTWo3ObnXqvhoG54+BhHwbXzEX3uD41TTy51JAtCUTadwjkJ9QlRv2ZyIqwVr7JGcJUqndNoz31Y+mOFFBC1UCgYACHpVFvCvAqIyn8G86asn5sz8wFPY7e4PQ/SXIBPE1MHcj8aisi2d5q3YZBokACVLKrIWr05tnssRZQEX8Z9U666do+3ZvYm1EaTAWvP0oGFqgW+5KCTL7Rdh3bpooOqArVKfNdiun0+aV6ABlPdC7lPhXCDnaldmSOYAaZIZs7QKBgF5OkUK5HWRefmNOQ8IWWfCt6hEOUPv29qgm6JKNXU/CobfBQjQIBcYyuZHZkvdFgMvMBZUu90QTus8WLqT01uAywG7yUHF70lHhuCQizY3URsXUBc1TvV8k52w3Cm64bnZiLQcpjEZIYiFB+a89plgesG8HHBwpH3Lk/9xfq2ahAoGALA/DBjdQ8Kz5vKrOaz1xdtjqgh4sXCql0XjEo5RazC1LGxdKZtZcjgMxlsEgaWf0Z4J9dD6JRrrZLMI5y3jMO33tZ9Q8avBEK29uNz8mjxEekvaxYr95bOX8KZ6TbtgAOPtoMp19gC+GCe2/08kgXnep6go4+zKMMColvnOAH0w=" // 必须，上一步中使用 RSA签名验签工具 生成的私钥
	client, err := alipay.New(appId, privateKey, false)
	if err != nil {
		logx.Info(err)
		return
	}

	var p = alipay.TradeWapPay{}
	p.NotifyURL = ""
	p.ReturnURL = "https://www.baidu.com"
	p.Subject = Num
	p.OutTradeNo = OutTradeNo
	p.TotalAmount = TotalAmount
	p.ProductCode = "QUICK_WAP_WAY"

	url, err := client.TradeWapPay(p)
	if err != nil {
		logx.Info(err)
	}

	// 这个 payURL 即是用于打开支付宝支付页面的 URL，可将输出的内容复制，到浏览器中访问该 URL 即可打开支付页面。
	payURL = url.String()
	fmt.Println(payURL)
	return payURL
}
