package helper

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"math/big"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const tokenCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const tokenSize = 22

func NewV4() []byte {
	bytes := make([]byte, 16)
	randSeed := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 16; i += 1 {
		bytes[i] = byte(randSeed.Uint32() % 256)
	}
	return bytes
}

func ByteToString(u []byte) string {
	buf := make([]byte, 36)

	hex.Encode(buf[0:8], u[0:4])
	buf[8] = '-'
	hex.Encode(buf[9:13], u[4:6])
	buf[13] = '-'
	hex.Encode(buf[14:18], u[6:8])
	buf[18] = '-'
	hex.Encode(buf[19:23], u[8:10])
	buf[23] = '-'
	hex.Encode(buf[24:], u[10:])

	return string(buf)
}

func GenerateToken(productID []byte, operationID string) string {
	if len(operationID) > 6 {
		operationID = operationID[:6]
	}
	var buffer bytes.Buffer
	buffer.Write(productID[:])
	buffer.WriteString(operationID)
	obj := NewV4()
	uuid := ByteToString(obj)
	i := new(big.Int)
	i.SetString(strings.Replace(uuid, "-", "", 4), 16)
	base := big.NewInt(int64(len(tokenCharset)))
	mod := new(big.Int)
	for i.String() != "0" {
		i.DivMod(i, base, mod)
		buffer.WriteByte(tokenCharset[mod.Int64()])
	}

	var res string
	size := tokenSize + len(productID) + len(operationID)
	n := size - buffer.Len()
	switch {
	case n < 0:
		res = buffer.String()[:size]
	case n > 0:
		buffer.Write(bytes.Repeat([]byte{tokenCharset[0]}, n))
		res = buffer.String()
	case n == 0:
		res = buffer.String()
	}

	return res
}

func SetCookie(token string, salt string) string {
	randomKey := rand.Uint64()
	timeStamp := time.Now().Unix() + YearTime
	timeStampStr := strconv.FormatInt(timeStamp, 10)
	randomKeyStr := strconv.FormatUint(randomKey, 10)
	key := randomKeyStr + timeStampStr + token + salt
	cipherStr := Md5(key)

	return fmt.Sprintf("%s-%s-%s-%s", randomKeyStr, timeStampStr, token, cipherStr)
}

func Md5(key string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(key))
	cipherStr := hex.EncodeToString(md5Ctx.Sum(nil))
	return cipherStr
}

func FormatLogPrint(logType string, format string, a ...interface{}) {
	printLog := logType + " " + time.Now().Format("2006/01/02 - 15:04:05") + " : " + format + "\r\n"
	log.Printf(printLog, a)
	_, _ = fmt.Fprintf(gin.DefaultWriter, printLog, a)
}
