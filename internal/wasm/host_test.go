package wasm

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tetratelabs/wazero/api"
)

// wasiAPI simulates the real WASI api
type wasiAPI struct {
}

func ArgsSizesGet(ctx api.Module, resultArgc, resultArgvBufSize uint32) uint32 {
	return 0
}

func (a *wasiAPI) ArgsSizesGet(ctx api.Module, resultArgc, resultArgvBufSize uint32) uint32 {
	return 0
}

func (a *wasiAPI) FdWrite(ctx api.Module, fd, iovs, iovsCount, resultSize uint32) uint32 {
	return 0
}

func TestNewHostModule(t *testing.T) {
	i32 := ValueTypeI32

	a := wasiAPI{}
	functionArgsSizesGet := "args_sizes_get"
	fnArgsSizesGet := reflect.ValueOf(a.ArgsSizesGet)
	functionFdWrite := "fd_write"
	fnFdWrite := reflect.ValueOf(a.FdWrite)

	tests := []struct {
		name, moduleName string
		nameToGoFunc     map[string]interface{}
		nameToMemory     map[string]*Memory
		expected         *Module
	}{
		{
			name:     "empty",
			expected: &Module{},
		},
		{
			name:       "only name",
			moduleName: "test",
			expected:   &Module{NameSection: &NameSection{ModuleName: "test"}},
		},
		{
			name:         "memory",
			nameToMemory: map[string]*Memory{"memory": {1, 2}},
			expected: &Module{
				MemorySection: &Memory{Min: 1, Max: 2},
				ExportSection: map[string]*Export{
					"memory": {Name: "memory", Type: ExternTypeMemory, Index: 0},
				},
			},
		},
		{
			name:       "two struct funcs",
			moduleName: "wasi_snapshot_preview1",
			nameToGoFunc: map[string]interface{}{
				functionArgsSizesGet: a.ArgsSizesGet,
				functionFdWrite:      a.FdWrite,
			},
			expected: &Module{
				TypeSection: []*FunctionType{
					{Params: []ValueType{i32, i32}, Results: []ValueType{i32}},
					{Params: []ValueType{i32, i32, i32, i32}, Results: []ValueType{i32}},
				},
				FunctionSection:     []Index{0, 1},
				HostFunctionSection: []*reflect.Value{&fnArgsSizesGet, &fnFdWrite},
				ExportSection: map[string]*Export{
					"args_sizes_get": {Name: "args_sizes_get", Type: ExternTypeFunc, Index: 0},
					"fd_write":       {Name: "fd_write", Type: ExternTypeFunc, Index: 1},
				},
				NameSection: &NameSection{
					ModuleName: "wasi_snapshot_preview1",
					FunctionNames: NameMap{
						{Index: 0, Name: "args_sizes_get"},
						{Index: 1, Name: "fd_write"},
					},
				},
			},
		},
		{
			name:       "one of each",
			moduleName: "env",
			nameToGoFunc: map[string]interface{}{
				functionArgsSizesGet: a.ArgsSizesGet,
			},
			nameToMemory: map[string]*Memory{
				"memory": {1, 1},
			},
			expected: &Module{
				TypeSection: []*FunctionType{
					{Params: []ValueType{i32, i32}, Results: []ValueType{i32}},
				},
				FunctionSection:     []Index{0},
				HostFunctionSection: []*reflect.Value{&fnArgsSizesGet},
				ExportSection: map[string]*Export{
					"args_sizes_get": {Name: "args_sizes_get", Type: ExternTypeFunc, Index: 0},
					"memory":         {Name: "memory", Type: ExternTypeMemory, Index: 0},
				},
				MemorySection: &Memory{Min: 1, Max: 1},
				NameSection: &NameSection{
					ModuleName: "env",
					FunctionNames: NameMap{
						{Index: 0, Name: "args_sizes_get"},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		tc := tt

		t.Run(tc.name, func(t *testing.T) {
			m, e := NewHostModule(tc.moduleName, tc.nameToGoFunc, tc.nameToMemory)
			require.NoError(t, e)
			requireHostModuleEquals(t, tc.expected, m)
		})
	}
}

func requireHostModuleEquals(t *testing.T, expected, actual *Module) {
	// `require.Equal(t, expected, actual)` fails reflect pointers don't match, so brute compare:
	require.Equal(t, expected.TypeSection, actual.TypeSection)
	require.Equal(t, expected.ImportSection, actual.ImportSection)
	require.Equal(t, expected.FunctionSection, actual.FunctionSection)
	require.Equal(t, expected.TableSection, actual.TableSection)
	require.Equal(t, expected.MemorySection, actual.MemorySection)
	require.Equal(t, expected.GlobalSection, actual.GlobalSection)
	require.Equal(t, expected.ExportSection, actual.ExportSection)
	require.Equal(t, expected.StartSection, actual.StartSection)
	require.Equal(t, expected.ElementSection, actual.ElementSection)
	require.Nil(t, actual.CodeSection) // Host functions are implemented in Go, not Wasm!
	require.Equal(t, expected.DataSection, actual.DataSection)
	require.Equal(t, expected.NameSection, actual.NameSection)

	// Special case because reflect.Value can't be compared with Equals
	require.Equal(t, len(expected.HostFunctionSection), len(actual.HostFunctionSection))
	for i := range expected.HostFunctionSection {
		require.Equal(t, (*expected.HostFunctionSection[i]).Type(), (*actual.HostFunctionSection[i]).Type())
	}
}

func TestNewHostModule_Errors(t *testing.T) {
	tests := []struct {
		name, moduleName string
		nameToGoFunc     map[string]interface{}
		nameToMemory     map[string]*Memory
		expectedErr      string
	}{
		{
			name:         "not a function",
			nameToGoFunc: map[string]interface{}{"fn": t},
			expectedErr:  "func[fn] kind != func: ptr",
		},
		{
			name:         "memory collides on func name",
			nameToGoFunc: map[string]interface{}{"fn": ArgsSizesGet},
			nameToMemory: map[string]*Memory{"fn": {1, 1}},
			expectedErr:  "memory[fn] exports the same name as a func",
		},
		{
			name:         "multiple memories",
			nameToMemory: map[string]*Memory{"memory": {1, 1}, "mem": {2, 2}},
			expectedErr:  "only one memory is allowed, but configured: mem, memory",
		},
		{
			name:         "memory max < min",
			nameToMemory: map[string]*Memory{"memory": {1, 0}},
			expectedErr:  "memory[memory] min 1 pages (64 Ki) > max 0 pages (0 Ki)",
		},
	}

	for _, tt := range tests {
		tc := tt

		t.Run(tc.name, func(t *testing.T) {
			_, e := NewHostModule(tc.moduleName, tc.nameToGoFunc, tc.nameToMemory)
			require.EqualError(t, e, tc.expectedErr)
		})
	}
}
