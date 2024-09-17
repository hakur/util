package tree

func NewAVLTreeNode[VT string | int | uint64 | uint32 | uint16 | uint8 | int64 | int32 | int16 | int8](value VT) (t *AVLTreeNode[VT]) {
	t = new(AVLTreeNode[VT])
	t.Value = value
	return
}

type AVLTreeNode[VT string | int | uint64 | uint32 | uint16 | uint8 | int64 | int32 | int16 | int8] struct {
	Parent *AVLTreeNode[VT]
	Left   *AVLTreeNode[VT]
	Right  *AVLTreeNode[VT]
	Value  VT
}

func (t *AVLTreeNode[VT]) RotateLeft() {
	if t.Parent != nil { // 普通节点失衡
		var isLeftNode bool
		if t.Parent.Value > t.Value {
			isLeftNode = true
		}
		rightNode := t.Right
		parentNode := t.Parent
		if isLeftNode {
			t.Parent.Left = rightNode
		} else {
			t.Parent.Right = rightNode
		}

		rightNode.Parent = parentNode

		t.Right = rightNode.Left
		if t.Right != nil {
			t.Right.Parent = t
		}
		rightNode.Left = t
		t.Parent = rightNode
	} else { // root 节点失衡
		tempLeft := t.Right.Left
		t.Right.Left = t
		t.Right.Parent = t.Parent
		t.Parent = t.Right
		t.Right = tempLeft
	}
}

func (t *AVLTreeNode[VT]) RotateRight() {
	if t.Parent != nil { // 普通节点失衡
		var isLeftNode bool
		if t.Parent.Value > t.Value {
			isLeftNode = true
		}
		leftNode := t.Left
		parentNode := t.Parent

		if isLeftNode {
			t.Parent.Left = leftNode
		} else {
			t.Parent.Right = leftNode
		}
		leftNode.Parent = parentNode

		t.Left = leftNode.Right
		if t.Left != nil {
			t.Left.Parent = t
		}
		leftNode.Right = t
		t.Parent = leftNode
	} else { // root 节点失衡
		tempRight := t.Left.Right
		t.Left.Right = t
		t.Left.Parent = t.Parent
		t.Parent = t.Left
		t.Left = tempRight
	}
}

// NewAVLTree create new avl tree for fast search in memory, not for disk data search
// NewAVLTree 新建一个二叉平衡树来获取在内存中快速搜索，而不是为了磁盘数据搜索
func NewAVLTree[VT string | int | uint64 | uint32 | uint16 | uint8 | int64 | int32 | int16 | int8]() (t *AVLTree[VT]) {
	t = new(AVLTree[VT])
	return
}

type AVLTree[VT string | int | uint64 | uint32 | uint16 | uint8 | int64 | int32 | int16 | int8] struct {
	Root *AVLTreeNode[VT]
}

// BalanceFactor get balance factor of node
// BalanceFactor 计算节点平衡因子
func (t *AVLTree[VT]) BalanceFactor(node *AVLTreeNode[VT]) int {
	if node == nil {
		return 0
	}
	// 平衡因子BF = 左子树深度－右子树深度。
	return t.GetNodeHeight(node.Left) - t.GetNodeHeight(node.Right)
}

func (t *AVLTree[VT]) GetNodeHeight(node *AVLTreeNode[VT]) int {
	if node == nil {
		return 0
	}

	leftHeight := t.GetNodeHeight(node.Left)
	rightHeight := t.GetNodeHeight(node.Right)
	var height int = rightHeight
	if leftHeight > rightHeight {
		height = leftHeight
	}
	return height + 1
}

func (t *AVLTree[VT]) rotate(node *AVLTreeNode[VT]) {
	var nodeBF int = t.BalanceFactor(node)
	if nodeBF < 2 && nodeBF > -2 {
		return
	}

	var leftBF, rightBF int = t.BalanceFactor(node.Left), t.BalanceFactor(node.Right)

	if nodeBF == -2 && rightBF == -1 { // RR 型, 失衡节点平衡因子=-2，失衡节点右子节点平衡因子=-1,直接使用左旋来平衡
		node.RotateLeft()
		if node == t.Root {
			t.Root = node.Parent
		}
	} else if nodeBF == 2 && leftBF == 1 { // LL 型，失衡节点平衡因子=2，失衡节点左子节点平衡因子=1,直接使用右旋来平衡
		node.RotateRight()
		if node == t.Root {
			t.Root = node.Parent
		}
	} else if nodeBF == 2 && leftBF == -1 { // LR 型，失衡节点平衡因子=2，失衡节点左子节点平衡因子=-1，先左旋转左子节点，再右旋失衡节点
		node.Left.RotateLeft()
		node.RotateRight()
		if node == t.Root {
			t.Root = node.Parent
		}
	} else if nodeBF == -2 && rightBF == 1 { // RL 型，失衡节点平衡因子=-2，失衡节点右子节点平衡因子=1，先右旋转右子节点，再左旋失衡节点
		node.Right.RotateRight()
		node.RotateLeft()
		if node == t.Root {
			t.Root = node.Parent
		}
	}
}

// Insert insert new value to avl tree, current not support insert repeated value
// Insert 将新值插入到avl树，当前不支持重复值插入
func (t *AVLTree[VT]) Insert(value VT) {
	if t.Root == nil {
		t.Root = NewAVLTreeNode(value)
		return
	}
	insertedNode := t.insert(t.Root, value)
	if insertedNode == nil {
		return
	}
	// 插入节点后，如果右多个祖节点失衡，那么调整最近的一个祖节点即可。删除时则需要递归检查每一个parent节点是否失衡
	t.recursiveRotate(insertedNode, true)
}

func (t *AVLTree[VT]) recursiveRotate(node *AVLTreeNode[VT], stopAtFirst bool) {
	if node == nil {
		return
	}
	var nodeBF int = t.BalanceFactor(node)
	if nodeBF > 1 || nodeBF < -1 {
		t.rotate(node)
		if stopAtFirst {
			return
		}
	}
	if node.Parent != nil {
		t.recursiveRotate(node.Parent, stopAtFirst)
	}
}

func (t *AVLTree[VT]) insert(node *AVLTreeNode[VT], value VT) (insertedNode *AVLTreeNode[VT]) {
	// TODO 使用溢出页来存储重复的值
	if node.Value == value { // 已存在的值 就跳过
		return
	}

	if value > node.Value { // 走右节点
		if node.Right == nil {
			node.Right = NewAVLTreeNode(value)
			node.Right.Parent = node
			insertedNode = node.Right
		} else {
			return t.insert(node.Right, value)
		}
	} else if value < node.Value { // 走左节点
		if node.Left == nil {
			node.Left = NewAVLTreeNode(value)
			node.Left.Parent = node
			insertedNode = node.Left
		} else {
			return t.insert(node.Left, value)
		}
	}
	return insertedNode
}

func (t *AVLTree[VT]) Delete(value VT) {
	if t.Root == nil {
		return
	}

	node := t.Search(value)
	if node == nil {
		return
	}

	t.deleteNode(node)
}

func (t *AVLTree[VT]) deleteNode(node *AVLTreeNode[VT]) {
	if node == nil {
		return
	}
	var rotateNode *AVLTreeNode[VT] // 最后执行旋转检查的节点

	if node.Left == nil && node.Right == nil { // 如果被删除的节点没有左右子节点，则直接从这个节点的parent中删除
		if node.Parent == nil {
			if node == t.Root { // 删除root节点
				t.Root = nil
			}
		} else {
			if node.Parent.Value > node.Value {
				node.Parent.Left = nil
			} else {
				node.Parent.Right = nil
			}
			node.Parent = nil
			rotateNode = node.Parent // 检查parent节点是否失衡
		}
	} else if node.Left != nil && node.Right != nil { // 如果左右子节点均不为空，则执行节点值替换
		// https://time.geekbang.org/column/article/641400
		var bf = t.BalanceFactor(node) // 如果左子节点比右子节点更高，则走左子节点,如果相同高度则默认走左树
		var goLeft bool = true
		if bf < 0 {
			goLeft = false
		}

		var subNode *AVLTreeNode[VT]
		if goLeft {
			// 左子节点不空，则找出左子节点的最右节点的值替换到被删除节点的值上，并删除最右子节点，然后重新从最右子节点的parent开始平衡
			// 如果左子节点没有右节点,就让左子节点顶替
			subNode = t.getDirectionLeaf(node.Left, false)
		} else {
			// 右子节点不空，则找出右子节点的最左子节点的值替换到被删除节点的值上，并删除最左子节点，然后重新从最左子节点的parent开始平衡
			// 如果右子节点没有左子节点，就让右子节点顶替
			subNode = t.getDirectionLeaf(node.Right, true)
		}

		tempValue := subNode.Value
		rotateNode = subNode.Parent
		t.deleteNode(subNode)
		node.Value = tempValue
	} else if node.Left != nil || node.Right != nil { // 如果左右子节点存在一个，则执行节点替换
		if node.Left != nil && node.Right == nil { // 只有左节点
			if node.Parent != nil {
				if node.Parent.Value > node.Value { // 自身是parent的的左子节点
					node.Parent.Left = node.Left
					node.Left.Parent = node.Parent
				} else {
					node.Parent.Right = node.Left
					node.Right.Parent = node.Parent
				}
			}
			if node == t.Root {
				t.Root = node.Left
				t.Root.Parent = nil
			}
			node.Parent = nil
			rotateNode = node.Left
		} else if node.Right != nil && node.Left == nil {
			if node.Parent != nil {
				if node.Parent.Value > node.Value { // 自身是parent的的左子节点
					node.Parent.Left = node.Right
					node.Right.Parent = node.Parent
				} else {
					node.Parent.Right = node.Right
					node.Right.Parent = node.Parent
				}
			}
			if node == t.Root {
				t.Root = node.Right
				t.Root.Parent = nil
			}
			node.Parent = nil
			rotateNode = node.Right
		}
	}

	t.recursiveRotate(rotateNode, false)
	t.recursiveRotate(t.Root, false) // 手动平衡一下root节点，凸(艹皿艹 )
}

func (t *AVLTree[VT]) getDirectionLeaf(node *AVLTreeNode[VT], isLeft bool) (target *AVLTreeNode[VT]) {
	if isLeft {
		if node.Left != nil {
			return t.getDirectionLeaf(node.Left, isLeft)
		} else {
			return node
		}
	} else {
		if node.Right != nil {
			return t.getDirectionLeaf(node.Right, isLeft)
		} else {
			return node
		}
	}
}

// Search search value node on avl tree,
// Search 搜索值在avl树上的节点
func (t *AVLTree[VT]) Search(value VT) (target *AVLTreeNode[VT]) {
	return t.SearchFrom(t.Root, value)
}

func (t *AVLTree[VT]) SearchFrom(node *AVLTreeNode[VT], value VT) (target *AVLTreeNode[VT]) {
	if node == nil {
		return
	}

	if node.Value == value {
		return node
	}

	if value > node.Value {
		return t.SearchFrom(node.Right, value)
	} else {
		return t.SearchFrom(node.Left, value)
	}
}

func (t *AVLTree[VT]) Walk(startNode *AVLTreeNode[VT], callback func(node *AVLTreeNode[VT])) {
	if startNode == nil {
		startNode = t.Root
	}
	if startNode == nil {
		return
	}
	t.walk(startNode, callback)
}

func (t *AVLTree[VT]) walk(startNode *AVLTreeNode[VT], callback func(node *AVLTreeNode[VT])) {
	if startNode == nil {
		return
	}
	// https://baijiahao.baidu.com/s?id=1728183194697373686&wfr=spider&for=pc

	// 前序遍历 root->left->right
	// callback(startNode)
	// t.walk(startNode.Left, callback)
	// t.walk(startNode.Right, callback)

	// 中序 意义：从小到大遍历
	t.walk(startNode.Left, callback)
	callback(startNode)
	t.walk(startNode.Right, callback)

	// 倒序 意义：从大到小遍历
	// t.walk(startNode.Right, callback)
	// callback(startNode)
	// t.walk(startNode.Left, callback)

	// 后序
	// t.walk(startNode.Left, callback)
	// t.walk(startNode.Right, callback)
	// callback(startNode)
}
