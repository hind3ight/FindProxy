//
//	@Description
//	@return
//  @author hind3ight
//  @createdtime
//  @updatedtime

package ip

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func ParseCidr(cidr string) []string {
	ip := strings.Split(cidr, "/")[0]
	ipSegs := strings.Split(ip, ".")
	maskLen, _ := strconv.Atoi(strings.Split(cidr, "/")[1])
	seg3MinIp, seg3MaxIp := getIpSeg3Range(ipSegs, maskLen)
	seg4MinIp, seg4MaxIp := getIpSeg4Range(ipSegs, maskLen)
	ipPrefix := ipSegs[0] + "." + ipSegs[1] + "."

	var ipPool []string
	for i := seg3MinIp; i <= seg3MaxIp; i++ {
		for j := 0; j <= 255; j++ {
			if i == seg3MinIp && j < seg4MinIp {
				continue
			}

			if i == seg3MaxIp && j == seg4MaxIp {
				continue
			}

			one := ipPrefix + strconv.Itoa(i) + "." + strconv.Itoa(j)
			ipPool = append(ipPool, one)
		}

	}

	return ipPool
}

// 将补全的ip转为二进制
func parsePart(str string) string {
	tmp := strings.TrimLeft(str, "0")
	if tmp == "" {
		tmp = "0"
	}
	return tmp
}

func binaryToDecimal(num int) int {
	var remainder int
	index := 0
	decimalNum := 0
	for num != 0 {
		remainder = num % 10
		num = num / 10
		decimalNum = decimalNum + remainder*int(math.Pow(2, float64(index)))
		index++
	}
	return decimalNum
}

//将获取的IpAddr 字符串，转换为10进制
func IpToLong(ipAddress string) (string, uint64) {
	var ipSlice []string = strings.Split(ipAddress, ".")
	var res uint64 = 0

	for k, v := range ipSlice {
		ip, _ := strconv.ParseUint(v, 10, 64)
		i := uint64(k)
		res |= ip << ((3 - i) << 3)
	}
	return ipAddress, res
}

func parseIpSegSlice(ipSegSlice []string) (string, string) {
	var ipByteAllStr, ipByteStr string
	for _, ipSeg := range ipSegSlice {
		one := []rune("00000000")
		tmp1, _ := strconv.Atoi(ipSeg)
		ipByte := strconv.FormatInt(int64(tmp1), 2) // 转为二进制
		str := string(one[:8-len(ipByte)]) + ipByte // 补全位数
		ipByteStr += ipByte
		ipByteAllStr += str
	}
	fmt.Println(ipByteStr)
	fmt.Println(len(ipByteStr))
	return ipByteAllStr, ipByteStr
}

//计算得到CIDR地址范围内可拥有的主机数量
func getCidrHostNum(maskLen int) uint {
	cidrIpNum := uint(0)
	var i uint = uint(32 - maskLen - 1)
	for ; i >= 1; i-- {
		cidrIpNum += 1 << i
	}
	return cidrIpNum
}

//获取Cidr的掩码
func getCidrIpMask(maskLen int) string {
	// ^uint32(0)二进制为32个比特1，通过向左位移，得到CIDR掩码的二进制
	cidrMask := ^uint32(0) << uint(32-maskLen)
	fmt.Println(fmt.Sprintf("%b \n", cidrMask))
	//计算CIDR掩码的四个片段，将想要得到的片段移动到内存最低8位后，将其强转为8位整型，从而得到
	cidrMaskSeg1 := uint8(cidrMask >> 24)
	cidrMaskSeg2 := uint8(cidrMask >> 16)
	cidrMaskSeg3 := uint8(cidrMask >> 8)
	cidrMaskSeg4 := uint8(cidrMask & uint32(255))

	return fmt.Sprint(cidrMaskSeg1) + "." + fmt.Sprint(cidrMaskSeg2) + "." + fmt.Sprint(cidrMaskSeg3) + "." + fmt.Sprint(cidrMaskSeg4)
}

//得到第三段IP的区间（第一片段.第二片段.第三片段.第四片段）
func getIpSeg3Range(ipSegSlice []string, maskLen int) (int, int) {
	if maskLen > 24 {
		segIp, _ := strconv.Atoi(ipSegSlice[2])
		return segIp, segIp
	}
	ipSeg, _ := strconv.Atoi(ipSegSlice[2])
	return getIpSegRange(uint8(ipSeg), uint8(24-maskLen))
}

//得到第四段IP的区间（第一片段.第二片段.第三片段.第四片段）
func getIpSeg4Range(ipSegs []string, maskLen int) (int, int) {
	ipSeg, _ := strconv.Atoi(ipSegs[3])
	segMinIp, segMaxIp := getIpSegRange(uint8(ipSeg), uint8(32-maskLen))
	return segMinIp + 1, segMaxIp
}

//根据用户输入的基础IP地址和CIDR掩码计算一个IP片段的区间
func getIpSegRange(userSegIp, offset uint8) (int, int) {
	var ipSegMax uint8 = 255
	netSegIp := ipSegMax << offset
	segMinIp := netSegIp & userSegIp
	segMaxIp := userSegIp&(255<<offset) | ^(255 << offset)
	return int(segMinIp), int(segMaxIp)
}
