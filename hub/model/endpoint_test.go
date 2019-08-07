package model

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestString(t *testing.T) {
	Convey("Get string expression of an endpoint", t, func() {
		se := ServiceEndpoint{
			IP:   "192.168.0.1",
			Port: 8000,
		}
		So(se.String(), ShouldEqual, "192.168.0.1:8000")
	})
}
