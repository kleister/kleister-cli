package mocks

import "github.com/solderapp/solder-go/solder"
import "github.com/stretchr/testify/mock"

// ClientAPI is an autogenerated mock type
type ClientAPI struct {
	mock.Mock
}

// BuildDelete provides a mock function with given fields: _a0, _a1
func (_m *ClientAPI) BuildDelete(_a0 string, _a1 string) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// BuildGet provides a mock function with given fields: _a0, _a1
func (_m *ClientAPI) BuildGet(_a0 string, _a1 string) (*solder.Build, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *solder.Build
	if rf, ok := ret.Get(0).(func(string, string) *solder.Build); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*solder.Build)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BuildList provides a mock function with given fields: _a0
func (_m *ClientAPI) BuildList(_a0 string) ([]*solder.Build, error) {
	ret := _m.Called(_a0)

	var r0 []*solder.Build
	if rf, ok := ret.Get(0).(func(string) []*solder.Build); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*solder.Build)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BuildPatch provides a mock function with given fields: _a0, _a1
func (_m *ClientAPI) BuildPatch(_a0 string, _a1 *solder.Build) (*solder.Build, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *solder.Build
	if rf, ok := ret.Get(0).(func(string, *solder.Build) *solder.Build); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*solder.Build)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, *solder.Build) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BuildPost provides a mock function with given fields: _a0, _a1
func (_m *ClientAPI) BuildPost(_a0 string, _a1 *solder.Build) (*solder.Build, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *solder.Build
	if rf, ok := ret.Get(0).(func(string, *solder.Build) *solder.Build); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*solder.Build)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, *solder.Build) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BuildVersionAppend provides a mock function with given fields: _a0
func (_m *ClientAPI) BuildVersionAppend(_a0 solder.BuildVersionOptions) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(solder.BuildVersionOptions) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// BuildVersionDelete provides a mock function with given fields: _a0
func (_m *ClientAPI) BuildVersionDelete(_a0 solder.BuildVersionOptions) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(solder.BuildVersionOptions) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// BuildVersionList provides a mock function with given fields: _a0
func (_m *ClientAPI) BuildVersionList(_a0 solder.BuildVersionOptions) ([]*solder.Version, error) {
	ret := _m.Called(_a0)

	var r0 []*solder.Version
	if rf, ok := ret.Get(0).(func(solder.BuildVersionOptions) []*solder.Version); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*solder.Version)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(solder.BuildVersionOptions) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ClientDelete provides a mock function with given fields: _a0
func (_m *ClientAPI) ClientDelete(_a0 string) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ClientGet provides a mock function with given fields: _a0
func (_m *ClientAPI) ClientGet(_a0 string) (*solder.Client, error) {
	ret := _m.Called(_a0)

	var r0 *solder.Client
	if rf, ok := ret.Get(0).(func(string) *solder.Client); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*solder.Client)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ClientList provides a mock function with given fields:
func (_m *ClientAPI) ClientList() ([]*solder.Client, error) {
	ret := _m.Called()

	var r0 []*solder.Client
	if rf, ok := ret.Get(0).(func() []*solder.Client); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*solder.Client)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ClientPackAppend provides a mock function with given fields: _a0
func (_m *ClientAPI) ClientPackAppend(_a0 solder.ClientPackOptions) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(solder.ClientPackOptions) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ClientPackDelete provides a mock function with given fields: _a0
func (_m *ClientAPI) ClientPackDelete(_a0 solder.ClientPackOptions) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(solder.ClientPackOptions) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ClientPackList provides a mock function with given fields: _a0
func (_m *ClientAPI) ClientPackList(_a0 solder.ClientPackOptions) ([]*solder.Pack, error) {
	ret := _m.Called(_a0)

	var r0 []*solder.Pack
	if rf, ok := ret.Get(0).(func(solder.ClientPackOptions) []*solder.Pack); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*solder.Pack)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(solder.ClientPackOptions) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ClientPatch provides a mock function with given fields: _a0
func (_m *ClientAPI) ClientPatch(_a0 *solder.Client) (*solder.Client, error) {
	ret := _m.Called(_a0)

	var r0 *solder.Client
	if rf, ok := ret.Get(0).(func(*solder.Client) *solder.Client); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*solder.Client)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*solder.Client) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ClientPost provides a mock function with given fields: _a0
func (_m *ClientAPI) ClientPost(_a0 *solder.Client) (*solder.Client, error) {
	ret := _m.Called(_a0)

	var r0 *solder.Client
	if rf, ok := ret.Get(0).(func(*solder.Client) *solder.Client); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*solder.Client)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*solder.Client) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ForgeBuildAppend provides a mock function with given fields: _a0
func (_m *ClientAPI) ForgeBuildAppend(_a0 solder.ForgeBuildOptions) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(solder.ForgeBuildOptions) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ForgeBuildDelete provides a mock function with given fields: _a0
func (_m *ClientAPI) ForgeBuildDelete(_a0 solder.ForgeBuildOptions) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(solder.ForgeBuildOptions) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ForgeBuildList provides a mock function with given fields: _a0
func (_m *ClientAPI) ForgeBuildList(_a0 solder.ForgeBuildOptions) ([]*solder.Build, error) {
	ret := _m.Called(_a0)

	var r0 []*solder.Build
	if rf, ok := ret.Get(0).(func(solder.ForgeBuildOptions) []*solder.Build); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*solder.Build)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(solder.ForgeBuildOptions) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ForgeGet provides a mock function with given fields: _a0
func (_m *ClientAPI) ForgeGet(_a0 string) (*solder.Forge, error) {
	ret := _m.Called(_a0)

	var r0 *solder.Forge
	if rf, ok := ret.Get(0).(func(string) *solder.Forge); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*solder.Forge)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ForgeList provides a mock function with given fields:
func (_m *ClientAPI) ForgeList() ([]*solder.Forge, error) {
	ret := _m.Called()

	var r0 []*solder.Forge
	if rf, ok := ret.Get(0).(func() []*solder.Forge); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*solder.Forge)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ForgeRefresh provides a mock function with given fields:
func (_m *ClientAPI) ForgeRefresh() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// KeyDelete provides a mock function with given fields: _a0
func (_m *ClientAPI) KeyDelete(_a0 string) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// KeyGet provides a mock function with given fields: _a0
func (_m *ClientAPI) KeyGet(_a0 string) (*solder.Key, error) {
	ret := _m.Called(_a0)

	var r0 *solder.Key
	if rf, ok := ret.Get(0).(func(string) *solder.Key); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*solder.Key)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// KeyList provides a mock function with given fields:
func (_m *ClientAPI) KeyList() ([]*solder.Key, error) {
	ret := _m.Called()

	var r0 []*solder.Key
	if rf, ok := ret.Get(0).(func() []*solder.Key); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*solder.Key)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// KeyPatch provides a mock function with given fields: _a0
func (_m *ClientAPI) KeyPatch(_a0 *solder.Key) (*solder.Key, error) {
	ret := _m.Called(_a0)

	var r0 *solder.Key
	if rf, ok := ret.Get(0).(func(*solder.Key) *solder.Key); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*solder.Key)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*solder.Key) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// KeyPost provides a mock function with given fields: _a0
func (_m *ClientAPI) KeyPost(_a0 *solder.Key) (*solder.Key, error) {
	ret := _m.Called(_a0)

	var r0 *solder.Key
	if rf, ok := ret.Get(0).(func(*solder.Key) *solder.Key); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*solder.Key)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*solder.Key) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MinecraftBuildAppend provides a mock function with given fields: _a0
func (_m *ClientAPI) MinecraftBuildAppend(_a0 solder.MinecraftBuildOptions) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(solder.MinecraftBuildOptions) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MinecraftBuildDelete provides a mock function with given fields: _a0
func (_m *ClientAPI) MinecraftBuildDelete(_a0 solder.MinecraftBuildOptions) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(solder.MinecraftBuildOptions) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MinecraftBuildList provides a mock function with given fields: _a0
func (_m *ClientAPI) MinecraftBuildList(_a0 solder.MinecraftBuildOptions) ([]*solder.Build, error) {
	ret := _m.Called(_a0)

	var r0 []*solder.Build
	if rf, ok := ret.Get(0).(func(solder.MinecraftBuildOptions) []*solder.Build); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*solder.Build)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(solder.MinecraftBuildOptions) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MinecraftGet provides a mock function with given fields: _a0
func (_m *ClientAPI) MinecraftGet(_a0 string) (*solder.Minecraft, error) {
	ret := _m.Called(_a0)

	var r0 *solder.Minecraft
	if rf, ok := ret.Get(0).(func(string) *solder.Minecraft); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*solder.Minecraft)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MinecraftList provides a mock function with given fields:
func (_m *ClientAPI) MinecraftList() ([]*solder.Minecraft, error) {
	ret := _m.Called()

	var r0 []*solder.Minecraft
	if rf, ok := ret.Get(0).(func() []*solder.Minecraft); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*solder.Minecraft)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MinecraftRefresh provides a mock function with given fields:
func (_m *ClientAPI) MinecraftRefresh() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ModDelete provides a mock function with given fields: _a0
func (_m *ClientAPI) ModDelete(_a0 string) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ModGet provides a mock function with given fields: _a0
func (_m *ClientAPI) ModGet(_a0 string) (*solder.Mod, error) {
	ret := _m.Called(_a0)

	var r0 *solder.Mod
	if rf, ok := ret.Get(0).(func(string) *solder.Mod); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*solder.Mod)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ModList provides a mock function with given fields:
func (_m *ClientAPI) ModList() ([]*solder.Mod, error) {
	ret := _m.Called()

	var r0 []*solder.Mod
	if rf, ok := ret.Get(0).(func() []*solder.Mod); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*solder.Mod)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ModPatch provides a mock function with given fields: _a0
func (_m *ClientAPI) ModPatch(_a0 *solder.Mod) (*solder.Mod, error) {
	ret := _m.Called(_a0)

	var r0 *solder.Mod
	if rf, ok := ret.Get(0).(func(*solder.Mod) *solder.Mod); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*solder.Mod)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*solder.Mod) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ModPost provides a mock function with given fields: _a0
func (_m *ClientAPI) ModPost(_a0 *solder.Mod) (*solder.Mod, error) {
	ret := _m.Called(_a0)

	var r0 *solder.Mod
	if rf, ok := ret.Get(0).(func(*solder.Mod) *solder.Mod); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*solder.Mod)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*solder.Mod) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ModUserAppend provides a mock function with given fields: _a0
func (_m *ClientAPI) ModUserAppend(_a0 solder.ModUserOptions) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(solder.ModUserOptions) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ModUserDelete provides a mock function with given fields: _a0
func (_m *ClientAPI) ModUserDelete(_a0 solder.ModUserOptions) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(solder.ModUserOptions) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ModUserList provides a mock function with given fields: _a0
func (_m *ClientAPI) ModUserList(_a0 solder.ModUserOptions) ([]*solder.User, error) {
	ret := _m.Called(_a0)

	var r0 []*solder.User
	if rf, ok := ret.Get(0).(func(solder.ModUserOptions) []*solder.User); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*solder.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(solder.ModUserOptions) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PackClientAppend provides a mock function with given fields: _a0
func (_m *ClientAPI) PackClientAppend(_a0 solder.PackClientOptions) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(solder.PackClientOptions) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PackClientDelete provides a mock function with given fields: _a0
func (_m *ClientAPI) PackClientDelete(_a0 solder.PackClientOptions) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(solder.PackClientOptions) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PackClientList provides a mock function with given fields: _a0
func (_m *ClientAPI) PackClientList(_a0 solder.PackClientOptions) ([]*solder.Client, error) {
	ret := _m.Called(_a0)

	var r0 []*solder.Client
	if rf, ok := ret.Get(0).(func(solder.PackClientOptions) []*solder.Client); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*solder.Client)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(solder.PackClientOptions) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PackDelete provides a mock function with given fields: _a0
func (_m *ClientAPI) PackDelete(_a0 string) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PackGet provides a mock function with given fields: _a0
func (_m *ClientAPI) PackGet(_a0 string) (*solder.Pack, error) {
	ret := _m.Called(_a0)

	var r0 *solder.Pack
	if rf, ok := ret.Get(0).(func(string) *solder.Pack); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*solder.Pack)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PackList provides a mock function with given fields:
func (_m *ClientAPI) PackList() ([]*solder.Pack, error) {
	ret := _m.Called()

	var r0 []*solder.Pack
	if rf, ok := ret.Get(0).(func() []*solder.Pack); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*solder.Pack)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PackPatch provides a mock function with given fields: _a0
func (_m *ClientAPI) PackPatch(_a0 *solder.Pack) (*solder.Pack, error) {
	ret := _m.Called(_a0)

	var r0 *solder.Pack
	if rf, ok := ret.Get(0).(func(*solder.Pack) *solder.Pack); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*solder.Pack)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*solder.Pack) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PackPost provides a mock function with given fields: _a0
func (_m *ClientAPI) PackPost(_a0 *solder.Pack) (*solder.Pack, error) {
	ret := _m.Called(_a0)

	var r0 *solder.Pack
	if rf, ok := ret.Get(0).(func(*solder.Pack) *solder.Pack); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*solder.Pack)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*solder.Pack) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PackUserAppend provides a mock function with given fields: _a0
func (_m *ClientAPI) PackUserAppend(_a0 solder.PackUserOptions) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(solder.PackUserOptions) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PackUserDelete provides a mock function with given fields: _a0
func (_m *ClientAPI) PackUserDelete(_a0 solder.PackUserOptions) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(solder.PackUserOptions) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PackUserList provides a mock function with given fields: _a0
func (_m *ClientAPI) PackUserList(_a0 solder.PackUserOptions) ([]*solder.User, error) {
	ret := _m.Called(_a0)

	var r0 []*solder.User
	if rf, ok := ret.Get(0).(func(solder.PackUserOptions) []*solder.User); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*solder.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(solder.PackUserOptions) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ProfileGet provides a mock function with given fields:
func (_m *ClientAPI) ProfileGet() (*solder.Profile, error) {
	ret := _m.Called()

	var r0 *solder.Profile
	if rf, ok := ret.Get(0).(func() *solder.Profile); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*solder.Profile)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ProfilePatch provides a mock function with given fields: _a0
func (_m *ClientAPI) ProfilePatch(_a0 *solder.Profile) (*solder.Profile, error) {
	ret := _m.Called(_a0)

	var r0 *solder.Profile
	if rf, ok := ret.Get(0).(func(*solder.Profile) *solder.Profile); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*solder.Profile)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*solder.Profile) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserDelete provides a mock function with given fields: _a0
func (_m *ClientAPI) UserDelete(_a0 string) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UserGet provides a mock function with given fields: _a0
func (_m *ClientAPI) UserGet(_a0 string) (*solder.User, error) {
	ret := _m.Called(_a0)

	var r0 *solder.User
	if rf, ok := ret.Get(0).(func(string) *solder.User); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*solder.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserList provides a mock function with given fields:
func (_m *ClientAPI) UserList() ([]*solder.User, error) {
	ret := _m.Called()

	var r0 []*solder.User
	if rf, ok := ret.Get(0).(func() []*solder.User); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*solder.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserModAppend provides a mock function with given fields: _a0
func (_m *ClientAPI) UserModAppend(_a0 solder.UserModOptions) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(solder.UserModOptions) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UserModDelete provides a mock function with given fields: _a0
func (_m *ClientAPI) UserModDelete(_a0 solder.UserModOptions) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(solder.UserModOptions) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UserModList provides a mock function with given fields: _a0
func (_m *ClientAPI) UserModList(_a0 solder.UserModOptions) ([]*solder.Mod, error) {
	ret := _m.Called(_a0)

	var r0 []*solder.Mod
	if rf, ok := ret.Get(0).(func(solder.UserModOptions) []*solder.Mod); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*solder.Mod)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(solder.UserModOptions) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserPackAppend provides a mock function with given fields: _a0
func (_m *ClientAPI) UserPackAppend(_a0 solder.UserPackOptions) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(solder.UserPackOptions) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UserPackDelete provides a mock function with given fields: _a0
func (_m *ClientAPI) UserPackDelete(_a0 solder.UserPackOptions) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(solder.UserPackOptions) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UserPackList provides a mock function with given fields: _a0
func (_m *ClientAPI) UserPackList(_a0 solder.UserPackOptions) ([]*solder.Pack, error) {
	ret := _m.Called(_a0)

	var r0 []*solder.Pack
	if rf, ok := ret.Get(0).(func(solder.UserPackOptions) []*solder.Pack); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*solder.Pack)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(solder.UserPackOptions) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserPatch provides a mock function with given fields: _a0
func (_m *ClientAPI) UserPatch(_a0 *solder.User) (*solder.User, error) {
	ret := _m.Called(_a0)

	var r0 *solder.User
	if rf, ok := ret.Get(0).(func(*solder.User) *solder.User); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*solder.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*solder.User) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserPost provides a mock function with given fields: _a0
func (_m *ClientAPI) UserPost(_a0 *solder.User) (*solder.User, error) {
	ret := _m.Called(_a0)

	var r0 *solder.User
	if rf, ok := ret.Get(0).(func(*solder.User) *solder.User); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*solder.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*solder.User) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// VersionBuildAppend provides a mock function with given fields: _a0
func (_m *ClientAPI) VersionBuildAppend(_a0 solder.VersionBuildOptions) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(solder.VersionBuildOptions) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// VersionBuildDelete provides a mock function with given fields: _a0
func (_m *ClientAPI) VersionBuildDelete(_a0 solder.VersionBuildOptions) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(solder.VersionBuildOptions) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// VersionBuildList provides a mock function with given fields: _a0
func (_m *ClientAPI) VersionBuildList(_a0 solder.VersionBuildOptions) ([]*solder.Build, error) {
	ret := _m.Called(_a0)

	var r0 []*solder.Build
	if rf, ok := ret.Get(0).(func(solder.VersionBuildOptions) []*solder.Build); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*solder.Build)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(solder.VersionBuildOptions) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// VersionDelete provides a mock function with given fields: _a0, _a1
func (_m *ClientAPI) VersionDelete(_a0 string, _a1 string) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// VersionGet provides a mock function with given fields: _a0, _a1
func (_m *ClientAPI) VersionGet(_a0 string, _a1 string) (*solder.Version, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *solder.Version
	if rf, ok := ret.Get(0).(func(string, string) *solder.Version); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*solder.Version)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// VersionList provides a mock function with given fields: _a0
func (_m *ClientAPI) VersionList(_a0 string) ([]*solder.Version, error) {
	ret := _m.Called(_a0)

	var r0 []*solder.Version
	if rf, ok := ret.Get(0).(func(string) []*solder.Version); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*solder.Version)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// VersionPatch provides a mock function with given fields: _a0, _a1
func (_m *ClientAPI) VersionPatch(_a0 string, _a1 *solder.Version) (*solder.Version, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *solder.Version
	if rf, ok := ret.Get(0).(func(string, *solder.Version) *solder.Version); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*solder.Version)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, *solder.Version) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// VersionPost provides a mock function with given fields: _a0, _a1
func (_m *ClientAPI) VersionPost(_a0 string, _a1 *solder.Version) (*solder.Version, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *solder.Version
	if rf, ok := ret.Get(0).(func(string, *solder.Version) *solder.Version); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*solder.Version)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, *solder.Version) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
