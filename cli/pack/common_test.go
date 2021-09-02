package pack

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tarantool/cartridge-cli/cli/context"
)

func TestGetPackageFullname(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	var ctx context.Ctx

	// TODO: internal error on bad type

	// w/o suffix
	ctx.Project.Name = "myapp"
	ctx.Pack.VersionRelease = "1.2.3-4"
	ctx.Pack.Suffix = ""

	ctx.Pack.Type = TgzType
	assert.Equal("myapp-1.2.3-4.tar.gz", getPackageFullname(&ctx))

	ctx.Pack.Type = RpmType
	assert.Equal("myapp-1.2.3-4.rpm", getPackageFullname(&ctx))

	ctx.Pack.Type = DebType
	assert.Equal("myapp-1.2.3-4.deb", getPackageFullname(&ctx))

	// w/ suffix
	ctx.Project.Name = "myapp"
	ctx.Pack.VersionRelease = "1.2.3-4"
	ctx.Pack.Suffix = "dev"

	ctx.Pack.Type = TgzType
	assert.Equal("myapp-1.2.3-4-dev.tar.gz", getPackageFullname(&ctx))

	ctx.Pack.Type = RpmType
	assert.Equal("myapp-1.2.3-4-dev.rpm", getPackageFullname(&ctx))
	ctx.Pack.Type = DebType
	assert.Equal("myapp-1.2.3-4-dev.deb", getPackageFullname(&ctx))
}

func TestGetImageTags(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	var ctx context.Ctx

	// TODO: internal error on bad type

	// VersionRelease
	ctx.Project.Name = "myapp"
	ctx.Pack.VersionRelease = "1.2.3-4"
	ctx.Pack.Suffix = ""
	ctx.Pack.ImageTags = []string{}

	assert.ElementsMatch([]string{"myapp:1.2.3-4"}, getImageTags(&ctx))

	// VersionRelease + Suffix
	ctx.Project.Name = "myapp"
	ctx.Pack.VersionRelease = "1.2.3-4"
	ctx.Pack.Suffix = "dev"
	ctx.Pack.ImageTags = []string{}

	assert.ElementsMatch([]string{"myapp:1.2.3-4-dev"}, getImageTags(&ctx))

	// ImageTags
	ctx.Project.Name = "myapp"
	ctx.Pack.VersionRelease = ""
	ctx.Pack.Suffix = ""
	ctx.Pack.ImageTags = []string{"my-first-image", "my-lovely-image"}

	assert.ElementsMatch([]string{"my-first-image", "my-lovely-image"}, getImageTags(&ctx))
}

func TestGetTarantoolMinVersion(t *testing.T) {
	assert := assert.New(t)

	expectedTarantoolVersions := map[string]map[string]string{
		"3.0.0-beta2": {RpmType: "3.0.0~beta2-1", DebType: "3.0.0~beta2-1"},
		"3.0.0-rc1":   {RpmType: "3.0.0~rc1-1", DebType: "3.0.0~rc1-1"},
		"3.0.0-rc2":   {RpmType: "3.0.0~rc2-1", DebType: "3.0.0~rc2-1"},

		"2.10.0-beta1-0-g7da4b1438": {RpmType: "2.10.0~beta1-1", DebType: "2.10.0~beta1-1"},
		"2.10.0-beta1":              {RpmType: "2.10.0~beta1-1", DebType: "2.10.0~beta1-1"},
		"2.10.0":                    {RpmType: "2.10.0-1", DebType: "2.10.0-1"},

		"2.9.0-alpha1-0-g7da4b1438": {RpmType: "2.9.0~alpha1-1", DebType: "2.9.0~alpha1-1"},
		"2.9.0-alpha1":              {RpmType: "2.9.0~alpha1-1", DebType: "2.9.0~alpha1-1"},
		"2.9.0":                     {RpmType: "2.9.0-1", DebType: "2.9.0-1"},

		"2.8.2-0-gfc96d10f5": {RpmType: "2.8.2", DebType: "2.8.2-0-gfc96d10f5"},
		"2.8.2":              {RpmType: "2.8.2", DebType: "2.8.2"},
	}

	for version, expectedByType := range expectedTarantoolVersions {
		for packType, expected := range expectedByType {
			tarantoolVersion, err := GetTarantoolMinVersion(version, packType)
			assert.Nil(err)
			assert.Equal(expected, tarantoolVersion)
		}
	}
}
