package util

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// AuthCode Comsenz discuz classical symmetric encryption/decryption function, example see authcode_test.go#TestAuthCode()
// AuthCode 康盛Discuz! 经典对称加解密函数,用法参考 authcode_test.go#TestAuthCode()
func AuthCode(str string, operation string, key string, expiry int64) string {
	// 动态密匙长度，相同的明文会生成不同密文就是依靠动态密匙
	// 加入随机密钥，可以令密文无任何规律，即便是原文和密钥完全相同，加密结果也会每次不同，增大破解难度。
	// 取值越大，密文变动规律越大，密文变化 = 16 的 ckeyLength 次方
	// 当此值为 0 时，则不产生随机密钥
	ckeyLength := 4

	// 密匙
	if key == "" {
		key = "#@!.5ebcQJx2Lz6GmcsqNiNHW.!@#"
	}
	key = Md5(key)

	// 密匙a会参与加解密
	keya := Md5(key[:16])
	// 密匙b会用来做数据完整性验证
	keyb := Md5(key[16:])
	// 密匙c用于变化生成的密文
	keyc := ""
	if ckeyLength != 0 {
		if operation == "DECODE" {
			keyc = str[:ckeyLength]
		} else {
			sTime := Md5(time.Now().String())
			sLen := 32 - ckeyLength
			keyc = sTime[sLen:]
		}
	}
	// 参与运算的密匙
	cryptkey := fmt.Sprintf("%s%s", keya, Md5(keya+keyc))
	ckkeyLength := len(cryptkey)
	// 明文，前10位用来保存时间戳，解密时验证数据有效性，10到26位用来保存$keyb(密匙b)，解密时会通过这个密匙验证数据完整性
	// 如果是解码的话，会从第$ckeyLength位开始，因为密文前$ckeyLength位保存 动态密匙，以保证解密正确
	if operation == "DECODE" {
		str = strings.Replace(str, "*", "+", -1)
		str = strings.Replace(str, "_", "/", -1)
		strByte, err := base64.StdEncoding.DecodeString(str[ckeyLength:])
		if err != nil {
			return ""
		}
		str = string(strByte)
	} else {
		if expiry != 0 {
			expiry = expiry + time.Now().Unix()
		}
		tmpMd5 := Md5(str + keyb)
		str = fmt.Sprintf("%010d%s%s", expiry, tmpMd5[:16], str)
	}
	stringLength := len(str)
	resdata := make([]byte, 0, stringLength)
	var rndkey, box [256]int
	// 产生密匙簿
	j := 0
	a := 0
	i := 0
	tmp := 0
	for i = 0; i < 256; i++ {
		rndkey[i] = int(cryptkey[i%ckkeyLength])
		box[i] = i
	}
	// 用固定的算法，打乱密匙簿，增加随机性，好像很复杂，实际上并不会增加密文的强度
	for i = 0; i < 256; i++ {
		j = (j + box[i] + rndkey[i]) % 256
		tmp = box[i]
		box[i] = box[j]
		box[j] = tmp
	}
	// 核心加解密部分
	a = 0
	j = 0
	tmp = 0
	for i = 0; i < stringLength; i++ {
		a = ((a + 1) % 256)
		j = ((j + box[a]) % 256)
		tmp = box[a]
		box[a] = box[j]
		box[j] = tmp
		// 从密匙簿得出密匙进行异或，再转成字符
		resdata = append(resdata, byte(int(str[i])^box[(box[a]+box[j])%256]))
	}
	result := string(resdata)
	if operation == "DECODE" {
		// substr($result, 0, 10) == 0 验证数据有效性
		// substr($result, 0, 10) - time() > 0 验证数据有效性
		// substr($result, 10, 16) == substr(md5(substr($result, 26).$keyb), 0, 16) 验证数据完整性
		// 验证数据有效性，请看未加密明文的格式
		frontTen, _ := strconv.ParseInt(result[:10], 10, 0)
		if (frontTen == 0 || frontTen-time.Now().Unix() > 0) && result[10:26] == Md5(result[26:] + keyb)[:16] {
			return result[26:]
		}
		return ""
	}
	// 把动态密匙保存在密文里，这也是为什么同样的明文，生产不同密文后能解密的原因
	// 因为加密后的密文可能是一些特殊字符，复制过程可能会丢失，所以用base64编码
	result = keyc + base64.StdEncoding.EncodeToString([]byte(result))
	result = strings.Replace(result, "+", "*", -1)
	result = strings.Replace(result, "/", "_", -1)
	return result
}
