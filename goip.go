package goip

import (
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jiazhoulvke/goutil"
)

type ipDB struct {
	path   string
	length uint32
	index  []byte
}

var ipdb ipDB

//Location 地理位置
type Location struct {
	Country  string
	Province string
	City     string
	District string
}

//SetDBPath 设置ip数据库文件所在位置
func SetDBPath(p string) error {
	if !goutil.IsExist(p) {
		return fmt.Errorf("ip database file [%s] is not found", p)
	}
	f, err := os.Open(p)
	if err != nil {
		return fmt.Errorf("open database error: %v", err)
	}
	defer f.Close()
	bs := make([]byte, 4, 4)
	_, err = f.Read(bs)
	if err != nil {
		return fmt.Errorf("read database error: %v", err)
	}
	length := binary.BigEndian.Uint32(bs)
	if length < 4 {
		return fmt.Errorf("database broken")
	}
	index := make([]byte, length-4, length-4)
	_, err = f.Read(index)
	if err != nil {
		return fmt.Errorf("read database error: %v", err)
	}
	ipdb = ipDB{
		path:   p,
		length: length,
		index:  index,
	}
	return nil
}

//Length length
func Length() int {
	return len(ipdb.index)
}

//Find 查询IP所在地址
func Find(ip string) (Location, error) {
	ipdot := strings.Split(ip, ".")
	var location Location
	if len(ipdot) != 4 {
		return location, fmt.Errorf("ip format error")
	}
	ipInt := [4]int{}
	for i := 0; i < 4; i++ {
		n, err := strconv.Atoi(ipdot[i])
		if err != nil || n > 255 || n < 0 {
			return location, fmt.Errorf("ip format error")
		}
		ipInt[i] = n
	}
	tmpOffset := ipInt[0] * 4
	start := binary.LittleEndian.Uint32(ipdb.index[tmpOffset : tmpOffset+4])
	maxCompLen := ipdb.length - 1024 - 4
	ipn, err := IPv4ToInt(ip)
	if err != nil {
		return location, err
	}
	var indexOffset uint32
	var indexLength uint32
	found := false
	nip := make([]byte, 4, 4)
	binary.BigEndian.PutUint32(nip, uint32(ipn))
	for start = start*8 + 1024; start < maxCompLen; start += 8 {
		a := ipdb.index[start : start+4]
		if string(a) >= string(nip) {
			//if strings.Compare(string(a), string(nip)) != -1 {
			indexOffset = binary.LittleEndian.Uint32([]byte{ipdb.index[start+4], ipdb.index[start+5], ipdb.index[start+6], byte(0)})
			indexLength = uint32(uint8(ipdb.index[start+7]))
			found = true
			break
		}
	}
	if !found {
		return location, fmt.Errorf("not found")
	}
	f, err := os.Open(ipdb.path)
	if err != nil {
		return location, err
	}
	defer f.Close()
	_, err = f.Seek(int64(ipdb.length+indexOffset-1024), 0)
	if err != nil {
		return location, err
	}
	bs := make([]byte, indexLength, indexLength)
	_, err = f.Read(bs)
	if err != nil {
		return location, err
	}
	data := strings.Split(string(bs), "\t")
	location.Country = data[0]
	location.Province = data[1]
	location.City = data[2]
	location.District = data[3]
	return location, nil
}

//IPv4ToInt 转换ipv4地址为int
func IPv4ToInt(ip string) (int, error) {
	numbers := strings.Split(ip, ".")
	if len(numbers) != 4 {
		return 0, fmt.Errorf("ip format error")
	}
	var nList [4]int
	var err error
	var result int
	for i := 0; i < 4; i++ {
		nList[i], err = strconv.Atoi(numbers[i])
		if err != nil {
			return 0, err
		}
	}
	result = (nList[0] << 24) | (nList[1] << 16) | (nList[2] << 8) | nList[3]
	return result, nil
}

//IntToIPv4 转换int为ipv4地址
func IntToIPv4(ipInt int) string {
	ip1 := (ipInt >> 24) & 0xff
	ip2 := (ipInt >> 16) & 0xff
	ip3 := (ipInt >> 8) & 0xff
	ip4 := ipInt & 0xff
	return fmt.Sprintf("%d.%d.%d.%d", ip1, ip2, ip3, ip4)
}
