package commands

// UndoIdentity undo identity
var UndoIdentity int64 = 0x800001

// UndoCommandDef undo operation
var UndoCommandDef CommandWordDef = CommandWordDef{
	Identity: UndoIdentity,
	Name:     "undo",
	Helper:   "undo operation",
	Type:     "name",
}

// ShowCommandDef show information
var ShowCommandDef CommandWordDef = CommandWordDef{
	Name:   "show",
	Helper: "show operation",
	Type:   "name",
}

// ResetCommandDef reset command define
var ResetCommandDef CommandWordDef = CommandWordDef{
	Name:   "reset",
	Helper: "reset operation",
	Type:   "name",
}

// DebugCommandDef debug command define
var DebugCommandDef CommandWordDef = CommandWordDef{
	Name:   "debug",
	Helper: "debug operation",
	Type:   "name",
}
