package evolution

// ProgNil-TreeNil
var ProgNil = Program{
	T:  TreeNil(),
	ID: "ProgNil",
}

var ProgBadTree = Program{
	T:  TreeBad(),
	ID: "BadTree",
}

// ProgX-TreeNil
var ProgX = Program{
	T:  TreeT_X(),
	ID: "ProgX",
}

// Prog1-TreeT_1
var Prog1 = Program{
	T:  TreeT_1(),
	ID: "Prog1",
}

// ProgTreeT_NT_T_0 | x*4
var ProgTreeT_NT_T_0 = Program{
	T:  TreeT_NT_T_0(),
	ID: "ProgTreeT_NT_T_0",
}

// ProgTreeT_NT_T_0 | x*5
var ProgTreeXby5 = Program{
	T:  TreeXby5(),
	ID: "TreeXby5",
}

// ProgTreeT_NT_T_0 + x+8
var ProgTreeT_NT_T_1 = Program{
	T:  TreeT_NT_T_1(),
	ID: "ProgTreeT_NT_T_0",
}

// ProgTreeT_NT_T_4 = x * x
var ProgTreeT_NT_T_4 = Program{
	T:  TreeT_NT_T_4(),
	ID: "ProgTreeT_NT_T_4",
}

// ProgTreeT_NT_T_NT_T_0 | x - x * 4
var ProgTreeT_NT_T_NT_T_0 = Program{
	T:  TreeT_NT_T_NT_T_0(),
	ID: "ProgTreeT_NT_T_NT_T_0",
}

// ProgTreeT_NT_T_NT_T_NT_T_0 | 4 - 0 + 4 + 8
var ProgTreeT_NT_T_NT_T_NT_T_0 = Program{
	T:  TreeT_NT_T_NT_T_NT_T_0(),
	ID: "ProgTreeT_NT_T_NT_T_NT_T_0",
}

// ProgTreeXXXX4 = x * x * x * x + 4
var ProgTreeXXXX4 = Program{
	T:  TreeT_NT_T_NT_T_NT_T_NT_T_1(),
	ID: "ProgTreeXXXX4",
}

// ProgXXXX = x * x * x * x
var ProgXXXX = Program{
	T:  TreeXXXX(),
	ID: "ProgTreeT_NT_T_NT_T_NT_T_2",
}

// ProgTreeT_NT_T_NT_T_NT_T_NT_T_NT_T_0
var ProgTreeT_NT_T_NT_T_NT_T_NT_T_NT_T_0 = Program{
	T:  TreeT_NT_T_NT_T_NT_T_NT_T_NT_T_0(),
	ID: "ProgTreeT_NT_T_NT_T_NT_T_NT_T_NT_T_0",
}

// ProgTreeT_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_0
var ProgTreeT_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_0 = Program{
	T:  TreeT_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_0(),
	ID: "ProgTreeT_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_0",
}

// ProgTreeT_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_0
var ProgTreeT_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_0 = Program{
	T:  TreeT_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_0(),
	ID: "ProgTreeT_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_NT_T_0",
}
