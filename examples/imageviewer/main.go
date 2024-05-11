package main

import (
	"fmt"
	"image"
	"image/draw"
	"log"
	"os"

	_ "image/jpeg"
	_ "image/png"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"

	"github.com/nfnt/resize"
	"github.com/pkg/errors"
	"github.com/rajveermalviya/go-wayland/examples/imageviewer/internal/swizzle"
	"github.com/rajveermalviya/go-wayland/examples/imageviewer/internal/tempfile"
	"github.com/rajveermalviya/go-wayland/wayland/client"
	xdg_shell "github.com/rajveermalviya/go-wayland/wayland/stable/xdg-shell"
	"golang.org/x/sys/unix"
)

// Global app state
type appState struct {
	appID         string
	title         string
	pImage        *image.RGBA
	width, height int32
	frame         *image.RGBA
	exit          bool
	skipDraw      bool

	display    *client.Display
	registry   *client.Registry
	shm        *client.Shm
	compositor *client.Compositor
	xdgWmBase  *xdg_shell.WmBase
	seat       *client.Seat

	surface     *client.Surface
	xdgSurface  *xdg_shell.Surface
	xdgTopLevel *xdg_shell.Toplevel

	keyboard *client.Keyboard
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) != 2 {
		return errors.Errorf("usage: %s file.jpg", os.Args[0])
	}

	fileName := os.Args[1]

	// Read the image file to *image.RGBA
	pImage, err := rgbaImageFromFile(fileName)
	if err != nil {
		return errors.WithStack(err)
	}

	// Resize again, for first frame
	frameImage := resize.Resize(0, uint(pImage.Rect.Dy()), pImage, resize.NearestNeighbor).(*image.RGBA)
	frameRect := frameImage.Bounds()

	app := &appState{
		// Set the title to `cat.jpg - imageviewer`
		title: fileName + " - imageviewer",
		appID: "osd_imageviewer",
		// Keep proxy image in cache, for use in resizing
		pImage: pImage,
		width:  int32(frameRect.Dx()),
		height: int32(frameRect.Dy()),
		frame:  frameImage,
	}

	if err := app.initWindow(); err != nil {
		return err
	}

	// Start the dispatch loop
	for !app.exit {
		app.dispatch()
	}

	app.cleanup()
	return nil
}

func (app *appState) initWindow() error {
	display, err := client.Connect("")
	if err != nil {
		return errors.Errorf("unable to connect to wayland server: %v", err)
	}
	app.display = display
	display.SetErrorHandler(app.HandleDisplayError)

	registry, err := app.display.GetRegistry()
	if err != nil {
		return errors.WithStack(err)
	}
	app.registry = registry

	registry.SetGlobalHandler(app.HandleRegistryGlobal)
	app.displayRoundTrip() // Wait for interfaces to register
	app.displayRoundTrip() // Wait for handler events

	// Create a wl_surface for toplevel window
	surface, err := app.compositor.CreateSurface()
	if err != nil {
		return errors.WithStack(err)
	}
	app.surface = surface

	// attach wl_surface to xdg_wmbase to get toplevel handle
	xdgSurface, err := app.xdgWmBase.GetXdgSurface(surface)
	if err != nil {
		return errors.WithStack(err)
	}
	app.xdgSurface = xdgSurface
	xdgSurface.SetConfigureHandler(app.HandleSurfaceConfigure)

	xdgTopLevel, err := xdgSurface.GetToplevel()
	if err != nil {
		return errors.WithStack(err)
	}
	app.xdgTopLevel = xdgTopLevel
	xdgTopLevel.SetConfigureHandler(app.HandleToplevelConfigure) // for window resize
	xdgTopLevel.SetCloseHandler(app.HandleToplevelClose)

	if err := xdgTopLevel.SetTitle(app.title); err != nil {
		return errors.WithStack(err)
	}
	if err := xdgTopLevel.SetAppId(app.appID); err != nil {
		return errors.WithStack(err)
	}
	if err := app.surface.Commit(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (app *appState) dispatch() {
	app.display.Context().Dispatch()
}

func (app *appState) context() *client.Context {
	return app.display.Context()
}

func (app *appState) HandleRegistryGlobal(e client.RegistryGlobalEvent) {
	switch e.Interface {
	case "wl_compositor":
		compositor := client.NewCompositor(app.context())
		err := app.registry.Bind(e.Name, e.Interface, e.Version, compositor)
		if err != nil {
			log.Fatalf("unable to bind wl_compositor interface: %v", err)
		}
		app.compositor = compositor
	case "wl_shm":
		shm := client.NewShm(app.context())
		err := app.registry.Bind(e.Name, e.Interface, e.Version, shm)
		if err != nil {
			log.Fatalf("unable to bind wl_shm interface: %v", err)
		}
		app.shm = shm
	case "xdg_wm_base":
		xdgWmBase := xdg_shell.NewWmBase(app.context())
		err := app.registry.Bind(e.Name, e.Interface, e.Version, xdgWmBase)
		if err != nil {
			log.Fatalf("unable to bind xdg_wm_base interface: %v", err)
		}
		app.xdgWmBase = xdgWmBase
		xdgWmBase.SetPingHandler(app.HandleWmBasePing)
	case "wl_seat":
		seat := client.NewSeat(app.context())
		err := app.registry.Bind(e.Name, e.Interface, e.Version, seat)
		if err != nil {
			log.Fatalf("unable to bind wl_seat interface: %v", err)
		}
		app.seat = seat
		seat.SetCapabilitiesHandler(app.HandleSeatCapabilities)
	}
}

func (app *appState) HandleSurfaceConfigure(e xdg_shell.SurfaceConfigureEvent) {
	if err := app.xdgSurface.AckConfigure(e.Serial); err != nil {
		log.Fatal("unable to ack xdg surface configure")
	}

	if app.skipDraw {
		return
	}

	buffer := app.drawFrame()
	if err := app.surface.Attach(buffer, 0, 0); err != nil {
		log.Fatalf("unable to attach buffer to surface: %v", err)
	}
	if err := app.surface.Commit(); err != nil {
		log.Fatalf("unable to commit surface state: %v", err)
	}
}

func (app *appState) HandleToplevelConfigure(e xdg_shell.ToplevelConfigureEvent) {
	width := e.Width
	height := e.Height

	if width == 0 || height == 0 {
		// Compositor is deferring to us
		return
	}

	if width == app.width && height == app.height {
		// No need to resize
		return
	}

	// Resize the proxy image to new frame size
	// and set it to frame image
	app.frame = resize.Resize(uint(width), uint(height), app.pImage, resize.Bilinear).(*image.RGBA)

	// Update app size
	app.width = width
	app.height = height
}

func (app *appState) drawFrame() *client.Buffer {
	logPrintln("drawing frame")
	app.skipDraw = true

	stride := app.width * 4
	size := stride * app.height

	file, err := tempfile.Create(int64(size))
	if err != nil {
		log.Fatalf("unable to create a temporary file: %v", err)
	}
	defer func() {
		if err2 := file.Close(); err2 != nil {
			logPrintf("unable to close file: %v", err2)
		}
	}()

	data, err := unix.Mmap(int(file.Fd()), 0, int(size), unix.PROT_READ|unix.PROT_WRITE, unix.MAP_SHARED)
	if err != nil {
		log.Fatalf("unable to create mapping: %v", err)
	}
	defer func() {
		if err2 := unix.Munmap(data); err2 != nil {
			logPrintf("unable to delete mapping: %v", err2)
		}
	}()

	pool, err := app.shm.CreatePool(int(file.Fd()), size)
	if err != nil {
		log.Fatalf("unable to create shm pool: %v", err)
	}
	defer func() {
		if err2 := pool.Destroy(); err2 != nil {
			logPrintf("unable to destroy shm pool: %v", err2)
		}
	}()

	buf, err := pool.CreateBuffer(0, app.width, app.height, stride, uint32(client.ShmFormatArgb8888))
	if err != nil {
		log.Fatalf("unable to create client.Buffer from shm pool: %v", err)
	}

	// Convert RGBA to BGRA
	copy(data, app.frame.Pix)
	swizzle.BGRA(data)

	buf.SetReleaseHandler(func(_ client.BufferReleaseEvent) {
		if err := buf.Destroy(); err != nil {
			logPrintf("unable to destroy buffer: %v", err)
		}
	})

	logPrintln("drawing frame complete")
	return buf
}

func (app *appState) HandleSeatCapabilities(e client.SeatCapabilitiesEvent) {
	haveKeyboard := (e.Capabilities * uint32(client.SeatCapabilityKeyboard)) != 0
	if haveKeyboard && app.keyboard == nil {
		app.attachKeyboard()
	} else if !haveKeyboard && app.keyboard != nil {
		app.releaseKeyboard()
	}
}

// HandleDisplayError handles client.Display errors
func (*appState) HandleDisplayError(e client.DisplayErrorEvent) {
	// Just log.Fatal for now
	log.Fatalf("display error event: %v", e)
}

func (app *appState) HandleWmBasePing(e xdg_shell.WmBasePingEvent) {
	app.xdgWmBase.Pong(e.Serial)
}

func (app *appState) HandleToplevelClose(xdg_shell.ToplevelCloseEvent) {
	app.exit = true
}

func (app *appState) displayRoundTrip() {
	callback, err := app.display.Sync()
	if err != nil {
		log.Fatalf("unable to get sync callback: %v", err)
	}
	defer func() {
		if err := callback.Destroy(); err != nil {
			logPrintln("unable to destroy callback:", err)
		}
	}()

	done := false
	callback.SetDoneHandler(func(client.CallbackDoneEvent) {
		done = true
	})

	// Wait for callback to return
	for !done {
		app.dispatch()
	}
}

func (app *appState) cleanup() {
	if app.keyboard != nil {
		app.releaseKeyboard()
	}

	if app.xdgTopLevel != nil {
		if err := app.xdgTopLevel.Destroy(); err != nil {
			logPrintln("unable to destroy xdg_toplevel:", err)
		}
		app.xdgTopLevel = nil
	}

	if app.xdgSurface != nil {
		if err := app.xdgSurface.Destroy(); err != nil {
			logPrintln("unable to destroy xdg_surface:", err)
		}
		app.xdgSurface = nil
	}

	if app.surface != nil {
		if err := app.surface.Destroy(); err != nil {
			logPrintln("unable to destroy wl_surface:", err)
		}
		app.surface = nil
	}

	// Release wl_seat handlers
	if app.seat != nil {
		if err := app.seat.Release(); err != nil {
			logPrintln("unable to destroy wl_seat:", err)
		}
		app.seat = nil
	}

	// Release xdg_wmbase
	if app.xdgWmBase != nil {
		if err := app.xdgWmBase.Destroy(); err != nil {
			logPrintln("unable to destroy xdg_wm_base:", err)
		}
		app.xdgWmBase = nil
	}

	if app.shm != nil {
		if err := app.shm.Destroy(); err != nil {
			logPrintln("unable to destroy wl_shm:", err)
		}
		app.shm = nil
	}

	if app.compositor != nil {
		if err := app.compositor.Destroy(); err != nil {
			logPrintln("unable to destroy wl_compositor:", err)
		}
		app.compositor = nil
	}

	if app.registry != nil {
		if err := app.registry.Destroy(); err != nil {
			logPrintln("unable to destroy wl_registry:", err)
		}
		app.registry = nil
	}

	if app.display != nil {
		if err := app.display.Destroy(); err != nil {
			logPrintln("unable to destroy wl_display:", err)
		}
	}

	if err := app.context().Close(); err != nil {
		logPrintln("unable to close wayland context:", err)
	}
}

func (app *appState) attachKeyboard() {
	keyboard, err := app.seat.GetKeyboard()
	if err != nil {
		log.Fatal("unable to register keyboard interface")
	}
	app.keyboard = keyboard
	keyboard.SetKeyHandler(app.HandleKeyboardKey)
}

func (app *appState) releaseKeyboard() {
	if err := app.keyboard.Release(); err != nil {
		logPrintln("unable to release keyboard interface")
	}
	app.keyboard = nil
}

func (app *appState) HandleKeyboardKey(e client.KeyboardKeyEvent) {
	// close on "esc"
	if e.Key == 1 {
		app.exit = true
	}
}

func rgbaImageFromFile(filePath string) (*image.RGBA, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	rgbaImage, ok := img.(*image.RGBA)
	if !ok {
		// Convert to RGBA if not already RGBA
		rect := img.Bounds()
		rgbaImage = image.NewRGBA(rect)
		draw.Draw(rgbaImage, rect, img, rect.Min, draw.Src)
	}

	return rgbaImage, nil
}

var logDisabled = os.Getenv("LOG_DISABLED") == "1"

func logPrintln(v ...interface{}) {
	if !logDisabled {
		log.Println(v...)
	}
}

func logPrintf(format string, v ...interface{}) {
	if !logDisabled {
		log.Printf(format, v...)
	}
}
