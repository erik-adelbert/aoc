package rotate

func CopyRotateLeft(a []int) {
	n := a[0]
	copy(a, a[1:])
	a[len(a)-1] = n
}

func DKReverse(a []int) {
	for l, r := 0, len(a)-1; l < r; l, r = l+1, r-1 {
		a[l], a[r] = a[r], a[l]
	}
}

func DKRotateLeft(a []int) {
	DKReverse(a[:1])
	DKReverse(a[1:])
	DKReverse(a)
}
