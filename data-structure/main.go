package main

import "fmt"

func main() {

	arr := []int{10, 5, 24, 30, 60, 40, 45, 15, 27, 49, 23, 42, 56, 12, 8, 55, 2, 9}
	fmt.Println(arr)
	t := creatTree(arr)
	preorder(t[0])
	fmt.Println()
	inorder(t[0])
	fmt.Println()
	afterorder(t[0])

}

type tree struct {
	data  int
	left  *tree
	right *tree
}

//创建二叉树
func creatTree(arr []int) []tree {
	d := make([]tree, 0)
	for i, ar := range arr {
		d = append(d, tree{})
		d[i].data = ar
	}
	for i := 0; i < len(arr)/2; i++ {
		d[i].left = &d[i*2+1]
		if i*2+2 < len(d) {
			d[i].right = &d[i*2+2]
		}
	}
	fmt.Println(d)
	return d
}

//前序遍历
func preorder(root tree) {
	//fmt.Print(root.data, " ")
	if root.left != nil {
		preorder(*root.left)
	}
	if root.right != nil {
		preorder(*root.right)
	}
}

//中序遍历
func inorder(root tree) {
	if root.left != nil {
		inorder(*root.left)
	}
	fmt.Print(root.data, " ")
	if root.right != nil {
		inorder(*root.right)
	}
}

//后序遍历
func afterorder(root tree) {
	if root.left != nil {
		afterorder(*root.left)
	}
	if root.right != nil {
		afterorder(*root.right)
	}
	fmt.Print(root.data, " ")
}
