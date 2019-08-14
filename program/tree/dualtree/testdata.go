package dualtree

//TERMINALS
var X1 = NodeType{kind: 0, value: "x", arity: 0}
var Const4 = NodeType{kind: 0, value: "4", arity: 0}
var Const8 = NodeType{kind: 0, value: "8", arity: 0}

// NON-TERMINALS
var Sub = NodeType{kind: 1, value: "-", arity: 2}
var Add = NodeType{kind: 1, value: "+", arity: 2}
var Mult = NodeType{kind: 1, value: "*", arity: 2}
var Sin = NodeType{kind: 1, value: "Sin", arity: 1}

// SAMPLE TREES

// TreeNil = 0
var TreeNil = func() *DualTree {
	t := DualTree{}
	return &t
}

// Tree0 = x
var Tree0 = func() *DualTree {
	t := DualTree{}
	t.root = X1.ToDualTreeNode(0)
	return &t
}

// Tree1 = x * 4
var Tree1 = func() *DualTree {
	t := DualTree{}
	t.root = Mult.ToDualTreeNode(0)
	t.root.left = X1.ToDualTreeNode(1)
	t.root.right = Const4.ToDualTreeNode(2)
	return &t
}

// Tree2 = x - x * 4
var Tree2 = func() *DualTree {
	t := DualTree{}
	t.root = Mult.ToDualTreeNode(0)
	t.root.right = Const4.ToDualTreeNode(4)
	t.root.left = Sub.ToDualTreeNode(1)
	t.root.left.left = X1.ToDualTreeNode(2)
	t.root.left.right = X1.ToDualTreeNode(3)
	return &t
}

// Tree3 = Sin(4 - x)
var Tree3 = func() *DualTree {
	t := DualTree{}
	t.root = Sin.ToDualTreeNode(0)
	t.root.left = Sub.ToDualTreeNode(1)
	t.root.left.left = X1.ToDualTreeNode(2)
	t.root.left.right = Const4.ToDualTreeNode(3)
	return &t
}

// Tree4 = Sin(x)
var Tree4 = func() *DualTree {
	t := DualTree{}
	t.root = Sin.ToDualTreeNode(0)
	t.root.left = X1.ToDualTreeNode(1)
	return &t
}

// Tree5 = Sin
var Tree5 = func() *DualTree {
	t := DualTree{}
	t.root = Sin.ToDualTreeNode(0)
	return &t
}

// Tree6 = x +
var Tree6 = func() *DualTree {
	t := DualTree{}
	t.root = Add.ToDualTreeNode(0)
	t.root.left = X1.ToDualTreeNode(1)
	return &t
}

// Tree7 =  +
var Tree7 = func() *DualTree {
	t := DualTree{}
	t.root = Add.ToDualTreeNode(0)
	return &t
}

// Tree8 =  x * x
var Tree8 = func() *DualTree {
	t := DualTree{}
	t.root = Mult.ToDualTreeNode(0)
	t.root.left = X1.ToDualTreeNode(1)
	t.root.right = X1.ToDualTreeNode(2)
	return &t
}