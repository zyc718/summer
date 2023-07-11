package example

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/assert/v2"
)

// 九九乘法表

func NineTable() {
	start := time.Now()
	for i := 1; i < 10; i++ {
		for j := 1; j <= i; j++ {
			fmt.Printf("%v x %v = %v", j, i, j*i)
			fmt.Printf(" | ")
		}
		fmt.Println()
	}
	tc := time.Since(start)

	fmt.Printf("过程时间是%v\n", tc)
}

// 最大公约数 最小公倍数

func MinMaxCommonDivisor() {
	start := time.Now()
	var a, b = 24, 10
	var c, d = a, b
	for a != b {
		if a > b {
			a = a - b
		} else if a < b {
			b = b - a
		}
	}
	fmt.Printf("%d 和 %d 最大的公约数是 %d \n", c, d, a)
	// 求出最大公约数。用两个数的乘积除以最大公约数 就是最小公倍数
	fmt.Printf("%d 和 %d 最小的公倍数是 %d \n", c, d, c*d/a)

	tc := time.Since(start)

	fmt.Printf("过程时间是%v\n", tc)
}

// 判断是不是回文
// 回文数的概念：即是给定一个数，这个数顺读和逆读都是一样的。例如：121，1221是回文数，123，1231不是回文数。

func Huiwen() {
	start := time.Now()
	str := "1221"
	var flag bool
	j := len(str) - 1
	for i := 0; i <= len(str)/2-1; i++ {
		if str[i] != str[j] {
			flag = false
			break
		}
		flag = true
		j--
	}
	fmt.Printf("输出的结果是%v\n", flag)
	tc := time.Since(start)
	fmt.Printf("过程时间是%v\n", tc)
}

// 水仙花数是指一个 3 位数，它的每个位上的数字的 3次幂之和等于它本身（例如：1^3 + 5^3+ 3^3 = 153）

func ShuiXianHua(num int) (result bool) {
	start := time.Now()
	// 分离出百位
	bai := num / 100
	// 分离出十位
	shi := (num / 10) % 10
	// 分离出个位
	ge := num % 10
	tc := time.Since(start)
	fmt.Printf("过程时间是%v\n", tc)
	fmt.Printf("输出的结果是%v\n", num == (bai*bai*bai+shi*shi*shi+ge*ge*ge))
	return num == (bai*bai*bai + shi*shi*shi + ge*ge*ge)

}

// 求同构书
// 正整数n若是它平方数的尾部，则称n为同构数。
// 例如：5的平方数是25，且5出现在25的右侧，那么5就是一个同构数。

func TongGou() {
	start := time.Now()
	// 分离出个位
	k := 10
	for i := 1; i <= 10000; i++ {
		// 小于10 求的是个位数，大于10 小于100 取得是最后两位数 以此类推
		if i == k {
			k *= 10
		}
		j := i * i
		if j%k == i {
			fmt.Printf("%v 是同构数，平方数 %v \n", i, j)
		}
	}
	tc := time.Since(start)
	fmt.Printf("过程时间是%v\n", tc)
}

// 翻转一个 int 类型的slice

func Reverse() {
	start := time.Now()

	var s = []int{1, 2, 3, 4, 5, 6, 7, 8}
	for i := 0; i < len(s)/2; i++ {
		s[i], s[len(s)-i-1] = s[len(s)-i-1], s[i]
	}
	fmt.Printf("输出的结果是%v\n", s)
	tc := time.Since(start)
	fmt.Printf("过程时间是%v\n", tc)
}

type StudentNode struct {
	No   int
	Data string
	Next *StudentNode
}

// 给链表插入一个节点 在链表的尾部插入

func PushNode(head *StudentNode, new *StudentNode) {
	// 创建一个辅助节点
	tmp := head
	// 找到最后一个节点
	for {
		if tmp.Next == nil {
			break
		}
		tmp = tmp.Next
	}

	// 将新节点加入到链表的最后
	tmp.Next = new

}

//5,4,2,1
// 根据插入方法，根据no的编号从大到小插入

func PushAfterNode(head *StudentNode, new *StudentNode) {
	tmp := head
	for {
		if tmp.No > new.No {
			if tmp.Next == nil {
				tmp.Next = new
				break
			}
		} else if tmp.No == new.No {
			break
		} else if tmp.No < new.No {
			// 找前一个
			previousNode := GetPrevious(head, tmp)
			previousNode.Next = new
			new.Next = tmp
			break
		}

		tmp = tmp.Next
	}

}

func GetPrevious(head, node *StudentNode) *StudentNode {
	tmp := head
	for {
		if head == node {
			return head
		}
		if tmp.Next == node {
			return tmp
		}
		tmp = tmp.Next
	}
}

// 替换字符串 转换大小写

func LowerToUpper() {
	start := time.Now()
	str1 := "sssnnnkabckoossabc"
	old := "abc"
	news := "xyz"
	var res []byte
	len1 := len(str1)
	lenOld := len(old)
	for i := 0; i < len(str1); i++ {
		if str1[i] == old[0] && i+lenOld-1 < len1 && str1[i:i+lenOld] == old {
			i += lenOld - 1
			for j := 0; j < len(news); j++ {
				res = append(res, news[j])
			}
		} else {
			res = append(res, str1[i])
		}

	}
	//str := strings.Replace(str1, "abc", "XYZ", -1)
	//str := strings.ToUpper(strings.Join(strings.Split("sssnnnkabckooss", "abc"), "XYZ"))
	fmt.Printf("输出的结果是%v\n", strings.ToUpper(string(res)))
	tc := time.Since(start)
	fmt.Printf("过程时间是%v\n", tc)
}

// 删除切片中的某个元素
func DeleteSlice() {
	start := time.Now()

	var s = []int{1, 2, 3, 4}
	var num = 4

	fmt.Printf("输出的结果是%v", s[4:])
	for i := 0; i < len(s); i++ {
		if assert.IsEqual(num, s[i]) {
			s = append(s[:i], s[i+1:]...)
		}
	}

	fmt.Printf("结果是%v\n", s)

	tc := time.Since(start)
	fmt.Printf("过程时间是%v\n", tc)
}
