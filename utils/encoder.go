package utils

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wumansgy/goEncrypt/des"
	"strconv"
)

func Packet(b []byte) ([]byte, error) {
	headerBytes := []byte("GDMP") //添加报文头:GDMP
	// 消息编码base64

	eb, err := des.DesCbcEncrypt(b, []byte("12345678"), []byte("00000000"))
	if err != nil {
		return nil, err
	}
	esb := base64.StdEncoding.EncodeToString(eb) //将加密结果进行base64编码,形成可见字符
	dataLenBytes := []byte(FormatNumberString(strconv.Itoa(len(esb))))
	newBytes := append(headerBytes, dataLenBytes...)
	newBytes = append(newBytes, []byte(esb)...)
	return newBytes, nil
}

// UPacket 对加密的des报文进行解密
func UPacket(b []byte) ([]byte, error) {
	if len(b) < 8 {
		return nil, errors.New("message not intact") //消息不完整
	}
	dataLenBytes := b[4:8] //取出报文头部:GDMP
	dateLen, err := strconv.Atoi(string(dataLenBytes))
	if err != nil {
		return nil, err
	}
	db := b[8:]
	if len(db) != dateLen {
		return nil, errors.New("消息长度错误")
	}
	decode, err := base64.StdEncoding.DecodeString(string(db)) //将消息进行base64解码
	if err != nil {
		return nil, err
	}
	rs, err := des.DesCbcDecrypt(decode, []byte("12345678"), []byte("00000000")) //将消息进行des解密,得到真实消息
	if err != nil {
		return nil, err
	}
	return rs, nil
}
func FormatNumberString(input string) string {
	if len(input) < 4 {
		// 在首部添加 '0'，直到字符串长度达到 3
		for len(input) < 4 {
			input = "0" + input
		}
	}
	return input
}

func ExtractCmdValue(jsonStr string) (string, error) {
	// 创建一个空的map来存储解析后的JSON数据
	var data map[string]interface{}
	// 解析JSON字符串
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		return "", err
	}
	// 获取Command字段的值
	cmdValue, ok := data["Cmd"].(string)
	if !ok {
		//fmt.Printf("the command value is not string : %s", err.Error())
		return "", fmt.Errorf("cmd的值不是字符串类型")
	}
	return cmdValue, nil
}
