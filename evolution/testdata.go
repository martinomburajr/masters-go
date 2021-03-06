package evolution

//TERMINALS
var X1 = SymbolicExpression{kind: 0, value: "x", arity: 0}
var Const0 = SymbolicExpression{kind: 0, value: "0", arity: 0}
var Const1 = SymbolicExpression{kind: 0, value: "1", arity: 0}
var Const2 = SymbolicExpression{kind: 0, value: "2", arity: 0}
var Const3 = SymbolicExpression{kind: 0, value: "3", arity: 0}
var Const4 = SymbolicExpression{kind: 0, value: "4", arity: 0}
var Const5 = SymbolicExpression{kind: 0, value: "5", arity: 0}
var Const6 = SymbolicExpression{kind: 0, value: "6", arity: 0}
var Const7 = SymbolicExpression{kind: 0, value: "7", arity: 0}
var Const8 = SymbolicExpression{kind: 0, value: "8", arity: 0}
var Const9 = SymbolicExpression{kind: 0, value: "9", arity: 0}

// NON-TERMINALS
var Sub = SymbolicExpression{kind: 1, value: "-", arity: 2}
var Add = SymbolicExpression{kind: 1, value: "+", arity: 2}
var Mult = SymbolicExpression{kind: 1, value: "*", arity: 2}
var Sin = SymbolicExpression{kind: 1, value: "sin", arity: 1}

// SAMPLE TREES

var TreeBad = func() *DualTree {
	t := DualTree{}
	t.root = Mult.ToDualTreeNode(RandString(5))
	return &t
}

// TreeNil = 0
var TreeNil = func() *DualTree {
	t := DualTree{}
	return &t
}

// TreeT_X = x
var TreeT_X = func() *DualTree {
	t := DualTree{}
	t.root = X1.ToDualTreeNode("0")
	return &t
}

var TreeT_10 = func() *DualTree {
	t := DualTree{}
	t.root = Const10.ToDualTreeNode("0")
	return &t
}

var TreeT_100 = func() *DualTree {
	t := DualTree{}
	t.root = Const100.ToDualTreeNode("0")
	return &t
}

var TreeT_1000 = func() *DualTree {
	t := DualTree{}
	t.root = Const1000.ToDualTreeNode("0")
	return &t
}

var TreeT_10000 = func() *DualTree {
	t := DualTree{}
	t.root = Const10000.ToDualTreeNode("0")
	return &t
}

// TreeT_1 = 0
var TreeT_1 = func() *DualTree {
	t := DualTree{}
	t.root = Const0.ToDualTreeNode("0")
	return &t
}

// TreeT_NT_T_4 = x * x
var TreeT_NT_T_4 = func() *DualTree {
	t := DualTree{}
	t.root = Mult.ToDualTreeNode(RandString(5))
	t.root.left = X1.ToDualTreeNode(RandString(5))
	t.root.right = X1.ToDualTreeNode(RandString(5))
	return &t
}

var Const10 = SymbolicExpression{kind: 0, value: "10", arity: 0}
var Const100 = SymbolicExpression{kind: 0, value: "100", arity: 0}
var Const1000 = SymbolicExpression{kind: 0, value: "1000", arity: 0}
var Const10000 = SymbolicExpression{kind: 0, value: "10000", arity: 0}
var Const100000 = SymbolicExpression{kind: 0, value: "100000", arity: 0}

var Tree_X10 = func() *DualTree {
	t := DualTree{}
	t.root = Mult.ToDualTreeNode(RandString(5))
	t.root.left = X1.ToDualTreeNode(RandString(5))
	t.root.right = Const10.ToDualTreeNode(RandString(5))
	return &t
}

var Tree_X100 = func() *DualTree {
	t := DualTree{}
	t.root = Mult.ToDualTreeNode(RandString(5))
	t.root.left = X1.ToDualTreeNode(RandString(5))
	t.root.right = Const100.ToDualTreeNode(RandString(5))
	return &t
}

var Tree_X1000 = func() *DualTree {
	t := DualTree{}
	t.root = Mult.ToDualTreeNode(RandString(5))
	t.root.left = X1.ToDualTreeNode(RandString(5))
	t.root.right = Const1000.ToDualTreeNode(RandString(5))
	return &t
}

var Tree_X10000 = func() *DualTree {
	t := DualTree{}
	t.root = Mult.ToDualTreeNode(RandString(5))
	t.root.left = X1.ToDualTreeNode(RandString(5))
	t.root.right = Const10000.ToDualTreeNode(RandString(5))
	return &t
}

var Tree_X100000 = func() *DualTree {
	t := DualTree{}
	t.root = Mult.ToDualTreeNode(RandString(5))
	t.root.left = X1.ToDualTreeNode(RandString(5))
	t.root.right = Const100000.ToDualTreeNode(RandString(5))
	return &t
}

var Tree_100000XXX = func() *DualTree {
	t := DualTree{}
	t.root = Mult.ToDualTreeNode(RandString(5))
	t.root.left = Mult.ToDualTreeNode(RandString(5))
	t.root.left.left = Const100000.ToDualTreeNode(RandString(5))
	t.root.left.right = X1.ToDualTreeNode(RandString(5))
	t.root.right = Mult.ToDualTreeNode(RandString(5))
	t.root.right.left = X1.ToDualTreeNode(RandString(5))
	t.root.right.right = X1.ToDualTreeNode(RandString(5))
	return &t
}

// TreeT_NT_T_0 = x * 4
var TreeT_NT_T_0 = func() *DualTree {
	t := DualTree{}
	t.root = Mult.ToDualTreeNode("1")
	t.root.left = X1.ToDualTreeNode("2")
	t.root.right = Const4.ToDualTreeNode("3")
	return &t
}

// TreeXby5 = x * 5
var TreeXby5 = func() *DualTree {
	t := DualTree{}
	t.root = Mult.ToDualTreeNode("1fas")
	t.root.left = X1.ToDualTreeNode("2ds")
	t.root.right = Const5.ToDualTreeNode("3dd")
	return &t
}

// TreeT_NT_T_5 = x - x
var TreeT_NT_T_5 = func() *DualTree {
	t := DualTree{}
	t.root = Sub.ToDualTreeNode("2")
	t.root.left = X1.ToDualTreeNode("3")
	t.root.right = X1.ToDualTreeNode("4")
	return &t
}

// TreeT_NT_T_1 = x + 8
var TreeT_NT_T_1 = func() *DualTree {
	t := DualTree{}
	t.root = Add.ToDualTreeNode("1234")
	t.root.left = X1.ToDualTreeNode("12345")
	t.root.right = Const8.ToDualTreeNode("456")
	return &t
}

// TreeT_NT_T_2 = 4 - 8
var TreeT_NT_T_2 = func() *DualTree {
	t := DualTree{}
	t.root = Sub.ToDualTreeNode(RandString(5))
	t.root.left = Const4.ToDualTreeNode(RandString(5))
	t.root.right = Const8.ToDualTreeNode(RandString(5))
	return &t
}

// TreeT_NT_T_3 = 8 * 8
var TreeT_NT_T_3 = func() *DualTree {
	t := DualTree{}
	t.root = Mult.ToDualTreeNode(RandString(5))
	t.root.left = Const8.ToDualTreeNode(RandString(5))
	t.root.right = Const8.ToDualTreeNode(RandString(5))
	return &t
}

// TreeT_NT_T_NT_T_0 = x - x * 4
var TreeT_NT_T_NT_T_0 = func() *DualTree {
	t := DualTree{}
	t.root = Mult.ToDualTreeNode("0")
	t.root.right = Const4.ToDualTreeNode("1")
	t.root.left = Sub.ToDualTreeNode("2")
	t.root.left.left = X1.ToDualTreeNode("3")
	t.root.left.right = X1.ToDualTreeNode("4")
	return &t
}

// TreeT_NT_T_NT_T_0 = (((x + (x * 4)) - 9) * 0)
var TreeXAddXMult4Sub9Mult0 = func() *DualTree {
	t := DualTree{}
	t.root = Mult.ToDualTreeNode("0")
	t.root.right = Const0.ToDualTreeNode("1")
	t.root.left = Sub.ToDualTreeNode("2")
	t.root.left.right = Const9.ToDualTreeNode("4")
	t.root.left.left = Add.ToDualTreeNode("43")
	t.root.left.left.left = X1.ToDualTreeNode("3f3")
	t.root.left.left.right = Mult.ToDualTreeNode("4sd3")
	t.root.left.left.right.left = X1.ToDualTreeNode("4fd3")
	t.root.left.left.right.right = Const4.ToDualTreeNode("4x3")
	return &t
}

// TreeT_NT_T_NT_T_1 = x + 8 * 4
var TreeT_NT_T_NT_T_1 = func() *DualTree {
	t := DualTree{}
	t.root = Mult.ToDualTreeNode(RandString(5))
	t.root.right = Const4.ToDualTreeNode(RandString(5))
	t.root.left = Add.ToDualTreeNode(RandString(5))
	t.root.left.left = X1.ToDualTreeNode(RandString(5))
	t.root.left.right = Const8.ToDualTreeNode(RandString(5))
	return &t
}

// TreeT_NT_T_NT_T_2 = x - 0 * 4
var TreeT_NT_T_NT_T_2 = func() *DualTree {
	t := DualTree{}
	t.root = Mult.ToDualTreeNode(RandString(5))
	t.root.right = Const4.ToDualTreeNode(RandString(5))
	t.root.left = Sub.ToDualTreeNode(RandString(5))
	t.root.left.left = X1.ToDualTreeNode(RandString(5))
	t.root.left.right = Const0.ToDualTreeNode(RandString(5))
	return &t
}

// TreeT_NT_T_NT_T_3 = x - 0 * 0
var TreeT_NT_T_NT_T_3 = func() *DualTree {
	t := DualTree{}
	t.root = Mult.ToDualTreeNode("0")
	t.root.right = Const0.ToDualTreeNode("2")
	t.root.left = Sub.ToDualTreeNode("1")
	t.root.left.left = X1.ToDualTreeNode("ll")
	t.root.left.right = Const0.ToDualTreeNode("32")
	return &t
}

// TreeT_NT_T_NT_T_4 = 4 - 0 + 0
var TreeT_NT_T_NT_T_4 = func() *DualTree {
	t := DualTree{}
	t.root = Add.ToDualTreeNode(RandString(5))
	t.root.right = Const0.ToDualTreeNode(RandString(5))
	t.root.left = Sub.ToDualTreeNode(RandString(5))
	t.root.left.left = Const4.ToDualTreeNode(RandString(5))
	t.root.left.right = Const0.ToDualTreeNode(RandString(5))
	return &t
}

// TreeT_NT_T_NT_T_NT_T_0 = 4 - 0 + 4 + 8
var TreeT_NT_T_NT_T_NT_T_0 = func() *DualTree {
	t := DualTree{}
	t.root = Add.ToDualTreeNode(RandString(5))
	t.root.right = Add.ToDualTreeNode(RandString(5))
	t.root.right.left = Const4.ToDualTreeNode(RandString(5))
	t.root.right.right = Const8.ToDualTreeNode(RandString(5))

	t.root.left = Sub.ToDualTreeNode(RandString(5))
	t.root.left.left = Const4.ToDualTreeNode(RandString(5))
	t.root.left.right = Const0.ToDualTreeNode(RandString(5))
	return &t
}

// TreeT_NT_T_NT_T_NT_T_1 = x * x + 4 + 8
var TreeT_NT_T_NT_T_NT_T_1 = func() *DualTree {
	t := DualTree{}
	t.root = Add.ToDualTreeNode(RandString(5))
	t.root.right = Add.ToDualTreeNode(RandString(5))
	t.root.right.left = Const4.ToDualTreeNode(RandString(5))
	t.root.right.right = Const8.ToDualTreeNode(RandString(5))

	t.root.left = Mult.ToDualTreeNode(RandString(5))
	t.root.left.left = X1.ToDualTreeNode(RandString(5))
	t.root.left.right = X1.ToDualTreeNode(RandString(5))
	return &t
}

// TreeXXXX = x * x * x * x
var TreeXXXX = func() *DualTree {
	t := DualTree{}
	t.root = Mult.ToDualTreeNode(RandString(5))
	t.root.right = Mult.ToDualTreeNode(RandString(5))
	t.root.right.left = X1.ToDualTreeNode(RandString(5))
	t.root.right.right = X1.ToDualTreeNode(RandString(5))

	t.root.left = Mult.ToDualTreeNode(RandString(5))
	t.root.left.left = X1.ToDualTreeNode(RandString(5))
	t.root.left.right = X1.ToDualTreeNode(RandString(5))
	return &t
}

// TreeT_NT_T_NT_T_NT_T_0 = x * x * x * x * x
var TreeT_NT_T_NT_T_NT_T_NT_T_0 = func() *DualTree {
	t := DualTree{}
	t.root = Mult.ToDualTreeNode(RandString(5))
	t.root.right = Mult.ToDualTreeNode(RandString(5))
	t.root.right.left = X1.ToDualTreeNode(RandString(5))
	t.root.right.right = X1.ToDualTreeNode(RandString(5))

	t.root.left = Mult.ToDualTreeNode(RandString(5))
	t.root.left.right = X1.ToDualTreeNode(RandString(5))

	t.root.left.left = Mult.ToDualTreeNode(RandString(5))
	t.root.left.left.left = X1.ToDualTreeNode(RandString(5))
	t.root.left.left.right = X1.ToDualTreeNode(RandString(5))
	return &t
}

// TreeT_NT_T_NT_T_NT_T_NT_T_1 = x * x * x * x + 4
var TreeT_NT_T_NT_T_NT_T_NT_T_1 = func() *DualTree {
	t := DualTree{}
	t.root = Mult.ToDualTreeNode(RandString(5))
	t.root.right = Add.ToDualTreeNode(RandString(5))
	t.root.right.left = X1.ToDualTreeNode(RandString(5))
	t.root.right.right = Const4.ToDualTreeNode(RandString(5))

	t.root.left = Mult.ToDualTreeNode(RandString(5))
	t.root.left.right = X1.ToDualTreeNode(RandString(5))

	t.root.left.left = Mult.ToDualTreeNode(RandString(5))
	t.root.left.left.left = X1.ToDualTreeNode(RandString(5))
	t.root.left.left.right = X1.ToDualTreeNode(RandString(5))
	return &t
}

// TreeT_NT_T_NT_T_NT_T_NT_T_NT_T_0 = 0 + 1 + 2 + 3 + 4 + 5
var TreeT_NT_T_NT_T_NT_T_NT_T_NT_T_0 = func() *DualTree {
	t := DualTree{}
	t.root = Add.ToDualTreeNode(RandString(5))

	t.root.left = Add.ToDualTreeNode(RandString(5))
	t.root.left.left = Add.ToDualTreeNode(RandString(5))
	t.root.left.left.left = Const0.ToDualTreeNode(RandString(5))
	t.root.left.left.right = Const1.ToDualTreeNode(RandString(5))

	t.root.left.right = Add.ToDualTreeNode(RandString(5))
	t.root.left.right.left = Const2.ToDualTreeNode(RandString(5))
	t.root.left.right.right = Const3.ToDualTreeNode(RandString(5))

	t.root.right = Add.ToDualTreeNode(RandString(5))
	t.root.right.left = Const4.ToDualTreeNode(RandString(5))
	t.root.right.right = Const5.ToDualTreeNode(RandString(5))

	return &t
}

// TreeT_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_0 = 0 + 1 + 2 + 3 + 4 + 5 + 6
var TreeT_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_0 = func() *DualTree {
	t := DualTree{}
	t.root = Add.ToDualTreeNode(RandString(5))

	t.root.left = Add.ToDualTreeNode(RandString(5))
	t.root.left.left = Add.ToDualTreeNode(RandString(5))
	t.root.left.left.left = Const0.ToDualTreeNode(RandString(5))
	t.root.left.left.right = Const1.ToDualTreeNode(RandString(5))

	t.root.left.right = Add.ToDualTreeNode(RandString(5))
	t.root.left.right.left = Const2.ToDualTreeNode(RandString(5))
	t.root.left.right.right = Const3.ToDualTreeNode(RandString(5))

	t.root.right = Add.ToDualTreeNode(RandString(5))
	t.root.right.right = Const6.ToDualTreeNode(RandString(5))

	t.root.right.left = Add.ToDualTreeNode(RandString(5))
	t.root.right.left.left = Const4.ToDualTreeNode(RandString(5))
	t.root.right.left.right = Const5.ToDualTreeNode(RandString(5))
	return &t
}

// TreeT_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_0 = 0 + 1 + 2 + 3 + 4 + 5 + 6 + 7
var TreeT_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_0 = func() *DualTree {
	t := DualTree{}
	t.root = Add.ToDualTreeNode(RandString(5))

	t.root.left = Add.ToDualTreeNode(RandString(5))
	t.root.left.left = Add.ToDualTreeNode(RandString(5))
	t.root.left.left.left = Const0.ToDualTreeNode(RandString(5))
	t.root.left.left.right = Const1.ToDualTreeNode(RandString(5))

	t.root.left.right = Add.ToDualTreeNode(RandString(5))
	t.root.left.right.left = Const2.ToDualTreeNode(RandString(5))
	t.root.left.right.right = Const3.ToDualTreeNode(RandString(5))

	t.root.right = Add.ToDualTreeNode(RandString(5))
	t.root.right.right = Add.ToDualTreeNode(RandString(5))
	t.root.right.right.left = Const6.ToDualTreeNode(RandString(5))
	t.root.right.right.right = Const7.ToDualTreeNode(RandString(5))

	t.root.right.left = Add.ToDualTreeNode(RandString(5))
	t.root.right.left.left = Const4.ToDualTreeNode(RandString(5))
	t.root.right.left.right = Const5.ToDualTreeNode(RandString(5))
	return &t
}

// TreeVine_D3 = sin(sin(sin(x)))
var TreeVine_D3 = func() *DualTree {
	t := DualTree{}
	t.root = Sin.ToDualTreeNode(RandString(5))
	t.root.left = Sin.ToDualTreeNode(RandString(5))
	t.root.left.left = Sin.ToDualTreeNode(RandString(5))
	t.root.left.left.left = X1.ToDualTreeNode(RandString(5))
	return &t
}

// TreeVine_D4 = sin(sin(sin(sin(x)))))
var TreeVine_D4 = func() *DualTree {
	t := DualTree{}
	t.root = Sin.ToDualTreeNode(RandString(5))
	t.root.left = Sin.ToDualTreeNode(RandString(5))
	t.root.left.left = Sin.ToDualTreeNode(RandString(5))
	t.root.left.left.left = Sin.ToDualTreeNode(RandString(5))
	t.root.left.left.left.left = X1.ToDualTreeNode(RandString(5))
	return &t
}

// TreeVine_D5_R = It looks like  triangle with depth 5 on the right side.
var TreeVine_D5_R = func() *DualTree {
	t := DualTree{}
	t.root = Add.ToDualTreeNode(RandString(5))
	t.root.left = Sin.ToDualTreeNode(RandString(5))
	t.root.left.left = Sin.ToDualTreeNode(RandString(5))
	t.root.left.left.left = Sin.ToDualTreeNode(RandString(5))
	t.root.left.left.left.left = X1.ToDualTreeNode(RandString(5))

	t.root.right = Sin.ToDualTreeNode(RandString(5))
	t.root.right.right = Sin.ToDualTreeNode(RandString(5))
	t.root.right.right.right = Sin.ToDualTreeNode(RandString(5))
	t.root.right.right.right.right = Sin.ToDualTreeNode(RandString(5))
	t.root.right.right.right.right.right = X1.ToDualTreeNode(RandString(5))
	return &t
}

// TreeVine_D5_R = It looks like  triangle with depth 6 on the right side.
var TreeVine_D6_R = func() *DualTree {
	t := DualTree{}
	t.root = Add.ToDualTreeNode("0")
	t.root.left = Sin.ToDualTreeNode("1")
	t.root.left.left = Sin.ToDualTreeNode("2")
	t.root.left.left.left = Sin.ToDualTreeNode("3")
	t.root.left.left.left.left = X1.ToDualTreeNode("4")

	t.root.right = Sin.ToDualTreeNode("5")
	t.root.right.right = Sin.ToDualTreeNode("6")
	t.root.right.right.right = Sin.ToDualTreeNode("7")
	t.root.right.right.right.right = Sin.ToDualTreeNode("8")
	t.root.right.right.right.right.right = Sin.ToDualTreeNode("9")
	t.root.right.right.right.right.right.left = X1.ToDualTreeNode("10")
	return &t
}

// Tree3 = Sin(4 - x)
var Tree3 = func() *DualTree {
	t := DualTree{}
	t.root = Sin.ToDualTreeNode(RandString(5))
	t.root.left = Sub.ToDualTreeNode(RandString(5))
	t.root.left.left = X1.ToDualTreeNode(RandString(5))
	t.root.left.right = Const4.ToDualTreeNode(RandString(5))
	return &t
}

// Tree4 = Sin(x)
var Tree4 = func() *DualTree {
	t := DualTree{}
	t.root = Sin.ToDualTreeNode(RandString(5))
	t.root.left = X1.ToDualTreeNode(RandString(5))
	return &t
}

// Tree5 = Sin
var Tree5 = func() *DualTree {
	t := DualTree{}
	t.root = Sin.ToDualTreeNode(RandString(5))
	return &t
}

// Tree6 = x +
var Tree6 = func() *DualTree {
	t := DualTree{}
	t.root = Add.ToDualTreeNode(RandString(5))
	t.root.left = X1.ToDualTreeNode(RandString(5))
	return &t
}

// Tree7 =  +
var Tree7 = func() *DualTree {
	t := DualTree{}
	t.root = Add.ToDualTreeNode(RandString(5))
	return &t
}

// Tree8 =  x * x
var Tree8 = func() *DualTree {
	t := DualTree{}
	t.root = Mult.ToDualTreeNode(RandString(5))
	t.root.left = X1.ToDualTreeNode(RandString(5))
	t.root.right = X1.ToDualTreeNode(RandString(5))
	return &t
}
