package media_test

import (
	"image"
	"testing"

	"github.com/pion/mediadevices/pkg/driver"
	"github.com/pion/mediadevices/pkg/io/video"
	"github.com/pion/mediadevices/pkg/prop"
	"go.viam.com/test"

	"github.com/edaniels/gostream/media"
)

func TestReaderClose(t *testing.T) {
	d := newFakeDriver("/dev/fake")

	vrc1 := media.NewVideoReadCloser(d, newFakeReader())
	vrc2 := media.NewVideoReadCloser(d, newFakeReader())

	if closedCount := d.(*fakeDriver).closedCount; closedCount != 0 {
		t.Fatalf("expected driver to be open, but was closed %d times", closedCount)
	}

	test.That(t, vrc1.Close(), test.ShouldHaveSameTypeAs, &media.DriverInUseError{})
	test.That(t, d.(*fakeDriver).closedCount, test.ShouldEqual, 0)

	test.That(t, vrc2.Close(), test.ShouldBeNil)
	test.That(t, d.(*fakeDriver).closedCount, test.ShouldEqual, 1)
}

// fakeDriver is a driver has a label and keeps track of how many times it is closed.
type fakeDriver struct {
	label       string
	closedCount int
}

func (d *fakeDriver) Open() error              { return nil }
func (d *fakeDriver) Properties() []prop.Media { return []prop.Media{} }
func (d *fakeDriver) ID() string               { return d.label }
func (d *fakeDriver) Info() driver.Info        { return driver.Info{Label: d.label} }
func (d *fakeDriver) Status() driver.State     { return "FakeState" }

func (d *fakeDriver) Close() error {
	d.closedCount++
	return nil
}

func newFakeDriver(label string) driver.Driver {
	return &fakeDriver{label: label}
}

// fakeReader is a reader that always returns a pixel-sized canvas.
type fakeReader struct{}

func (r *fakeReader) Read() (img image.Image, release func(), err error) {
	return image.NewNRGBA(image.Rect(0, 0, 1, 1)), func() {}, nil
}

func newFakeReader() video.Reader {
	return &fakeReader{}
}
