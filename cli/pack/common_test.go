package pack

import (
	"fmt"
	"runtime"
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
	assert.Equal(fmt.Sprintf("myapp-1.2.3-4.%s.rpm", runtime.GOARCH), getPackageFullname(&ctx))

	ctx.Pack.Type = DebType
	assert.Equal(fmt.Sprintf("myapp_1.2.3-4-1_%s.deb", runtime.GOARCH), getPackageFullname(&ctx))

	// w/ suffix
	ctx.Project.Name = "myapp"
	ctx.Pack.VersionRelease = "1.2.3-4"
	ctx.Pack.Suffix = "dev"

	ctx.Pack.Type = TgzType
	assert.Equal("myapp-1.2.3-4-dev.tar.gz", getPackageFullname(&ctx))

	ctx.Pack.Type = RpmType
	assert.Equal(fmt.Sprintf("myapp-1.2.3-4-dev.%s.rpm", runtime.GOARCH), getPackageFullname(&ctx))
	ctx.Pack.Type = DebType
	assert.Equal(fmt.Sprintf("myapp_1.2.3-4-1-dev_%s.deb", runtime.GOARCH), getPackageFullname(&ctx))
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
