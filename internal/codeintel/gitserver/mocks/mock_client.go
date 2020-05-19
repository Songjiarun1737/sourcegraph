// Code generated by github.com/efritz/go-mockgen 0.1.0; DO NOT EDIT.

package mocks

import (
	db "github.com/sourcegraph/sourcegraph/internal/codeintel/db"
	gitserver "github.com/sourcegraph/sourcegraph/internal/codeintel/gitserver"
	"sync"
)

// MockClient is a mock impelementation of the Client interface (from the
// package github.com/sourcegraph/sourcegraph/internal/codeintel/gitserver)
// used for unit testing.
type MockClient struct {
	// CommitsNearFunc is an instance of a mock function object controlling
	// the behavior of the method CommitsNear.
	CommitsNearFunc *ClientCommitsNearFunc
	// DirectoryChildrenFunc is an instance of a mock function object
	// controlling the behavior of the method DirectoryChildren.
	DirectoryChildrenFunc *ClientDirectoryChildrenFunc
	// HeadFunc is an instance of a mock function object controlling the
	// behavior of the method Head.
	HeadFunc *ClientHeadFunc
	// RepoStatusFunc is an instance of a mock function object controlling
	// the behavior of the method RepoStatus.
	RepoStatusFunc *ClientRepoStatusFunc
}

// NewMockClient creates a new mock of the Client interface. All methods
// return zero values for all results, unless overwritten.
func NewMockClient() *MockClient {
	return &MockClient{
		CommitsNearFunc: &ClientCommitsNearFunc{
			defaultHook: func(db.DB, int, string) (map[string][]string, error) {
				return nil, nil
			},
		},
		DirectoryChildrenFunc: &ClientDirectoryChildrenFunc{
			defaultHook: func(db.DB, int, string, []string) (map[string][]string, error) {
				return nil, nil
			},
		},
		HeadFunc: &ClientHeadFunc{
			defaultHook: func(db.DB, int) (string, error) {
				return "", nil
			},
		},
		RepoStatusFunc: &ClientRepoStatusFunc{
			defaultHook: func(int, string) (bool, bool, error) {
				return false, false, nil
			},
		},
	}
}

// NewMockClientFrom creates a new mock of the MockClient interface. All
// methods delegate to the given implementation, unless overwritten.
func NewMockClientFrom(i gitserver.Client) *MockClient {
	return &MockClient{
		CommitsNearFunc: &ClientCommitsNearFunc{
			defaultHook: i.CommitsNear,
		},
		DirectoryChildrenFunc: &ClientDirectoryChildrenFunc{
			defaultHook: i.DirectoryChildren,
		},
		HeadFunc: &ClientHeadFunc{
			defaultHook: i.Head,
		},
		RepoStatusFunc: &ClientRepoStatusFunc{
			defaultHook: i.RepoStatus,
		},
	}
}

// ClientCommitsNearFunc describes the behavior when the CommitsNear method
// of the parent MockClient instance is invoked.
type ClientCommitsNearFunc struct {
	defaultHook func(db.DB, int, string) (map[string][]string, error)
	hooks       []func(db.DB, int, string) (map[string][]string, error)
	history     []ClientCommitsNearFuncCall
	mutex       sync.Mutex
}

// CommitsNear delegates to the next hook function in the queue and stores
// the parameter and result values of this invocation.
func (m *MockClient) CommitsNear(v0 db.DB, v1 int, v2 string) (map[string][]string, error) {
	r0, r1 := m.CommitsNearFunc.nextHook()(v0, v1, v2)
	m.CommitsNearFunc.appendCall(ClientCommitsNearFuncCall{v0, v1, v2, r0, r1})
	return r0, r1
}

// SetDefaultHook sets function that is called when the CommitsNear method
// of the parent MockClient instance is invoked and the hook queue is empty.
func (f *ClientCommitsNearFunc) SetDefaultHook(hook func(db.DB, int, string) (map[string][]string, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// CommitsNear method of the parent MockClient instance inovkes the hook at
// the front of the queue and discards it. After the queue is empty, the
// default hook function is invoked for any future action.
func (f *ClientCommitsNearFunc) PushHook(hook func(db.DB, int, string) (map[string][]string, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultDefaultHook with a function that returns
// the given values.
func (f *ClientCommitsNearFunc) SetDefaultReturn(r0 map[string][]string, r1 error) {
	f.SetDefaultHook(func(db.DB, int, string) (map[string][]string, error) {
		return r0, r1
	})
}

// PushReturn calls PushDefaultHook with a function that returns the given
// values.
func (f *ClientCommitsNearFunc) PushReturn(r0 map[string][]string, r1 error) {
	f.PushHook(func(db.DB, int, string) (map[string][]string, error) {
		return r0, r1
	})
}

func (f *ClientCommitsNearFunc) nextHook() func(db.DB, int, string) (map[string][]string, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *ClientCommitsNearFunc) appendCall(r0 ClientCommitsNearFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of ClientCommitsNearFuncCall objects
// describing the invocations of this function.
func (f *ClientCommitsNearFunc) History() []ClientCommitsNearFuncCall {
	f.mutex.Lock()
	history := make([]ClientCommitsNearFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// ClientCommitsNearFuncCall is an object that describes an invocation of
// method CommitsNear on an instance of MockClient.
type ClientCommitsNearFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 db.DB
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 int
	// Arg2 is the value of the 3rd argument passed to this method
	// invocation.
	Arg2 string
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 map[string][]string
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c ClientCommitsNearFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1, c.Arg2}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c ClientCommitsNearFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1}
}

// ClientDirectoryChildrenFunc describes the behavior when the
// DirectoryChildren method of the parent MockClient instance is invoked.
type ClientDirectoryChildrenFunc struct {
	defaultHook func(db.DB, int, string, []string) (map[string][]string, error)
	hooks       []func(db.DB, int, string, []string) (map[string][]string, error)
	history     []ClientDirectoryChildrenFuncCall
	mutex       sync.Mutex
}

// DirectoryChildren delegates to the next hook function in the queue and
// stores the parameter and result values of this invocation.
func (m *MockClient) DirectoryChildren(v0 db.DB, v1 int, v2 string, v3 []string) (map[string][]string, error) {
	r0, r1 := m.DirectoryChildrenFunc.nextHook()(v0, v1, v2, v3)
	m.DirectoryChildrenFunc.appendCall(ClientDirectoryChildrenFuncCall{v0, v1, v2, v3, r0, r1})
	return r0, r1
}

// SetDefaultHook sets function that is called when the DirectoryChildren
// method of the parent MockClient instance is invoked and the hook queue is
// empty.
func (f *ClientDirectoryChildrenFunc) SetDefaultHook(hook func(db.DB, int, string, []string) (map[string][]string, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// DirectoryChildren method of the parent MockClient instance inovkes the
// hook at the front of the queue and discards it. After the queue is empty,
// the default hook function is invoked for any future action.
func (f *ClientDirectoryChildrenFunc) PushHook(hook func(db.DB, int, string, []string) (map[string][]string, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultDefaultHook with a function that returns
// the given values.
func (f *ClientDirectoryChildrenFunc) SetDefaultReturn(r0 map[string][]string, r1 error) {
	f.SetDefaultHook(func(db.DB, int, string, []string) (map[string][]string, error) {
		return r0, r1
	})
}

// PushReturn calls PushDefaultHook with a function that returns the given
// values.
func (f *ClientDirectoryChildrenFunc) PushReturn(r0 map[string][]string, r1 error) {
	f.PushHook(func(db.DB, int, string, []string) (map[string][]string, error) {
		return r0, r1
	})
}

func (f *ClientDirectoryChildrenFunc) nextHook() func(db.DB, int, string, []string) (map[string][]string, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *ClientDirectoryChildrenFunc) appendCall(r0 ClientDirectoryChildrenFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of ClientDirectoryChildrenFuncCall objects
// describing the invocations of this function.
func (f *ClientDirectoryChildrenFunc) History() []ClientDirectoryChildrenFuncCall {
	f.mutex.Lock()
	history := make([]ClientDirectoryChildrenFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// ClientDirectoryChildrenFuncCall is an object that describes an invocation
// of method DirectoryChildren on an instance of MockClient.
type ClientDirectoryChildrenFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 db.DB
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 int
	// Arg2 is the value of the 3rd argument passed to this method
	// invocation.
	Arg2 string
	// Arg3 is the value of the 4th argument passed to this method
	// invocation.
	Arg3 []string
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 map[string][]string
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c ClientDirectoryChildrenFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1, c.Arg2, c.Arg3}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c ClientDirectoryChildrenFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1}
}

// ClientHeadFunc describes the behavior when the Head method of the parent
// MockClient instance is invoked.
type ClientHeadFunc struct {
	defaultHook func(db.DB, int) (string, error)
	hooks       []func(db.DB, int) (string, error)
	history     []ClientHeadFuncCall
	mutex       sync.Mutex
}

// Head delegates to the next hook function in the queue and stores the
// parameter and result values of this invocation.
func (m *MockClient) Head(v0 db.DB, v1 int) (string, error) {
	r0, r1 := m.HeadFunc.nextHook()(v0, v1)
	m.HeadFunc.appendCall(ClientHeadFuncCall{v0, v1, r0, r1})
	return r0, r1
}

// SetDefaultHook sets function that is called when the Head method of the
// parent MockClient instance is invoked and the hook queue is empty.
func (f *ClientHeadFunc) SetDefaultHook(hook func(db.DB, int) (string, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// Head method of the parent MockClient instance inovkes the hook at the
// front of the queue and discards it. After the queue is empty, the default
// hook function is invoked for any future action.
func (f *ClientHeadFunc) PushHook(hook func(db.DB, int) (string, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultDefaultHook with a function that returns
// the given values.
func (f *ClientHeadFunc) SetDefaultReturn(r0 string, r1 error) {
	f.SetDefaultHook(func(db.DB, int) (string, error) {
		return r0, r1
	})
}

// PushReturn calls PushDefaultHook with a function that returns the given
// values.
func (f *ClientHeadFunc) PushReturn(r0 string, r1 error) {
	f.PushHook(func(db.DB, int) (string, error) {
		return r0, r1
	})
}

func (f *ClientHeadFunc) nextHook() func(db.DB, int) (string, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *ClientHeadFunc) appendCall(r0 ClientHeadFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of ClientHeadFuncCall objects describing the
// invocations of this function.
func (f *ClientHeadFunc) History() []ClientHeadFuncCall {
	f.mutex.Lock()
	history := make([]ClientHeadFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// ClientHeadFuncCall is an object that describes an invocation of method
// Head on an instance of MockClient.
type ClientHeadFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 db.DB
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 int
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 string
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c ClientHeadFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c ClientHeadFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1}
}

// ClientRepoStatusFunc describes the behavior when the RepoStatus method of
// the parent MockClient instance is invoked.
type ClientRepoStatusFunc struct {
	defaultHook func(int, string) (bool, bool, error)
	hooks       []func(int, string) (bool, bool, error)
	history     []ClientRepoStatusFuncCall
	mutex       sync.Mutex
}

// RepoStatus delegates to the next hook function in the queue and stores
// the parameter and result values of this invocation.
func (m *MockClient) RepoStatus(v0 int, v1 string) (bool, bool, error) {
	r0, r1, r2 := m.RepoStatusFunc.nextHook()(v0, v1)
	m.RepoStatusFunc.appendCall(ClientRepoStatusFuncCall{v0, v1, r0, r1, r2})
	return r0, r1, r2
}

// SetDefaultHook sets function that is called when the RepoStatus method of
// the parent MockClient instance is invoked and the hook queue is empty.
func (f *ClientRepoStatusFunc) SetDefaultHook(hook func(int, string) (bool, bool, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// RepoStatus method of the parent MockClient instance inovkes the hook at
// the front of the queue and discards it. After the queue is empty, the
// default hook function is invoked for any future action.
func (f *ClientRepoStatusFunc) PushHook(hook func(int, string) (bool, bool, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultDefaultHook with a function that returns
// the given values.
func (f *ClientRepoStatusFunc) SetDefaultReturn(r0 bool, r1 bool, r2 error) {
	f.SetDefaultHook(func(int, string) (bool, bool, error) {
		return r0, r1, r2
	})
}

// PushReturn calls PushDefaultHook with a function that returns the given
// values.
func (f *ClientRepoStatusFunc) PushReturn(r0 bool, r1 bool, r2 error) {
	f.PushHook(func(int, string) (bool, bool, error) {
		return r0, r1, r2
	})
}

func (f *ClientRepoStatusFunc) nextHook() func(int, string) (bool, bool, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *ClientRepoStatusFunc) appendCall(r0 ClientRepoStatusFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of ClientRepoStatusFuncCall objects describing
// the invocations of this function.
func (f *ClientRepoStatusFunc) History() []ClientRepoStatusFuncCall {
	f.mutex.Lock()
	history := make([]ClientRepoStatusFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// ClientRepoStatusFuncCall is an object that describes an invocation of
// method RepoStatus on an instance of MockClient.
type ClientRepoStatusFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 int
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 string
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 bool
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 bool
	// Result2 is the value of the 3rd result returned from this method
	// invocation.
	Result2 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c ClientRepoStatusFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c ClientRepoStatusFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1, c.Result2}
}
