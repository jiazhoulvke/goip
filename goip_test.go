package goip

import (
	"os"
	"path/filepath"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFind(t *testing.T) {
	Convey("测试goip", t, func() {
		dbPath := filepath.Join(os.ExpandEnv("$GOPATH"), "src", "github.com", "jiazhoulvke", "goip", "17monipdb.dat")
		err := SetDBPath(dbPath)
		So(err, ShouldBeNil)
		for _, ip := range []string{"8.8.8.8", "114.114.114.114", "223.5.5.5", "192.168.1.1", "0.0.0.0", "127.0.0.1"} {
			location, err := Find(ip)
			So(err, ShouldBeNil)
			t.Log("location:", location)
		}
	})

}
