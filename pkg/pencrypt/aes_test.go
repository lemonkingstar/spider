package pencrypt

import (
	"fmt"
	"testing"
)

/*
AES（Advanced Encryption Standard）高级加密标准，是流行的对称加密算法。

AES 有 5 种加密模式，分别是：
- 电子密码本模式（ECB，Electronic Code Book）
- 加密块链模式（CBC，Cipher Block Chaining），如果明文长度不是分组长度 16 字节的整数倍需要进行填充
- 计数模式（CTR，Counter）
- 密码反馈模式（CFB，Cipher FeedBack）
- 输出反馈模式（OFB，Output FeedBack）

AES 是对称分组加密算法，每组长度为 128bits，即 16 字节。

AES 秘钥的长度只能是16、24 或 32 字节，分别对应三种加密模式 AES-128、AES-192 和 AES-256，三者的区别是加密轮数不同。

| AES		分组长度(字节)	密钥长度(字节)	加密轮数
| AES-128	16				16				10
| AES-192	16				24				12
| AES-256	16				32				14

 */

const key = "2a08f271128e0f40b63e75e2ad4db451"

func TestAesEncrypt(t *testing.T) {
	text := "i love china! i love chinese food!"
	aes := NewAesEncryptor([]byte(text), true)
	aes.SetKey([]byte(key))
	b, err := aes.Encrypt()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(b)
	fmt.Println(string(b))

	// to base64
	str, err := aes.Encrypt2Base64()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(str)

	// decrypt
	aes = NewAesEncryptor(b, false)
	aes.SetKey([]byte(key))
	b2, err := aes.Decrypt()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b2))

	// from base64
	aes = NewAesEncryptor([]byte(str), false)
	aes.SetKey([]byte(key))
	b3, err := aes.Decrypt5Base64()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b3))
}
