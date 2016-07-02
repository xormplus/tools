package main

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/oskca/sciter"
	"github.com/oskca/sciter/window"
)

var tool *Tool
var w *window.Window
var ret *sciter.Value

//sciter.SW_RESIZEABLE|sciter.SW_CONTROLS|sciter.SW_TITLEBAR|sciter.SW_MAIN|sciter.SW_ENABLE_DEBUG
func main() {

	ret = sciter.NewValue()
	var err error
	tool = new(Tool)

	w, err = window.New(sciter.SW_MAIN|sciter.SW_GLASSY, &sciter.Rect{0, 0, 0, 0})

	if err != nil {
		log.Fatal("Create Window Error: ", err)
	}
	w.SetTitle("xorm tools")

	w.LoadFile("tools.html")
	//	err = w.LoadHtml(indexhtml, "/")
	if err != nil {
		log.Fatal("Create Window Error: ", err)
	}

	setEventHandler(w)

	w.Show()
	w.Run()
}

func setEventHandler(w *window.Window) {
	var err error
	w.DefineFunction("getDir",
		func(args ...*sciter.Value) *sciter.Value {
			var bits int

			returnCmd := sciter.NewValue()
			returnCmd.Set("cmd", sciter.NewValue("done"))

			tool.inputDir = args[0].Get("path1").String()
			tool.outputDir = args[0].Get("path2").String()

			switch args[0].Get("radioGroup").Int() {
			case 1:
				aes := &AesEncrypt{PubKey: args[0].Get("passwd").String()}
				DoAes(aes)
			case 2:
				des := &DesEncrypt{PubKey: args[0].Get("passwd").String()}
				DoDes(des)
			case 3:
				tripleDes := &TripleDesEncrypt{PubKey: args[0].Get("passwd").String()}
				DoTripleDes(tripleDes)
			case 4:

				rsa := &RsaEncrypt{}
				if args[0].Get("rsaMode").Int() == 1 {
					rsa.EncryptMode = MODE_PUBKEY_ENCRYPT
					rsa.DecryptMode = MODE_PRIKEY_DECRYPT
				} else {
					rsa.EncryptMode = MODE_PRIKEY_ENCRYPT
					rsa.DecryptMode = MODE_PUBKEY_DECRYPT
				}
				bits, err = strconv.Atoi(args[0].Get("bitwd").String())
				if err != nil {
					AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [错误日志]  内容：[" + err.Error() + "]</div>")
					return ret
				}
				rsa.Bits = bits

				DoRsa(rsa)
			case 5:
				rsa := &RsaEncrypt{}
				if args[0].Get("rsaMode").Int() == 1 {
					rsa.EncryptMode = MODE_PUBKEY_ENCRYPT
					rsa.DecryptMode = MODE_PRIKEY_DECRYPT
				} else {
					rsa.EncryptMode = MODE_PRIKEY_ENCRYPT
					rsa.DecryptMode = MODE_PUBKEY_DECRYPT
				}
				bits, err = strconv.Atoi(args[0].Get("bitwd").String())
				if err != nil {
					PutKey("pubkey", "<div style=\"color:#FF0000\">"+NowTime()+"  [错误日志]  内容：["+err.Error()+"]</div>")
					return ret
				}
				rsa.Bits = bits
				DoGenRsaPerm(rsa)

			case 6:
				rsa := &RsaEncrypt{}
				rsa.EncryptMode = MODE_PUBKEY_ENCRYPT
				rsa.DecryptMode = MODE_PRIKEY_DECRYPT
				DoRsaModePubkeyEncrypt(rsa)

			case 7:
				rsa := &RsaEncrypt{}
				rsa.EncryptMode = MODE_PRIKEY_ENCRYPT
				rsa.DecryptMode = MODE_PUBKEY_DECRYPT
				DoRsaModePrikeyEncrypt(rsa)
			}

			return returnCmd
		})
}

func DoGenRsaPerm(rsaEncrypt *RsaEncrypt) {
	root, err := w.GetRootElement()
	if err != nil {
		PutKey("pubkey", "<div style=\"color:#FF0000\">"+NowTime()+"  [错误日志]  内容：["+err.Error()+"]</div>")
		root.CallFunction("enable", ret)
		root.Update(false)
		return
	}

	ClearKey()

	err = rsaEncrypt.GenRsaKey(tool)
	if err != nil {
		PutKey("pubkey", "<div style=\"color:#FF0000\">"+NowTime()+"  [错误日志]  内容：["+err.Error()+"]</div>")
		root.CallFunction("enable", ret)
		root.Update(false)
		return
	}
	PutKey("pubkey", rsaEncrypt.PubKey)
	PutKey("prikey", rsaEncrypt.PriKey)

	root.CallFunction("enable", ret)
	root.Update(false)
	return
}

func DoRsaModePubkeyEncrypt(rsaEncrypt *RsaEncrypt) {
	root, err := w.GetRootElement()
	if err != nil {
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [错误日志]  内容：[" + err.Error() + "]</div>")
		root.CallFunction("enable", ret)
		root.Update(false)
		return
	}
	ClearMsg()

	resultElement, err := root.SelectById("inputPubkey")
	if err != nil {
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [错误日志]  内容：[" + err.Error() + "]</div>")
		root.CallFunction("enable", ret)
		root.Update(false)
		return
	}
	pubkey, err := resultElement.Text()
	if err != nil {
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [错误日志]  内容：[" + err.Error() + "]</div>")
		root.CallFunction("enable", ret)
		root.Update(false)
		return
	}

	//	log.Println(pubkey)
	pubkey = strings.Replace(pubkey, "-----BEGIN PUBLIC KEY-----", "", -1)
	pubkey = strings.Replace(pubkey, "-----END PUBLIC KEY-----", "", -1)
	pubkey = strings.Replace(pubkey, " ", "\n", -1)
	rsaEncrypt.PubKey = "-----BEGIN PUBLIC KEY-----" + pubkey + "-----END PUBLIC KEY-----"
	//	log.Println(rsaEncrypt.PubKey)

	err = filepath.Walk(tool.inputDir, rsaEncrypt.walkFunc)
	if err != nil {
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [错误日志]  内容：[" + err.Error() + "]</div>")
		root.CallFunction("enable", ret)
		root.Update(false)
		return
	}

	root.CallFunction("enable", ret)
	root.Update(false)
	return
}

func DoRsaModePrikeyEncrypt(rsaEncrypt *RsaEncrypt) {
	root, err := w.GetRootElement()
	if err != nil {
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [错误日志]  内容：[" + err.Error() + "]</div>")
		root.CallFunction("enable", ret)
		root.Update(false)
		return
	}
	ClearMsg()

	resultElement, err := root.SelectById("inputPrikey")
	if err != nil {
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [错误日志]  内容：[" + err.Error() + "]</div>")
		root.CallFunction("enable", ret)
		root.Update(false)
		return
	}
	prikey, err := resultElement.Text()
	if err != nil {
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [错误日志]  内容：[" + err.Error() + "]</div>")
		root.CallFunction("enable", ret)
		root.Update(false)
		return
	}

	prikey = strings.Replace(prikey, "-----BEGIN PRIVATE KEY-----", "", -1)
	prikey = strings.Replace(prikey, "-----END PRIVATE KEY-----", "", -1)
	prikey = strings.Replace(prikey, " ", "\n", -1)
	rsaEncrypt.PriKey = "-----BEGIN PRIVATE KEY-----" + prikey + "-----END PRIVATE KEY-----"

	err = filepath.Walk(tool.inputDir, rsaEncrypt.walkFunc)
	if err != nil {
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [错误日志]  内容：[" + err.Error() + "]</div>")
		root.CallFunction("enable", ret)
		root.Update(false)
		return
	}

	root.CallFunction("enable", ret)
	root.Update(false)
	return
}

func DoAes(aesEncrypt *AesEncrypt) {

	root, err := w.GetRootElement()
	if err != nil {
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [错误日志]  内容：[" + err.Error() + "]</div>")
		root.CallFunction("enable", ret)
		root.Update(false)
		return
	}

	ClearMsg()

	err = filepath.Walk(tool.inputDir, aesEncrypt.walkFunc)
	if err != nil {
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [错误日志]  内容：[" + err.Error() + "]</div>")
		root.CallFunction("enable", ret)
		root.Update(false)
		return
	}

	root.CallFunction("enable", ret)
	root.Update(false)
	return
}

func DoDes(desEncrypt *DesEncrypt) {

	root, err := w.GetRootElement()
	if err != nil {
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [错误日志]  内容：[" + err.Error() + "]</div>")
		root.CallFunction("enable", ret)
		root.Update(false)
		return
	}

	ClearMsg()

	err = filepath.Walk(tool.inputDir, desEncrypt.walkFunc)
	if err != nil {
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [错误日志]  内容：[" + err.Error() + "]</div>")
		root.CallFunction("enable", ret)
		root.Update(false)
		return
	}

	root.CallFunction("enable", ret)
	root.Update(false)
	return
}

func DoTripleDes(tripleDesEncrypt *TripleDesEncrypt) {

	root, err := w.GetRootElement()
	if err != nil {
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [错误日志]  内容：[" + err.Error() + "]</div>")
		root.CallFunction("enable", ret)
		root.Update(false)
		return
	}

	ClearMsg()

	err = filepath.Walk(tool.inputDir, tripleDesEncrypt.walkFunc)
	if err != nil {
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [错误日志]  内容：[" + err.Error() + "]</div>")
		root.CallFunction("enable", ret)
		root.Update(false)
		return
	}

	root.CallFunction("enable", ret)
	root.Update(false)
	return
}

func DoRsa(rsaEncrypt *RsaEncrypt) {

	root, err := w.GetRootElement()
	if err != nil {
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [错误日志]  内容：[" + err.Error() + "]</div>")
		root.CallFunction("enable", ret)
		root.Update(false)
		return
	}

	ClearMsg()

	AppendMsg("<div style=\"color:#FF8C00\">" + NowTime() + "  [生成秘钥中...]</div>")

	err = rsaEncrypt.GenRsaKey(tool)
	if err != nil {
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [错误日志]  内容：[" + err.Error() + "]</div>")
		root.CallFunction("enable", ret)
		root.Update(false)
		return
	}
	AppendMsg("<div style=\"color:#43CD80\">" + NowTime() + "  [生成秘钥完成]</div>")

	err = filepath.Walk(tool.inputDir, rsaEncrypt.walkFunc)
	if err != nil {
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [错误日志]  内容：[" + err.Error() + "]</div>")
		root.CallFunction("enable", ret)
		root.Update(false)
		return
	}

	root.CallFunction("enable", ret)
	root.Update(false)
	return
}

func NowTime() string {
	return time.Now().Format("2006/01/02 15:04:05.9999")
}

func AppendMsg(msg string) error {
	root, err := w.GetRootElement()
	if err != nil {
		return err
	}

	resultElement, err := root.SelectById("result")
	if err != nil {
		return err
	}
	err = resultElement.SetHtml(msg, sciter.SIH_APPEND_AFTER_LAST)
	if err != nil {
		return err
	}
	root.Update(true)
	return nil
}

func PutKey(id string, key string) error {
	root, err := w.GetRootElement()
	if err != nil {
		return err
	}

	resultElement, err := root.SelectById(id)
	if err != nil {
		return err
	}
	err = resultElement.SetText(key)
	if err != nil {
		return err
	}
	root.Update(true)
	return nil
}

func ClearKey() error {
	root, err := w.GetRootElement()
	if err != nil {
		return err
	}

	resultElement1, err := root.SelectById("pubkey")
	if err != nil {
		return err
	}
	err = resultElement1.SetText("")
	if err != nil {
		return err
	}
	resultElement2, err := root.SelectById("prikey")
	if err != nil {
		return err
	}
	err = resultElement2.SetText("")
	if err != nil {
		return err
	}
	root.Update(true)
	return nil
}

func ClearMsg() error {
	root, err := w.GetRootElement()
	if err != nil {
		return err
	}

	resultElement, err := root.SelectById("result")
	if err != nil {
		return err
	}
	err = resultElement.SetHtml("<div id=\"result\" class=\"list\"></div>", sciter.SOH_REPLACE)
	if err != nil {
		return err
	}
	root.Update(true)
	return nil
}

func (rsaEncrypt *RsaEncrypt) walkFunc(path string, info os.FileInfo, err error) error {
	if err != nil {
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [错误日志]  内容：[" + err.Error() + "]</div>")
		return err
	}

	if info.IsDir() {
		return nil
	}

	size := info.Size()

	if size > 5*1024*1024 {
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [文件超大]  文件：[" + path + "]</div>")
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [错误日志]  内容：[文件大于5MB，请切分为小配置文件]</div>")
		return nil
	}

	file, err := os.Open(path)

	if err != nil {
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [加密失败]  文件：[" + path + "]</div>")
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [错误日志]  内容：[" + err.Error() + "]</div>")
		return err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)

	if err != nil {
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [加密失败]  文件：[" + path + "]</div>")
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [错误日志]  内容：[" + err.Error() + "]</div>")
		return err
	}

	bytes, err = rsaEncrypt.Encrypt(string(bytes))

	if err != nil {
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [加密失败]  文件：[" + path + "]</div>")
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [错误日志]  内容：[" + err.Error() + "]</div>")
		return err
	}
	err = AppendMsg("<div style=\"color:#00008B\">" + NowTime() + "  [加密完成]  文件：[" + path + "]</div>")

	start := len(tool.inputDir)
	end := len(path)

	out := tool.outputDir + Substr(path, start, end)
	efile, err := os.Create(out)
	if err != nil {
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [加密失败]  文件：[" + path + "]</div>")
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [错误日志]  内容：[" + err.Error() + "]</div>")
		return err
	}
	defer efile.Close()

	_, err = efile.WriteString(base64.StdEncoding.EncodeToString(bytes))
	if err != nil {
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [加密失败]  文件：[" + path + "]</div>")
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [错误日志]  内容：[" + err.Error() + "]</div>")
		return err
	}

	return nil
}

func (aesEncrypt *AesEncrypt) walkFunc(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if info.IsDir() {
		return nil
	}

	size := info.Size()

	if size > 5*1024*1024 {
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [文件超大]  文件：[" + path + "]</div>")
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [错误日志]  内容：[文件大于5MB，请切分为小配置文件]</div>")
		return nil
	}

	file, err := os.Open(path)

	if err != nil {
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [加密失败]  文件：[" + path + "]</div>")
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [错误日志]  内容：[" + err.Error() + "]</div>")
		return err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)

	if err != nil {
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [加密失败]  文件：[" + path + "]</div>")
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [错误日志]  内容：[" + err.Error() + "]</div>")
		return err
	}

	bytes, err = aesEncrypt.Encrypt(string(bytes))

	if err != nil {
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [加密失败]  文件：[" + path + "]</div>")
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [错误日志]  内容：[" + err.Error() + "]</div>")
		return err
	}
	err = AppendMsg("<div style=\"color:#00008B\">" + NowTime() + "  [加密完成]  文件：[" + path + "]</div>")

	start := len(tool.inputDir)
	end := len(path)

	out := tool.outputDir + Substr(path, start, end)
	efile, err := os.Create(out)
	if err != nil {
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [加密失败]  文件：[" + path + "]</div>")
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [错误日志]  内容：[" + err.Error() + "]</div>")
		return err
	}
	defer efile.Close()

	_, err = efile.WriteString(base64.StdEncoding.EncodeToString(bytes))
	if err != nil {
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [加密失败]  文件：[" + path + "]</div>")
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [错误日志]  内容：[" + err.Error() + "]</div>")
		return err
	}

	return nil
}

func (desEncrypt *DesEncrypt) walkFunc(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if info.IsDir() {
		return nil
	}

	size := info.Size()

	if size > 5*1024*1024 {
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [文件超大]  文件：[" + path + "]</div>")
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [错误日志]  内容：[文件大于5MB，请切分为小配置文件]</div>")
		return nil
	}

	file, err := os.Open(path)

	if err != nil {
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [加密失败]  文件：[" + path + "]</div>")
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [错误日志]  内容：[" + err.Error() + "]</div>")
		return err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)

	if err != nil {
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [加密失败]  文件：[" + path + "]</div>")
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [错误日志]  内容：[" + err.Error() + "]</div>")
		return err
	}

	bytes, err = desEncrypt.Encrypt(string(bytes))

	if err != nil {
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [加密失败]  文件：[" + path + "]</div>")
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [错误日志]  内容：[" + err.Error() + "]</div>")
		return err
	}
	err = AppendMsg("<div style=\"color:#00008B\">" + NowTime() + "  [加密完成]  文件：[" + path + "]</div>")

	start := len(tool.inputDir)
	end := len(path)

	out := tool.outputDir + Substr(path, start, end)
	efile, err := os.Create(out)
	if err != nil {
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [加密失败]  文件：[" + path + "]</div>")
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [错误日志]  内容：[" + err.Error() + "]</div>")
		return err
	}
	defer efile.Close()

	_, err = efile.WriteString(base64.StdEncoding.EncodeToString(bytes))
	if err != nil {
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [加密失败]  文件：[" + path + "]</div>")
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [错误日志]  内容：[" + err.Error() + "]</div>")
		return err
	}

	return nil
}

func (tripleDesEncrypt *TripleDesEncrypt) walkFunc(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if info.IsDir() {
		return nil
	}

	size := info.Size()

	if size > 5*1024*1024 {
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [文件超大]  文件：[" + path + "]</div>")
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [错误日志]  内容：[文件大于5MB，请切分为小配置文件]</div>")
		return nil
	}

	file, err := os.Open(path)

	if err != nil {
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [加密失败]  文件：[" + path + "]</div>")
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [错误日志]  内容：[" + err.Error() + "]</div>")
		return err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)

	if err != nil {
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [加密失败]  文件：[" + path + "]</div>")
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [错误日志]  内容：[" + err.Error() + "]</div>")
		return err
	}

	bytes, err = tripleDesEncrypt.Encrypt(string(bytes))

	if err != nil {
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [加密失败]  文件：[" + path + "]</div>")
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [错误日志]  内容：[" + err.Error() + "]</div>")
		return err
	}
	err = AppendMsg("<div style=\"color:#00008B\">" + NowTime() + "  [加密完成]  文件：[" + path + "]</div>")

	start := len(tool.inputDir)
	end := len(path)

	out := tool.outputDir + Substr(path, start, end)
	efile, err := os.Create(out)
	if err != nil {
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [加密失败]  文件：[" + path + "]</div>")
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [错误日志]  内容：[" + err.Error() + "]</div>")
		return err
	}
	defer efile.Close()

	_, err = efile.WriteString(base64.StdEncoding.EncodeToString(bytes))
	if err != nil {
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [加密失败]  文件：[" + path + "]</div>")
		AppendMsg("<div style=\"color:#FF0000\">" + NowTime() + "  [错误日志]  内容：[" + err.Error() + "]</div>")
		return err
	}

	return nil
}

func Substr(str string, start, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}
	return string(rs[start:end])
}
