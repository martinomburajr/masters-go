package evolution

import (
	"fmt"
	"testing"
)

var TestSectionDivider = func(title string, t *testing.T) {
	t.Log(fmt.Sprintf("############################## %s ##########################", title))
}

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

// TreeNil = 0
var TreeNil = func() *DualTree {
	t := DualTree{}
	return &t
}

// TreeT_0 = x
var TreeT_0 = func() *DualTree {
	t := DualTree{}
	t.root = X1.ToDualTreeNode(0)
	return &t
}

// TreeT_1 = 0
var TreeT_1 = func() *DualTree {
	t := DualTree{}
	t.root = Const0.ToDualTreeNode(0)
	return &t
}

// TreeT_NT_T_4 = x * x
var TreeT_NT_T_4 = func() *DualTree {
	t := DualTree{}
	t.root = Mult.ToDualTreeNode(0)
	t.root.left = X1.ToDualTreeNode(1)
	t.root.right = X1.ToDualTreeNode(2)
	return &t
}

// TreeT_NT_T_0 = x * 4
var TreeT_NT_T_0 = func() *DualTree {
	t := DualTree{}
	t.root = Mult.ToDualTreeNode(0)
	t.root.left = X1.ToDualTreeNode(1)
	t.root.right = Const4.ToDualTreeNode(2)
	return &t
}

// TreeT_NT_T_1 = x + 8
var TreeT_NT_T_1 = func() *DualTree {
	t := DualTree{}
	t.root = Add.ToDualTreeNode(0)
	t.root.left = X1.ToDualTreeNode(1)
	t.root.right = Const8.ToDualTreeNode(2)
	return &t
}

// TreeT_NT_T_2 = 4 - 8
var TreeT_NT_T_2 = func() *DualTree {
	t := DualTree{}
	t.root = Sub.ToDualTreeNode(0)
	t.root.left = Const4.ToDualTreeNode(1)
	t.root.right = Const8.ToDualTreeNode(2)
	return &t
}

// TreeT_NT_T_3 = 8 * 8
var TreeT_NT_T_3 = func() *DualTree {
	t := DualTree{}
	t.root = Mult.ToDualTreeNode(0)
	t.root.left = Const8.ToDualTreeNode(1)
	t.root.right = Const8.ToDualTreeNode(2)
	return &t
}

// TreeT_NT_T_NT_T_0 = x - x * 4
var TreeT_NT_T_NT_T_0 = func() *DualTree {
	t := DualTree{}
	t.root = Mult.ToDualTreeNode(0)
	t.root.right = Const4.ToDualTreeNode(4)
	t.root.left = Sub.ToDualTreeNode(1)
	t.root.left.left = X1.ToDualTreeNode(2)
	t.root.left.right = X1.ToDualTreeNode(3)
	return &t
}

// TreeT_NT_T_NT_T_1 = x + 8 * 4
var TreeT_NT_T_NT_T_1 = func() *DualTree {
	t := DualTree{}
	t.root = Mult.ToDualTreeNode(0)
	t.root.right = Const4.ToDualTreeNode(4)
	t.root.left = Add.ToDualTreeNode(1)
	t.root.left.left = X1.ToDualTreeNode(2)
	t.root.left.right = Const8.ToDualTreeNode(3)
	return &t
}

// TreeT_NT_T_NT_T_2 = x - 0 * 4
var TreeT_NT_T_NT_T_2 = func() *DualTree {
	t := DualTree{}
	t.root = Mult.ToDualTreeNode(0)
	t.root.right = Const4.ToDualTreeNode(4)
	t.root.left = Sub.ToDualTreeNode(1)
	t.root.left.left = X1.ToDualTreeNode(2)
	t.root.left.right = Const0.ToDualTreeNode(3)
	return &t
}

// TreeT_NT_T_NT_T_3 = x - 0 * 0
var TreeT_NT_T_NT_T_3 = func() *DualTree {
	t := DualTree{}
	t.root = Mult.ToDualTreeNode(0)
	t.root.right = Const0.ToDualTreeNode(4)
	t.root.left = Sub.ToDualTreeNode(1)
	t.root.left.left = X1.ToDualTreeNode(2)
	t.root.left.right = Const0.ToDualTreeNode(3)
	return &t
}

// TreeT_NT_T_NT_T_4 = 4 - 0 + 0
var TreeT_NT_T_NT_T_4 = func() *DualTree {
	t := DualTree{}
	t.root = Add.ToDualTreeNode(0)
	t.root.right = Const0.ToDualTreeNode(4)
	t.root.left = Sub.ToDualTreeNode(1)
	t.root.left.left = Const4.ToDualTreeNode(2)
	t.root.left.right = Const0.ToDualTreeNode(3)
	return &t
}

// TreeT_NT_T_NT_T_NT_T_0 = 4 - 0 + 4 + 8
var TreeT_NT_T_NT_T_NT_T_0 = func() *DualTree {
	t := DualTree{}
	t.root = Add.ToDualTreeNode(0)
	t.root.right = Add.ToDualTreeNode(5)
	t.root.right.left = Const4.ToDualTreeNode(4)
	t.root.right.right = Const8.ToDualTreeNode(7)

	t.root.left = Sub.ToDualTreeNode(1)
	t.root.left.left = Const4.ToDualTreeNode(2)
	t.root.left.right = Const0.ToDualTreeNode(3)
	return &t
}

// TreeT_NT_T_NT_T_NT_T_1 = x * x + 4 + 8
var TreeT_NT_T_NT_T_NT_T_1 = func() *DualTree {
	t := DualTree{}
	t.root = Add.ToDualTreeNode(0)
	t.root.right = Add.ToDualTreeNode(5)
	t.root.right.left = Const4.ToDualTreeNode(4)
	t.root.right.right = Const8.ToDualTreeNode(7)

	t.root.left = Mult.ToDualTreeNode(1)
	t.root.left.left = X1.ToDualTreeNode(2)
	t.root.left.right = X1.ToDualTreeNode(3)
	return &t
}

// TreeT_NT_T_NT_T_NT_T_2 = x * x * x * x
var TreeT_NT_T_NT_T_NT_T_2 = func() *DualTree {
	t := DualTree{}
	t.root = Mult.ToDualTreeNode(0)
	t.root.right = Mult.ToDualTreeNode(5)
	t.root.right.left = X1.ToDualTreeNode(4)
	t.root.right.right = X1.ToDualTreeNode(7)

	t.root.left = Mult.ToDualTreeNode(1)
	t.root.left.left = X1.ToDualTreeNode(2)
	t.root.left.right = X1.ToDualTreeNode(3)
	return &t
}

// TreeT_NT_T_NT_T_NT_T_0 = x * x * x * x * x
var TreeT_NT_T_NT_T_NT_T_NT_T_0 = func() *DualTree {
	t := DualTree{}
	t.root = Mult.ToDualTreeNode(0)
	t.root.right = Mult.ToDualTreeNode(5)
	t.root.right.left = X1.ToDualTreeNode(4)
	t.root.right.right = X1.ToDualTreeNode(7)

	t.root.left = Mult.ToDualTreeNode(1)
	t.root.left.right = X1.ToDualTreeNode(3)

	t.root.left.left = Mult.ToDualTreeNode(2)
	t.root.left.left.left = X1.ToDualTreeNode(2)
	t.root.left.left.right = X1.ToDualTreeNode(2)
	return &t
}

// TreeT_NT_T_NT_T_NT_T_NT_T_1 = x * x * x * x + 4
var TreeT_NT_T_NT_T_NT_T_NT_T_1 = func() *DualTree {
	t := DualTree{}
	t.root = Mult.ToDualTreeNode(0)
	t.root.right = Add.ToDualTreeNode(5)
	t.root.right.left = X1.ToDualTreeNode(4)
	t.root.right.right = Const4.ToDualTreeNode(7)

	t.root.left = Mult.ToDualTreeNode(1)
	t.root.left.right = X1.ToDualTreeNode(3)

	t.root.left.left = Mult.ToDualTreeNode(2)
	t.root.left.left.left = X1.ToDualTreeNode(2)
	t.root.left.left.right = X1.ToDualTreeNode(2)
	return &t
}

// TreeT_NT_T_NT_T_NT_T_NT_T_NT_T_0 = 0 + 1 + 2 + 3 + 4 + 5
var TreeT_NT_T_NT_T_NT_T_NT_T_NT_T_0 = func() *DualTree {
	t := DualTree{}
	t.root = Add.ToDualTreeNode(0)

	t.root.left = Add.ToDualTreeNode(1)
	t.root.left.left = Add.ToDualTreeNode(1)
	t.root.left.left.left = Const0.ToDualTreeNode(1)
	t.root.left.left.right = Const1.ToDualTreeNode(1)

	t.root.left.right = Add.ToDualTreeNode(1)
	t.root.left.right.left = Const2.ToDualTreeNode(1)
	t.root.left.right.right = Const3.ToDualTreeNode(1)

	t.root.right = Add.ToDualTreeNode(1)
	t.root.right.left = Const4.ToDualTreeNode(1)
	t.root.right.right = Const5.ToDualTreeNode(1)

	return &t
}

// TreeT_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_0 = 0 + 1 + 2 + 3 + 4 + 5 + 6
var TreeT_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_0 = func() *DualTree {
	t := DualTree{}
	t.root = Add.ToDualTreeNode(0)

	t.root.left = Add.ToDualTreeNode(1)
	t.root.left.left = Add.ToDualTreeNode(1)
	t.root.left.left.left = Const0.ToDualTreeNode(1)
	t.root.left.left.right = Const1.ToDualTreeNode(1)

	t.root.left.right = Add.ToDualTreeNode(1)
	t.root.left.right.left = Const2.ToDualTreeNode(1)
	t.root.left.right.right = Const3.ToDualTreeNode(1)

	t.root.right = Add.ToDualTreeNode(1)
	t.root.right.right = Const6.ToDualTreeNode(1)

	t.root.right.left = Add.ToDualTreeNode(1)
	t.root.right.left.left = Const4.ToDualTreeNode(1)
	t.root.right.left.right = Const5.ToDualTreeNode(1)
	return &t
}

// TreeT_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_0 = 0 + 1 + 2 + 3 + 4 + 5 + 6 + 7
var TreeT_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_0 = func() *DualTree {
	t := DualTree{}
	t.root = Add.ToDualTreeNode(0)

	t.root.left = Add.ToDualTreeNode(1)
	t.root.left.left = Add.ToDualTreeNode(1)
	t.root.left.left.left = Const0.ToDualTreeNode(1)
	t.root.left.left.right = Const1.ToDualTreeNode(1)

	t.root.left.right = Add.ToDualTreeNode(1)
	t.root.left.right.left = Const2.ToDualTreeNode(1)
	t.root.left.right.right = Const3.ToDualTreeNode(1)

	t.root.right = Add.ToDualTreeNode(1)
	t.root.right.right = Add.ToDualTreeNode(1)
	t.root.right.right.left = Const6.ToDualTreeNode(1)
	t.root.right.right.right = Const7.ToDualTreeNode(1)

	t.root.right.left = Add.ToDualTreeNode(1)
	t.root.right.left.left = Const4.ToDualTreeNode(1)
	t.root.right.left.right = Const5.ToDualTreeNode(1)
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
