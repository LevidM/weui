package php

import (
	"archive/zip"
	"bufio"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"html"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"net"
	"net/http"
	"net/smtp"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)
/**
 * @Description: 当前时间戳
 * @return int64
 */
func Time() int64 {
	return time.Now().Unix()
}
/**
 * @Description: 格式化时间类型
 * @param format Y-m-d H:i:s
 * @param ts 时间类型 2021-01-02 13:45:34 +0800 CST
 * @return string 时间字符串
 */
func Date(format string, ts ...time.Time) string {
	patterns := []string{
		// 年
		"Y", "2006", // 4 位数字完整表示的年份
		"y", "06", // 2 位数字表示的年份

		// 月
		"m", "01", // 数字表示的月份，有前导零
		"n", "1", // 数字表示的月份，没有前导零
		"M", "Jan", // 三个字母缩写表示的月份
		"F", "January", // 月份，完整的文本格式，例如 January 或者 March

		// 日
		"d", "02", // 月份中的第几天，有前导零的 2 位数字
		"j", "2", // 月份中的第几天，没有前导零

		"D", "Mon", // 星期几，文本表示，3 个字母
		"l", "Monday", // 星期几，完整的文本格式;L的小写字母

		// 时间
		"g", "3", // 小时，12 小时格式，没有前导零
		"G", "15", // 小时，24 小时格式，没有前导零
		"h", "03", // 小时，12 小时格式，有前导零
		"H", "15", // 小时，24 小时格式，有前导零

		"a", "pm", // 小写的上午和下午值
		"A", "PM", // 小写的上午和下午值

		"i", "04", // 有前导零的分钟数
		"s", "05", // 秒数，有前导零
	}
	replacer := strings.NewReplacer(patterns...)
	format = replacer.Replace(format)

	t := time.Now()
	if len(ts) > 0 {
		t = ts[0]
	}
	return t.Format(format)
}
/**
 * @Description: 时间戳转换成时间类型
 * @param t 时间戳
 * @return time.Time 2021-01-02 13:45:34 +0800 CST
 */
func Unix2Time(t int64)time.Time{
	const timeLayout = "2006-01-02 15:04:05"
	str:= time.Unix(t, 0).Format(timeLayout)
	return Timestr2Time(str)
}
/**
 * @Description: 时间字符串转换时间类型
 * @param str 时间字符串
 * @return time.Time 2021-01-02 13:45:34 +0800 CST
 */
func Timestr2Time(str string)time.Time{
	const Layout = "2006-01-02 15:04:05"//时间常量
	loc, _ := time.LoadLocation("Asia/Shanghai")
	times,_ := time.ParseInLocation(Layout,str,loc)
	return times
}
/**
 * @Description: 时间字符串转换时间戳
 * @param s 时间字符类型
 * @return int64 时间戳
 */
func Str2Time(s string)int64{
	const timeLayout = "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")
	tmp, _ := time.ParseInLocation(timeLayout, s, loc)
	timestamp := tmp.Unix()
	return timestamp
}
/**
 * @Description: 友好时间显示
 * @param t 时间戳
 * @return string
 */
func TimeLine(t int64)string{
	now :=time.Now().Unix()
	var xx string
	if now<=t{
		xx= Date("Y-m-d H:i:s",Unix2Time(t))
	}else{
		t= now-t
		f:=map[int]string{
			31536000 :"年",
			2592000:"个月",
			604800 :"星期",
			86400 :"天",
			3600:"小时",
			60 :"分钟",
			1 :"秒"}
		var keys []int
		for k := range f {
			keys = append(keys, k)
		}
		sort.Sort(sort.Reverse(sort.IntSlice(keys)))
		for _,v:=range keys{
			k1:=int64(v)
			x:= t/k1
			if x!=0 {
				x1 := strconv.FormatInt(x,10)
				xx= x1+f[v]+"前"
				break
			}
		}}
	return xx
}
/**
 * @Description: 格式化字节
 * @param sizes
 * @return string
 */
func FileCount(sizes uint64)(string){
	a:=[...]string{"B", "KB", "MB", "GB", "TB", "PB"}
	pos:=0
	s:=float64(sizes)
	for s>=1024 {
		s =s/1024
		pos++
	}
	c := strconv.FormatFloat(s,'f',2,64)
	return c+" "+a[pos];
}
/**
 * @Description: 生成随机id
 * @param prefix 前缀字符
 * @return string
 */
func Uniqid(prefix string) string {
	now := time.Now()
	return fmt.Sprintf("%s%08x%05x", prefix, now.Unix(), now.UnixNano()%0x100000)
}
/**
 * @Description: 结构体转换字典map
 * @param obj 结构体
 * @return map[string]interface{}
 */
func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}
/**
 * @Description: 类型判断
 * @param a
 * @return reflect.Type
 */
func Typeof(a interface{})reflect.Type{
	return reflect.TypeOf(a)
}
/**
 * @Description: 读取文件
 * @param path 路径
 * @return string
 */
func ReadFile(path string)string{
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	bufReader := bufio.NewReader(file)
	var i = 0
	s:=""
	for{
		i++
		line,err := bufReader.ReadString(';')
		s+=line
		if err == io.EOF {
			break
		}
	}
	return s
}
/**
 * @Description: 判断元素是否在数组,切片,字典中
 * @param needle 元素
 * @param haystack 数组
 * @return bool
 */
func InArray(needle interface{}, haystack interface{}) bool {
	val := reflect.ValueOf(haystack)
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			if reflect.DeepEqual(needle, val.Index(i).Interface()) {
				return true
			}
		}
	case reflect.Map:
		for _, k := range val.MapKeys() {
			if reflect.DeepEqual(needle, val.MapIndex(k).Interface()) {
				return true
			}
		}
	default:
		panic("haystack: haystack type muset be slice, array or map")
	}

	return false
}
/**
 * @Description: 四舍五入保留小数位数
 * @param value 浮点数
 * @param n 小数位数 0是整数
 * @return float64
 */
func Round(value float64,n int) float64 {
	n10 := math.Pow10(n)
	return math.Trunc((value+0.5/n10)*n10) / n10
}
/**
 * @Description: 将utf-8编码的字符串转换为GBK编码
 * @param str 要转换字符串
 * @return string
 */
func Utf2Gbk(str string) string {
	ret, _ := simplifiedchinese.GBK.NewEncoder().String(str)
	return ret
	b, _ := simplifiedchinese.GBK.NewEncoder().Bytes([]byte(str))
	return string(b)
}
/**
 * @Description: 将GBK编码的字符串转换为utf-8编码
 * @param gbkStr
 * @return string
 */
func Gbk2Utf8(gbkStr string) string {
	ret, _ := simplifiedchinese.GBK.NewDecoder().String(gbkStr)
	return ret
	b, _ := simplifiedchinese.GBK.NewDecoder().Bytes([]byte(gbkStr))
	return string(b)
}
/**
 * @Description: 文件修改时间,返回时间戳
 * @param file
 * @return int64 时间戳
 * @return error
 */
func FileMtime(file string) (int64, error) {
	f, err := os.Stat(file)
	if err != nil {
		return 0, err
	}
	return f.ModTime().Unix(), nil
}
/**
 * @Description: 返回文件大小字节
 * @param file 文件名
 * @return uint64 字节大小
 * @return error
 */
func FileSize(file string) (uint64, error) {
	f, err := os.Stat(file)
	if err != nil {
		return 0, err
	}
	return uint64(f.Size()), nil
}
/**
 * @Description: 写文件
 * @param filename 文件名
 * @param data []byte("我爱你12345")
 * @return error
 */
func WriteFile(filename string, data []byte) error {
	os.MkdirAll(path.Dir(filename), os.ModePerm)
	return ioutil.WriteFile(filename, data, 0655)
}
/**
 * @Description: 判断是否文件,不存在也返回否
 * @param filePath
 * @return bool
 */
func IsFile(filePath string) bool {
	f, e := os.Stat(filePath)
	if e != nil {
		return false
	}
	return !f.IsDir()
}
/**
 * @Description: 判断是否目录
 * @param dir
 * @return bool
 */
func IsDir(dir string) bool {
	f, e := os.Stat(dir)
	if e != nil {
		return false
	}
	return f.IsDir()
}
/**
 * @Description: 判断文件或目录是否存在
 * @param path
 * @return bool
 */
func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}
/**
 * @Description: 创建目录 php.CreateDir("./112/qq","112/vv")
 * @param dirs 多个目录,支持多级别
 * @return err
 */
func CreateDir(dirs ...string) (err error) {
	for _, v := range dirs {
		exist:= IsExist(v)
		if !exist {
			err = os.MkdirAll(v, os.ModePerm)
			CheckErr(err)
		}
	}
	return err
}
/**
 * @Description: 获取go路径,返回切片
 * @return []string
 */
func GetGopath() []string {
	gopath := os.Getenv("GOPATH")
	var paths []string
	if runtime.GOOS == "windows" {
		gopath = strings.Replace(gopath, "\\", "/", -1)
		paths = strings.Split(gopath, ";")
	} else {
		paths = strings.Split(gopath, ":")
	}
	return paths
}
//是否邮箱
func IsEmail(email string) bool {
	return regexp.MustCompile(`(?i)[A-Z0-9._%+-]+@(?:[A-Z0-9-]+\.)+[A-Z]{2,6}`).MatchString(email)
}
//是否url
func IsUrl(url string) bool {
	return regexp.MustCompile(`(ftp|http|https):\/\/(\w+:{0,1}\w*@)?(\S+)(:[0-9]+)?(\/|\/([\w#!:.?+=&%@!\-\/]))?`).MatchString(url)
}
/**
 * @Description: 输出csv
 * @param filename 文件名
 * @param data
	data:= [][]string{
		{"1", "test1", "李白"},
		{"2", "test2", "2015-12-26"},
		{"3", "test3", "test3-1"},
	}
 */
func WriteCsv(filename string,data [][]string){
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.WriteString("\xEF\xBB\xBF")
	w := csv.NewWriter(f)
	w.WriteAll(data)
}
/**
 * @Description: 读取csv返回二维切片
 * @param filename
 * @return [][]string
 */
func ReadCsv(filename string)[][]string{
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	reader := csv.NewReader(f)
	var s [][]string
	for {
		record,err := reader.Read()
		if err == io.EOF {
			break
		}
		s=append(s,record)
	}
	return s
}
/**
 * @Description: 获取exe路径
 * @return string C:\www\go\bin
 */
func GetPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))
	return path[:index]
}
/**
 * @Description: 数据反顺序输出
 * @param s 切片
 * @return []interface{}
 */
func ArrayReverse(s []interface{}) []interface{} {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]    }
	return s
}
/**
 * @Description: 合并两个数组
 * @param ss 多个数组
 * @return []string
 */
func ArrayMerge(ss ...[]string) []string {
	n := 0
	for _, v := range ss {
		n += len(v)
	}
	s := make([]string, 0, n)
	for _, v := range ss {
		s = append(s, v...)
	}
	return s
}
/**
 * @Description: 返回value组成新数组
 * @param elements 字典
 * @return []interface{} 新数组
 */
func ArrayValues(elements map[interface{}]interface{}) []interface{} {
	i, vals := 0, make([]interface{}, len(elements))
	for _, val := range elements {
		vals[i] = val
		i++
	}
	return vals
}
/**
 * @Description: 返回key组成新数组
 * @param elements
 * @return []interface{}
 */
func ArrayKeys(elements map[interface{}]interface{}) []interface{} {
	i, keys := 0, make([]interface{}, len(elements))
	for key, _ := range elements {
		keys[i] = key
		i++
	}
	return keys
}
/**
 * @Description: 切片转换字符串	s:=[]interface{}{"我们","爱你"}
	v:=php.Array2String(s,"")
 * @param array 切片
 * @param ss 分隔符
 * @return string
 */
func Array2String(array []interface{},ss string) string {
	return strings.Replace(strings.Trim(fmt.Sprint(array), "[]"), " ", ss, -1)
}
/**
 * @Description: 产生随机数,最小,最大
 * @param min 最小数
 * @param max 最大数
 * @return int64 随机
 */
func MtRand(min, max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(max-min+1) + min
}
/**
 * @Description: md5加密
 * @param str
 * @return string
 */
func Md5(str string) string {
	h := md5.New()
	io.WriteString(h, str)
	return fmt.Sprintf("%x", h.Sum(nil))
}
/**
 * @Description: 生成sha1
 * @param str
 * @return string
 */
func Sha1(str string) string {
	hash := sha1.New()
	io.WriteString(hash, str)
	return fmt.Sprintf("%x", hash.Sum(nil))
}
/**
 * @Description: 创建密码散列,password 密码
 * @param password 要加密字符串
 * @return string
 * @return error
 */
func PasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
/**
 * @Description: 验证密码,密码和密码散列
 * @param password 密码
 * @param hash 上面生成的散列密码
 * @return bool
 */
func PasswordVerify(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
/**
 * @Description: 四舍五入,分割字符
 * @param number 浮点数
 * @param decimals 保留位数2
 * @param decPoint .
 * @param thousandsSep , 分隔符
 * @return string
 */
func NumberFormat(number float64, decimals uint, decPoint, thousandsSep string) string {
	neg := false
	if number < 0 {
		number = -number
		neg = true
	}
	dec := int(decimals)
	// Will round off
	str := fmt.Sprintf("%."+strconv.Itoa(dec)+"F", number)
	prefix, suffix := "", ""
	if dec > 0 {
		prefix = str[:len(str)-(dec+1)]
		suffix = str[len(str)-dec:]	} else {
		prefix = str
	}
	sep := []byte(thousandsSep)
	n, l1, l2 := 0, len(prefix), len(sep)
	// thousands sep num
	c := (l1 - 1) / 3
	tmp := make([]byte, l2*c+l1)
	pos := len(tmp) - 1
	for i := l1 - 1; i >= 0; i, n, pos = i-1, n+1, pos-1 {
		if l2 > 0 && n > 0 && n%3 == 0 {
			for j := range sep {
				tmp[pos] = sep[l2-j-1]
				pos--
			}
		}
		tmp[pos] = prefix[i]	}
	s := string(tmp)
	if dec > 0 {
		s += decPoint + suffix
	}
	if neg {
		s = "-" + s
	}

	return s
}
/**
 * @Description: 去除字符中HTML标记
 * @param content 字符串
 * @return string
 */
func StripTags(content string) string {
	re := regexp.MustCompile(`<(.|\n)*?>`)
	return re.ReplaceAllString(content,"")
}
/**
 * @Description: 返回当前 Unix 时间戳和微秒数
 * @return float64
 */
func MicroTime() float64 {
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	micSeconds := float64(now.Nanosecond()) / 1000000000
	return float64(now.Unix()) + micSeconds
}
/**
 * @Description: 判断是否为空
 * @param val 变量
 * @return bool
 */
func Empty(val interface{}) bool {
	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.String, reflect.Array:
		return v.Len() == 0
	case reflect.Map, reflect.Slice:
		return v.Len() == 0 || v.IsNil()
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}

	return reflect.DeepEqual(val, reflect.Zero(v.Type()).Interface())
}
/**
 * @Description: 显示错误
 * @param err
 */
func CheckErr(err error) {
	if err != nil{
		panic(err)
	}
}
/**
 * @Description: 计算字符个数
 * @param str 字符串
 * @return int
 */
func Len(str string)int{
	return len([]rune(str))
}
/**
 * @Description: cli实现提示输入 qq:=php.Ask("请输入",nil)
 * @param str 提示字符串
 * @param check nil
 * @return string
 */
func Ask(str string,check func(string) error)string{
	if check == nil {
		check = func(in string) error {
			if len(in) > 0 {
				return nil
			} else {
				return fmt.Errorf("Cannot be empty")
			}
		}
	}
	input := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf(str+"")
		line, _, err := input.ReadLine()
		for err == io.EOF {
			<-time.After(time.Millisecond)
			line, _, err = input.ReadLine()
		}
		if err != nil {
			fmt.Printf("<warn>%s \n\n", err)
		} else if err = check(string(line)); err != nil {
			fmt.Printf("<warn>%s \n\n", err)
		} else {
			return string(line)
		}
	}
}
//cli选择
var RenderChooseQuestion = func(question string) string {
	return question + "\n"
}
var RenderChooseOption = func(key, value string, size int) string {
	return fmt.Sprintf("%-"+fmt.Sprintf("%d", size+1)+"s %s\n", key+")", value)
}
var RenderChooseQuery = func() string {
	return "Choose: "
}
/**
 * @Description: cli选择
 * @param question 提示内容
 * @param choices  map[string]string{
						"1":  "苹果",
						"2": "橘子",
						"3":   "西瓜",
					}
 * @return string 返回key
 */
func Choose(question string, choices map[string]string) string {
	options := RenderChooseQuestion(question)
	keys := []string{}
	max := 0
	for k, _ := range choices {
		if l := len(k); l > max {
			max = l
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		options += RenderChooseOption(k, choices[k], max)
	}
	options += RenderChooseQuery()
	return Ask(options,func(in string) error {
		if _, ok := choices[in]; ok {
			return nil
		} else {
			return fmt.Errorf("Choose one of: %s", strings.Join(keys, ", "))
		}
	})
}
var ConfirmRejection = "<warn>Please respond with \"y\" or \"n\"\n\n"
var ConfirmYesRegex = regexp.MustCompile(`^(?i)y(es)?$`)
var ConfirmNoRegex = regexp.MustCompile(`^(?i)no?$`)
/**
 * @Description: cli选择yes/no
 * @param question 提示内容支持回复 yes y/no n
 * @return bool
 */
func Confirm(question string) bool {
	cb := func(value string) error {return nil}
	for {
		res := Ask(question, cb)
		if ConfirmYesRegex.MatchString(res) {
			return true
		} else if ConfirmNoRegex.MatchString(res) {
			return false
		} else {
			fmt.Printf(ConfirmRejection)
		}
	}
}
/**
 * @Description: 设置字体颜色,仅用于linux,cmder,标准cmd不支持
 * @param str 字符串
 * @param color1
 * @param extraArgs 忽略
 * @return string php.ColorLinux("测试","red")
 */
func ColorLinux(str string, color1 string, extraArgs ...interface{}) string {
	//闪烁效果
	var isBlink int64 = 0
	var color int
	var weight int=1
	m:= map[string]int{"green":32, "red":31,"yellow":33,"black":30,"white":37,"blue":34,"zi":35,"qing":36}
	if v, ok := m[color1]; ok {
		color=v
	}else{
		color=31
	}
	if len(extraArgs) > 0 {
		isBlink = reflect.ValueOf(extraArgs[0]).Int()
	}

	//下划线效果
	var isUnderLine int64 = 0
	if len(extraArgs) > 1 {
		isUnderLine = reflect.ValueOf(extraArgs[1]).Int()
	}
	var mo []string
	if isBlink > 0 {
		mo = append(mo, "05")
	}
	if isUnderLine > 0 {
		mo = append(mo, "04")
	}
	if weight > 0 {
		mo = append(mo, fmt.Sprintf("%d", weight))
	}
	if len(mo) <= 0 {
		mo = append(mo, "0")
	}
	buf := bytes.Buffer{}
	buf.WriteString("\033[")
	buf.WriteString(strings.Join(mo, ";"))
	buf.WriteString(";")
	buf.WriteString(fmt.Sprintf("%d", color))
	buf.WriteString("m")
	buf.WriteString(str)
	buf.WriteString("\033[0m")
	return buf.String()
}
/**
 * @Description: 支持windows设置颜色
 * @param s 字符
 * @param i1 颜色
 */
/*func Color(s interface{}, i1 string){
	c:=	map[string]int{
		"black":0,
		"blue":1,
		"green":2,
		"cyan":3,//青色
		"red":4,
		"purple":5,//紫色
		"yellow":6,//
		"white":15,
		"gray":8,
		"qing":3,
		"zi":5,
	}
	var i int
	if v, ok := c[i1]; ok {
		i=v
	}else{
		i=4
	}
	if runtime.GOOS=="windows"{
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	proc := kernel32.NewProc("SetConsoleTextAttribute")
	handle, _, _ := proc.Call(uintptr(syscall.Stdout), uintptr(i))
	fmt.Print(s)
	handle, _, _ = proc.Call(uintptr(syscall.Stdout), uintptr(7))
	CloseHandle := kernel32.NewProc("CloseHandle")
	CloseHandle.Call(handle)
	}
}*/
/**
 * @Description: 获取本地IP 返回切片
 * @return ips
 */
func Getip() (string) {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("net.Interfaces failed, err:", err.Error())
	}

	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()

			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						return ipnet.IP.String()
					}
				}
			}
		}

	}
	return ""
}
/**
 * @Description: 自动访问网址 php.Openurl("https://www.baidu.com")
 * @param uri
 */
func Openurl(url string) {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default:
		cmd = "xdg-open"
	}
	args = append(args, url)
	cmds:=exec.Command(cmd, args...)
	if runtime.GOOS=="windows" {
		//cmds.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	}
	cmds.Start()
}
/**
 * @Description: 发送邮件不使用任何扩展
	from :="logwwwove@qq.com"
	to1:="yobybxy@163.com,18291448834@163.com"
	secret := "abc"
	host :="smtp.qq.com"
	port := 25
	subject:="主题"
	body:="内容是测试"
	err:=php.SendEmail(from,to1,subject,body,secret,host,port)
	php.CheckErr(err)
 * @param from 发送人邮箱
 * @param to1 收件人邮箱,多个用,隔开"yoby21bxy@163.com,182914114811834@163.com"
 * @param subject 标题
 * @param body 内容
 * @param secret 密钥,qq邮箱授权码密码
 * @param host 主机地址
 * @param port 端口25
 * @return err
 */
func SendEmail(from string,to1 string,subject string,body string,secret string,host string,port int)(err error){
	to2:=strings.Split(to1,",")
	to:=to2[0]
	auth := smtp.PlainAuth("", from, secret, host)
	msg := []byte("To: "+to+"\r\n" +
		"Subject: "+subject+"\r\n" +
		"\r\n" +
		body+"\r\n")
	err = smtp.SendMail(host+":"+strconv.Itoa(port), auth, from, to2, msg)
	return err
}
/**
 * @Description: int转换字符串
 * @param i int
 * @return s 字符串
 */
func Int2String(i interface{})(s string){
	ty:=reflect.TypeOf(i).String()
	if ty=="int"{
		ii:=i.(int)
		s=strconv.Itoa(ii)
	}else if ty=="int64"{
		ii:=i.(int64)
		s= strconv.FormatInt(ii, 10)
	}else if ty=="uint64"{
		ii:=i.(uint64)
		s= strconv.FormatUint(ii, 10)
	}
	return s
}
func String2Int(i string,ty string)(s interface{}){
	if ty=="int"{
		s, _ = strconv.Atoi(i)
	}else if ty=="int32"{
		ii,_:= strconv.ParseInt(i, 10, 64)
		s=int32(ii)
	}else if ty=="int64"{
		s,_= strconv.ParseInt(i, 10, 64)
	}else if ty=="uint"{
		s,_= strconv.ParseUint(i, 10, 64)
	}
	return s
}
/**
 * @Description: 浮点数转换字符串
 * @param i 浮点数
 * @param ty 类型 64表示64位浮点数 32表示32位浮点数
 * @return s
 */
func Float2String(i float64,ty int)(s string){
	if ty==64{
		s= strconv.FormatFloat(i, 'E', -1, 64)
	}else{
		s= strconv.FormatFloat(i, 'E', -1, 32)
	}
	return s
}
/**
 * @Description: 字符串转换浮点数
 * @param i 浮点字符串
 * @param ty 32或64浮点
 * @return s
 */
func String2Float(i string,ty int)(s1 interface{}){
	s, _:= strconv.ParseFloat(i, ty)
	if ty==32{
		s1=float32(s)
	}else{
		s1=s
	}
	return s1
}
/**
 * @Description: 压缩文件zip Zip([]string{"./1.txt","qq/"},"./1.zip")
 * @param files 可以是文件或目录
 * @param dest zip文件名包括扩展名
 * @return error
 */
func Zip(files []string, dest string)error{
	d, _ := os.Create(dest)
	defer d.Close()
	w := zip.NewWriter(d)
	defer w.Close()
	for _, file := range files {
		f, err:= os.Open(file)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		err = zips(f, "", w)
		if err != nil {
			return err
		}
	}
	return nil
}
func zips(file *os.File, prefix string, zw *zip.Writer) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		if len(prefix) == 0 {
			prefix = info.Name()
		} else {
			prefix = prefix + "/" + info.Name()
		}
		fileInfos, err := file.Readdir(-1)
		if err != nil {
			return err
		}
		for _, fi := range fileInfos {
			f, err := os.Open(file.Name() + "/" + fi.Name())
			if err != nil {
				return err
			}
			err = zips(f, prefix, zw)
			if err != nil {
				return err
			}
		}
	} else {
		header, err := zip.FileInfoHeader(info)
		if len(prefix) == 0 {
			header.Name = header.Name
		} else {
			header.Name = prefix + "/" + header.Name
		}
		if err != nil {
			return err
		}
		writer, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, file)
		file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
/**
 * @Description: 解压缩zip Unzip("./1.zip","./ss")
 * @param archive 要解压zip路径
 * @param target 目标文件夹可以生成
 * @return error
 */
func Unzip(archive, target string)error{
	reader, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(target, 0755); err != nil {
		return err
	}
	for _, file := range reader.File {
		path := filepath.Join(target, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}
		dir := filepath.Dir(path)
		if len(dir) > 0 {
			if _, err = os.Stat(dir); os.IsNotExist(err) {
				err = os.MkdirAll(dir, 0755)
				if err != nil {
					return err
				}
			}
		}
		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()
		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}
	}
	return nil
}
/**
 * @Description: 随机字符串
 * @param l 指定长度
 * @return string
 */
func  RandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
/**
 * @Description: 随机数字
 * @param width 指定长度
 * @return string
 */
func RandomNumber(width int) string {
	numeric := [10]byte{0,1,2,3,4,5,6,7,8,9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&sb, "%d", numeric[ rand.Intn(r) ])
	}
	return sb.String()
}

/**
 * @Description: 编码成base64
str:="http://www.baidu.com-+_"
str=php.Base64Encode(str,false)
 * @param s 要编码字符串
 * @param isurl 是否url false或true
 * @return s1 字符串
 */
func Base64Encode(s string,isurl bool) (s1 string){
	if isurl==true{
		s1 = base64.URLEncoding.EncodeToString([]byte(s))
	}else{
		s1=base64.StdEncoding.EncodeToString([]byte(s))
	}
	return s1
}
/**
 * @Description: 解码base64 str=php.Base64Decode(str,false)
 * @param s 要解码字符串
 * @param isurl 是否url
 * @return string
 */
func Base64Decode(s string,isurl bool) (string){
	var s1 []byte
	x := len(s) * 3 % 4
	switch {
	case x == 2:
		s += "=="
	case x == 1:
		s += "="
	}
		if isurl==true{
		s1,_= base64.URLEncoding.DecodeString(s)
	}else{
		s1,_=base64.StdEncoding.DecodeString(s)
	}
	return string(s1)
}
/**
 * @Description: 解码emoji网页上显示
 * @param s
 * @return string
 */
func EmojiDecode(s string) string {
	re := regexp.MustCompile("\\[[\\\\u0-9a-zA-Z]+\\]")
	//提取emoji数据表达式
	reg := regexp.MustCompile("\\[\\\\u|]")
	src := re.FindAllString(s, -1)
	for i := 0; i < len(src); i++ {
		e := reg.ReplaceAllString(src[i], "")
		p, err := strconv.ParseInt(e, 16, 32)
		if err == nil {
			s = strings.Replace(s, src[i], string(rune(p)), -1)
		}
	}
	return s
}
/**
 * @Description: 编码emoji成unicode
 * @param s
 * @return string
 */
func EmojiEncode(s string) string {
	ret := ""
	rs := []rune(s)
	for i := 0; i < len(rs); i++ {
		if len(string(rs[i])) == 4 {
			u := `[\u` + strconv.FormatInt(int64(rs[i]), 16) + `]`
			ret += u

		} else {
			ret += string(rs[i])
		}
	}
	return ret
}
/**
 * @Description: 判断是否微信浏览器,必须启动http使用
 * @param r *http.Request
 * @return bool
 */
func IsWeixin(r *http.Request)bool{
	if strings.Index(r.UserAgent(),"icroMessenger")==-1{
		return false
	}else{
		return true
	}
}
/**
 * @Description: 隐藏手机中间四位或电话中间四位
 * @param phone 手机号
 * @return str
 */
func HideTel(phone string) (str string) {
	re := regexp.MustCompile(`0[0-9]{2,3}[-]?[2-9][0-9]{6,7}[-]?[0-9]?`)
	is := re.Match([]byte(phone))
	if is == true {
		re = regexp.MustCompile(`(0[0-9]{2,3}[-]?[2-9])[0-9]{3,4}([0-9]{3}[-]?[0-9]?)`)
		str = re.ReplaceAllString(phone, "$1****$2")
	} else {
		re = regexp.MustCompile(`(1[34578]{1}[0-9])[0-9]{4}([0-9]{4})`)
		str = re.ReplaceAllString(phone, "$1****$2")
	}
	return
}
/**
 * @Description: emoji编码成实体直接输出不需要转码
 * @param s
 * @return ss
 */
func Emoji(s string) (ss string) {
	s1 := strings.Split(s, "")
	for _, v := range s1 {
		if len(v) >= 4 {
			vv := []rune(v)
			k := int(vv[0])
			ss += "&#" + strconv.Itoa(k) + ";"
		} else {
			ss += v
		}
	}
	return
}
/**
 * @Description: 支持中英文字符串截取
 * @param s 原字符串
 * @param begin 开始位置
 * @param leng 长度
 * @return str
 */
func Cutstr(s string,begin int,leng int)(str string){
	s0 := []rune(s)
	l:=len(s0)
	if begin<0{
		begin=0
	}
	if begin>=l{
		begin=l
	}
	end:=begin+leng
	if end>l{
		end=l
	}
	str=string(s0[begin:end])
	return
}
/**
 * @Description: 对称加密解密函数
 * @param text 要加密或解密字符串
 * @param false解码 true加密
密钥 字符串
s:=php.Authcode("1234==+wo我们",true,"abc")
s=php.Authcode(s,false,"abc")
 * @return string
 */
func Authcode(text string, params ...interface{}) string {
	l := len(params)

	isEncode := false
	key := ""
	expiry := 0
	cKeyLen := 4

	if l > 0 {
		isEncode = params[0].(bool)
	}

	if l > 1 {
		key = params[1].(string)
	}

	if l > 2 {
		expiry = params[2].(int)
		if expiry < 0 {
			expiry = 0
		}
	}

	if l > 3 {
		cKeyLen = params[3].(int)
		if cKeyLen < 0 {
			cKeyLen = 0
		}
	}
	if cKeyLen > 32 {
		cKeyLen = 32
	}

	timestamp := time.Now().Unix()

	// md5加密key
	mKey := Md5(key)

	// 参与加密的
	keyA := Md5(mKey[0:16])
	// 用于验证数据有效性的
	keyB := Md5(mKey[16:])
	// 动态部分
	var keyC string
	if cKeyLen > 0 {
		if isEncode {
			// 加密的时候，动态获取一个秘钥
			keyC = Md5(fmt.Sprint(timestamp))[32-cKeyLen:]
		} else {
			// 解密的时候从头部获取动态秘钥部分
			keyC = text[0:cKeyLen]
		}
	}

	// 加入了动态的秘钥
	cryptKey := keyA + Md5(keyA+keyC)
	// 秘钥长度
	keyLen := len(cryptKey)
	if isEncode {
		// 加密 前10位是过期验证字符串 10-26位字符串验证
		var d int64
		if expiry > 0 {
			d = timestamp + int64(expiry)
		}
		text = fmt.Sprintf("%010d%s%s", d, Md5(text + keyB)[0:16], text)
	} else {
		// 解密
		text = string(Base64Decode(text[cKeyLen:],false))
	}

	// 字符串长度
	textLen := len(text)
	if textLen <= 0 {
		return ""
	}

	// 密匙簿
	box := RangeArray(0, 256)

	// 对称算法
	var rndKey []int
	cryptKeyB := []byte(cryptKey)
	for i := 0; i < 256; i++ {
		pos := i % keyLen
		rndKey = append(rndKey, int(cryptKeyB[pos]))
	}

	j := 0
	for i := 0; i < 256; i++ {
		j = (j + box[i] + rndKey[i]) % 256
		box[i], box[j] = box[j], box[i]
	}

	textB := []byte(text)
	a := 0
	j = 0
	var result []byte
	for i := 0; i < textLen; i++ {
		a = (a + 1) % 256
		j = (j + box[a]) % 256
		box[a], box[j] = box[j], box[a]
		result = append(result, byte(int(textB[i])^(box[(box[a]+box[j])%256])))
	}

	if isEncode {
		return keyC + strings.Replace(Base64Encode(string(result),false), "=", "", -1)
	}

	// 获取前10位，判断过期时间
	d, _ := strconv.ParseInt(string(result[0:10]), 10, 0)
	if (d == 0 || d-timestamp > 0) && string(result[10:26]) == Md5(string(result[26:]) + keyB)[0:16] {
		return string(result[26:])
	}

	return ""
}
/**
 * @Description: 生成序列数组
 * @param m 开始值
 * @param n 结束值
 * @return slice [0 1 2 3 4 5 6 7 8 9]
 */
func RangeArray(m, n int) (b []int) {
	if m >= n || m < 0 {
		return b
	}

	c := make([]int, 0, n-m)
	for i := m; i < n; i++ {
		c = append(c, i)
	}

	return c
}
/**
 * @Description: 生成随机颜色
 * @return string
 */
func RandColor()string{
	str:="abcdef0123456789"
	var s string
	for i:=0;i<6;i++{
		n:=MtRand(0,15)
		time.Sleep(1)
		s+=Cutstr(str,n,1)
	}
	return "#"+s
}
/**
 * @Description: 图片转换成base64
filename:="1.jpg"
s:=php.Img2Base64(filename)
 * @param filename
 * @return s
 */
func Img2Base64(filename string)(s string){
	ext:= filepath.Ext(filename)
	ext=strings.TrimLeft(ext,".")
	srcByte, _ := ioutil.ReadFile(filename)
	s=Base64Encode(string(srcByte),false)
	s="data:image/"+ext+";base64,"+s;
	return
}
/**
 * @Description: base64还原成图片
 * @param path 当前 . qq/ss 最后不要带/
 * @param data 上传的base64
 * @return ps 路径
 */
func Base642Img(path,data string)(ps string){
	re := regexp.MustCompile(`^(data:\s*image\/(\w+);base64,)`)
	r:=re.FindStringSubmatch(data)
	ext:=	r[2]
	bs:=strings.Replace(data,r[1],"",-1)
	CreateDir(path)
	ps=path+"/"+Sha1(data)+"."+ext
	bs = Base64Decode(bs,false)
	ioutil.WriteFile(ps, []byte(bs), 0666)
	return ps
}
/**
 * @Description: 中文字符串长度
 * @param str
 * @return int
 */
func StringLen(str string) int {
	return utf8.RuneCountInString(str)
}
/**
 * @Description: html字符串转换实体
 * @param s
 * @return string
 */
func HtmlEncode(s string)string{
	return html.EscapeString(s)
}
/**
 * @Description: html实体字符串还原
 * @param s
 * @return string
 */
func HtmlDecode(s string)string{
	return html.UnescapeString(s)
}
/**
 * @Description: 网址编码
 * @param str
 * @return string
 */
func UrlEncode(str string) string {
	return url.QueryEscape(str)
}

/**
 * @Description: 网址解码
 * @param str
 * @return string
 * @return error
 */
func UrlDecode(str string) (string, error) {
	return url.QueryUnescape(str)
}
//判断是否gbk编码,需要先判断是否utf8才可以
func Isgbk(s string) bool {
	if Isutf8(s){
		return false
	}
	data:=[]byte(s)
	length := len(data)
	var i int = 0
	for i < length {
		//fmt.Printf("for %x\n", data[i])
		if data[i] <= 0xff {
			//编码小于等于127,只有一个字节的编码，兼容ASCII吗
			i++
			continue
		} else {
			//大于127的使用双字节编码
			if  data[i] >= 0x81 &&
				data[i] <= 0xfe &&
				data[i + 1] >= 0x40 &&
				data[i + 1] <= 0xfe &&
				data[i + 1] != 0xf7 {
				i += 2
				continue
			} else {
				return false
			}
		}
	}
	return true
}
//是否utf8
func Isutf8(s string)bool{
	return utf8.ValidString(s)
}
/**
 * @Description: GET请求,支持gzip
 * @param u 网址
 * @return string
 */
func Get(u string)string{
	rs,_:=http.Get(u)
	defer rs.Body.Close()
	body, _ := io.ReadAll(rs.Body)
t:=string(body)
	if Isgbk(t){
		return Gbk2Utf8(t)
	}
	return t
}
//table转换成数组
func TableArr(s string)[][]string{
	re := regexp.MustCompile("<table[^>]*?>")
	s=re.ReplaceAllString(s, "")
	re = regexp.MustCompile("<tbody[^>]*?>")
	s=re.ReplaceAllString(s, "")
	re = regexp.MustCompile("<tr[^>]*?>")
	s=re.ReplaceAllString(s, "")
	re = regexp.MustCompile("<td[^>]*?>")
	s=re.ReplaceAllString(s, "")
	s=strings.Replace(s,"</tr>","{tr}",-1)
	s=strings.Replace(s,"</td>","{td}",-1)
	re = regexp.MustCompile("<[/!]*?[^<>]*?>")
	s=re.ReplaceAllString(s, "")
	re = regexp.MustCompile("([rn])[s]+")
	s=re.ReplaceAllString(s, "")
	re = regexp.MustCompile("&nbsp;")
	s=re.ReplaceAllString(s, "")
	re = regexp.MustCompile("</tbody>")
	s=re.ReplaceAllString(s, "")
	re = regexp.MustCompile("</table>")
	s=re.ReplaceAllString(s, "")
	re = regexp.MustCompile(`\s{2,}`)
	s=re.ReplaceAllString(s, "")
	s=strings.Replace(s," ","",-1)
	s=strings.Replace(s,"	","",-1)
	s=strings.Replace(s,"\r","",-1)
	s=strings.Replace(s,"\t","",-1)
	s=strings.Replace(s,"\n","",-1)
arr:=strings.Split(s,"{tr}")
arr=arr[:len(arr)-1]
var arr1 [][]string
	for _, v := range arr {
		arr2:=strings.Split(v,"{td}")
		arr2=arr2[:len(arr2)-1]
		arr1=append(arr1,arr2)
	}
	return arr1
}

var aes128=[]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
var aes192=[]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
}
var aes256=[]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
}
/**
 * @Description: Aes加密
 * @param text
 * @param key ,分别代表AES-128, AES-192和 AES-256
 * @return string
 * @return error
 */
func AesEn(text string, k string) (string, error) {
	var key []byte
	switch k {
	case "aes128":
		key=aes128
	case "aes192":
		key=aes192
	case "aes256":
		key=aes256
	}
	var iv = key[:aes.BlockSize]
	encrypted := make([]byte, len(text))
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	encrypter := cipher.NewCFBEncrypter(block, iv)
	encrypter.XORKeyStream(encrypted, []byte(text))
	return hex.EncodeToString(encrypted), nil
}
//aes解密 支持aes128 aes192 aes256
func AesDe(encrypted string, k string) (string, error) {
	var key []byte
	switch k {
	case "aes128":
		key=aes128
	case "aes192":
		key=aes192
	case "aes256":
		key=aes256
	}
	var err error
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	src, err := hex.DecodeString(encrypted)
	if err != nil {
		return "", err
	}
	var iv = key[:aes.BlockSize]
	decrypted := make([]byte, len(src))
	var block cipher.Block
	block, err = aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	decrypter := cipher.NewCFBDecrypter(block, iv)
	decrypter.XORKeyStream(decrypted, src)
	return string(decrypted), nil
}
//根据年月日判断星期 2021-03-17
func GetWeekday(ri string)string{
	var weekday = [7]string{"周日", "周一", "周二", "周三", "周四", "周五", "周六"}
	t:=Timestr2Time(ri+" 00:00:00")
	n:=int(t.Weekday())
	return weekday[n]
}
var dict=map[string]string{"腌":"yan","嗄":"a","迫":"po","捱":"ai","艾":"ai","瑷":"ai","嗌":"ai","犴":"an","鳌":"ao","廒":"ao","拗":"niu","岙":"ao","鏊":"ao","扒":"ba","岜":"ba","耙":"pa","鲅":"ba","癍":"ban","膀":"pang","磅":"bang","炮":"pao","曝":"pu","刨":"pao","瀑":"pu","陂":"bei","埤":"pi","鹎":"bei","邶":"bei","孛":"bei","鐾":"bei","鞴":"bei","畚":"ben","甏":"beng","舭":"bi","秘":"mi","辟":"pi","泌":"mi","裨":"bi","濞":"bi","庳":"bi","嬖":"bi","畀":"bi","筚":"bi","箅":"bi","襞":"bi","跸":"bi","笾":"bian","扁":"bian","碥":"bian","窆":"bian","便":"bian","弁":"bian","缏":"bian","骠":"biao","杓":"shao","飚":"biao","飑":"biao","瘭":"biao","髟":"biao","玢":"bin","豳":"bin","镔":"bin","膑":"bin","屏":"ping","泊":"bo","逋":"bu","晡":"bu","钸":"bu","醭":"bu","埔":"pu","瓿":"bu","礤":"ca","骖":"can","藏":"cang","艚":"cao","侧":"ce","喳":"zha","刹":"sha","瘥":"chai","禅":"chan","廛":"chan","镡":"tan","澶":"chan","躔":"chan","阊":"chang","鲳":"chang","长":"chang","苌":"chang","氅":"chang","鬯":"chang","焯":"chao","朝":"chao","车":"che","琛":"chen","谶":"chen","榇":"chen","蛏":"cheng","埕":"cheng","枨":"cheng","塍":"cheng","裎":"cheng","螭":"chi","眵":"chi","墀":"chi","篪":"chi","坻":"di","瘛":"chi","种":"zhong","重":"zhong","仇":"chou","帱":"chou","俦":"chou","雠":"chou","臭":"chou","楮":"chu","畜":"chu","嘬":"zuo","膪":"chuai","巛":"chuan","椎":"zhui","呲":"ci","兹":"zi","伺":"si","璁":"cong","楱":"cou","攒":"zan","爨":"cuan","隹":"zhui","榱":"cui","撮":"cuo","鹾":"cuo","嗒":"da","哒":"da","沓":"ta","骀":"tai","绐":"dai","埭":"dai","甙":"dai","弹":"dan","澹":"dan","叨":"dao","纛":"dao","簦":"deng","提":"ti","翟":"zhai","绨":"ti","丶":"dian","佃":"dian","簟":"dian","癜":"dian","调":"tiao","铞":"diao","佚":"yi","堞":"die","瓞":"die","揲":"die","垤":"die","疔":"ding","岽":"dong","硐":"dong","恫":"dong","垌":"dong","峒":"dong","芏":"du","煅":"duan","碓":"dui","镦":"dui","囤":"tun","铎":"duo","缍":"duo","驮":"tuo","沲":"tuo","柁":"tuo","哦":"o","恶":"e","轭":"e","锷":"e","鹗":"e","阏":"e","诶":"ea","鲕":"er","珥":"er","佴":"er","番":"fan","彷":"pang","霏":"fei","蜚":"fei","鲱":"fei","芾":"fei","瀵":"fen","鲼":"fen","否":"fou","趺":"fu","桴":"fu","莩":"fu","菔":"fu","幞":"fu","郛":"fu","绂":"fu","绋":"fu","祓":"fu","砩":"fu","黻":"fu","罘":"fu","蚨":"fu","脯":"pu","滏":"fu","黼":"fu","鲋":"fu","鳆":"fu","咖":"ka","噶":"ga","轧":"zha","陔":"gai","戤":"gai","扛":"kang","戆":"gang","筻":"gang","槔":"gao","藁":"gao","缟":"gao","咯":"ge","仡":"yi","搿":"ge","塥":"ge","鬲":"ge","哿":"ge","句":"ju","缑":"gou","鞲":"gou","笱":"gou","遘":"gou","瞽":"gu","罟":"gu","嘏":"gu","牿":"gu","鲴":"gu","栝":"kuo","莞":"guan","纶":"lun","涫":"guan","涡":"wo","呙":"guo","馘":"guo","猓":"guo","咳":"ke","氦":"hai","颔":"han","吭":"keng","颃":"hang","巷":"xiang","蚵":"ke","翮":"he","吓":"xia","桁":"heng","泓":"hong","蕻":"hong","黉":"hong","後":"hou","唿":"hu","煳":"hu","浒":"hu","祜":"hu","岵":"hu","鬟":"huan","圜":"huan","郇":"xun","锾":"huan","逭":"huan","咴":"hui","虺":"hui","会":"hui","溃":"kui","哕":"hui","缋":"hui","锪":"huo","蠖":"huo","缉":"ji","稽":"ji","赍":"ji","丌":"ji","咭":"ji","亟":"ji","殛":"ji","戢":"ji","嵴":"ji","蕺":"ji","系":"xi","蓟":"ji","霁":"ji","荠":"qi","跽":"ji","哜":"ji","鲚":"ji","洎":"ji","芰":"ji","茄":"qie","珈":"jia","迦":"jia","笳":"jia","葭":"jia","跏":"jia","郏":"jia","恝":"jia","铗":"jia","袷":"qia","蛱":"jia","角":"jiao","挢":"jiao","岬":"jia","徼":"jiao","湫":"qiu","敫":"jiao","瘕":"jia","浅":"qian","蒹":"jian","搛":"jian","湔":"jian","缣":"jian","犍":"jian","鹣":"jian","鲣":"jian","鞯":"jian","蹇":"jian","謇":"jian","硷":"jian","枧":"jian","戬":"jian","谫":"jian","囝":"jian","裥":"jian","笕":"jian","翦":"jian","趼":"jian","楗":"jian","牮":"jian","踺":"jian","茳":"jiang","礓":"jiang","耩":"jiang","降":"jiang","绛":"jiang","洚":"jiang","鲛":"jiao","僬":"jiao","鹪":"jiao","艽":"jiao","茭":"jiao","嚼":"jiao","峤":"qiao","觉":"jiao","校":"xiao","噍":"jiao","醮":"jiao","疖":"jie","喈":"jie","桔":"ju","拮":"jie","桀":"jie","颉":"jie","婕":"jie","羯":"jie","鲒":"jie","蚧":"jie","骱":"jie","衿":"jin","馑":"jin","卺":"jin","廑":"jin","堇":"jin","槿":"jin","靳":"jin","缙":"jin","荩":"jin","赆":"jin","妗":"jin","旌":"jing","腈":"jing","憬":"jing","肼":"jing","迳":"jing","胫":"jing","弪":"jing","獍":"jing","扃":"jiong","鬏":"jiu","疚":"jiu","僦":"jiu","桕":"jiu","疽":"ju","裾":"ju","苴":"ju","椐":"ju","锔":"ju","琚":"ju","鞫":"ju","踽":"ju","榉":"ju","莒":"ju","遽":"ju","倨":"ju","钜":"ju","犋":"ju","屦":"ju","榘":"ju","窭":"ju","讵":"ju","醵":"ju","苣":"ju","圈":"quan","镌":"juan","蠲":"juan","锩":"juan","狷":"juan","桊":"juan","鄄":"juan","獗":"jue","攫":"jue","孓":"jue","橛":"jue","珏":"jue","桷":"jue","劂":"jue","爝":"jue","镢":"jue","觖":"jue","筠":"jun","麇":"jun","捃":"jun","浚":"jun","喀":"ka","卡":"ka","佧":"ka","胩":"ka","锎":"kai","蒈":"kai","剀":"kai","垲":"kai","锴":"kai","戡":"kan","莰":"kan","闶":"kang","钪":"kang","尻":"kao","栲":"kao","柯":"ke","疴":"ke","钶":"ke","颏":"ke","珂":"ke","髁":"ke","壳":"ke","岢":"ke","溘":"ke","骒":"ke","缂":"ke","氪":"ke","锞":"ke","裉":"ken","倥":"kong","崆":"kong","箜":"kong","芤":"kou","眍":"kou","筘":"kou","刳":"ku","堀":"ku","喾":"ku","侉":"kua","蒯":"kuai","哙":"kuai","狯":"kuai","郐":"kuai","匡":"kuang","夼":"kuang","邝":"kuang","圹":"kuang","纩":"kuang","贶":"kuang","岿":"kui","悝":"kui","睽":"kui","逵":"kui","馗":"kui","夔":"kui","喹":"kui","隗":"wei","暌":"kui","揆":"kui","蝰":"kui","跬":"kui","喟":"kui","聩":"kui","篑":"kui","蒉":"kui","愦":"kui","锟":"kun","醌":"kun","琨":"kun","髡":"kun","悃":"kun","阃":"kun","蛞":"kuo","砬":"la","落":"luo","剌":"la","瘌":"la","涞":"lai","崃":"lai","铼":"lai","赉":"lai","濑":"lai","斓":"lan","镧":"lan","谰":"lan","漤":"lan","罱":"lan","稂":"lang","阆":"lang","莨":"liang","蒗":"lang","铹":"lao","痨":"lao","醪":"lao","栳":"lao","铑":"lao","耢":"lao","勒":"le","仂":"le","叻":"le","泐":"le","鳓":"le","了":"le","镭":"lei","嫘":"lei","缧":"lei","檑":"lei","诔":"lei","耒":"lei","酹":"lei","塄":"leng","愣":"leng","藜":"li","骊":"li","黧":"li","缡":"li","嫠":"li","鲡":"li","蓠":"li","澧":"li","锂":"li","醴":"li","鳢":"li","俪":"li","砺":"li","郦":"li","詈":"li","猁":"li","溧":"li","栎":"li","轹":"li","傈":"li","坜":"li","苈":"li","疠":"li","疬":"li","篥":"li","粝":"li","跞":"li","俩":"liang","裢":"lian","濂":"lian","臁":"lian","奁":"lian","蠊":"lian","琏":"lian","蔹":"lian","裣":"lian","楝":"lian","潋":"lian","椋":"liang","墚":"liang","寮":"liao","鹩":"liao","蓼":"liao","钌":"liao","廖":"liao","尥":"liao","洌":"lie","捩":"lie","埒":"lie","躐":"lie","鬣":"lie","辚":"lin","遴":"lin","啉":"lin","瞵":"lin","懔":"lin","廪":"lin","蔺":"lin","膦":"lin","酃":"ling","柃":"ling","鲮":"ling","呤":"ling","镏":"liu","旒":"liu","骝":"liu","鎏":"liu","锍":"liu","碌":"lu","鹨":"liu","茏":"long","栊":"long","泷":"long","砻":"long","癃":"long","垅":"long","偻":"lou","蝼":"lou","蒌":"lou","耧":"lou","嵝":"lou","露":"lu","瘘":"lou","噜":"lu","轳":"lu","垆":"lu","胪":"lu","舻":"lu","栌":"lu","镥":"lu","绿":"lv","辘":"lu","簏":"lu","潞":"lu","辂":"lu","渌":"lu","氇":"lu","捋":"lv","稆":"lv","率":"lv","闾":"lv","栾":"luan","銮":"luan","滦":"luan","娈":"luan","脔":"luan","锊":"lve","猡":"luo","椤":"luo","脶":"luo","镙":"luo","倮":"luo","蠃":"luo","瘰":"luo","珞":"luo","泺":"luo","荦":"luo","雒":"luo","呒":"mu","抹":"mo","唛":"mai","杩":"ma","么":"me","埋":"mai","荬":"mai","脉":"mai","劢":"mai","颟":"man","蔓":"man","鳗":"man","鞔":"man","螨":"man","墁":"man","缦":"man","熳":"man","镘":"man","邙":"mang","硭":"mang","旄":"mao","茆":"mao","峁":"mao","泖":"mao","昴":"mao","耄":"mao","瑁":"mao","懋":"mao","瞀":"mao","麽":"me","没":"mei","嵋":"mei","湄":"mei","猸":"mei","镅":"mei","鹛":"mei","浼":"mei","钔":"men","瞢":"meng","甍":"meng","礞":"meng","艨":"meng","黾":"mian","鳘":"min","溟":"ming","暝":"ming","模":"mo","谟":"mo","嫫":"mo","镆":"mo","瘼":"mo","耱":"mo","貊":"mo","貘":"mo","牟":"mou","鍪":"mou","蛑":"mou","侔":"mou","毪":"mu","坶":"mu","仫":"mu","唔":"wu","那":"na","镎":"na","哪":"na","呢":"ne","肭":"na","艿":"nai","鼐":"nai","萘":"nai","柰":"nai","蝻":"nan","馕":"nang","攮":"nang","曩":"nang","猱":"nao","铙":"nao","硇":"nao","蛲":"nao","垴":"nao","坭":"ni","猊":"ni","铌":"ni","鲵":"ni","祢":"mi","睨":"ni","慝":"te","伲":"ni","鲇":"nian","鲶":"nian","埝":"nian","嬲":"niao","茑":"niao","脲":"niao","啮":"nie","陧":"nie","颞":"nie","臬":"nie","蘖":"nie","甯":"ning","聍":"ning","狃":"niu","侬":"nong","耨":"nou","孥":"nu","胬":"nu","钕":"nv","恧":"nv","褰":"qian","掮":"qian","荨":"xun","钤":"qian","箝":"qian","鬈":"quan","缱":"qian","肷":"qian","纤":"xian","茜":"qian","慊":"qian","椠":"qian","戗":"qiang","镪":"qiang","锖":"qiang","樯":"qiang","嫱":"qiang","雀":"que","缲":"qiao","硗":"qiao","劁":"qiao","樵":"qiao","谯":"qiao","鞒":"qiao","愀":"qiao","鞘":"qiao","郄":"xi","箧":"qie","亲":"qin","覃":"tan","溱":"qin","檎":"qin","锓":"qin","嗪":"qin","螓":"qin","揿":"qin","吣":"qin","圊":"qing","鲭":"qing","檠":"qing","黥":"qing","謦":"qing","苘":"qing","磬":"qing","箐":"qing","綮":"qi","茕":"qiong","邛":"dao","蛩":"tun","筇":"qiong","跫":"qiong","銎":"qiong","楸":"qiu","俅":"qiu","赇":"qiu","逑":"qiu","犰":"qiu","蝤":"qiu","巯":"qiu","鼽":"qiu","糗":"qiu","区":"qu","祛":"qu","麴":"qu","诎":"qu","衢":"qu","癯":"qu","劬":"qu","璩":"qu","氍":"qu","朐":"qu","磲":"qu","鸲":"qu","蕖":"qu","蠼":"qu","蘧":"qu","阒":"qu","颧":"quan","荃":"quan","铨":"quan","辁":"quan","筌":"quan","绻":"quan","畎":"quan","阕":"que","悫":"que","髯":"ran","禳":"rang","穰":"rang","仞":"ren","妊":"ren","轫":"ren","衽":"ren","狨":"rong","肜":"rong","蝾":"rong","嚅":"ru","濡":"ru","薷":"ru","襦":"ru","颥":"ru","洳":"ru","溽":"ru","蓐":"ru","朊":"ruan","蕤":"rui","枘":"rui","箬":"ruo","挲":"suo","脎":"sa","塞":"sai","鳃":"sai","噻":"sai","毵":"san","馓":"san","糁":"san","霰":"xian","磉":"sang","颡":"sang","缫":"sao","鳋":"sao","埽":"sao","瘙":"sao","色":"se","杉":"shan","鲨":"sha","痧":"sha","裟":"sha","铩":"sha","唼":"sha","酾":"shai","栅":"zha","跚":"shan","芟":"shan","埏":"shan","钐":"shan","舢":"shan","剡":"yan","鄯":"shan","疝":"shan","蟮":"shan","墒":"shang","垧":"shang","绱":"shang","蛸":"shao","筲":"shao","苕":"tiao","召":"zhao","劭":"shao","猞":"she","畲":"she","折":"zhe","滠":"she","歙":"xi","厍":"she","莘":"shen","娠":"shen","诜":"shen","什":"shen","谂":"shen","渖":"shen","矧":"shen","胂":"shen","椹":"shen","省":"sheng","眚":"sheng","嵊":"sheng","嘘":"xu","蓍":"shi","鲺":"shi","识":"shi","拾":"shi","埘":"shi","莳":"shi","炻":"shi","鲥":"shi","豕":"shi","似":"si","噬":"shi","贳":"shi","铈":"shi","螫":"shi","筮":"shi","殖":"zhi","熟":"shu","艏":"shou","菽":"shu","摅":"shu","纾":"shu","毹":"shu","疋":"shu","数":"shu","属":"shu","术":"shu","澍":"shu","沭":"shu","丨":"shu","腧":"shu","说":"shuo","妁":"shuo","蒴":"shuo","槊":"shuo","搠":"shuo","鸶":"si","澌":"si","缌":"si","锶":"si","厶":"si","蛳":"si","驷":"si","泗":"si","汜":"si","兕":"si","姒":"si","耜":"si","笥":"si","忪":"song","淞":"song","崧":"song","凇":"song","菘":"song","竦":"song","溲":"sou","飕":"sou","蜩":"tiao","萜":"tie","汀":"ting","葶":"ting","莛":"ting","梃":"ting","佟":"tong","酮":"tong","仝":"tong","茼":"tong","砼":"tong","钭":"dou","酴":"tu","钍":"tu","堍":"tu","抟":"tuan","忒":"te","煺":"tui","暾":"tun","氽":"tun","乇":"tuo","砣":"tuo","沱":"tuo","跎":"tuo","坨":"tuo","橐":"tuo","酡":"tuo","鼍":"tuo","庹":"tuo","拓":"tuo","柝":"tuo","箨":"tuo","腽":"wa","崴":"wai","芄":"wan","畹":"wan","琬":"wan","脘":"wan","菀":"wan","尢":"you","辋":"wang","魍":"wang","逶":"wei","葳":"wei","隈":"wei","惟":"wei","帏":"wei","圩":"wei","囗":"wei","潍":"wei","嵬":"wei","沩":"wei","涠":"wei","尾":"wei","玮":"wei","炜":"wei","韪":"wei","洧":"wei","艉":"wei","鲔":"wei","遗":"yi","尉":"wei","軎":"wei","璺":"wen","阌":"wen","蓊":"weng","蕹":"weng","渥":"wo","硪":"wo","龌":"wo","圬":"wu","吾":"wu","浯":"wu","鼯":"wu","牾":"wu","迕":"wu","庑":"wu","痦":"wu","芴":"wu","杌":"wu","焐":"wu","阢":"wu","婺":"wu","鋈":"wu","樨":"xi","栖":"qi","郗":"xi","蹊":"qi","淅":"xi","熹":"xi","浠":"xi","僖":"xi","穸":"xi","螅":"xi","菥":"xi","舾":"xi","矽":"xi","粞":"xi","硒":"xi","醯":"xi","欷":"xi","鼷":"xi","檄":"xi","隰":"xi","觋":"xi","屣":"xi","葸":"xi","蓰":"xi","铣":"xi","饩":"xi","阋":"xi","禊":"xi","舄":"xi","狎":"xia","硖":"xia","柙":"xia","暹":"xian","莶":"xian","祆":"xian","籼":"xian","跹":"xian","鹇":"xian","痫":"xian","猃":"xian","燹":"xian","蚬":"xian","筅":"xian","冼":"xian","岘":"xian","骧":"xiang","葙":"xiang","芗":"xiang","缃":"xiang","庠":"xiang","鲞":"xiang","蟓":"xiang","削":"xue","枵":"xiao","绡":"xiao","筱":"xiao","邪":"xie","勰":"xie","缬":"xie","血":"xue","榭":"xie","瀣":"xie","薤":"xie","燮":"xie","躞":"xie","廨":"xie","绁":"xie","渫":"xie","榍":"xie","獬":"xie","昕":"xin","忻":"xin","囟":"xin","陉":"jing","荥":"ying","饧":"tang","硎":"xing","荇":"xing","芎":"xiong","馐":"xiu","庥":"xiu","鸺":"xiu","貅":"xiu","髹":"xiu","宿":"xiu","岫":"xiu","溴":"xiu","吁":"xu","盱":"xu","顼":"xu","糈":"xu","醑":"xu","洫":"xu","溆":"xu","蓿":"xu","萱":"xuan","谖":"xuan","儇":"xuan","煊":"xuan","痃":"xuan","铉":"xuan","泫":"xuan","碹":"xuan","楦":"xuan","镟":"xuan","踅":"xue","泶":"xue","鳕":"xue","埙":"xun","曛":"xun","窨":"xun","獯":"xun","峋":"xun","洵":"xun","恂":"xun","浔":"xun","鲟":"xun","蕈":"xun","垭":"ya","岈":"ya","琊":"ya","痖":"ya","迓":"ya","砑":"ya","咽":"yan","鄢":"yan","菸":"yan","崦":"yan","铅":"qian","芫":"yuan","兖":"yan","琰":"yan","罨":"yan","厣":"yan","焱":"yan","酽":"yan","谳":"yan","鞅":"yang","炀":"yang","蛘":"yang","约":"yue","珧":"yao","轺":"yao","繇":"yao","鳐":"yao","崾":"yao","钥":"yao","曜":"yao","铘":"ye","烨":"ye","邺":"ye","靥":"ye","晔":"ye","猗":"yi","铱":"yi","欹":"qi","黟":"yi","怡":"yi","沂":"yi","圯":"yi","荑":"yi","诒":"yi","眙":"yi","嶷":"yi","钇":"yi","舣":"yi","酏":"yi","熠":"yi","弋":"yi","懿":"yi","镒":"yi","峄":"yi","怿":"yi","悒":"yi","佾":"yi","殪":"yi","挹":"yi","埸":"yi","劓":"yi","镱":"yi","瘗":"yi","癔":"yi","翊":"yi","蜴":"yi","氤":"yin","堙":"yin","洇":"yin","鄞":"yin","狺":"yin","夤":"yin","圻":"qi","饮":"yin","吲":"yin","胤":"yin","茚":"yin","璎":"ying","撄":"ying","嬴":"ying","滢":"ying","潆":"ying","蓥":"ying","瘿":"ying","郢":"ying","媵":"ying","邕":"yong","镛":"yong","墉":"yong","慵":"yong","痈":"yong","鳙":"yong","饔":"yong","喁":"yong","俑":"yong","莸":"you","猷":"you","疣":"you","蚰":"you","蝣":"you","莜":"you","牖":"you","铕":"you","卣":"you","宥":"you","侑":"you","蚴":"you","釉":"you","馀":"yu","萸":"yu","禺":"yu","妤":"yu","欤":"yu","觎":"yu","窬":"yu","蝓":"yu","嵛":"yu","舁":"yu","雩":"yu","龉":"yu","伛":"yu","圉":"yu","庾":"yu","瘐":"yu","窳":"yu","俣":"yu","毓":"yu","峪":"yu","煜":"yu","燠":"yu","蓣":"yu","饫":"yu","阈":"yu","鬻":"yu","聿":"yu","钰":"yu","鹆":"yu","蜮":"yu","眢":"yuan","箢":"yuan","员":"yuan","沅":"yuan","橼":"yuan","塬":"yuan","爰":"yuan","螈":"yuan","鼋":"yuan","掾":"yuan","垸":"yuan","瑗":"yuan","刖":"yue","瀹":"yue","樾":"yue","龠":"yue","氲":"yun","昀":"yun","郧":"yun","狁":"yun","郓":"yun","韫":"yun","恽":"yun","扎":"zha","拶":"za","咋":"za","仔":"zai","昝":"zan","瓒":"zan","奘":"zang","唣":"zao","择":"ze","迮":"ze","赜":"ze","笮":"ze","箦":"ze","舴":"ze","昃":"ze","缯":"zeng","罾":"zeng","齄":"zha","柞":"zha","痄":"zha","瘵":"zhai","旃":"zhan","璋":"zhang","漳":"zhang","嫜":"zhang","鄣":"zhang","仉":"zhang","幛":"zhang","着":"zhe","啁":"zhou","爪":"zhao","棹":"zhao","笊":"zhao","摺":"zhe","磔":"zhe","这":"zhe","柘":"zhe","桢":"zhen","蓁":"zhen","祯":"zhen","浈":"zhen","畛":"zhen","轸":"zhen","稹":"zhen","圳":"zhen","徵":"zhi","钲":"zheng","卮":"zhi","胝":"zhi","祗":"zhi","摭":"zhi","絷":"zhi","埴":"zhi","轵":"zhi","黹":"zhi","帙":"zhi","轾":"zhi","贽":"zhi","陟":"zhi","忮":"zhi","彘":"zhi","膣":"zhi","鸷":"zhi","骘":"zhi","踬":"zhi","郅":"zhi","觯":"zhi","锺":"zhong","螽":"zhong","舯":"zhong","碡":"zhou","绉":"zhou","荮":"zhou","籀":"zhou","酎":"zhou","洙":"zhu","邾":"zhu","潴":"zhu","槠":"zhu","橥":"zhu","舳":"zhu","瘃":"zhu","渚":"zhu","麈":"zhu","箸":"zhu","炷":"zhu","杼":"zhu","翥":"zhu","疰":"zhu","颛":"zhuan","赚":"zhuan","馔":"zhuan","僮":"tong","缒":"zhui","肫":"zhun","窀":"zhun","涿":"zhuo","倬":"zhuo","濯":"zhuo","诼":"zhuo","禚":"zhuo","浞":"zhuo","谘":"zi","淄":"zi","髭":"zi","孳":"zi","粢":"zi","趑":"zi","觜":"zui","缁":"zi","鲻":"zi","嵫":"zi","笫":"zi","耔":"zi","腙":"zong","偬":"zong","诹":"zou","陬":"zou","鄹":"zou","驺":"zou","鲰":"zou","菹":"ju","镞":"zu","躜":"zuan","缵":"zuan","蕞":"zui","撙":"zun","胙":"zuo","阿":"a","柏":"bai","蚌":"beng","薄":"bo","堡":"bao","呗":"bei","贲":"ben","臂":"bi","瘪":"bie","槟":"bin","剥":"bo","伯":"bo","卜":"bu","参":"can","嚓":"ca","差":"cha","孱":"chan","绰":"chuo","称":"cheng","澄":"cheng","大":"da","单":"dan","得":"de","的":"de","地":"di","都":"dou","读":"du","度":"du","蹲":"dun","佛":"fo","伽":"jia","盖":"gai","镐":"hao","给":"gei","呱":"gua","氿":"jiu","桧":"hui","掴":"guo","蛤":"ha","还":"hai","和":"he","核":"he","哼":"heng","鹄":"hu","划":"hua","夹":"jia","贾":"jia","芥":"jie","劲":"jin","荆":"jing","颈":"jing","貉":"he","吖":"a","啊":"a","锕":"a","哎":"ai","哀":"ai","埃":"ai","唉":"ai","欸":"ai","锿":"ai","挨":"ai","皑":"ai","癌":"ai","毐":"ai","矮":"ai","蔼":"ai","霭":"ai","砹":"ai","爱":"ai","隘":"ai","碍":"ai","嗳":"ai","嫒":"ai","叆":"ai","暧":"ai","安":"an","桉":"an","氨":"an","庵":"an","谙":"an","鹌":"an","鞍":"an","俺":"an","埯":"an","唵":"an","铵":"an","揞":"an","岸":"an","按":"an","胺":"an","案":"an","暗":"an","黯":"an","玵":"an","肮":"ang","昂":"ang","盎":"ang","凹":"ao","敖":"ao","遨":"ao","嗷":"ao","獒":"ao","熬":"ao","聱":"ao","螯":"ao","翱":"ao","謷":"ao","鏖":"ao","袄":"ao","媪":"ao","坳":"ao","傲":"ao","奥":"ao","骜":"ao","澳":"ao","懊":"ao","八":"ba","巴":"ba","叭":"ba","芭":"ba","疤":"ba","捌":"ba","笆":"ba","粑":"ba","拔":"ba","茇":"ba","妭":"ba","菝":"ba","跋":"ba","魃":"ba","把":"ba","靶":"ba","坝":"ba","爸":"ba","罢":"ba","霸":"ba","灞":"ba","吧":"ba","钯":"ba","掰":"bai","白":"bai","百":"bai","佰":"bai","捭":"bai","摆":"bai","败":"bai","拜":"bai","稗":"bai","扳":"ban","攽":"ban","班":"ban","般":"ban","颁":"ban","斑":"ban","搬":"ban","瘢":"ban","阪":"ban","坂":"ban","板":"ban","版":"ban","钣":"ban","舨":"ban","办":"ban","半":"ban","伴":"ban","拌":"ban","绊":"ban","瓣":"ban","扮":"ban","邦":"bang","帮":"bang","梆":"bang","浜":"bang","绑":"bang","榜":"bang","棒":"bang","傍":"bang","谤":"bang","蒡":"bang","镑":"bang","包":"bao","苞":"bao","孢":"bao","胞":"bao","龅":"bao","煲":"bao","褒":"bao","雹":"bao","饱":"bao","宝":"bao","保":"bao","鸨":"bao","葆":"bao","褓":"bao","报":"bao","抱":"bao","趵":"bao","豹":"bao","鲍":"bao","暴":"bao","爆":"bao","枹":"bao","杯":"bei","卑":"bei","悲":"bei","碑":"bei","北":"bei","贝":"bei","狈":"bei","备":"bei","背":"bei","钡":"bei","倍":"bei","悖":"bei","被":"bei","辈":"bei","惫":"bei","焙":"bei","蓓":"bei","碚":"bei","褙":"bei","别":"bei","蹩":"bei","椑":"bei","奔":"ben","倴":"ben","犇":"ben","锛":"ben","本":"ben","苯":"ben","坌":"ben","笨":"ben","崩":"beng","绷":"beng","嘣":"beng","甭":"beng","泵":"beng","迸":"beng","镚":"beng","蹦":"beng","屄":"bi","逼":"bi","荸":"bi","鼻":"bi","匕":"bi","比":"bi","吡":"bi","沘":"bi","妣":"bi","彼":"bi","秕":"bi","笔":"bi","俾":"bi","鄙":"bi","币":"bi","必":"bi","毕":"bi","闭":"bi","庇":"bi","诐":"bi","苾":"bi","荜":"bi","毖":"bi","哔":"bi","陛":"bi","毙":"bi","铋":"bi","狴":"bi","萆":"bi","梐":"bi","敝":"bi","婢":"bi","赑":"bi","愎":"bi","弼":"bi","蓖":"bi","痹":"bi","滗":"bi","碧":"bi","蔽":"bi","馝":"bi","弊":"bi","薜":"bi","篦":"bi","壁":"bi","避":"bi","髀":"bi","璧":"bi","芘":"bi","边":"bian","砭":"bian","萹":"bian","编":"bian","煸":"bian","蝙":"bian","鳊":"bian","鞭":"bian","贬":"bian","匾":"bian","褊":"bian","藊":"bian","卞":"bian","抃":"bian","苄":"bian","汴":"bian","忭":"bian","变":"bian","遍":"bian","辨":"bian","辩":"bian","辫":"bian","标":"biao","骉":"biao","彪":"biao","摽":"biao","膘":"biao","飙":"biao","镖":"biao","瀌":"biao","镳":"biao","表":"biao","婊":"biao","裱":"biao","鳔":"biao","憋":"bie","鳖":"bie","宾":"bin","彬":"bin","傧":"bin","滨":"bin","缤":"bin","濒":"bin","摈":"bin","殡":"bin","髌":"bin","鬓":"bin","冰":"bing","兵":"bing","丙":"bing","邴":"bing","秉":"bing","柄":"bing","饼":"bing","炳":"bing","禀":"bing","并":"bing","病":"bing","摒":"bing","拨":"bo","波":"bo","玻":"bo","钵":"bo","饽":"bo","袯":"bo","菠":"bo","播":"bo","驳":"bo","帛":"bo","勃":"bo","钹":"bo","铂":"bo","亳":"bo","舶":"bo","脖":"bo","博":"bo","鹁":"bo","渤":"bo","搏":"bo","馎":"bo","箔":"bo","膊":"bo","踣":"bo","馞":"bo","礴":"bo","跛":"bo","檗":"bo","擘":"bo","簸":"bo","啵":"bo","蕃":"bo","哱":"bo","卟":"bu","补":"bu","捕":"bu","哺":"bu","不":"bu","布":"bu","步":"bu","怖":"bu","钚":"bu","部":"bu","埠":"bu","簿":"bu","擦":"ca","猜":"cai","才":"cai","材":"cai","财":"cai","裁":"cai","采":"cai","彩":"cai","睬":"cai","踩":"cai","菜":"cai","蔡":"cai","餐":"can","残":"can","蚕":"can","惭":"can","惨":"can","黪":"can","灿":"can","粲":"can","璨":"can","穇":"can","仓":"cang","伧":"cang","苍":"cang","沧":"cang","舱":"cang","操":"cao","糙":"cao","曹":"cao","嘈":"cao","漕":"cao","槽":"cao","螬":"cao","草":"cao","册":"ce","厕":"ce","测":"ce","恻":"ce","策":"ce","岑":"cen","涔":"cen","噌":"ceng","层":"ceng","嶒":"ceng","蹭":"ceng","叉":"cha","杈":"cha","插":"cha","馇":"cha","锸":"cha","茬":"cha","茶":"cha","搽":"cha","嵖":"cha","猹":"cha","槎":"cha","碴":"cha","察":"cha","檫":"cha","衩":"cha","镲":"cha","汊":"cha","岔":"cha","侘":"cha","诧":"cha","姹":"cha","蹅":"cha","拆":"chai","钗":"chai","侪":"chai","柴":"chai","豺":"chai","虿":"chai","茝":"chai","觇":"chan","掺":"chan","搀":"chan","襜":"chan","谗":"chan","婵":"chan","馋":"chan","缠":"chan","蝉":"chan","潺":"chan","蟾":"chan","巉":"chan","产":"chan","浐":"chan","谄":"chan","铲":"chan","阐":"chan","蒇":"chan","骣":"chan","冁":"chan","忏":"chan","颤":"chan","羼":"chan","韂":"chan","伥":"chang","昌":"chang","菖":"chang","猖":"chang","娼":"chang","肠":"chang","尝":"chang","常":"chang","偿":"chang","徜":"chang","嫦":"chang","厂":"chang","场":"chang","昶":"chang","惝":"chang","敞":"chang","怅":"chang","畅":"chang","倡":"chang","唱":"chang","裳":"chang","抄":"chao","怊":"chao","钞":"chao","超":"chao","晁":"chao","巢":"chao","嘲":"chao","潮":"chao","吵":"chao","炒":"chao","耖":"chao","砗":"che","扯":"che","彻":"che","坼":"che","掣":"che","撤":"che","澈":"che","瞮":"che","抻":"chen","郴":"chen","嗔":"chen","瞋":"chen","臣":"chen","尘":"chen","辰":"chen","沉":"chen","忱":"chen","陈":"chen","宸":"chen","晨":"chen","谌":"chen","碜":"chen","衬":"chen","龀":"chen","趁":"chen","柽":"cheng","琤":"cheng","撑":"cheng","瞠":"cheng","成":"cheng","丞":"cheng","呈":"cheng","诚":"cheng","承":"cheng","城":"cheng","铖":"cheng","程":"cheng","惩":"cheng","酲":"cheng","橙":"cheng","逞":"cheng","骋":"cheng","秤":"cheng","铛":"cheng","樘":"cheng","吃":"chi","哧":"chi","鸱":"chi","蚩":"chi","笞":"chi","嗤":"chi","痴":"chi","媸":"chi","魑":"chi","池":"chi","弛":"chi","驰":"chi","迟":"chi","茌":"chi","持":"chi","踟":"chi","尺":"chi","齿":"chi","侈":"chi","耻":"chi","豉":"chi","褫":"chi","彳":"chi","叱":"chi","斥":"chi","赤":"chi","饬":"chi","炽":"chi","翅":"chi","敕":"chi","啻":"chi","傺":"chi","匙":"chi","冲":"chong","充":"chong","忡":"chong","茺":"chong","舂":"chong","憧":"chong","艟":"chong","虫":"chong","崇":"chong","宠":"chong","铳":"chong","抽":"chou","瘳":"chou","惆":"chou","绸":"chou","畴":"chou","酬":"chou","稠":"chou","愁":"chou","筹":"chou","踌":"chou","丑":"chou","瞅":"chou","出":"chu","初":"chu","樗":"chu","刍":"chu","除":"chu","厨":"chu","锄":"chu","滁":"chu","蜍":"chu","雏":"chu","橱":"chu","躇":"chu","蹰":"chu","杵":"chu","础":"chu","储":"chu","楚":"chu","褚":"chu","亍":"chu","处":"chu","怵":"chu","绌":"chu","搐":"chu","触":"chu","憷":"chu","黜":"chu","矗":"chu","揣":"chuai","搋":"chuai","膗":"chuai","踹":"chuai","川":"chuan","氚":"chuan","穿":"chuan","舡":"chuan","船":"chuan","遄":"chuan","椽":"chuan","舛":"chuan","喘":"chuan","串":"chuan","钏":"chuan","疮":"chuang","窗":"chuang","床":"chuang","闯":"chuang","创":"chuang","怆":"chuang","吹":"chui","炊":"chui","垂":"chui","陲":"chui","捶":"chui","棰":"chui","槌":"chui","锤":"chui","春":"chun","瑃":"chun","椿":"chun","蝽":"chun","纯":"chun","莼":"chun","唇":"chun","淳":"chun","鹑":"chun","醇":"chun","蠢":"chun","踔":"chuo","戳":"chuo","啜":"chuo","惙":"chuo","辍":"chuo","龊":"chuo","歠":"chuo","疵":"ci","词":"ci","茈":"ci","茨":"ci","祠":"ci","瓷":"ci","辞":"ci","慈":"ci","磁":"ci","雌":"ci","鹚":"ci","糍":"ci","此":"ci","泚":"ci","跐":"ci","次":"ci","刺":"ci","佽":"ci","赐":"ci","匆":"cong","苁":"cong","囱":"cong","枞":"cong","葱":"cong","骢":"cong","聪":"cong","从":"cong","丛":"cong","淙":"cong","悰":"cong","琮":"cong","凑":"cou","辏":"cou","腠":"cou","粗":"cu","徂":"cu","殂":"cu","促":"cu","猝":"cu","蔟":"cu","醋":"cu","踧":"cu","簇":"cu","蹙":"cu","蹴":"cu","汆":"cuan","撺":"cuan","镩":"cuan","蹿":"cuan","窜":"cuan","篡":"cuan","崔":"cui","催":"cui","摧":"cui","璀":"cui","脆":"cui","萃":"cui","啐":"cui","淬":"cui","悴":"cui","毳":"cui","瘁":"cui","粹":"cui","翠":"cui","村":"cun","皴":"cun","存":"cun","忖":"cun","寸":"cun","吋":"cun","搓":"cuo","磋":"cuo","蹉":"cuo","嵯":"cuo","矬":"cuo","痤":"cuo","脞":"cuo","挫":"cuo","莝":"cuo","厝":"cuo","措":"cuo","锉":"cuo","错":"cuo","酇":"cuo","咑":"da","垯":"da","耷":"da","搭":"da","褡":"da","达":"da","怛":"da","妲":"da","荙":"da","笪":"da","答":"da","跶":"da","靼":"da","瘩":"da","鞑":"da","打":"da","呆":"dai","歹":"dai","逮":"dai","傣":"dai","代":"dai","岱":"dai","迨":"dai","玳":"dai","带":"dai","殆":"dai","贷":"dai","待":"dai","怠":"dai","袋":"dai","叇":"dai","戴":"dai","黛":"dai","襶":"dai","呔":"dai","丹":"dan","担":"dan","眈":"dan","耽":"dan","郸":"dan","聃":"dan","殚":"dan","瘅":"dan","箪":"dan","儋":"dan","胆":"dan","疸":"dan","掸":"dan","亶":"dan","旦":"dan","但":"dan","诞":"dan","萏":"dan","啖":"dan","淡":"dan","惮":"dan","蛋":"dan","氮":"dan","赕":"dan","当":"dang","裆":"dang","挡":"dang","档":"dang","党":"dang","谠":"dang","凼":"dang","砀":"dang","宕":"dang","荡":"dang","菪":"dang","刀":"dao","忉":"dao","氘":"dao","舠":"dao","导":"dao","岛":"dao","捣":"dao","倒":"dao","捯":"dao","祷":"dao","蹈":"dao","到":"dao","盗":"dao","悼":"dao","道":"dao","稻":"dao","焘":"dao","锝":"de","嘚":"de","德":"de","扽":"den","灯":"deng","登":"deng","噔":"deng","蹬":"deng","等":"deng","戥":"deng","邓":"deng","僜":"deng","凳":"deng","嶝":"deng","磴":"deng","瞪":"deng","镫":"deng","低":"di","羝":"di","堤":"di","嘀":"di","滴":"di","狄":"di","迪":"di","籴":"di","荻":"di","敌":"di","涤":"di","笛":"di","觌":"di","嫡":"di","镝":"di","氐":"di","邸":"di","诋":"di","抵":"di","底":"di","柢":"di","砥":"di","骶":"di","玓":"di","弟":"di","帝":"di","递":"di","娣":"di","第":"di","谛":"di","蒂":"di","棣":"di","睇":"di","缔":"di","碲":"di","嗲":"dia","掂":"dian","滇":"dian","颠":"dian","巅":"dian","癫":"dian","典":"dian","点":"dian","碘":"dian","踮":"dian","电":"dian","甸":"dian","阽":"dian","坫":"dian","店":"dian","玷":"dian","垫":"dian","钿":"dian","淀":"dian","惦":"dian","奠":"dian","殿":"dian","靛":"dian","刁":"diao","叼":"diao","汈":"diao","凋":"diao","貂":"diao","碉":"diao","雕":"diao","鲷":"diao","屌":"diao","吊":"diao","钓":"diao","窎":"diao","掉":"diao","铫":"diao","爹":"die","跌":"die","迭":"die","谍":"die","耋":"die","喋":"die","牒":"die","叠":"die","碟":"die","嵽":"die","蝶":"die","蹀":"die","鲽":"die","仃":"ding","叮":"ding","玎":"ding","盯":"ding","町":"ding","耵":"ding","顶":"ding","酊":"ding","鼎":"ding","订":"ding","钉":"ding","定":"ding","啶":"ding","腚":"ding","碇":"ding","锭":"ding","丢":"diu","铥":"diu","东":"dong","冬":"dong","咚":"dong","氡":"dong","鸫":"dong","董":"dong","懂":"dong","动":"dong","冻":"dong","侗":"dong","栋":"dong","胨":"dong","洞":"dong","胴":"dong","兜":"dou","蔸":"dou","篼":"dou","抖":"dou","陡":"dou","蚪":"dou","斗":"dou","豆":"dou","逗":"dou","痘":"dou","窦":"dou","督":"du","嘟":"du","毒":"du","独":"du","渎":"du","椟":"du","犊":"du","牍":"du","黩":"du","髑":"du","厾":"du","笃":"du","堵":"du","赌":"du","睹":"du","杜":"du","肚":"du","妒":"du","渡":"du","镀":"du","蠹":"du","端":"duan","短":"duan","段":"duan","断":"duan","缎":"duan","椴":"duan","锻":"duan","簖":"duan","堆":"dui","队":"dui","对":"dui","兑":"dui","怼":"dui","憝":"dui","吨":"dun","惇":"dun","敦":"dun","墩":"dun","礅":"dun","盹":"dun","趸":"dun","沌":"dun","炖":"dun","砘":"dun","钝":"dun","盾":"dun","顿":"dun","遁":"dun","多":"duo","咄":"duo","哆":"duo","掇":"duo","裰":"duo","夺":"duo","踱":"duo","朵":"duo","垛":"duo","哚":"duo","躲":"duo","亸":"duo","剁":"duo","舵":"duo","堕":"duo","惰":"duo","跺":"duo","屙":"e","婀":"e","讹":"e","囮":"e","俄":"e","莪":"e","峨":"e","娥":"e","锇":"e","鹅":"e","蛾":"e","额":"e","厄":"e","扼":"e","苊":"e","呃":"e","垩":"e","饿":"e","鄂":"e","谔":"e","萼":"e","遏":"e","愕":"e","腭":"e","颚":"e","噩":"e","鳄":"e","恩":"en","蒽":"en","摁":"en","鞥":"eng","儿":"er","而":"er","鸸":"er","尔":"er","耳":"er","迩":"er","饵":"er","洱":"er","铒":"er","二":"er","贰":"er","发":"fa","乏":"fa","伐":"fa","罚":"fa","垡":"fa","阀":"fa","筏":"fa","法":"fa","砝":"fa","珐":"fa","帆":"fan","幡":"fan","藩":"fan","翻":"fan","凡":"fan","矾":"fan","钒":"fan","烦":"fan","樊":"fan","燔":"fan","繁":"fan","蹯":"fan","蘩":"fan","反":"fan","返":"fan","犯":"fan","饭":"fan","泛":"fan","范":"fan","贩":"fan","畈":"fan","梵":"fan","方":"fang","邡":"fang","坊":"fang","芳":"fang","枋":"fang","钫":"fang","防":"fang","妨":"fang","肪":"fang","房":"fang","鲂":"fang","仿":"fang","访":"fang","纺":"fang","舫":"fang","放":"fang","飞":"fei","妃":"fei","非":"fei","菲":"fei","啡":"fei","绯":"fei","扉":"fei","肥":"fei","淝":"fei","腓":"fei","匪":"fei","诽":"fei","悱":"fei","棐":"fei","斐":"fei","榧":"fei","翡":"fei","篚":"fei","吠":"fei","肺":"fei","狒":"fei","废":"fei","沸":"fei","费":"fei","痱":"fei","镄":"fei","分":"fen","芬":"fen","吩":"fen","纷":"fen","氛":"fen","酚":"fen","坟":"fen","汾":"fen","棼":"fen","焚":"fen","鼢":"fen","粉":"fen","份":"fen","奋":"fen","忿":"fen","偾":"fen","粪":"fen","愤":"fen","丰":"feng","风":"feng","沣":"feng","枫":"feng","封":"feng","砜":"feng","疯":"feng","峰":"feng","烽":"feng","葑":"feng","锋":"feng","蜂":"feng","酆":"feng","冯":"feng","逢":"feng","缝":"feng","讽":"feng","唪":"feng","凤":"feng","奉":"feng","俸":"feng","缶":"fou","夫":"fu","呋":"fu","肤":"fu","麸":"fu","跗":"fu","稃":"fu","孵":"fu","敷":"fu","弗":"fu","伏":"fu","凫":"fu","扶":"fu","芙":"fu","孚":"fu","拂":"fu","苻":"fu","服":"fu","怫":"fu","茯":"fu","氟":"fu","俘":"fu","浮":"fu","符":"fu","匐":"fu","涪":"fu","艴":"fu","幅":"fu","辐":"fu","蜉":"fu","福":"fu","蝠":"fu","抚":"fu","甫":"fu","拊":"fu","斧":"fu","府":"fu","俯":"fu","釜":"fu","辅":"fu","腑":"fu","腐":"fu","父":"fu","讣":"fu","付":"fu","负":"fu","妇":"fu","附":"fu","咐":"fu","阜":"fu","驸":"fu","赴":"fu","复":"fu","副":"fu","赋":"fu","傅":"fu","富":"fu","腹":"fu","缚":"fu","赙":"fu","蝮":"fu","覆":"fu","馥":"fu","袱":"fu","旮":"ga","嘎":"ga","钆":"ga","尜":"ga","尕":"ga","尬":"ga","该":"gai","垓":"gai","荄":"gai","赅":"gai","改":"gai","丐":"gai","钙":"gai","溉":"gai","概":"gai","甘":"gan","玕":"gan","肝":"gan","坩":"gan","苷":"gan","矸":"gan","泔":"gan","柑":"gan","竿":"gan","酐":"gan","疳":"gan","尴":"gan","杆":"gan","秆":"gan","赶":"gan","敢":"gan","感":"gan","澉":"gan","橄":"gan","擀":"gan","干":"gan","旰":"gan","绀":"gan","淦":"gan","骭":"gan","赣":"gan","冈":"gang","冮":"gang","刚":"gang","肛":"gang","纲":"gang","钢":"gang","缸":"gang","罡":"gang","岗":"gang","港":"gang","杠":"gang","皋":"gao","高":"gao","羔":"gao","睾":"gao","膏":"gao","篙":"gao","糕":"gao","杲":"gao","搞":"gao","槁":"gao","稿":"gao","告":"gao","郜":"gao","诰":"gao","锆":"gao","戈":"ge","圪":"ge","纥":"ge","疙":"ge","哥":"ge","胳":"ge","鸽":"ge","袼":"ge","搁":"ge","割":"ge","歌":"ge","革":"ge","阁":"ge","格":"ge","隔":"ge","嗝":"ge","膈":"ge","骼":"ge","镉":"ge","舸":"ge","葛":"ge","个":"ge","各":"ge","虼":"ge","硌":"ge","铬":"ge","根":"gen","跟":"gen","哏":"gen","亘":"gen","艮":"gen","茛":"gen","庚":"geng","耕":"geng","浭":"geng","赓":"geng","羹":"geng","埂":"geng","耿":"geng","哽":"geng","绠":"geng","梗":"geng","鲠":"geng","更":"geng","工":"gong","弓":"gong","公":"gong","功":"gong","攻":"gong","肱":"gong","宫":"gong","恭":"gong","蚣":"gong","躬":"gong","龚":"gong","塨":"gong","觥":"gong","巩":"gong","汞":"gong","拱":"gong","珙":"gong","共":"gong","贡":"gong","供":"gong","勾":"gou","佝":"gou","沟":"gou","钩":"gou","篝":"gou","苟":"gou","岣":"gou","狗":"gou","枸":"gou","构":"gou","购":"gou","诟":"gou","垢":"gou","够":"gou","彀":"gou","媾":"gou","觏":"gou","估":"gu","咕":"gu","沽":"gu","孤":"gu","姑":"gu","轱":"gu","鸪":"gu","菰":"gu","菇":"gu","蛄":"gu","蓇":"gu","辜":"gu","酤":"gu","觚":"gu","毂":"gu","箍":"gu","古":"gu","谷":"gu","汩":"gu","诂":"gu","股":"gu","骨":"gu","牯":"gu","钴":"gu","羖":"gu","蛊":"gu","鼓":"gu","榾":"gu","鹘":"gu","臌":"gu","瀔":"gu","固":"gu","故":"gu","顾":"gu","梏":"gu","崮":"gu","雇":"gu","锢":"gu","痼":"gu","瓜":"gua","刮":"gua","胍":"gua","鸹":"gua","剐":"gua","寡":"gua","卦":"gua","诖":"gua","挂":"gua","褂":"gua","乖":"guai","拐":"guai","怪":"guai","关":"guan","观":"guan","官":"guan","倌":"guan","蒄":"guan","棺":"guan","瘝":"guan","鳏":"guan","馆":"guan","管":"guan","贯":"guan","冠":"guan","掼":"guan","惯":"guan","祼":"guan","盥":"guan","灌":"guan","瓘":"guan","鹳":"guan","罐":"guan","琯":"guan","光":"guang","咣":"guang","胱":"guang","广":"guang","犷":"guang","桄":"guang","逛":"guang","归":"gui","圭":"gui","龟":"gui","妫":"gui","规":"gui","皈":"gui","闺":"gui","硅":"gui","瑰":"gui","鲑":"gui","宄":"gui","轨":"gui","庋":"gui","匦":"gui","诡":"gui","鬼":"gui","姽":"gui","癸":"gui","晷":"gui","簋":"gui","柜":"gui","炅":"gui","刿":"gui","刽":"gui","贵":"gui","桂":"gui","跪":"gui","鳜":"gui","衮":"gun","绲":"gun","辊":"gun","滚":"gun","磙":"gun","鲧":"gun","棍":"gun","埚":"guo","郭":"guo","啯":"guo","崞":"guo","聒":"guo","锅":"guo","蝈":"guo","国":"guo","帼":"guo","虢":"guo","果":"guo","椁":"guo","蜾":"guo","裹":"guo","过":"guo","哈":"ha","铪":"ha","孩":"hai","骸":"hai","胲":"hai","海":"hai","醢":"hai","亥":"hai","骇":"hai","害":"hai","嗐":"hai","嗨":"hai","顸":"han","蚶":"han","酣":"han","憨":"han","鼾":"han","邗":"han","邯":"han","含":"han","函":"han","晗":"han","焓":"han","涵":"han","韩":"han","寒":"han","罕":"han","喊":"han","蔊":"han","汉":"han","汗":"han","旱":"han","捍":"han","悍":"han","菡":"han","焊":"han","撖":"han","撼":"han","翰":"han","憾":"han","瀚":"han","夯":"hang","杭":"hang","绗":"hang","航":"hang","沆":"hang","蒿":"hao","薅":"hao","嚆":"hao","蚝":"hao","毫":"hao","嗥":"hao","豪":"hao","壕":"hao","嚎":"hao","濠":"hao","好":"hao","郝":"hao","号":"hao","昊":"hao","耗":"hao","浩":"hao","皓":"hao","滈":"hao","颢":"hao","灏":"hao","诃":"he","呵":"he","喝":"he","嗬":"he","禾":"he","合":"he","何":"he","劾":"he","河":"he","曷":"he","阂":"he","盍":"he","荷":"he","菏":"he","盒":"he","涸":"he","颌":"he","阖":"he","贺":"he","赫":"he","褐":"he","鹤":"he","壑":"he","黑":"hei","嘿":"hei","痕":"hen","很":"hen","狠":"hen","恨":"hen","亨":"heng","恒":"heng","珩":"heng","横":"heng","衡":"heng","蘅":"heng","啈":"heng","轰":"hong","訇":"hong","烘":"hong","薨":"hong","弘":"hong","红":"hong","闳":"hong","宏":"hong","荭":"hong","虹":"hong","竑":"hong","洪":"hong","鸿":"hong","哄":"hong","讧":"hong","吽":"hong","齁":"hou","侯":"hou","喉":"hou","猴":"hou","瘊":"hou","骺":"hou","篌":"hou","糇":"hou","吼":"hou","后":"hou","郈":"hou","厚":"hou","垕":"hou","逅":"hou","候":"hou","堠":"hou","鲎":"hou","乎":"hu","呼":"hu","忽":"hu","轷":"hu","烀":"hu","惚":"hu","滹":"hu","囫":"hu","狐":"hu","弧":"hu","胡":"hu","壶":"hu","斛":"hu","葫":"hu","猢":"hu","湖":"hu","瑚":"hu","鹕":"hu","槲":"hu","蝴":"hu","糊":"hu","醐":"hu","觳":"hu","虎":"hu","唬":"hu","琥":"hu","互":"hu","户":"hu","冱":"hu","护":"hu","沪":"hu","枑":"hu","怙":"hu","戽":"hu","笏":"hu","瓠":"hu","扈":"hu","鹱":"hu","花":"hua","砉":"hua","华":"hua","哗":"hua","骅":"hua","铧":"hua","猾":"hua","滑":"hua","化":"hua","画":"hua","话":"hua","桦":"hua","婳":"hua","觟":"hua","怀":"huai","徊":"huai","淮":"huai","槐":"huai","踝":"huai","耲":"huai","坏":"huai","欢":"huan","獾":"huan","环":"huan","洹":"huan","桓":"huan","萑":"huan","寰":"huan","缳":"huan","缓":"huan","幻":"huan","奂":"huan","宦":"huan","换":"huan","唤":"huan","涣":"huan","浣":"huan","患":"huan","焕":"huan","痪":"huan","豢":"huan","漶":"huan","鲩":"huan","擐":"huan","肓":"huang","荒":"huang","塃":"huang","慌":"huang","皇":"huang","黄":"huang","凰":"huang","隍":"huang","喤":"huang","遑":"huang","徨":"huang","湟":"huang","惶":"huang","媓":"huang","煌":"huang","锽":"huang","潢":"huang","璜":"huang","蝗":"huang","篁":"huang","艎":"huang","磺":"huang","癀":"huang","蟥":"huang","簧":"huang","鳇":"huang","恍":"huang","晃":"huang","谎":"huang","幌":"huang","滉":"huang","皝":"huang","灰":"hui","诙":"hui","挥":"hui","恢":"hui","晖":"hui","辉":"hui","麾":"hui","徽":"hui","隳":"hui","回":"hui","茴":"hui","洄":"hui","蛔":"hui","悔":"hui","毁":"hui","卉":"hui","汇":"hui","讳":"hui","荟":"hui","浍":"hui","诲":"hui","绘":"hui","恚":"hui","贿":"hui","烩":"hui","彗":"hui","晦":"hui","秽":"hui","惠":"hui","喙":"hui","慧":"hui","蕙":"hui","蟪":"hui","珲":"hun","昏":"hun","荤":"hun","阍":"hun","惛":"hun","婚":"hun","浑":"hun","馄":"hun","混":"hun","魂":"hun","诨":"hun","溷":"hun","耠":"huo","劐":"huo","豁":"huo","活":"huo","火":"huo","伙":"huo","钬":"huo","夥":"huo","或":"huo","货":"huo","获":"huo","祸":"huo","惑":"huo","霍":"huo","镬":"huo","攉":"huo","藿":"huo","嚯":"huo","讥":"ji","击":"ji","叽":"ji","饥":"ji","玑":"ji","圾":"ji","芨":"ji","机":"ji","乩":"ji","肌":"ji","矶":"ji","鸡":"ji","剞":"ji","唧":"ji","积":"ji","笄":"ji","屐":"ji","姬":"ji","基":"ji","犄":"ji","嵇":"ji","畸":"ji","跻":"ji","箕":"ji","齑":"ji","畿":"ji","墼":"ji","激":"ji","羁":"ji","及":"ji","吉":"ji","岌":"ji","汲":"ji","级":"ji","极":"ji","即":"ji","佶":"ji","笈":"ji","急":"ji","疾":"ji","棘":"ji","集":"ji","蒺":"ji","楫":"ji","辑":"ji","嫉":"ji","瘠":"ji","藉":"ji","籍":"ji","几":"ji","己":"ji","虮":"ji","挤":"ji","脊":"ji","掎":"ji","戟":"ji","麂":"ji","计":"ji","记":"ji","伎":"ji","纪":"ji","技":"ji","忌":"ji","际":"ji","妓":"ji","季":"ji","剂":"ji","迹":"ji","济":"ji","既":"ji","觊":"ji","继":"ji","偈":"ji","祭":"ji","悸":"ji","寄":"ji","寂":"ji","绩":"ji","暨":"ji","稷":"ji","鲫":"ji","髻":"ji","冀":"ji","骥":"ji","加":"jia","佳":"jia","枷":"jia","浃":"jia","痂":"jia","家":"jia","袈":"jia","嘉":"jia","镓":"jia","荚":"jia","戛":"jia","颊":"jia","甲":"jia","胛":"jia","钾":"jia","假":"jia","价":"jia","驾":"jia","架":"jia","嫁":"jia","稼":"jia","戋":"jian","尖":"jian","奸":"jian","歼":"jian","坚":"jian","间":"jian","肩":"jian","艰":"jian","监":"jian","兼":"jian","菅":"jian","笺":"jian","缄":"jian","煎":"jian","拣":"jian","茧":"jian","柬":"jian","俭":"jian","捡":"jian","检":"jian","减":"jian","剪":"jian","睑":"jian","简":"jian","碱":"jian","见":"jian","件":"jian","饯":"jian","建":"jian","荐":"jian","贱":"jian","剑":"jian","健":"jian","舰":"jian","涧":"jian","渐":"jian","谏":"jian","践":"jian","锏":"jian","毽":"jian","腱":"jian","溅":"jian","鉴":"jian","键":"jian","僭":"jian","箭":"jian","江":"jiang","将":"jiang","姜":"jiang","豇":"jiang","浆":"jiang","僵":"jiang","缰":"jiang","疆":"jiang","讲":"jiang","奖":"jiang","桨":"jiang","蒋":"jiang","匠":"jiang","酱":"jiang","犟":"jiang","糨":"jiang","交":"jiao","郊":"jiao","浇":"jiao","娇":"jiao","姣":"jiao","骄":"jiao","胶":"jiao","椒":"jiao","蛟":"jiao","焦":"jiao","跤":"jiao","蕉":"jiao","礁":"jiao","佼":"jiao","狡":"jiao","饺":"jiao","绞":"jiao","铰":"jiao","矫":"jiao","皎":"jiao","脚":"jiao","搅":"jiao","剿":"jiao","缴":"jiao","叫":"jiao","轿":"jiao","较":"jiao","教":"jiao","窖":"jiao","酵":"jiao","侥":"jiao","阶":"jie","皆":"jie","接":"jie","秸":"jie","揭":"jie","嗟":"jie","街":"jie","孑":"jie","节":"jie","讦":"jie","劫":"jie","杰":"jie","诘":"jie","洁":"jie","结":"jie","捷":"jie","睫":"jie","截":"jie","碣":"jie","竭":"jie","姐":"jie","解":"jie","介":"jie","戒":"jie","届":"jie","界":"jie","疥":"jie","诫":"jie","借":"jie","巾":"jin","斤":"jin","今":"jin","金":"jin","津":"jin","矜":"jin","筋":"jin","襟":"jin","仅":"jin","紧":"jin","锦":"jin","谨":"jin","尽":"jin","进":"jin","近":"jin","晋":"jin","烬":"jin","浸":"jin","禁":"jin","觐":"jin","噤":"jin","茎":"jing","京":"jing","泾":"jing","经":"jing","菁":"jing","惊":"jing","晶":"jing","睛":"jing","粳":"jing","兢":"jing","精":"jing","鲸":"jing","井":"jing","阱":"jing","刭":"jing","景":"jing","儆":"jing","警":"jing","径":"jing","净":"jing","痉":"jing","竞":"jing","竟":"jing","敬":"jing","靖":"jing","静":"jing","境":"jing","镜":"jing","迥":"jiong","炯":"jiong","窘":"jiong","纠":"jiu","鸠":"jiu","究":"jiu","赳":"jiu","阄":"jiu","揪":"jiu","啾":"jiu","九":"jiu","久":"jiu","玖":"jiu","灸":"jiu","韭":"jiu","酒":"jiu","旧":"jiu","臼":"jiu","咎":"jiu","柩":"jiu","救":"jiu","厩":"jiu","就":"jiu","舅":"jiu","鹫":"jiu","军":"jun","均":"jun","君":"jun","钧":"jun","菌":"jun","皲":"jun","俊":"jun","郡":"jun","峻":"jun","骏":"jun","竣":"jun","拘":"ju","狙":"ju","居":"ju","驹":"ju","掬":"ju","雎":"ju","鞠":"ju","局":"ju","菊":"ju","焗":"ju","橘":"ju","咀":"ju","沮":"ju","矩":"ju","举":"ju","龃":"ju","巨":"ju","拒":"ju","具":"ju","炬":"ju","俱":"ju","剧":"ju","据":"ju","距":"ju","惧":"ju","飓":"ju","锯":"ju","聚":"ju","踞":"ju","捐":"juan","涓":"juan","娟":"juan","鹃":"juan","卷":"juan","倦":"juan","绢":"juan","眷":"juan","隽":"juan","撅":"jue","噘":"jue","决":"jue","诀":"jue","抉":"jue","绝":"jue","掘":"jue","崛":"jue","厥":"jue","谲":"jue","蕨":"jue","爵":"jue","蹶":"jue","矍":"jue","倔":"jue","咔":"ka","开":"kai","揩":"kai","凯":"kai","铠":"kai","慨":"kai","楷":"kai","忾":"kai","刊":"kan","勘":"kan","龛":"kan","堪":"kan","坎":"kan","侃":"kan","砍":"kan","槛":"kan","看":"kan","瞰":"kan","康":"kang","慷":"kang","糠":"kang","亢":"kang","伉":"kang","抗":"kang","炕":"kang","考":"kao","拷":"kao","烤":"kao","铐":"kao","犒":"kao","靠":"kao","苛":"ke","轲":"ke","科":"ke","棵":"ke","搕":"ke","嗑":"ke","稞":"ke","窠":"ke","颗":"ke","磕":"ke","瞌":"ke","蝌":"ke","可":"ke","坷":"ke","渴":"ke","克":"ke","刻":"ke","恪":"ke","客":"ke","课":"ke","肯":"ken","垦":"ken","恳":"ken","啃":"ken","坑":"keng","铿":"keng","空":"kong","孔":"kong","恐":"kong","控":"kong","抠":"kou","口":"kou","叩":"kou","扣":"kou","寇":"kou","蔻":"kou","枯":"ku","哭":"ku","窟":"ku","骷":"ku","苦":"ku","库":"ku","绔":"ku","裤":"ku","酷":"ku","夸":"kua","垮":"kua","挎":"kua","胯":"kua","跨":"kua","块":"kuai","快":"kuai","侩":"kuai","脍":"kuai","筷":"kuai","宽":"kuan","髋":"kuan","款":"kuan","诓":"kuang","哐":"kuang","筐":"kuang","狂":"kuang","诳":"kuang","旷":"kuang","况":"kuang","矿":"kuang","框":"kuang","眶":"kuang","亏":"kui","盔":"kui","窥":"kui","葵":"kui","魁":"kui","傀":"kui","匮":"kui","馈":"kui","愧":"kui","坤":"kun","昆":"kun","鲲":"kun","捆":"kun","困":"kun","扩":"kuo","括":"kuo","阔":"kuo","廓":"kuo","垃":"la","拉":"la","啦":"la","邋":"la","旯":"la","喇":"la","腊":"la","蜡":"la","辣":"la","来":"lai","莱":"lai","徕":"lai","睐":"lai","赖":"lai","癞":"lai","籁":"lai","兰":"lan","岚":"lan","拦":"lan","栏":"lan","婪":"lan","阑":"lan","蓝":"lan","澜":"lan","褴":"lan","篮":"lan","览":"lan","揽":"lan","缆":"lan","榄":"lan","懒":"lan","烂":"lan","滥":"lan","啷":"lang","郎":"lang","狼":"lang","琅":"lang","廊":"lang","榔":"lang","锒":"lang","螂":"lang","朗":"lang","浪":"lang","捞":"lao","劳":"lao","牢":"lao","崂":"lao","老":"lao","佬":"lao","姥":"lao","唠":"lao","烙":"lao","涝":"lao","酪":"lao","雷":"lei","羸":"lei","垒":"lei","磊":"lei","蕾":"lei","儡":"lei","肋":"lei","泪":"lei","类":"lei","累":"lei","擂":"lei","嘞":"lei","棱":"leng","楞":"leng","冷":"leng","睖":"leng","厘":"li","狸":"li","离":"li","梨":"li","犁":"li","鹂":"li","喱":"li","蜊":"li","漓":"li","璃":"li","黎":"li","罹":"li","篱":"li","蠡":"li","礼":"li","李":"li","里":"li","俚":"li","逦":"li","哩":"li","娌":"li","理":"li","鲤":"li","力":"li","历":"li","厉":"li","立":"li","吏":"li","丽":"li","励":"li","呖":"li","利":"li","沥":"li","枥":"li","例":"li","戾":"li","隶":"li","荔":"li","俐":"li","莉":"li","莅":"li","栗":"li","砾":"li","蛎":"li","唳":"li","笠":"li","粒":"li","雳":"li","痢":"li","连":"lian","怜":"lian","帘":"lian","莲":"lian","涟":"lian","联":"lian","廉":"lian","鲢":"lian","镰":"lian","敛":"lian","脸":"lian","练":"lian","炼":"lian","恋":"lian","殓":"lian","链":"lian","良":"liang","凉":"liang","梁":"liang","粮":"liang","粱":"liang","两":"liang","魉":"liang","亮":"liang","谅":"liang","辆":"liang","靓":"liang","量":"liang","晾":"liang","踉":"liang","辽":"liao","疗":"liao","聊":"liao","僚":"liao","寥":"liao","撩":"liao","嘹":"liao","獠":"liao","潦":"liao","缭":"liao","燎":"liao","料":"liao","撂":"liao","瞭":"liao","镣":"liao","咧":"lie","列":"lie","劣":"lie","冽":"lie","烈":"lie","猎":"lie","裂":"lie","趔":"lie","拎":"lin","邻":"lin","林":"lin","临":"lin","淋":"lin","琳":"lin","粼":"lin","嶙":"lin","潾":"lin","霖":"lin","磷":"lin","鳞":"lin","麟":"lin","凛":"lin","檩":"lin","吝":"lin","赁":"lin","躏":"lin","伶":"ling","灵":"ling","苓":"ling","囹":"ling","泠":"ling","玲":"ling","瓴":"ling","铃":"ling","凌":"ling","陵":"ling","聆":"ling","菱":"ling","棂":"ling","蛉":"ling","翎":"ling","羚":"ling","绫":"ling","零":"ling","龄":"ling","岭":"ling","领":"ling","另":"ling","令":"ling","溜":"liu","熘":"liu","刘":"liu","浏":"liu","留":"liu","流":"liu","琉":"liu","硫":"liu","馏":"liu","榴":"liu","瘤":"liu","柳":"liu","绺":"liu","六":"liu","遛":"liu","龙":"long","咙":"long","珑":"long","胧":"long","聋":"long","笼":"long","隆":"long","窿":"long","陇":"long","拢":"long","垄":"long","娄":"lou","楼":"lou","髅":"lou","搂":"lou","篓":"lou","陋":"lou","镂":"lou","漏":"lou","喽":"lou","撸":"lu","卢":"lu","芦":"lu","庐":"lu","炉":"lu","泸":"lu","鸬":"lu","颅":"lu","鲈":"lu","卤":"lu","虏":"lu","掳":"lu","鲁":"lu","橹":"lu","录":"lu","赂":"lu","鹿":"lu","禄":"lu","路":"lu","箓":"lu","漉":"lu","戮":"lu","鹭":"lu","麓":"lu","峦":"luan","孪":"luan","挛":"luan","鸾":"luan","卵":"luan","乱":"luan","抡":"lun","仑":"lun","伦":"lun","囵":"lun","沦":"lun","轮":"lun","论":"lun","啰":"luo","罗":"luo","萝":"luo","逻":"luo","锣":"luo","箩":"luo","骡":"luo","螺":"luo","裸":"luo","洛":"luo","络":"luo","骆":"luo","摞":"luo","漯":"luo","驴":"lv","榈":"lv","吕":"lv","侣":"lv","旅":"lv","铝":"lv","屡":"lv","缕":"lv","膂":"lv","褛":"lv","履":"lv","律":"lv","虑":"lv","氯":"lv","滤":"lv","掠":"lve","略":"lve","妈":"ma","麻":"ma","蟆":"ma","马":"ma","犸":"ma","玛":"ma","码":"ma","蚂":"ma","骂":"ma","吗":"ma","嘛":"ma","霾":"mai","买":"mai","迈":"mai","麦":"mai","卖":"mai","霡":"mai","蛮":"man","馒":"man","瞒":"man","满":"man","曼":"man","谩":"man","幔":"man","漫":"man","慢":"man","牤":"mang","芒":"mang","忙":"mang","盲":"mang","氓":"mang","茫":"mang","莽":"mang","漭":"mang","蟒":"mang","猫":"mao","毛":"mao","矛":"mao","茅":"mao","牦":"mao","锚":"mao","髦":"mao","蝥":"mao","蟊":"mao","冇":"mao","卯":"mao","铆":"mao","茂":"mao","冒":"mao","贸":"mao","袤":"mao","帽":"mao","貌":"mao","玫":"mei","枚":"mei","眉":"mei","莓":"mei","梅":"mei","媒":"mei","楣":"mei","煤":"mei","酶":"mei","霉":"mei","每":"mei","美":"mei","镁":"mei","妹":"mei","昧":"mei","袂":"mei","寐":"mei","媚":"mei","魅":"mei","门":"men","扪":"men","闷":"men","焖":"men","懑":"men","们":"men","虻":"meng","萌":"meng","蒙":"meng","盟":"meng","檬":"meng","曚":"meng","朦":"meng","猛":"meng","锰":"meng","蜢":"meng","懵":"meng","孟":"meng","梦":"meng","咪":"mi","眯":"mi","弥":"mi","迷":"mi","猕":"mi","谜":"mi","醚":"mi","糜":"mi","麋":"mi","靡":"mi","米":"mi","弭":"mi","觅":"mi","密":"mi","幂":"mi","谧":"mi","蜜":"mi","眠":"mian","绵":"mian","棉":"mian","免":"mian","勉":"mian","娩":"mian","冕":"mian","渑":"mian","湎":"mian","缅":"mian","腼":"mian","面":"mian","喵":"miao","苗":"miao","描":"miao","瞄":"miao","秒":"miao","渺":"miao","藐":"miao","妙":"miao","庙":"miao","缥":"miao","咩":"mie","灭":"mie","蔑":"mie","篾":"mie","乜":"mie","民":"min","皿":"min","抿":"min","泯":"min","闽":"min","悯":"min","敏":"min","名":"ming","明":"ming","鸣":"ming","茗":"ming","冥":"ming","铭":"ming","瞑":"ming","螟":"ming","酩":"ming","命":"ming","谬":"miu","摸":"mo","馍":"mo","摹":"mo","膜":"mo","摩":"mo","磨":"mo","蘑":"mo","魔":"mo","末":"mo","茉":"mo","殁":"mo","沫":"mo","陌":"mo","莫":"mo","秣":"mo","蓦":"mo","漠":"mo","寞":"mo","墨":"mo","默":"mo","嬷":"mo","缪":"mou","哞":"mou","眸":"mou","谋":"mou","某":"mou","母":"mu","牡":"mu","亩":"mu","拇":"mu","姆":"mu","木":"mu","目":"mu","沐":"mu","苜":"mu","牧":"mu","钼":"mu","募":"mu","墓":"mu","幕":"mu","睦":"mu","慕":"mu","暮":"mu","穆":"mu","拿":"na","呐":"na","纳":"na","钠":"na","衲":"na","捺":"na","乃":"nai","奶":"nai","氖":"nai","奈":"nai","耐":"nai","囡":"nan","男":"nan","南":"nan","难":"nan","喃":"nan","楠":"nan","赧":"nan","腩":"nan","囔":"nang","囊":"nang","孬":"nao","呶":"nao","挠":"nao","恼":"nao","脑":"nao","瑙":"nao","闹":"nao","淖":"nao","讷":"ne","馁":"nei","内":"nei","嫩":"nen","恁":"nen","能":"neng","嗯":"ng","妮":"ni","尼":"ni","泥":"ni","怩":"ni","倪":"ni","霓":"ni","拟":"ni","你":"ni","旎":"ni","昵":"ni","逆":"ni","匿":"ni","腻":"ni","溺":"ni","拈":"nian","蔫":"nian","年":"nian","黏":"nian","捻":"nian","辇":"nian","撵":"nian","碾":"nian","廿":"nian","念":"nian","娘":"niang","酿":"niang","鸟":"niao","袅":"niao","尿":"niao","捏":"nie","聂":"nie","涅":"nie","嗫":"nie","镊":"nie","镍":"nie","蹑":"nie","孽":"nie","您":"nin","宁":"ning","咛":"ning","狞":"ning","柠":"ning","凝":"ning","拧":"ning","佞":"ning","泞":"ning","妞":"niu","牛":"niu","扭":"niu","忸":"niu","纽":"niu","钮":"niu","农":"nong","哝":"nong","浓":"nong","脓":"nong","弄":"nong","奴":"nu","驽":"nu","努":"nu","弩":"nu","怒":"nu","暖":"nuan","疟":"nue","虐":"nue","挪":"nuo","诺":"nuo","喏":"nuo","懦":"nuo","糯":"nuo","女":"nv","噢":"o","讴":"ou","瓯":"ou","欧":"ou","殴":"ou","鸥":"ou","呕":"ou","偶":"ou","藕":"ou","怄":"ou","趴":"pa","啪":"pa","葩":"pa","杷":"pa","爬":"pa","琶":"pa","帕":"pa","怕":"pa","拍":"pai","排":"pai","徘":"pai","牌":"pai","哌":"pai","派":"pai","湃":"pai","潘":"pan","攀":"pan","爿":"pan","盘":"pan","磐":"pan","蹒":"pan","蟠":"pan","判":"pan","盼":"pan","叛":"pan","畔":"pan","乓":"pang","滂":"pang","庞":"pang","旁":"pang","螃":"pang","耪":"pang","抛":"pao","咆":"pao","庖":"pao","袍":"pao","跑":"pao","泡":"pao","呸":"pei","胚":"pei","陪":"pei","培":"pei","赔":"pei","裴":"pei","沛":"pei","佩":"pei","配":"pei","喷":"pen","盆":"pen","抨":"peng","怦":"peng","砰":"peng","烹":"peng","嘭":"peng","朋":"peng","彭":"peng","棚":"peng","蓬":"peng","硼":"peng","鹏":"peng","澎":"peng","篷":"peng","膨":"peng","捧":"peng","碰":"peng","丕":"pi","批":"pi","纰":"pi","坯":"pi","披":"pi","砒":"pi","劈":"pi","噼":"pi","霹":"pi","皮":"pi","枇":"pi","毗":"pi","蚍":"pi","疲":"pi","啤":"pi","琵":"pi","脾":"pi","貔":"pi","匹":"pi","痞":"pi","癖":"pi","屁":"pi","睥":"pi","媲":"pi","僻":"pi","譬":"pi","偏":"pian","篇":"pian","翩":"pian","骈":"pian","蹁":"pian","片":"pian","骗":"pian","剽":"piao","漂":"piao","飘":"piao","瓢":"piao","殍":"piao","瞟":"piao","票":"piao","氕":"pie","瞥":"pie","撇":"pie","拼":"pin","姘":"pin","贫":"pin","频":"pin","嫔":"pin","颦":"pin","品":"pin","聘":"pin","乒":"ping","娉":"ping","平":"ping","评":"ping","坪":"ping","苹":"ping","凭":"ping","瓶":"ping","萍":"ping","钋":"po","坡":"po","泼":"po","颇":"po","婆":"po","鄱":"po","叵":"po","珀":"po","破":"po","粕":"po","魄":"po","剖":"pou","抔":"pou","扑":"pu","铺":"pu","噗":"pu","仆":"pu","匍":"pu","菩":"pu","葡":"pu","蒲":"pu","璞":"pu","圃":"pu","浦":"pu","普":"pu","谱":"pu","蹼":"pu","七":"qi","沏":"qi","妻":"qi","柒":"qi","凄":"qi","萋":"qi","戚":"qi","期":"qi","欺":"qi","嘁":"qi","漆":"qi","齐":"qi","芪":"qi","其":"qi","歧":"qi","祈":"qi","祇":"qi","脐":"qi","畦":"qi","跂":"qi","崎":"qi","骑":"qi","琪":"qi","棋":"qi","旗":"qi","鳍":"qi","麒":"qi","乞":"qi","岂":"qi","企":"qi","杞":"qi","启":"qi","起":"qi","绮":"qi","气":"qi","讫":"qi","迄":"qi","弃":"qi","汽":"qi","泣":"qi","契":"qi","砌":"qi","葺":"qi","器":"qi","憩":"qi","俟":"qi","掐":"qia","洽":"qia","恰":"qia","千":"qian","仟":"qian","阡":"qian","芊":"qian","迁":"qian","钎":"qian","牵":"qian","悭":"qian","谦":"qian","签":"qian","愆":"qian","前":"qian","虔":"qian","钱":"qian","钳":"qian","乾":"qian","潜":"qian","黔":"qian","遣":"qian","谴":"qian","欠":"qian","芡":"qian","倩":"qian","堑":"qian","嵌":"qian","歉":"qian","羌":"qiang","枪":"qiang","戕":"qiang","腔":"qiang","蜣":"qiang","锵":"qiang","墙":"qiang","蔷":"qiang","抢":"qiang","羟":"qiang","襁":"qiang","呛":"qiang","炝":"qiang","跄":"qiang","悄":"qiao","跷":"qiao","锹":"qiao","敲":"qiao","橇":"qiao","乔":"qiao","侨":"qiao","荞":"qiao","桥":"qiao","憔":"qiao","瞧":"qiao","巧":"qiao","俏":"qiao","诮":"qiao","峭":"qiao","窍":"qiao","翘":"qiao","撬":"qiao","切":"qie","且":"qie","妾":"qie","怯":"qie","窃":"qie","挈":"qie","惬":"qie","趄":"qie","锲":"qie","钦":"qin","侵":"qin","衾":"qin","芹":"qin","芩":"qin","秦":"qin","琴":"qin","禽":"qin","勤":"qin","擒":"qin","噙":"qin","寝":"qin","沁":"qin","青":"qing","轻":"qing","氢":"qing","倾":"qing","卿":"qing","清":"qing","蜻":"qing","情":"qing","晴":"qing","氰":"qing","擎":"qing","顷":"qing","请":"qing","庆":"qing","罄":"qing","穷":"qiong","穹":"qiong","琼":"qiong","丘":"qiu","秋":"qiu","蚯":"qiu","鳅":"qiu","囚":"qiu","求":"qiu","虬":"qiu","泅":"qiu","酋":"qiu","球":"qiu","遒":"qiu","裘":"qiu","岖":"qu","驱":"qu","屈":"qu","蛆":"qu","躯":"qu","趋":"qu","蛐":"qu","黢":"qu","渠":"qu","瞿":"qu","曲":"qu","取":"qu","娶":"qu","龋":"qu","去":"qu","趣":"qu","觑":"qu","悛":"quan","权":"quan","全":"quan","诠":"quan","泉":"quan","拳":"quan","痊":"quan","蜷":"quan","醛":"quan","犬":"quan","劝":"quan","券":"quan","炔":"que","缺":"que","瘸":"que","却":"que","确":"que","鹊":"que","阙":"que","榷":"que","逡":"qun","裙":"qun","群":"qun","蚺":"ran","然":"ran","燃":"ran","冉":"ran","苒":"ran","染":"ran","瓤":"rang","壤":"rang","攘":"rang","嚷":"rang","让":"rang","荛":"rao","饶":"rao","娆":"rao","桡":"rao","扰":"rao","绕":"rao","惹":"re","热":"re","人":"ren","壬":"ren","仁":"ren","忍":"ren","荏":"ren","稔":"ren","刃":"ren","认":"ren","任":"ren","纫":"ren","韧":"ren","饪":"ren","扔":"reng","仍":"reng","日":"ri","戎":"rong","茸":"rong","荣":"rong","绒":"rong","容":"rong","嵘":"rong","蓉":"rong","溶":"rong","榕":"rong","熔":"rong","融":"rong","冗":"rong","氄":"rong","柔":"rou","揉":"rou","糅":"rou","蹂":"rou","鞣":"rou","肉":"rou","如":"ru","茹":"ru","铷":"ru","儒":"ru","孺":"ru","蠕":"ru","汝":"ru","乳":"ru","辱":"ru","入":"ru","缛":"ru","褥":"ru","阮":"ruan","软":"ruan","蕊":"rui","蚋":"rui","锐":"rui","瑞":"rui","睿":"rui","闰":"run","润":"run","若":"ruo","偌":"ruo","弱":"ruo","仨":"sa","洒":"sa","撒":"sa","卅":"sa","飒":"sa","萨":"sa","腮":"sai","赛":"sai","三":"san","叁":"san","伞":"san","散":"san","桑":"sang","搡":"sang","嗓":"sang","丧":"sang","搔":"sao","骚":"sao","扫":"sao","嫂":"sao","臊":"sao","涩":"se","啬":"se","铯":"se","瑟":"se","穑":"se","森":"sen","僧":"seng","杀":"sha","沙":"sha","纱":"sha","砂":"sha","啥":"sha","傻":"sha","厦":"sha","歃":"sha","煞":"sha","霎":"sha","筛":"shai","晒":"shai","山":"shan","删":"shan","苫":"shan","衫":"shan","姗":"shan","珊":"shan","煽":"shan","潸":"shan","膻":"shan","闪":"shan","陕":"shan","讪":"shan","汕":"shan","扇":"shan","善":"shan","骟":"shan","缮":"shan","擅":"shan","膳":"shan","嬗":"shan","赡":"shan","鳝":"shan","伤":"shang","殇":"shang","商":"shang","觞":"shang","熵":"shang","晌":"shang","赏":"shang","上":"shang","尚":"shang","捎":"shao","烧":"shao","梢":"shao","稍":"shao","艄":"shao","勺":"shao","芍":"shao","韶":"shao","少":"shao","邵":"shao","绍":"shao","哨":"shao","潲":"shao","奢":"she","赊":"she","舌":"she","佘":"she","蛇":"she","舍":"she","设":"she","社":"she","射":"she","涉":"she","赦":"she","摄":"she","慑":"she","麝":"she","申":"shen","伸":"shen","身":"shen","呻":"shen","绅":"shen","砷":"shen","深":"shen","神":"shen","沈":"shen","审":"shen","哂":"shen","婶":"shen","肾":"shen","甚":"shen","渗":"shen","葚":"shen","蜃":"shen","慎":"shen","升":"sheng","生":"sheng","声":"sheng","昇":"sheng","牲":"sheng","笙":"sheng","甥":"sheng","绳":"sheng","圣":"sheng","胜":"sheng","晟":"sheng","剩":"sheng","尸":"shi","失":"shi","师":"shi","诗":"shi","虱":"shi","狮":"shi","施":"shi","湿":"shi","十":"shi","时":"shi","实":"shi","食":"shi","蚀":"shi","史":"shi","矢":"shi","使":"shi","始":"shi","驶":"shi","屎":"shi","士":"shi","氏":"shi","示":"shi","世":"shi","仕":"shi","市":"shi","式":"shi","势":"shi","事":"shi","侍":"shi","饰":"shi","试":"shi","视":"shi","拭":"shi","柿":"shi","是":"shi","适":"shi","恃":"shi","室":"shi","逝":"shi","轼":"shi","舐":"shi","弑":"shi","释":"shi","谥":"shi","嗜":"shi","誓":"shi","收":"shou","手":"shou","守":"shou","首":"shou","寿":"shou","受":"shou","狩":"shou","授":"shou","售":"shou","兽":"shou","绶":"shou","瘦":"shou","殳":"shu","书":"shu","抒":"shu","枢":"shu","叔":"shu","姝":"shu","殊":"shu","倏":"shu","梳":"shu","淑":"shu","舒":"shu","疏":"shu","输":"shu","蔬":"shu","秫":"shu","孰":"shu","赎":"shu","塾":"shu","暑":"shu","黍":"shu","署":"shu","蜀":"shu","鼠":"shu","薯":"shu","曙":"shu","戍":"shu","束":"shu","述":"shu","树":"shu","竖":"shu","恕":"shu","庶":"shu","墅":"shu","漱":"shu","刷":"shua","唰":"shua","耍":"shua","衰":"shuai","摔":"shuai","甩":"shuai","帅":"shuai","蟀":"shuai","闩":"shuan","拴":"shuan","栓":"shuan","涮":"shuan","双":"shuang","霜":"shuang","孀":"shuang","爽":"shuang","谁":"shui","水":"shui","税":"shui","睡":"shui","吮":"shun","顺":"shun","舜":"shun","瞬":"shun","烁":"shuo","铄":"shuo","朔":"shuo","硕":"shuo","司":"si","丝":"si","私":"si","咝":"si","思":"si","斯":"si","厮":"si","撕":"si","嘶":"si","死":"si","巳":"si","四":"si","寺":"si","祀":"si","饲":"si","肆":"si","嗣":"si","松":"song","嵩":"song","怂":"song","耸":"song","悚":"song","讼":"song","宋":"song","送":"song","诵":"song","颂":"song","搜":"sou","嗖":"sou","馊":"sou","艘":"sou","叟":"sou","擞":"sou","嗽":"sou","苏":"su","酥":"su","俗":"su","夙":"su","诉":"su","肃":"su","素":"su","速":"su","粟":"su","嗉":"su","塑":"su","溯":"su","簌":"su","酸":"suan","蒜":"suan","算":"suan","虽":"sui","睢":"sui","绥":"sui","隋":"sui","随":"sui","髓":"sui","岁":"sui","祟":"sui","遂":"sui","碎":"sui","隧":"sui","穗":"sui","孙":"sun","损":"sun","笋":"sun","隼":"sun","唆":"suo","梭":"suo","蓑":"suo","羧":"suo","缩":"suo","所":"suo","索":"suo","唢":"suo","琐":"suo","锁":"suo","他":"ta","它":"ta","她":"ta","铊":"ta","塌":"ta","塔":"ta","獭":"ta","挞":"ta","榻":"ta","踏":"ta","蹋":"ta","胎":"tai","台":"tai","邰":"tai","抬":"tai","苔":"tai","跆":"tai","太":"tai","汰":"tai","态":"tai","钛":"tai","泰":"tai","酞":"tai","贪":"tan","摊":"tan","滩":"tan","瘫":"tan","坛":"tan","昙":"tan","谈":"tan","痰":"tan","谭":"tan","潭":"tan","檀":"tan","坦":"tan","袒":"tan","毯":"tan","叹":"tan","炭":"tan","探":"tan","碳":"tan","汤":"tang","嘡":"tang","羰":"tang","唐":"tang","堂":"tang","棠":"tang","塘":"tang","搪":"tang","膛":"tang","镗":"tang","糖":"tang","螳":"tang","倘":"tang","淌":"tang","躺":"tang","烫":"tang","趟":"tang","涛":"tao","绦":"tao","掏":"tao","滔":"tao","韬":"tao","饕":"tao","逃":"tao","桃":"tao","陶":"tao","萄":"tao","淘":"tao","讨":"tao","套":"tao","特":"te","疼":"teng","腾":"teng","誊":"teng","滕":"teng","藤":"teng","剔":"ti","梯":"ti","踢":"ti","啼":"ti","题":"ti","醍":"ti","蹄":"ti","体":"ti","屉":"ti","剃":"ti","涕":"ti","悌":"ti","惕":"ti","替":"ti","天":"tian","添":"tian","田":"tian","恬":"tian","甜":"tian","填":"tian","忝":"tian","殄":"tian","舔":"tian","掭":"tian","佻":"tiao","挑":"tiao","条":"tiao","迢":"tiao","笤":"tiao","髫":"tiao","窕":"tiao","眺":"tiao","粜":"tiao","跳":"tiao","帖":"tie","贴":"tie","铁":"tie","餮":"tie","铤":"ting","厅":"ting","听":"ting","烃":"ting","廷":"ting","亭":"ting","庭":"ting","停":"ting","蜓":"ting","婷":"ting","霆":"ting","挺":"ting","艇":"ting","通":"tong","嗵":"tong","同":"tong","彤":"tong","桐":"tong","铜":"tong","童":"tong","潼":"tong","瞳":"tong","统":"tong","捅":"tong","桶":"tong","筒":"tong","恸":"tong","痛":"tong","偷":"tou","头":"tou","投":"tou","骰":"tou","透":"tou","凸":"tu","秃":"tu","突":"tu","图":"tu","荼":"tu","徒":"tu","途":"tu","涂":"tu","屠":"tu","土":"tu","吐":"tu","兔":"tu","菟":"tu","湍":"tuan","团":"tuan","疃":"tuan","彖":"tuan","推":"tui","颓":"tui","腿":"tui","退":"tui","蜕":"tui","褪":"tui","吞":"tun","屯":"tun","饨":"tun","豚":"tun","臀":"tun","托":"tuo","拖":"tuo","脱":"tuo","佗":"tuo","陀":"tuo","驼":"tuo","鸵":"tuo","妥":"tuo","椭":"tuo","唾":"tuo","挖":"wa","哇":"wa","洼":"wa","娲":"wa","蛙":"wa","娃":"wa","瓦":"wa","佤":"wa","袜":"wa","歪":"wai","外":"wai","弯":"wan","剜":"wan","湾":"wan","蜿":"wan","豌":"wan","丸":"wan","纨":"wan","完":"wan","玩":"wan","顽":"wan","烷":"wan","宛":"wan","挽":"wan","晚":"wan","惋":"wan","婉":"wan","绾":"wan","皖":"wan","碗":"wan","万":"wan","腕":"wan","汪":"wang","亡":"wang","王":"wang","网":"wang","枉":"wang","罔":"wang","往":"wang","惘":"wang","妄":"wang","忘":"wang","旺":"wang","望":"wang","危":"wei","威":"wei","偎":"wei","微":"wei","煨":"wei","薇":"wei","巍":"wei","韦":"wei","为":"wei","违":"wei","围":"wei","闱":"wei","桅":"wei","唯":"wei","帷":"wei","维":"wei","伟":"wei","伪":"wei","苇":"wei","纬":"wei","委":"wei","诿":"wei","娓":"wei","萎":"wei","猥":"wei","痿":"wei","卫":"wei","未":"wei","位":"wei","味":"wei","畏":"wei","胃":"wei","谓":"wei","喂":"wei","猬":"wei","渭":"wei","蔚":"wei","慰":"wei","魏":"wei","温":"wen","瘟":"wen","文":"wen","纹":"wen","闻":"wen","蚊":"wen","雯":"wen","刎":"wen","吻":"wen","紊":"wen","稳":"wen","问":"wen","汶":"wen","翁":"weng","嗡":"weng","瓮":"weng","挝":"wo","莴":"wo","倭":"wo","喔":"wo","窝":"wo","蜗":"wo","我":"wo","肟":"wo","沃":"wo","卧":"wo","握":"wo","幄":"wo","斡":"wo","乌":"wu","邬":"wu","污":"wu","巫":"wu","呜":"wu","钨":"wu","诬":"wu","屋":"wu","无":"wu","毋":"wu","芜":"wu","吴":"wu","梧":"wu","蜈":"wu","五":"wu","午":"wu","伍":"wu","仵":"wu","怃":"wu","忤":"wu","妩":"wu","武":"wu","侮":"wu","捂":"wu","鹉":"wu","舞":"wu","兀":"wu","勿":"wu","戊":"wu","务":"wu","坞":"wu","物":"wu","误":"wu","悟":"wu","晤":"wu","骛":"wu","雾":"wu","寤":"wu","鹜":"wu","夕":"xi","兮":"xi","西":"xi","吸":"xi","汐":"xi","希":"xi","昔":"xi","析":"xi","唏":"xi","牺":"xi","息":"xi","奚":"xi","悉":"xi","烯":"xi","惜":"xi","晰":"xi","稀":"xi","翕":"xi","犀":"xi","皙":"xi","锡":"xi","溪":"xi","熙":"xi","蜥":"xi","熄":"xi","嘻":"xi","膝":"xi","嬉":"xi","羲":"xi","蟋":"xi","曦":"xi","习":"xi","席":"xi","袭":"xi","媳":"xi","洗":"xi","玺":"xi","徙":"xi","喜":"xi","禧":"xi","戏":"xi","细":"xi","隙":"xi","呷":"xia","虾":"xia","瞎":"xia","匣":"xia","侠":"xia","峡":"xia","狭":"xia","遐":"xia","瑕":"xia","暇":"xia","辖":"xia","霞":"xia","黠":"xia","下":"xia","夏":"xia","罅":"xia","仙":"xian","先":"xian","氙":"xian","掀":"xian","酰":"xian","锨":"xian","鲜":"xian","闲":"xian","贤":"xian","弦":"xian","咸":"xian","涎":"xian","娴":"xian","衔":"xian","舷":"xian","嫌":"xian","显":"xian","险":"xian","跣":"xian","藓":"xian","苋":"xian","县":"xian","现":"xian","限":"xian","线":"xian","宪":"xian","陷":"xian","馅":"xian","羡":"xian","献":"xian","腺":"xian","乡":"xiang","相":"xiang","香":"xiang","厢":"xiang","湘":"xiang","箱":"xiang","襄":"xiang","镶":"xiang","详":"xiang","祥":"xiang","翔":"xiang","享":"xiang","响":"xiang","饷":"xiang","飨":"xiang","想":"xiang","向":"xiang","项":"xiang","象":"xiang","像":"xiang","橡":"xiang","肖":"xiao","枭":"xiao","哓":"xiao","骁":"xiao","逍":"xiao","消":"xiao","宵":"xiao","萧":"xiao","硝":"xiao","销":"xiao","箫":"xiao","潇":"xiao","霄":"xiao","魈":"xiao","嚣":"xiao","崤":"xiao","淆":"xiao","小":"xiao","晓":"xiao","孝":"xiao","哮":"xiao","笑":"xiao","效":"xiao","啸":"xiao","挟":"xie","些":"xie","楔":"xie","歇":"xie","蝎":"xie","协":"xie","胁":"xie","偕":"xie","斜":"xie","谐":"xie","揳":"xie","携":"xie","撷":"xie","鞋":"xie","写":"xie","泄":"xie","泻":"xie","卸":"xie","屑":"xie","械":"xie","亵":"xie","谢":"xie","邂":"xie","懈":"xie","蟹":"xie","心":"xin","芯":"xin","辛":"xin","欣":"xin","锌":"xin","新":"xin","歆":"xin","薪":"xin","馨":"xin","鑫":"xin","信":"xin","衅":"xin","星":"xing","猩":"xing","惺":"xing","腥":"xing","刑":"xing","邢":"xing","形":"xing","型":"xing","醒":"xing","擤":"xing","兴":"xing","杏":"xing","幸":"xing","性":"xing","姓":"xing","悻":"xing","凶":"xiong","兄":"xiong","匈":"xiong","讻":"xiong","汹":"xiong","胸":"xiong","雄":"xiong","熊":"xiong","休":"xiu","咻":"xiu","修":"xiu","羞":"xiu","朽":"xiu","秀":"xiu","袖":"xiu","绣":"xiu","锈":"xiu","嗅":"xiu","欻":"xu","戌":"xu","须":"xu","胥":"xu","虚":"xu","墟":"xu","需":"xu","魆":"xu","徐":"xu","许":"xu","诩":"xu","栩":"xu","旭":"xu","序":"xu","叙":"xu","恤":"xu","酗":"xu","勖":"xu","绪":"xu","续":"xu","絮":"xu","婿":"xu","蓄":"xu","煦":"xu","轩":"xuan","宣":"xuan","揎":"xuan","喧":"xuan","暄":"xuan","玄":"xuan","悬":"xuan","旋":"xuan","漩":"xuan","璇":"xuan","选":"xuan","癣":"xuan","炫":"xuan","绚":"xuan","眩":"xuan","渲":"xuan","靴":"xue","薛":"xue","穴":"xue","学":"xue","噱":"xue","雪":"xue","谑":"xue","勋":"xun","熏":"xun","薰":"xun","醺":"xun","旬":"xun","寻":"xun","巡":"xun","询":"xun","荀":"xun","循":"xun","训":"xun","讯":"xun","汛":"xun","迅":"xun","驯":"xun","徇":"xun","逊":"xun","殉":"xun","巽":"xun","丫":"ya","压":"ya","押":"ya","鸦":"ya","桠":"ya","鸭":"ya","牙":"ya","伢":"ya","芽":"ya","蚜":"ya","崖":"ya","涯":"ya","睚":"ya","衙":"ya","哑":"ya","雅":"ya","亚":"ya","讶":"ya","娅":"ya","氩":"ya","揠":"ya","呀":"ya","恹":"yan","胭":"yan","烟":"yan","焉":"yan","阉":"yan","淹":"yan","湮":"yan","嫣":"yan","延":"yan","闫":"yan","严":"yan","言":"yan","妍":"yan","岩":"yan","炎":"yan","沿":"yan","研":"yan","盐":"yan","阎":"yan","蜒":"yan","筵":"yan","颜":"yan","檐":"yan","奄":"yan","俨":"yan","衍":"yan","掩":"yan","郾":"yan","眼":"yan","偃":"yan","演":"yan","魇":"yan","鼹":"yan","厌":"yan","砚":"yan","彦":"yan","艳":"yan","晏":"yan","唁":"yan","宴":"yan","验":"yan","谚":"yan","堰":"yan","雁":"yan","焰":"yan","滟":"yan","餍":"yan","燕":"yan","赝":"yan","央":"yang","泱":"yang","殃":"yang","鸯":"yang","秧":"yang","扬":"yang","羊":"yang","阳":"yang","杨":"yang","佯":"yang","疡":"yang","徉":"yang","洋":"yang","仰":"yang","养":"yang","氧":"yang","痒":"yang","怏":"yang","样":"yang","恙":"yang","烊":"yang","漾":"yang","幺":"yao","夭":"yao","吆":"yao","妖":"yao","腰":"yao","邀":"yao","爻":"yao","尧":"yao","肴":"yao","姚":"yao","窑":"yao","谣":"yao","摇":"yao","徭":"yao","遥":"yao","瑶":"yao","杳":"yao","咬":"yao","舀":"yao","窈":"yao","药":"yao","要":"yao","鹞":"yao","耀":"yao","耶":"ye","掖":"ye","椰":"ye","噎":"ye","爷":"ye","揶":"ye","也":"ye","冶":"ye","野":"ye","业":"ye","叶":"ye","页":"ye","曳":"ye","夜":"ye","液":"ye","谒":"ye","腋":"ye","一":"yi","伊":"yi","衣":"yi","医":"yi","依":"yi","咿":"yi","揖":"yi","壹":"yi","漪":"yi","噫":"yi","仪":"yi","夷":"yi","饴":"yi","宜":"yi","咦":"yi","贻":"yi","姨":"yi","胰":"yi","移":"yi","痍":"yi","颐":"yi","疑":"yi","彝":"yi","乙":"yi","已":"yi","以":"yi","苡":"yi","矣":"yi","迤":"yi","蚁":"yi","倚":"yi","椅":"yi","旖":"yi","乂":"yi","亿":"yi","义":"yi","艺":"yi","刈":"yi","忆":"yi","议":"yi","屹":"yi","亦":"yi","异":"yi","抑":"yi","呓":"yi","邑":"yi","役":"yi","译":"yi","易":"yi","诣":"yi","绎":"yi","驿":"yi","轶":"yi","弈":"yi","奕":"yi","疫":"yi","羿":"yi","益":"yi","谊":"yi","逸":"yi","翌":"yi","肄":"yi","裔":"yi","意":"yi","溢":"yi","缢":"yi","毅":"yi","薏":"yi","翳":"yi","臆":"yi","翼":"yi","因":"yin","阴":"yin","茵":"yin","荫":"yin","音":"yin","姻":"yin","铟":"yin","喑":"yin","愔":"yin","吟":"yin","垠":"yin","银":"yin","淫":"yin","寅":"yin","龈":"yin","霪":"yin","尹":"yin","引":"yin","蚓":"yin","隐":"yin","瘾":"yin","印":"yin","英":"ying","莺":"ying","婴":"ying","嘤":"ying","罂":"ying","缨":"ying","樱":"ying","鹦":"ying","膺":"ying","鹰":"ying","迎":"ying","茔":"ying","荧":"ying","盈":"ying","莹":"ying","萤":"ying","营":"ying","萦":"ying","楹":"ying","蝇":"ying","赢":"ying","瀛":"ying","颍":"ying","颖":"ying","影":"ying","应":"ying","映":"ying","硬":"ying","哟":"yo","唷":"yo","佣":"yong","拥":"yong","庸":"yong","雍":"yong","壅":"yong","臃":"yong","永":"yong","甬":"yong","咏":"yong","泳":"yong","勇":"yong","涌":"yong","恿":"yong","蛹":"yong","踊":"yong","用":"yong","优":"you","攸":"you","忧":"you","呦":"you","幽":"you","悠":"you","尤":"you","由":"you","邮":"you","犹":"you","油":"you","铀":"you","鱿":"you","游":"you","友":"you","有":"you","酉":"you","莠":"you","黝":"you","又":"you","右":"you","幼":"you","佑":"you","柚":"you","囿":"you","诱":"you","鼬":"you","迂":"yu","纡":"yu","於":"yu","淤":"yu","瘀":"yu","于":"yu","余":"yu","盂":"yu","臾":"yu","鱼":"yu","竽":"yu","俞":"yu","狳":"yu","谀":"yu","娱":"yu","渔":"yu","隅":"yu","揄":"yu","逾":"yu","腴":"yu","渝":"yu","愉":"yu","瑜":"yu","榆":"yu","虞":"yu","愚":"yu","舆":"yu","与":"yu","予":"yu","屿":"yu","宇":"yu","羽":"yu","雨":"yu","禹":"yu","语":"yu","圄":"yu","玉":"yu","驭":"yu","芋":"yu","妪":"yu","郁":"yu","育":"yu","狱":"yu","浴":"yu","预":"yu","域":"yu","欲":"yu","谕":"yu","遇":"yu","喻":"yu","御":"yu","寓":"yu","裕":"yu","愈":"yu","誉":"yu","豫":"yu","鹬":"yu","鸢":"yuan","鸳":"yuan","冤":"yuan","渊":"yuan","元":"yuan","园":"yuan","垣":"yuan","袁":"yuan","原":"yuan","圆":"yuan","援":"yuan","媛":"yuan","缘":"yuan","猿":"yuan","源":"yuan","辕":"yuan","远":"yuan","苑":"yuan","怨":"yuan","院":"yuan","愿":"yuan","曰":"yue","月":"yue","岳":"yue","钺":"yue","阅":"yue","悦":"yue","跃":"yue","越":"yue","粤":"yue","晕":"yun","云":"yun","匀":"yun","芸":"yun","纭":"yun","耘":"yun","允":"yun","陨":"yun","殒":"yun","孕":"yun","运":"yun","酝":"yun","愠":"yun","韵":"yun","蕴":"yun","熨":"yun","匝":"za","咂":"za","杂":"za","砸":"za","灾":"zai","甾":"zai","哉":"zai","栽":"zai","载":"zai","宰":"zai","崽":"zai","再":"zai","在":"zai","糌":"zan","簪":"zan","咱":"zan","趱":"zan","暂":"zan","錾":"zan","赞":"zan","赃":"zang","脏":"zang","臧":"zang","驵":"zang","葬":"zang","遭":"zao","糟":"zao","凿":"zao","早":"zao","枣":"zao","蚤":"zao","澡":"zao","藻":"zao","皂":"zao","灶":"zao","造":"zao","噪":"zao","燥":"zao","躁":"zao","则":"ze","责":"ze","泽":"ze","啧":"ze","帻":"ze","仄":"ze","贼":"zei","怎":"zen","谮":"zen","增":"zeng","憎":"zeng","锃":"zeng","赠":"zeng","甑":"zeng","吒":"zha","挓":"zha","哳":"zha","揸":"zha","渣":"zha","楂":"zha","札":"zha","闸":"zha","铡":"zha","眨":"zha","砟":"zha","乍":"zha","诈":"zha","咤":"zha","炸":"zha","蚱":"zha","榨":"zha","拃":"zha","斋":"zhai","摘":"zhai","宅":"zhai","窄":"zhai","债":"zhai","砦":"zhai","寨":"zhai","沾":"zhan","毡":"zhan","粘":"zhan","詹":"zhan","谵":"zhan","瞻":"zhan","斩":"zhan","盏":"zhan","展":"zhan","崭":"zhan","搌":"zhan","辗":"zhan","占":"zhan","栈":"zhan","战":"zhan","站":"zhan","绽":"zhan","湛":"zhan","蘸":"zhan","张":"zhang","章":"zhang","獐":"zhang","彰":"zhang","樟":"zhang","蟑":"zhang","涨":"zhang","掌":"zhang","丈":"zhang","仗":"zhang","杖":"zhang","帐":"zhang","账":"zhang","胀":"zhang","障":"zhang","嶂":"zhang","瘴":"zhang","钊":"zhao","招":"zhao","昭":"zhao","找":"zhao","沼":"zhao","兆":"zhao","诏":"zhao","赵":"zhao","照":"zhao","罩":"zhao","肇":"zhao","蜇":"zhe","遮":"zhe","哲":"zhe","辄":"zhe","蛰":"zhe","谪":"zhe","辙":"zhe","者":"zhe","锗":"zhe","赭":"zhe","褶":"zhe","浙":"zhe","蔗":"zhe","鹧":"zhe","贞":"zhen","针":"zhen","侦":"zhen","珍":"zhen","帧":"zhen","胗":"zhen","真":"zhen","砧":"zhen","斟":"zhen","甄":"zhen","榛":"zhen","箴":"zhen","臻":"zhen","诊":"zhen","枕":"zhen","疹":"zhen","缜":"zhen","阵":"zhen","鸩":"zhen","振":"zhen","朕":"zhen","赈":"zhen","震":"zhen","镇":"zhen","争":"zheng","征":"zheng","怔":"zheng","峥":"zheng","狰":"zheng","睁":"zheng","铮":"zheng","筝":"zheng","蒸":"zheng","拯":"zheng","整":"zheng","正":"zheng","证":"zheng","郑":"zheng","诤":"zheng","政":"zheng","挣":"zheng","症":"zheng","之":"zhi","支":"zhi","只":"zhi","汁":"zhi","芝":"zhi","吱":"zhi","枝":"zhi","知":"zhi","肢":"zhi","织":"zhi","栀":"zhi","脂":"zhi","蜘":"zhi","执":"zhi","直":"zhi","侄":"zhi","值":"zhi","职":"zhi","植":"zhi","跖":"zhi","踯":"zhi","止":"zhi","旨":"zhi","址":"zhi","芷":"zhi","纸":"zhi","祉":"zhi","指":"zhi","枳":"zhi","咫":"zhi","趾":"zhi","酯":"zhi","至":"zhi","志":"zhi","豸":"zhi","帜":"zhi","制":"zhi","质":"zhi","炙":"zhi","治":"zhi","栉":"zhi","峙":"zhi","挚":"zhi","桎":"zhi","致":"zhi","秩":"zhi","掷":"zhi","痔":"zhi","窒":"zhi","蛭":"zhi","智":"zhi","痣":"zhi","滞":"zhi","置":"zhi","雉":"zhi","稚":"zhi","中":"zhong","忠":"zhong","终":"zhong","盅":"zhong","钟":"zhong","衷":"zhong","肿":"zhong","冢":"zhong","踵":"zhong","仲":"zhong","众":"zhong","舟":"zhou","州":"zhou","诌":"zhou","周":"zhou","洲":"zhou","粥":"zhou","妯":"zhou","轴":"zhou","肘":"zhou","纣":"zhou","咒":"zhou","宙":"zhou","胄":"zhou","昼":"zhou","皱":"zhou","骤":"zhou","帚":"zhou","朱":"zhu","侏":"zhu","诛":"zhu","茱":"zhu","珠":"zhu","株":"zhu","诸":"zhu","铢":"zhu","猪":"zhu","蛛":"zhu","竹":"zhu","竺":"zhu","逐":"zhu","烛":"zhu","躅":"zhu","主":"zhu","拄":"zhu","煮":"zhu","嘱":"zhu","瞩":"zhu","伫":"zhu","苎":"zhu","助":"zhu","住":"zhu","贮":"zhu","注":"zhu","驻":"zhu","柱":"zhu","祝":"zhu","著":"zhu","蛀":"zhu","铸":"zhu","筑":"zhu","抓":"zhua","跩":"zhuai","拽":"zhuai","专":"zhuan","砖":"zhuan","转":"zhuan","啭":"zhuan","撰":"zhuan","篆":"zhuan","妆":"zhuang","庄":"zhuang","桩":"zhuang","装":"zhuang","壮":"zhuang","状":"zhuang","撞":"zhuang","幢":"zhuang","追":"zhui","骓":"zhui","锥":"zhui","坠":"zhui","缀":"zhui","惴":"zhui","赘":"zhui","谆":"zhun","准":"zhun","拙":"zhuo","捉":"zhuo","桌":"zhuo","灼":"zhuo","茁":"zhuo","卓":"zhuo","斫":"zhuo","浊":"zhuo","酌":"zhuo","啄":"zhuo","擢":"zhuo","镯":"zhuo","孜":"zi","咨":"zi","姿":"zi","赀":"zi","资":"zi","辎":"zi","嗞":"zi","滋":"zi","锱":"zi","龇":"zi","子":"zi","姊":"zi","秭":"zi","籽":"zi","梓":"zi","紫":"zi","訾":"zi","滓":"zi","自":"zi","字":"zi","恣":"zi","眦":"zi","渍":"zi","宗":"zong","综":"zong","棕":"zong","踪":"zong","鬃":"zong","总":"zong","纵":"zong","粽":"zong","邹":"zou","走":"zou","奏":"zou","揍":"zou","租":"zu","足":"zu","卒":"zu","族":"zu","诅":"zu","阻":"zu","组":"zu","俎":"zu","祖":"zu","纂":"zuan","钻":"zuan","攥":"zuan","嘴":"zui","最":"zui","罪":"zui","醉":"zui","尊":"zun","遵":"zun","樽":"zun","鳟":"zun","昨":"zuo","左":"zuo","佐":"zuo","作":"zuo","坐":"zuo","阼":"zuo","怍":"zuo","祚":"zuo","唑":"zuo","座":"zuo","做":"zuo","酢":"zuo","斌":"bin","曾":"zeng","查":"zha","査":"zha","乘":"cheng","传":"chuan","丁":"ding","行":"xing","瑾":"jin","婧":"jing","恺":"kai","阚":"kan","奎":"kui","乐":"le","陆":"lu","逯":"lv","璐":"lu","淼":"miao","闵":"min","娜":"na","奇":"qi","琦":"qi","强":"qiang","邱":"qiu","芮":"rui","莎":"sha","盛":"sheng","石":"shi","祎":"yi","殷":"yin","瑛":"ying","昱":"yu","眃":"yun","琢":"zhuo","枰":"ping","玟":"min","珉":"min","珣":"xun","淇":"qi","缈":"miao","彧":"yu","祺":"qi","骞":"qian","垚":"yao","妸":"e","烜":"hui","祁":"qi","傢":"jia","珮":"pei","濮":"pu","屺":"qi","珅":"shen","缇":"ti","霈":"pei","晞":"xi","璠":"fan","骐":"qi","姞":"ji","偲":"cai","齼":"chu","宓":"mi","朴":"pu","萁":"qi","颀":"qi","阗":"tian","湉":"tian","翀":"chong","岷":"min","桤":"qi","囯":"guo","浛":"han","勐":"meng","苠":"min","岍":"qian","皞":"hao","岐":"qi","溥":"pu","锘":"muo","渼":"mei","燊":"shen","玚":"chang","亓":"qi","湋":"wei","涴":"wan","沤":"ou","胖":"pang","莆":"pu","扦":"qian","僳":"su","坍":"tan","锑":"ti","嚏":"ti","腆":"tian","丿":"pie","鼗":"tao","芈":"mi","匚":"fang","刂":"li","冂":"tong","亻":"dan","仳":"pi","俜":"ping","俳":"pai","倜":"ti","傥":"tang","傩":"nuo","佥":"qian","勹":"bao","亠":"tou","廾":"gong","匏":"pao","扌":"ti","拚":"pin","掊":"pou","搦":"nuo","擗":"pi","啕":"tao","嗦":"suo","嗍":"suo","辔":"pei","嘌":"piao","嗾":"sou","嘧":"mi","帔":"pei","帑":"tang","彡":"san","犭":"fan","狍":"pao","狲":"sun","狻":"jun","飧":"sun","夂":"zhi","饣":"shi","庀":"pi","忄":"shu","愫":"su","闼":"ta","丬":"jiang","氵":"san","汔":"qi","沔":"mian","汨":"mi","泮":"pan","洮":"tao","涑":"su","淠":"pi","湓":"pen","溻":"ta","溏":"tang","濉":"sui","宀":"bao","搴":"qian","辶":"zou","逄":"pang","逖":"ti","遢":"ta","邈":"miao","邃":"sui","彐":"ji","屮":"cao","娑":"suo","嫖":"piao","纟":"jiao","缗":"min","瑭":"tang","杪":"miao","桫":"suo","榀":"pin","榫":"sun","槭":"qi","甓":"pi","攴":"po","耆":"qi","牝":"pin","犏":"pian","氆":"pu","攵":"fan","肽":"tai","胼":"pian","脒":"mi","脬":"pao","旆":"pei","炱":"tai","燧":"sui","灬":"biao","礻":"shi","祧":"tiao","忑":"te","忐":"tan","愍":"min","肀":"yu","碛":"qi","眄":"mian","眇":"miao","眭":"sui","睃":"suo","瞍":"sou","畋":"tian","罴":"pi","蠓":"meng","蠛":"mie","笸":"po","筢":"pa","衄":"nv","艋":"meng","敉":"mi","糸":"mi","綦":"qi","醅":"pei","醣":"tang","趿":"ta","觫":"su","龆":"tiao","鲆":"ping","稣":"su","鲐":"tai","鲦":"tiao","鳎":"ta","髂":"qia","縻":"mi","裒":"pou","冫":"liang","冖":"tu","讠":"yan","谇":"sui","谝":"pian","谡":"su","卩":"dan","阝":"zuo","陴":"pi","邳":"pi","郫":"pi","郯":"tan","廴":"yin","凵":"qian","圮":"pi","堋":"peng","鼙":"pi","艹":"cao","芑":"qi","苤":"pie","荪":"sun","荽":"sui","葜":"qia","蒎":"pai","蔌":"su","蕲":"qi","薮":"sou","薹":"tai","蘼":"mi","钅":"jin","钷":"po","钽":"tan","铍":"pi","铴":"tang","铽":"te","锫":"pei","锬":"tan","锼":"sou","镤":"pu","镨":"pu","皤":"po","鹈":"ti","鹋":"miao","疒":"bing","疱":"pao","衤":"yi","袢":"pan","裼":"ti","襻":"pan","耥":"tang","耦":"ou","虍":"hu","蛴":"qi","蜞":"qi","蜱":"pi","螋":"sou","螗":"tang","螵":"piao","蟛":"peng",}
//获取汉字拼音,支持返回全拼音和拼音首字母,同时支持sp分隔符 php.GetPinyin(s,"")
func GetPinyin(s,sp string)(pinyin,shortpinyin string){
	n:=utf8.RuneCountInString(s)
	var list []string
	for i:=0;i<n;i++{
		list=append(list,Cutstr(s,i,1))
	}
	for _, v := range list {
		pin,isok :=dict[v]
		if isok{
			pinyin+=pin+sp
			shortpinyin+=Cutstr(pin,0,1)
		}else{
			pinyin+=v
			shortpinyin+=Cutstr(v,0,1)
		}
		
	}
	return pinyin,shortpinyin
}
